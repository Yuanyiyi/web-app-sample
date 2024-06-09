package grpc_server

import (
	"sync"

	"github.com/web-app-sample/pkg/utils/msg_queue"
)

const (
	GetResourceEvent       EventType = "GetResourceEvent"
	CreateRobotEvent       EventType = "CreateRobotEvent"
	DestroyRobotEvent      EventType = "DestroyRobotEvent"
	GetRobotListEvent      EventType = "GetRobotListEvent"
	GetAgentRobotListEvent EventType = "GetAgentRobotListEvent"
	AgentOffline           EventType = "AgentOffline"
	ProtoMsg               EventType = "ProtoMsg"
	MsgQueue               EventType = "MsgQueue"
)

type EventType string

type ServerEvent struct {
	Type      EventType
	AgentName string
	Data      any
}

type RobotInfo struct {
	Total   int32
	Running int32
}

var (
	publish        = msg_queue.GetNewMsgQueue()
	agentStreamMgr = sync.Map{}
)
