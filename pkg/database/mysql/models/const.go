package models

const (
	// 未开始、执行中、已终止、已完成
	NoStarted State = "no-started"
	Waiting   State = "waiting"
	Executing State = "executing"
	Running   State = "running"
	Completed State = "completed"
	Failed    State = "failed"
	WaitToEnd State = "wait_to_end"
	Stop      State = "stop"
	Creating  State = "creating"
	Deleting  State = "deleting"
	Deleted   State = "Deleted"
)

type State string
