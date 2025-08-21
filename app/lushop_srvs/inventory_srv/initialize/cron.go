package initialize

import (
	"inventorysrv/global"
	"inventorysrv/handler"
	"inventorysrv/model"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

func InitCron() {
	c := cron.New(cron.WithSeconds())
	// 每小时执行一次库存兜底检查
	_, err := c.AddFunc("0 0 */1 * * *", func() {
		var timeoutOrders []model.StockSellDetail
		// 查询1小时前创建且未归还的订单
		global.DB.Where("status = 1 AND created_at < ?", time.Now().Add(-1*time.Hour)).
			Find(&timeoutOrders)

		zap.S().Infof("发现%d个超时未支付订单，开始处理", len(timeoutOrders))
		for _, order := range timeoutOrders {
			if err := handler.ProcessOrderReback(order.OrderSn); err != nil {
				zap.S().Errorf("订单%s处理失败: %v", order.OrderSn, err)
				order.Status = 3 // 标记为超时未支付
				global.DB.Save(&order)
			}
		}
	})

	if err != nil {
		zap.S().Panic("定时任务初始化失败", err)
	}
	c.Start()
	zap.S().Info("定时任务服务已启动")
}
