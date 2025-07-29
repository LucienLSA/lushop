package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"ordersrv/global"
	"ordersrv/model"
	v2goodsproto "ordersrv/proto/goods"
	v2inventoryproto "ordersrv/proto/inventory"
	v2orderproto "ordersrv/proto/order"
	"time"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	v2orderproto.UnimplementedOrderServer
}

var Tracer = otel.Tracer(global.ServerConfig.JaegerInfo.TracerName)

// 获取用户的购物车列表
func (s *OrderServer) CartItemList(ctx context.Context, req *v2orderproto.UserInfo) (*v2orderproto.CartItemListResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "CartItemList",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	var shopCarts []model.ShoppingCart
	var rsp v2orderproto.CartItemListResponse
	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &v2orderproto.ShopCartInfoResponse{
			Id:      shopCart.ID,
			UserId:  shopCart.ID,
			GoodsId: shopCart.Goods,
			Nums:    shopCart.Nums,
			Checked: shopCart.Checked,
		})
	}
	return &rsp, nil
}

// 将商品添加到购物车
// 购物车中原本没有该商品。购物车中存在该商品
func (s *OrderServer) CreateCartItem(ctx context.Context, req *v2orderproto.CartItemRequest) (*v2orderproto.ShopCartInfoResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "CreateCartItem",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	var shopCart model.ShoppingCart
	result := global.DB.Where(&model.ShoppingCart{Goods: req.GoodsId, User: req.UserId}).First(&shopCart)
	if result.RowsAffected == 1 {
		// 如果记录存在，就合并购物车记录
		shopCart.Nums += req.Nums
	} else {
		// 插入操作
		shopCart.User = req.UserId
		shopCart.Goods = req.GoodsId
		shopCart.Nums = req.Nums
		shopCart.Checked = false
	}
	if result := global.DB.Save(&shopCart); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "新建购物车记录失败")
	}
	return &v2orderproto.ShopCartInfoResponse{Id: shopCart.ID}, nil
}

// 更新购物车记录，更新数量和选中状态
func (s *OrderServer) UpdateCartItem(ctx context.Context, req *v2orderproto.CartItemRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "UpdateCartItem",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	var shopCart model.ShoppingCart
	// result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart, req.Id)
	result := global.DB.First(&shopCart, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	// 验证记录是否属于当前用户
	if shopCart.User != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "无权操作此购物车记录")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	if result := global.DB.Save(&shopCart); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "更新购物车记录失败")
	}
	return &emptypb.Empty{}, nil
}

// 删除购物车清单
func (s *OrderServer) DeleteCartItem(ctx context.Context, req *v2orderproto.CartItemRequest) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "DeleteCartItem",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

// 订单业务
// 建立订单表用于后续消息传递
type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
}

// 实例
var orderInfo = OrderListener{}

// 用于trace上下文传递的消息结构体
// 注意：model.OrderInfo需可序列化
// 实现了 trace 上下文的跨进程传递：
// handler 层用 otel.GetTextMapPropagator().Inject 把 trace 上下文写入消息体。
// 事务监听器用 otel.GetTextMapPropagator().Extract 恢复上下文，再创建父 span 和子 span。
// 这样 Jaeger UI 里能看到 handler 和事务监听器的 trace 在同一条链路下，trace 串联完整。
type OrderMessage struct {
	OrderInfo model.OrderInfo   `json:"order_info"`
	TraceMap  map[string]string `json:"trace_map"`
}

