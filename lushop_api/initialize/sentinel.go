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
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               global.ServerConfig.SentinelInfo.App.Name,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
			Threshold:              6,
			StatIntervalInMs:       global.ServerConfig.SentinelInfo.Stat.GlobalStatisticIntervalMsTotal,
		},
	})
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		// Statistic time span=5s, recoveryTimeout=3s, maxErrorCount=50
		{
			Resource:         global.ServerConfig.SentinelInfo.App.Name,
			Strategy:         circuitbreaker.ErrorRatio,
			RetryTimeoutMs:   3000, // 3s之后尝试恢复
			MinRequestAmount: 10,   // 静默数
			StatIntervalMs:   global.ServerConfig.SentinelInfo.Stat.GlobalStatisticIntervalMsTotal,
			Threshold:        0.4,
		},
	})
	if err != nil {
		otelzap.S().Fatalf("加载规则失败:%v", err)
	}
}
