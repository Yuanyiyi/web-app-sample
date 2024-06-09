package sentinel

import (
	"log"

	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/sirupsen/logrus"

	sentinel "github.com/alibaba/sentinel-golang/api"
)

var isInit bool = false

func init() {
	initSentinelSdk()
}

func initSentinelSdk() {
	if isInit {
		return
	}
	isInit = true
	err := sentinel.InitDefault()
	if err != nil {
		log.Fatal(err)
	}
}

func LimitingPolicy(resourceName string, threshold float64) {
	if threshold < 1 {
		threshold = 20
	}
	// 限制 请求创建ds的频率
	_, err := flow.LoadRules([]*flow.Rule{
		{
			Resource:               resourceName,
			Threshold:              threshold,
			TokenCalculateStrategy: flow.Direct,
			ControlBehavior:        flow.Reject,
		},
	})
	if err != nil {
		logrus.Warnf("Failed to configure current limiting policy, error: %s", err.Error())
		return
	}
	// 熔断规则, 10s、大于50个请求 最大rt 100ms 窗口3s
	_, err = circuitbreaker.LoadRules([]*circuitbreaker.Rule{
		{
			Resource:         resourceName,
			Strategy:         circuitbreaker.SlowRequestRatio,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 50,
			StatIntervalMs:   10000,
			MaxAllowedRtMs:   100,
			Threshold:        0.3,
		},
	})
	if err != nil {
		logrus.Warnf("Circuit breaker policy configuration failed, error: %s", err.Error())
		return
	}
}
