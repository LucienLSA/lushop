package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"ordersrv/global"
	"ordersrv/model"
	proto_goods "ordersrv/proto/gen/goods"
	proto_inventory "ordersrv/proto/gen/inventory"
	proto_order "ordersrv/proto/gen/order"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderServer struct {
	proto_order.UnimplementedOrderServer
}

// 购物车
func (s *OrderServer) CartItemList(ctx context.Context, req *proto_order.UserInfo) (*proto_order.CartItemListResponse, error) {
	// 获取用户的购物车列表
	var shopCarts []model.ShoppingCart
	var rsp proto_order.CartItemListResponse
	result := global.DB.Where(&model.ShoppingCart{User: req.Id}).Find(&shopCarts)
	if result.Error != nil {
		return nil, result.Error
	} else {
		rsp.Total = int32(result.RowsAffected)
	}
	for _, shopCart := range shopCarts {
		rsp.Data = append(rsp.Data, &proto_order.ShopCartInfoResponse{
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
func (s *OrderServer) CreateCartItem(ctx context.Context, req *proto_order.CartItemRequest) (*proto_order.ShopCartInfoResponse, error) {
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
	global.DB.Save(&shopCart)
	return &proto_order.ShopCartInfoResponse{Id: shopCart.ID}, nil
}
func (s *OrderServer) UpdateCartItem(ctx context.Context, req *proto_order.CartItemRequest) (*emptypb.Empty, error) {
	var shopCart model.ShoppingCart
	result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).First(&shopCart, req.Id)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	shopCart.Checked = req.Checked
	if req.Nums > 0 {
		shopCart.Nums = req.Nums
	}
	global.DB.Save(&shopCart)
	return &emptypb.Empty{}, nil
}

func (s *OrderServer) DeleteCartItem(ctx context.Context, req *proto_order.CartItemRequest) (*emptypb.Empty, error) {
	result := global.DB.Where("goods=? and user=?", req.GoodsId, req.UserId).Delete(&model.ShoppingCart{})
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "购物车记录不存在")
	}
	return &emptypb.Empty{}, nil
}

type OrderListener struct {
	Code        codes.Code
	Detail      string
	ID          int32
	OrderAmount float32
}

func (o *OrderListener) ExecuteLocalTransaction(msg *primitive.Message) primitive.LocalTransactionState {
	var orderInfo model.OrderInfo
	_ = json.Unmarshal(msg.Body, &orderInfo)

	var goodsIds []int32
	var shopCarts []model.ShoppingCart
	goodsNumsMap := make(map[int32]int32)
	result := global.DB.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Find(&shopCarts)
	if result.RowsAffected == 0 {
		// return nil, status.Errorf(codes.InvalidArgument, "未选中结算的商品")
		o.Code = codes.InvalidArgument
		o.Detail = "未选中结算的商品"
		return primitive.RollbackMessageState
	}

	for _, shopCart := range shopCarts {
		goodsIds = append(goodsIds, shopCart.Goods)
		goodsNumsMap[shopCart.Goods] = shopCart.Nums
	}

	// 跨服务调用 商品
	goods, err := global.GoodsSrvClient.BatchGetGoods(context.Background(), &proto_goods.BatchGoodsIdInfo{
		Id: goodsIds,
	})
	if err != nil {
		o.Code = codes.Internal
		o.Detail = "批量查询商品信息失败"
		// return nil, status.Errorf(codes.Internal, "批量查询商品信息失败")
		return primitive.RollbackMessageState
	}
	var orderAmount float32
	var orderGoods []*model.OrderGoods
	var goodsInvInfo []*proto_inventory.GoodsInvInfo
	for _, good := range goods.Data {
		orderAmount += good.ShopPrice * float32(goodsNumsMap[good.Id])
		orderGoods = append(orderGoods, &model.OrderGoods{
			Goods:      good.Id,
			GoodsName:  good.Name,
			GoodsImage: good.GoodsFrontImage,
			GoodsPrice: good.ShopPrice,
			Nums:       goodsNumsMap[good.Id],
		})
		goodsInvInfo = append(goodsInvInfo, &proto_inventory.GoodsInvInfo{
			GoodsId: good.Id,
			Num:     goodsNumsMap[good.Id],
		})
	}
	// 跨服务 库存扣减
	_, err = global.InventorySrvClient.Sell(context.Background(),
		&proto_inventory.SellInfo{OrderSn: orderInfo.OrderSn, GoodsInfo: goodsInvInfo})
	if err != nil {
		// 当遇到网络问题时，如何避免误判，对于库存服务中的sell逻辑，判断返回的状态码信息，确定是什么原因
		// 遇到该状态码时，则返回Commit
		o.Code = codes.ResourceExhausted
		o.Detail = "扣减库存失败"
		return primitive.RollbackMessageState
		// return nil, status.Errorf(codes.ResourceExhausted, "扣减库存失败")
	}

	// 生产订单表
	// 20250406xxx时间戳方式生成订单号
	tx := global.DB.Begin()
	orderInfo.OrderMount = orderAmount
	result = tx.Save(&orderInfo)
	if result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "创建订单失败"
		return primitive.CommitMessageState
		// return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	o.OrderAmount = orderAmount
	o.ID = orderInfo.ID
	for _, orderGood := range orderGoods {
		orderGood.Order = orderInfo.ID
	}
	// 批量插入orderGoods
	if result := tx.CreateInBatches(orderGoods, 100); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "批量插入订单商品失败"
		return primitive.CommitMessageState
		// return nil, status.Errorf(codes.Internal, "创建订单失败")
	}
	if result := tx.Where(&model.ShoppingCart{User: orderInfo.User, Checked: true}).Delete(&model.ShoppingCart{}); result.RowsAffected == 0 {
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "删除购物车记录失败"
		return primitive.CommitMessageState
		// return nil, status.Errorf(codes.Internal, "创建订单失败")
	}

	// 提交事务之前发送延迟消息
	p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.226.140:9876"}))
	if err != nil {
		panic(err)
	}
	if err = p.Start(); err != nil {
		panic(err)
	}
	msg = primitive.NewMessage("order_timeout", msg.Body)
	msg.WithDelayTimeLevel(16)
	_, err = p.SendSync(context.Background(), msg)
	if err != nil {
		zap.S().Errorf("发送延迟消息失败:%s\n", err)
		tx.Rollback()
		o.Code = codes.Internal
		o.Detail = "发送延迟消息失败"
		return primitive.CommitMessageState
	}
	// if err = p.Shutdown(); err != nil {
	// 	panic(err)
	// }
	// 提交事务
	tx.Commit()
	o.Code = codes.OK
	return primitive.RollbackMessageState
}

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
func (s *OrderServer) CreateOrder(ctx context.Context, req *proto_order.OrderRequest) (*proto_order.OrderInfoResponse, error) {

	orderListener := &OrderListener{}
	p, err := rocketmq.NewTransactionProducer(
		&OrderListener{},
		producer.WithNameServer([]string{"192.168.226.140:9876"}),
	)
	if err != nil {
		zap.S().Errorf("生成producer失败: %s", err.Error())
		return nil, err
	}
	if err = p.Start(); err != nil {
		zap.S().Errorf("启动producer失败: %s", err.Error())
		return nil, err
	}

	order := model.OrderInfo{
		OrderSn:      GenerateOrderSn(req.UserId),
		Address:      req.Address,
		SignerName:   req.Name,
		SingerMobile: req.Mobile,
		Post:         req.Post,
		User:         req.UserId,
	}

	jsonString, _ := json.Marshal(order)

	_, err = p.SendMessageInTransaction(context.Background(),
		primitive.NewMessage("order_reback", jsonString))
	if err != nil {
		fmt.Printf("发送失败:%s\n", err)
		return nil, status.Error(codes.Internal, "发送消息失败")
	}
	if orderListener.Code != codes.OK {
		return nil, status.Error(orderListener.Code, orderListener.Detail)
	}
	return &proto_order.OrderInfoResponse{
		Id:      orderListener.ID,
		OrderSn: order.OrderSn,
		Total:   orderListener.OrderAmount,
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
		if order.Status != "TRADE_SUCCESS" {
			tx := global.DB.Begin()
			// 归还库存 ， 发送一个普通消息到 order_reback的 topic中
			// 修改订单的状态为已结束
			order.Status = "TRADE_CLOSED"
			tx.Save(&order)
			p, err := rocketmq.NewProducer(producer.WithNameServer([]string{"192.168.226.140:9876"}))
			if err != nil {
				panic(err)
			}
			if err = p.Start(); err != nil {
				panic(err)
			}
			_, err = p.SendSync(context.Background(), primitive.NewMessage("order_reback", msgs[i].Body))
			if err != nil {
				tx.Rollback()
				fmt.Printf("发送失败:%s\n", err)
				return consumer.ConsumeRetryLater, nil
			}
			// if err = p.Shutdown(); err != nil {
			// 	panic(err)
			// }
		}
	}
	return consumer.ConsumeSuccess, nil
}