// 执行本地事务的监听器
func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	// 解析消息体，恢复trace上下文
	var orderMsg OrderMessage
	_ = json.Unmarshal(msg.Body, &orderMsg)
	parentCtx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.MapCarrier(orderMsg.TraceMap))

	ctx, parentSpan := Tracer.Start(parentCtx, "OrderLocalTransaction",
		trace.WithAttributes(
			attribute.String("order_sn", orderMsg.OrderInfo.OrderSn),
			attribute.Int64("user_id", int64(orderMsg.OrderInfo.User)),
		),
	)
	defer parentSpan.End()

	// 购物车清单检查
	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	result := global.DB.Where(&model.ShoppingCart{User: orderMsg.OrderInfo.User, Checked: true}).Find(&shopCarts)
	if result.RowsAffected == 0 {
		o.Code = codes.InvalidArgument
		o.Detail = "未选中结算的商品"
		return primitive.RollbackMessageState
	}

	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 子span
	goodsSpanCtx, goodsSpan := Tracer.Start(ctx, "BatchGetGoods")
	// 批量商品查询
	goods, err := global.GoodsSrvClient.BatchGetGoods(goodsSpanCtx, &v2goodsproto.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	goodsSpan.End()
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		return primitive.RollbackMessageState
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*v2inventoryproto.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &v2inventoryproto.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}

	// 子Span
	invSpanCtx, invSpan := Tracer.Start(ctx, "InventorySell")
	// 库存扣减
	_, err = global.InventorySrvClient.Sell(invSpanCtx,
		&v2inventoryproto.SellInfo{OrderSn: orderMsg.OrderInfo.OrderSn, GoodsInfo: goodsInvInfo})
	invSpan.End()
	if err != nil {
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
	}
	// 事务保存订单信息
	tx := global.DB.Begin()
	orderMsg.OrderInfo.OrderMount = orderAmount
	result = tx.Save(&orderMsg.OrderInfo)
	if result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
	}

	o.OrderAmount = orderAmount
	o.ID = orderMsg.OrderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderMsg.OrderInfo.ID
	}
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
	}
	if result := tx.Where(&model.ShoppingCart{User: orderMsg.OrderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
	}

	// p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.226.140:9876"}))
	// if err != nil {
	// 	panic(err)
	// }
	// if err = p.Start(); err != nil {
	// 	panic(err)
	// }

	// 生产者发出延时消息order_timeout
	msg = primitive.NewMessage(global.ServerConfig.RocketMQConfig.TopicTimeOut, msg.Body)
	msg.WithDelayTimeLevel(3)
	_, err = global.MQSendClient.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延迟消息失败:%s\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延迟消息失败"
		return primitive.CommitMessageState
	}
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

// 事务消息回查
func (o *OrderListener) CheckLocalTransaction(msg *primitive.MessageExt) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)
	result := global.DB.Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).First(&orderInfo)
	if result.RowsAffected == 0 {
		// 这里并不能确定库存是否扣减，所以需要在库存服务中保证幂等性，扣减库存成功
		return primitive.CommitMessageState
	}
	// 查询到了库存，逻辑就不用归还
	return primitive.RollbackMessageState
}

