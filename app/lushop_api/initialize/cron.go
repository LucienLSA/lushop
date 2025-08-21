package initialize

import (
	"context"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"lushopapi/global"
)

// InitCron 初始化定时任务
func InitCron() {
	c := cron.New(cron.WithSeconds())

	// 添加黑名单清理任务 (每天凌晨2点执行)
	_, err := c.AddFunc("0 0 2 * * ?", cleanBlacklist)
	if err != nil {
		zap.S().Errorf("添加黑名单清理任务失败: %v", err)
		return
	}

	c.Start()
	zap.S().Info("定时任务服务启动成功")
}

// cleanBlacklist 清理过期的JWT黑名单
func cleanBlacklist() {
	ctx := context.Background()
	cursor := uint64(0)
	count := 100 // 每次扫描数量
	deleted := 0

	for {
		// 扫描黑名单键
		keys, newCursor, err := global.RedisClient.Scan(ctx, cursor, "jwt_blacklist:*", int64(count)).Result()
		if err != nil {
			zap.S().Errorf("扫描黑名单失败: %v", err)
			break
		}

		// 批量删除过期键
		if len(keys) > 0 {
			// 检查每个键的过期时间
			for _, key := range keys {
				ttl, err := global.RedisClient.TTL(ctx, key).Result()
				if err != nil || ttl <= 0 {
					// 已过期或获取TTL失败，直接删除
					if delErr := global.RedisClient.Del(ctx, key).Err(); delErr == nil {
						deleted++
					}
				}
			}
		}
		cursor = newCursor
		if cursor == 0 {
			break // 扫描完成
		}
	}
	zap.S().Infof("黑名单清理完成，共删除过期键 %d 个", deleted)
}
