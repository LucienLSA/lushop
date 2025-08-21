package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"inventorysrv/global"
	"inventorysrv/model"
	proto "inventorysrv/proto"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"

	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
)

type InventoryServer struct {
	proto.UnimplementedInventoryServer
	// db *gorm.DB
}

var _ proto.InventoryServer = &InventoryServer{}

// func NewInventoryServer(db *gorm.DB) *InventoryServer {
// 	return &InventoryServer{db: db}
// }

// 设置库存
func (v *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	// invDB := v.db.WithContext(ctx)
	// invDB := initialize.NewDBClient(ctx)
	//只有是主键的情况才能直接用id   goodsid不是主键
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	if result := global.DB.Save(&inv); result.Error != nil {
		return nil, status.Errorf(codes.Internal, "设置库存失败")
	}
	return &emptypb.Empty{}, nil
}

// 库存详情查询
func (v *InventoryServer) InvDetail(ctx context.Context, req *proto.GoodsInvInfo) (*proto.GoodsInvInfo, error) {
	var inv model.Inventory
	// invDB := NewInventoryDB(ctx)
	result := global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	if result.RowsAffected == 0 {
		return nil, status.Errorf(codes.NotFound, "库存信息不存在")
	}
	return &proto.GoodsInvInfo{
		GoodsId: inv.Goods,
		Num:     inv.Stocks,
	}, nil
}

// 扣减库存
// 存在并发问题，本地事务
// var m sync.Mutex 手动加入锁
// func (v *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
// 	tx := global.DB.Begin()
// 	// 通过下面的tx事务操作，默认不会自动提交
// 	for _, goodInfo := range req.GoodsInfo {
// 		var inv model.Inventory
// 		// if result := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
// 		// 	tx.Rollback() // 回滚之前的操作
// 		// 	return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
// 		// }
// 		for {
// 			if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
// 				tx.Rollback() // 回滚之前的操作
// 				return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
// 			}
// 			// 判断库存是否充足
// 			if inv.Stocks < goodInfo.Num {
// 				tx.Rollback()
// 				return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
// 			}
// 			// 扣减
// 			inv.Stocks -= goodInfo.Num
// 			result := tx.Model(&model.Inventory{}).Select("Stocks", "Version").
// 				Where("goods = ? and version=?", goodInfo.GoodsId, inv.Version).
// 				Updates(model.Inventory{Stocks: inv.Stocks, Version: inv.Version + 1})
// 			if result.RowsAffected == 0 {
// 				zap.S().Info("库存扣减失败，乐观锁冲突，重试中")
// 			} else {
// 				break
// 			}
// 		}

//			// tx.Save(&inv)
//		}
//		if err := tx.Commit().Error; err != nil {
//			return nil, status.Errorf(codes.Internal, "提交事务失败")
//		} // 必须手动提交
//		return &emptypb.Empty{}, nil
//	}
//
// 使用go-redsync 分布式锁包 实现库存扣减
func (v *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// client := goredislib.NewClient(&goredislib.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// 	// Addr: "localhost:6379",
	// })
	// 创建 Redis 连接池和 RedSync 分布式锁实例
	pool := goredis.NewPool(global.Rdb)
	// 支持多个 Redis 实例（避免单点故障）
	// 配置多个 Redis 实例连接池
	// pools := []redsync.Pool{
	//     goredis.NewPool(global.Rdb1),
	//     goredis.NewPool(global.Rdb2),
	//     goredis.NewPool(global.Rdb3),
	// }

	// // 创建 RedSync 实例
	// rs := redsync.New(pools...)
	rs := redsync.New(pool)
	// 开始数据库事务
	tx := global.DB.Begin()
	// 记录订单号和状态
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var details []model.GoodsDetail
	// 将请求中的商品信息转换为 GoodsDetail 列表
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})
		// 处理每个商品
		var inv model.Inventory
		// 为每个商品创建分布式锁（基于商品ID）
		mutex := rs.NewMutex(
			fmt.Sprintf("goods_%d", goodInfo.GoodsId),
			redsync.WithExpiry(5*time.Second),            // 设置锁的TTL为5秒
			redsync.WithTries(3),                         // 设置最大重试次数为3次
			redsync.WithRetryDelay(100*time.Millisecond), // 设置重试间隔为100毫秒
		)
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}
		// 确保锁最终会被释放
		defer func() {
			if ok, err := mutex.Unlock(); !ok || err != nil {
				zap.S().Errorw("释放redis分布式锁异常: %v", err)
			}
		}()
		// 获取锁后检查库存是否存在
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减库存并保存，会出现数据不一致
		// inv.Stocks -= goodInfo.Num
		// tx.Save(&inv)
		result := tx.Model(&model.Inventory{}).
			Where("goods = ? AND stocks >= ?", goodInfo.GoodsId, goodInfo.Num).
			Update("stocks", gorm.Expr("stocks - ?", goodInfo.Num))

		if result.Error != nil {
			tx.Rollback()
			return nil, status.Errorf(codes.Internal, "库存扣减失败: %v", result.Error)
		}

	}
	// 保存库存扣减记录
	sellDetail.Detail = details
	// 写sellDetail表
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "保存库存扣减历史记录失败")
	}
	// 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	} // 必须手动提交
	return &emptypb.Empty{}, nil
}

