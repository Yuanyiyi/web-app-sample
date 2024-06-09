package grpc_server

import (
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/web-app-sample/pkg/database/mysql/db_connection"
	"github.com/web-app-sample/pkg/database/mysql/models"
	"github.com/web-app-sample/pkg/pb/web_app_sample"
	"github.com/web-app-sample/pkg/utils/msg_queue"
	"github.com/web-app-sample/pkg/utils/startenv"
)

type RouterMapInterface struct {
	taskRecords *models.HotfixTaskGormDB
}

func NewRouterMapInterface() *RouterMapInterface {
	db := db_connection.GetGormDB()
	return &RouterMapInterface{
		taskRecords: models.NewHotfixTaskGormDB(db),
	}
}

func (r *RouterMapInterface) RouterMap(event *ServerEvent, eventType EventType) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("RouterMap panic error: %s", err1)
		}
	}()
	if event == nil {
		return
	}
	logrus.WithFields(logrus.Fields{"msg": event}).Debug("RouterMap")
	switch eventType {
	case ProtoMsg:
		err = r.RouterMapOfProtoMessage(event)
	case MsgQueue:
		err = r.RouterMapOfMsgQueue(event)
	default:
		logrus.WithFields(logrus.Fields{
			"msg": event,
		}).Info("routerMap default")
	}
	return
}

func (r *RouterMapInterface) RouterMapOfProtoMessage(event *ServerEvent) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("RouterMapOfProtoMessage panic error: %s", err1)
		}
	}()
	//logrus.WithFields(logrus.Fields{"msg": event.Data}).Info("RouterMapOfProtoMessage")
	logrus.WithField("msg", event.Data).Info("RouterMapOfProtoMessage")
	if event.Data == nil {
		return nil
	}
	pMsg := event.Data.(proto.Message)
	switch msg := pMsg.(type) {
	case *web_app_sample.CreateRobotRsp:
		event.Type = CreateRobotEvent
		err = r.CreateRobotRspHandle(event, msg)
	case *web_app_sample.DestroyRobotRsp:
		event.Type = DestroyRobotEvent
		err = r.DestroyRobotRspHandle(event, msg)
	case *web_app_sample.GetResourceRsp:
		event.Type = GetResourceEvent
		err = r.GetResourceRspHandle(event, msg)
	case *web_app_sample.GetRobotListRsp:
		event.Type = GetRobotListEvent
		err = r.GetRobotListRspHandle(event, msg)
	default:
		logrus.WithFields(logrus.Fields{
			"msg": pMsg,
		}).Info("RouterMapOfProtoMessage default")
	}
	publish.Publish(&msg_queue.MsgStruct{
		MsgType: msg_queue.GrpcServerMessage,
		MsgData: event,
	})
	return
}

func (r *RouterMapInterface) RouterMapOfMsgQueue(event *ServerEvent) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("RouterMapOfMsgQueue panic error: %s", err1)
		}
	}()
	logrus.WithFields(logrus.Fields{"msg": event.Data}).Debug("RouterMapOfMsgQueue")
	switch event.Type {
	case GetResourceEvent:
		r.GetAllAgentResourceReqHandle()
	case GetRobotListEvent:
		r.GetAllAgentRobotListReqHandle()
	case CreateRobotEvent:
		err = r.CreateRobotReqHandle(event)
	case DestroyRobotEvent:
		err = r.DestroyRobotReqHandle(event)
	default:
		logrus.WithFields(logrus.Fields{
			"msg": event,
		}).Info("RouterMapOfMsgQueue default")
	}
	return
}

func (r *RouterMapInterface) CreateRobotRspHandle(event *ServerEvent, rsp *web_app_sample.CreateRobotRsp) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("CreateRobotRspHandle panic error: %s", err1)
		}
	}()
	logrus.WithFields(logrus.Fields{"msg": rsp}).Debug("CreateRobotRspHandle")
	return nil
}

func (r *RouterMapInterface) DestroyRobotRspHandle(event *ServerEvent, rsp *web_app_sample.DestroyRobotRsp) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("DestroyRobotRspHandle panic error: %s", err1)
		}
	}()
	logrus.WithFields(logrus.Fields{"msg": rsp}).Debug("DestroyRobotRspHandle")
	return
}

func (r *RouterMapInterface) GetResourceRspHandle(event *ServerEvent, rsp *web_app_sample.GetResourceRsp) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("GetResourceRspHandle panic error: %s", err1)
		}
	}()
	logrus.WithFields(logrus.Fields{"msg": rsp}).Debug("GetResourceRspHandle")
	return
}

func (r *RouterMapInterface) GetRobotListRspHandle(event *ServerEvent, rsp *web_app_sample.GetRobotListRsp) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("GetRobotListRspHandle panic error: %s", err1)
		}
	}()

	return
}

func (r *RouterMapInterface) GetAllAgentResourceReqHandle() {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("GetAllAgentResourceReqHandle panic error: %s", err1)
		}
	}()
	logrus.WithFields(logrus.Fields{}).Debug("GetAllAgentResourceReqHandle")
	agentStreamMgr.Range(func(key, value any) bool {
		// 获取agent所在机器上的总资源
		pevent, _ := anypb.New(&web_app_sample.GetResourceReq{AutoAgent: key.(string)})
		agentStream := value.(*AgentStream)
		agentStream.SendStreamPassiveEvent(&web_app_sample.StreamPassiveEvent{PEvent: pevent})
		return true
	})
}

func (r *RouterMapInterface) GetAllAgentRobotListReqHandle() {
	logrus.WithFields(logrus.Fields{}).Debug("GetAllAgentRobotListReqHandle")
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("GetAllAgentRobotListReqHandle panic error: %s", err1)
		}
	}()
	agentStreamMgr.Range(func(key, value any) bool {
		// 获取agent上模拟的客户端数据
		pevent, _ := anypb.New(&web_app_sample.GetRobotListReq{AutoAgent: key.(string)})
		agentStream := value.(*AgentStream)
		agentStream.SendStreamPassiveEvent(&web_app_sample.StreamPassiveEvent{PEvent: pevent})

		return true
	})
}

func (r *RouterMapInterface) CreateRobotReqHandle(event *ServerEvent) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("CreateRobotReqHandle panic error: %s", err1)
		}
	}()
	return
}

func (r *RouterMapInterface) DestroyRobotReqHandle(event *ServerEvent) (err error) {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("DestroyRobotReqHandle panic error: %s", err1)
		}
	}()
	return
}

func IsTest() bool {
	if startenv.GetEnvironment() == "auto-test" {
		return true
	}
	return false
}
