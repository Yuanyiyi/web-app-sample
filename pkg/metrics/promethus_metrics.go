package metrics

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	HotfixManagerSubsystem = "hotfix_manager"
	TaskCount              = "task_count"
	TaskExecTimeout        = "task_exec_timeout"
	TaskCompleted          = "task_completed"
	TaskExecuting          = "task_executing"
	TaskNoStarted          = "task_no_started"
	TaskFail               = "task_failed"
)

var (
	RecordMap = CollectorMap{Record: sync.Map{}}

	// 定义
	//GameServerStateCount = prometheus.NewGaugeVec(
	//	prometheus.GaugeOpts{
	//		Name: "gameservers_state_count",
	//		Help: "The number of gameservers per state",
	//	},
	//	[]string{"state"},
	//)

	TaskTaskCountCnt   = newGaugeFunc(TaskCount, "task request count")
	TaskExecTimeoutCnt = newGaugeFunc(TaskExecTimeout, "the task exec timeout")
	TaskCompletedCnt   = newGaugeFunc(TaskCompleted, "task completed count")
	TaskExecutingCnt   = newGaugeFunc(TaskExecuting, "task executing count")
	TaskNoStartedCnt   = newGaugeFunc(TaskNoStarted, "task no_started count")
	TaskFailCnt        = newGaugeFunc(TaskFail, "task failed count")

	CommonCollectors = []prometheus.Collector{
		TaskTaskCountCnt,
		TaskExecTimeoutCnt,
		TaskCompletedCnt,
		TaskExecutingCnt,
		TaskNoStartedCnt,
		TaskFailCnt,
	}
)

func init() {
	// 注册
	//prometheus.MustRegister(GameServerStateCount)
	prometheus.MustRegister(CommonCollectors...)
}

// 对每种类型进行计数，每上报一次刷新一次指标
func newGaugeFunc(name, describe string) prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Subsystem: HotfixManagerSubsystem,
		Name:      name,
		Help:      describe,
	}, func() float64 {
		return float64(RecordMap.GetAndResetValue(name))
	})
}

// 定时统计的指标
func newGaugeFuncCount(name, describe string) prometheus.GaugeFunc {
	return prometheus.NewGaugeFunc(prometheus.GaugeOpts{
		Subsystem: HotfixManagerSubsystem,
		Name:      name,
		Help:      describe,
	}, func() float64 {
		return float64(RecordMap.GetValue(name))
	})
}