// TCC 方案
// 尝试扣减库存，将库存预扣减 冻结库存，预留资源，防止超卖。
func (v *InventoryServer) TrySell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// client := goredislib.NewClient(&goredislib.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// 	// Addr: "localhost:6379",
	// })
	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		// inv.Stocks -= goodInfo.Num
		inv.Freeze += goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	} // 必须手动提交
	return &emptypb.Empty{}, nil
}

// 订单确认，正式扣减库存并解冻。
func (v *InventoryServer) ComfirmSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// client := goredislib.NewClient(&goredislib.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// 	// Addr: "localhost:6379",
	// })
	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Stocks -= goodInfo.Num
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	} // 必须手动提交
	return &emptypb.Empty{}, nil
}

// 订单取消，解冻冻结库存，库存回滚。
func (v *InventoryServer) CancelSell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// client := goredislib.NewClient(&goredislib.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// 	// Addr: "localhost:6379",
	// })
	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.InventoryNew
		mutex := rs.NewMutex(fmt.Sprintf("goods_%d", goodInfo.GoodsId))
		if err := mutex.Lock(); err != nil {
			return nil, status.Errorf(codes.Internal, "获取redis分布式锁异常")
		}

		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 判断库存是否充足
		if inv.Stocks < goodInfo.Num {
			tx.Rollback()
			return nil, status.Errorf(codes.ResourceExhausted, "库存不足")
		}
		// 扣减
		inv.Freeze -= goodInfo.Num
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	} // 必须手动提交
	return &emptypb.Empty{}, nil
}

// 库存归还
// 1. 订单的超时归还
// 2. 订单创建失败，归还之前扣减的库存
// 3. 手动归还
// 这里的归还方案废除，由下面的AutoReback重构
func (v *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 增加库存数量
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 必须手动提交
	return &emptypb.Empty{}, nil
}

// 新增公共库存归还函数
func ProcessOrderReback(orderSn string) error {
	// 独立开启事务
	tx := global.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// 将inv的库存加回去，将selldetail的status设置为2，在事务中执行
	// 通过 StockSellDetail 表（库存扣减明细表）来判断该订单是否已经归还过库存。
	// 只处理 Status=1（已扣减未归还）的订单，防止重复归还。
	var sellDetail model.StockSellDetail
	if result := tx.Where("order_sn = ? AND status = 1", orderSn).First(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil // 无需处理的订单
	}

	// 批量更新库存
	// 如果查询到逐个归还库存
	for _, item := range sellDetail.Detail {
		// 先查询Inventory表，使用update会有锁冲突，当多个并发进入mysql会自动锁住
		if err := tx.Model(&model.Inventory{}).
			Where("goods = ?", item.Goods).
			Update("stocks", gorm.Expr("stocks + ?", item.Num)).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	// 将该订单的扣减明细状态设为2（已归还）
	// 如果更新失败，回滚事务
	if err := tx.Model(&sellDetail).Update("status", 2).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// 自动归还库存，放在consumer监听库存中
// 幂等性：通过扣减明细表和状态字段，确保同一订单不会被重复归还库存。
// 事务性：归还库存和状态更新在同一事务内，保证数据一致性。
// 自动重试：遇到数据库等临时问题时，返回 ConsumeRetryLater，消息会自动重试，保证最终归还成功。
// 适用场景：订单取消、超时未支付等需要自动归还库存的场景。
func AutoReback(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range me {
		//既然是归还库存，那么我应该具体的知道每件商品应该归还多少,但是有一个问题是什么?重复归还的问题
		//所以说这个接口应该确保幂等性,你不能因为消息的重复发送导致一个订单的库存归还多次,没有扣减的库存你别归还
		//如何确保这些都没有问题，新建一张表，这张表记录了详细的订单扣减细节，以及归还细节
		var orderInfo OrderInfo
		err := json.Unmarshal(me[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败:%v\n", me[i].Body)
			//根据业务来   订单号都解析失败了，感觉是错误的信息
			//ConsumeRetryLater 保证下次还能执行
			//ConsumeSuccess 丢弃
			return consumer.ConsumeSuccess, nil
		}

		// 直接调用事务内聚的归还函数
		err = ProcessOrderReback(orderInfo.OrderSn)
		if err != nil {
			zap.S().Errorf("订单%s处理失败: %v", orderInfo.OrderSn, err)
			return consumer.ConsumeRetryLater, err
		}
		zap.S().Infof("订单%s库存归还成功", orderInfo.OrderSn)
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