// 新建订单
func (s *OrderServer) CreateOrder(ctx context.Context, req *v2orderproto.OrderRequest) (*v2orderproto.OrderInfoResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "CreateOrder",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	// 提取trace上下文
	traceMap := make(map[string]string)
	otel.GetTextMapPropagator().Inject(ctx, propagation.MapCarrier(traceMap))

	// 封装消息体
	msgBody, _ := json.Marshal(OrderMessage{
		OrderInfo: order,
		TraceMap:  traceMap,
	})

	// 生产者发送一个事务消息 库存回归reback的 topic中
	_, err := global.MQSendTranClient.SendMessageInTransaction(context.Background(),
		primitive.NewMessage(global.ServerConfig.RocketMQConfig.TopicReback, msgBody))
	if err != nil {
		fmt.Printf("发送失败:%s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}
	if orderInfo.Code != codes.OK {
		return nil, status.Error(orderInfo.Code, orderInfo.Detail)
	}
	return &v2orderproto.OrderInfoResponse{
		// Id:      orderInfo.ID,
		OrderSn: order.OrderSn,
		Total:   orderInfo.OrderAmount,
	}, nil
}

// 订单超时
func OrderTimeout(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	for i := range msgs {
		var orderInfo model.OrderInfo
		_ = json.Unmarshal(msgs[i].Body, &orderInfo)
		fmt.Printf("获取到订单的超时消息:%v", time.Now())
		// 查询订单的支付状态
		// 如果已支付什么都不做，如果未支付，应该归还库存
		var order model.OrderInfo
		result := global.DB.Model(model.OrderInfo{}).
			Where(model.OrderInfo{OrderSn: orderInfo.OrderSn}).
			First(&order)
		if result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		// 订单超时，自动归还库存
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			// p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.226.140:9876"}))
			// if err != nil {
			// 	panic(err)
			// }
			// if err = p.Start(); err != nil {
			// 	panic(err)
			// }
			// 归还库存，生产者发送一个普通消息到 reback的 topic中
			// 通知库存服务把这笔订单占用的库存释放回去
			_, err := global.MQSendClient.SendSync(context.Background(),
				primitive.NewMessage(global.ServerConfig.RocketMQConfig.TopicReback, msgs[i].Body))
			if err != nil {
				tx.Rollback()
				zap.S().Errorf("【超时归还】发送失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			// if err = p.Shutdown(); err != nil {
			// 	panic(err)
			// }
			// 修改订单的状态为已支付
			order.Status = "TRADE_CLOSED"
			if result := tx.Save(&order); result.Error != nil {
				tx.Rollback()
				zap.S().Errorf("【超时归还】修改支付失败: %s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
		}
	}
	return consumer.ConsumeSuccess, nil
}

// 订单列表
func (s *OrderServer) OrderList(ctx context.Context, req *v2orderproto.OrderFilterRequest) (*v2orderproto.OrderListResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "OrderList",
		trace.WithAttributes(
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	var orders []model.OrderInfo
	var rsp v2orderproto.OrderListResponse
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)
	// 分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &v2orderproto.OrderInfoResponse{
			UserId:  order.User,
			Id:      order.ID,
			OrderSn: order.OrderSn,
			PayType: order.PayType,
			Status:  order.Status,
			Post:    order.Post,
			Total:   order.OrderMount,
			Address: order.Address,
			Name:    order.SignerName,
			Mobile:  order.SingerMobile,
			AddTime: order.CreatedAt.Format("2006-01-02 15:03:05"),
		})
	}
	return &rsp, nil
}

// 获取订单详情
func (s *OrderServer) OrderDetail(ctx context.Context, req *v2orderproto.OrderRequest) (*v2orderproto.OrderInfoDetailResponse, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "OrderDetail",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()

	var order model.OrderInfo
	var rsp v2orderproto.OrderInfoDetailResponse
	result := global.DB.Where("id=? AND user=?", req.Id, req.UserId).First(&order)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo := v2orderproto.OrderInfoResponse{}
	orderInfo.Id = order.ID
	orderInfo.UserId = order.User
	orderInfo.OrderSn = order.OrderSn
	orderInfo.PayType = order.PayType
	orderInfo.Status = order.Status
	orderInfo.Post = order.Post
	orderInfo.Mobile = order.SingerMobile
	orderInfo.Name = order.SignerName
	orderInfo.Total = order.OrderMount
	orderInfo.Address = order.Address
	rsp.OrderInfo = &orderInfo
	var orderGoods []model.OrderGoods
	result = global.DB.Where(&model.OrderGoods{Order: order.ID}).Find(&orderGoods)
	if result.Error != nil {
		return nil, result.Error
	}
	for _, orderGoods := range orderGoods {
		rsp.Goods = append(rsp.Goods, &v2orderproto.OrderItemResponse{
			GoodsId:    orderGoods.Goods,
			GoodsName:  orderGoods.GoodsName,
			GoodsImage: orderGoods.GoodsImage,
			GoodsPrice: orderGoods.GoodsPrice,
			Nums:       orderGoods.Nums,
		})
	}
	return &rsp, nil
}

// 更新订单状态
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *v2orderproto.OrderStatus) (*emptypb.Empty, error) {
	md, _ := metadata.FromIncomingContext(ctx)
	_, span := Tracer.Start(ctx, "UpdateOrderStatus",
		trace.WithAttributes(
			attribute.Int64("id", int64(req.GetId())),
			attribute.StringSlice("client-id", md.Get("client-id")),
			attribute.StringSlice("user-id", md.Get("user-id")),
		),
	)
	defer span.End()
	if result := global.DB.Model(&model.OrderInfo{}).
		Where("order_sn = ?", req.OrderSn).
		Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}