// 订单列表
func (s *OrderServer) OrderList(ctx context.Context, req *proto_order.OrderFilterRequest) (*proto_order.OrderListResponse, error) {
	var orders []model.OrderInfo
	var rsp proto_order.OrderListResponse
	var total int64
	global.DB.Where(&model.OrderInfo{User: req.UserId}).Count(&total)
	rsp.Total = int32(total)
	// 分页
	global.DB.Scopes(Paginate(int(req.Pages), int(req.PagePerNums))).Where(&model.OrderInfo{User: req.UserId}).Find(&orders)
	for _, order := range orders {
		rsp.Data = append(rsp.Data, &proto_order.OrderInfoResponse{
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
func (s *OrderServer) OrderDetail(ctx context.Context, req *proto_order.OrderRequest) (*proto_order.OrderInfoDetailResponse, error) {
	var order model.OrderInfo
	var rsp proto_order.OrderInfoDetailResponse
	result := global.DB.Where("id=? AND user=?", req.Id, req.UserId).First(&order)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	orderInfo := proto_order.OrderInfoResponse{}
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
		rsp.Goods = append(rsp.Goods, &proto_order.OrderItemResponse{
			GoodsId:    orderGoods.Goods,
			GoodsName:  orderGoods.GoodsName,
			GoodsImage: orderGoods.GoodsImage,
			GoodsPrice: orderGoods.GoodsPrice,
			Nums:       orderGoods.Nums,
		})
	}
	return &rsp, nil
}
func (s *OrderServer) UpdateOrderStatus(ctx context.Context, req *proto_order.OrderStatus) (*emptypb.Empty, error) {
	if result := global.DB.Model(&model.OrderInfo{}).
		Where("order_sn = ?", req.OrderSn).
		Update("status", req.Status); result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "订单不存在")
	}
	return &emptypb.Empty{}, nil
}
