package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"inventorysrv/global"
	"inventorysrv/model"
	"inventorysrv/proto"

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

// func NewInventoryServer(db *gorm.DB) *InventoryServer {
// 	return &InventoryServer{db: db}
// }

// 设置库存
func (v *InventoryServer) SetInv(ctx context.Context, req *proto.GoodsInvInfo) (*emptypb.Empty, error) {
	var inv model.Inventory
	// invDB := v.db.WithContext(ctx)
	// invDB := initialize.NewDBClient(ctx)
	global.DB.Where(&model.Inventory{Goods: req.GoodsId}).First(&inv)
	inv.Goods = req.GoodsId
	inv.Stocks = req.Num
	global.DB.Save(&inv)
	return &emptypb.Empty{}, nil
}

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
// var m sync.Mutex
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
func (v *InventoryServer) Sell(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	// client := goredislib.NewClient(&goredislib.Options{
	// 	Addr: fmt.Sprintf("%s:%s", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	// 	// Addr: "localhost:6379",
	// })
	pool := goredis.NewPool(global.Rdb)
	rs := redsync.New(pool)
	tx := global.DB.Begin()
	sellDetail := model.StockSellDetail{
		OrderSn: req.OrderSn,
		Status:  1,
	}
	var details []model.GoodsDetail
	for _, goodInfo := range req.GoodsInfo {
		details = append(details, model.GoodsDetail{
			Goods: goodInfo.GoodsId,
			Num:   goodInfo.Num,
		})
		var inv model.Inventory
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
		tx.Save(&inv)
		if ok, err := mutex.Unlock(); !ok || err != nil {
			return nil, status.Errorf(codes.Internal, "释放redis分布式锁异常")
		}
	}
	sellDetail.Detail = details
	// 写sellDetail表
	if result := tx.Create(&sellDetail); result.RowsAffected == 0 {
		tx.Rollback()
		return nil, status.Error(codes.Internal, "保存库存扣减历史记录失败")
	}
	if err := tx.Commit().Error; err != nil {
		return nil, status.Errorf(codes.Internal, "提交事务失败")
	} // 必须手动提交
	return &emptypb.Empty{}, nil
}

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
func (v *InventoryServer) Reback(ctx context.Context, req *proto.SellInfo) (*emptypb.Empty, error) {
	tx := global.DB.Begin()
	for _, goodInfo := range req.GoodsInfo {
		var inv model.Inventory
		if result := global.DB.Where(&model.Inventory{Goods: goodInfo.GoodsId}).First(&inv); result.RowsAffected == 0 {
			tx.Rollback() // 回滚之前的操作
			return nil, status.Errorf(codes.InvalidArgument, "无库存信息")
		}
		// 扣减
		inv.Stocks += goodInfo.Num
		tx.Save(&inv)
	}
	tx.Commit() // 必须手动提交
	return &emptypb.Empty{}, nil
}

// 自动归还库存，放在consumer监听库存中
func AutoReback(ctx context.Context, me ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
	type OrderInfo struct {
		OrderSn string
	}
	for i := range me {
		var orderInfo OrderInfo
		err := json.Unmarshal(me[i].Body, &orderInfo)
		if err != nil {
			zap.S().Errorf("解析json失败:%v\n", me[i].Body)
			return consumer.ConsumeSuccess, nil
		}
		// 将inv的库存加回去，将selldetail的status设置为2，在事务中执行
		tx := global.DB.Begin()
		var sellDetail model.StockSellDetail
		if result := tx.Model(&model.StockSellDetail{}).Where(&model.
			StockSellDetail{OrderSn: orderInfo.OrderSn, Status: 1}).
			First(&sellDetail); result.RowsAffected == 0 {
			return consumer.ConsumeSuccess, nil
		}
		// 如果查询到逐个归还库存
		for _, orderGood := range sellDetail.Detail {
			// 先查询Inventory表，但是使用update会有锁冲突，并发情况下
			result := tx.Model(&model.Inventory{}).Where(&model.Inventory{Goods: orderGood.Goods}).
				Update("stocks", gorm.Expr("stock+?", orderGood.Num))
			if result.RowsAffected == 0 {
				tx.Rollback()
				return consumer.ConsumeRetryLater, nil
			}
		}
		sellDetail.Status = 2
		result := tx.Model(&model.StockSellDetail{}).Where(&model.StockSellDetail{OrderSn: orderInfo.OrderSn}).
			Update("status", 2)
		if result.RowsAffected == 0 {
			tx.Rollback()
			return consumer.ConsumeRetryLater, nil
		}
		tx.Commit()
		return consumer.ConsumeSuccess, nil
	}
	return consumer.ConsumeSuccess, nil
}
