package initialize

import (
	"lushopapi/global"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
)

func Sentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		otelzap.S().Fatalf("初始化sentinel 失败:%v", err)
	}
	// 限流 当请求量超过阈值时，Sentinel 会立即拒绝多余请求，避免服务过载。
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               global.ServerConfig.SentinelInfo.App.Name,
			TokenCalculateStrategy: flow.Direct, // 直接统计请求量（非冷启动模式）
			ControlBehavior:        flow.Reject, // Reject 表示超阈值时直接拒绝请求（快速失败）
			Threshold:              6,           // 每秒允许的最大请求数为 6（超过则触发限流）
			StatIntervalInMs: global.ServerConfig.SentinelInfo.
				Stat.GlobalStatisticIntervalMsTotal, // 统计时间窗口（从配置读取，单位毫秒）。
		},
	})
	if err != nil {
		otelzap.S().Fatalf("加载限流规则失败:%v", err)
	}
	// 熔断降级
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         global.ServerConfig.SentinelInfo.App.Name,
			Strategy:         circuitbreaker.ErrorRatio, // ErrorRatio 表示基于 错误比例 触发熔断。
			RetryTimeoutMs:   3000,                      // 熔断后 3 秒 尝试恢复（半开状态探测）。
			MinRequestAmount: 10,                        // 静默数
			StatIntervalMs: global.ServerConfig.SentinelInfo.
				Stat.GlobalStatisticIntervalMsTotal, // 统计时间窗口
			Threshold: 0.4, // 错误比例阈值为 0.4（即 40% 请求失败时熔断）。
		},
	})
	if err != nil {
		otelzap.S().Fatalf("加载熔断降级规则失败:%v", err)
	}
}
