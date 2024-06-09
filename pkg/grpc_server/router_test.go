package grpc_server

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/known/anypb"

	"github.com/web-app-sample/pkg/database/mysql/models"
	"github.com/web-app-sample/pkg/pb/web_app_sample"
)

var (
	robotIds = []string{"robot_id_1"}
	agentIds = []string{"agent_id_1"}
)

type RouterUtilTestSuite struct {
	models.ModelsTestSuite
	father    models.ModelsTestSuite
	routerMap *RouterMapInterface
}

func TestRouterUtilTestSuite(t *testing.T) {
	suite.Run(t, new(RouterUtilTestSuite))
}

func (t *RouterUtilTestSuite) SetupSuite() {
	t.father.SetupSuite()
	t.routerMap = NewRouterMapInterface()
}

func (t *RouterUtilTestSuite) steamEvent(data interface{}, event *ServerEvent) {
	var err error
	steamEvent := &web_app_sample.StreamActiveEvent{}
	switch data.(type) {
	case *web_app_sample.DestroyRobotRsp:
		steamEvent.AEvent, err = anypb.New(data.(*web_app_sample.DestroyRobotRsp))
	case *web_app_sample.CreateRobotRsp:
		steamEvent.AEvent, err = anypb.New(data.(*web_app_sample.CreateRobotRsp))
	case *web_app_sample.GetResourceRsp:
		steamEvent.AEvent, err = anypb.New(data.(*web_app_sample.GetResourceRsp))
	case *web_app_sample.GetRobotListRsp:
		steamEvent.AEvent, err = anypb.New(data.(*web_app_sample.GetRobotListRsp))
	}
	t.Nil(err)
	m, err := steamEvent.AEvent.UnmarshalNew()
	t.Nil(err)
	event.Data = m
}
func (t *RouterUtilTestSuite) insert() {
}
func (t *RouterUtilTestSuite) delete() {
}

func (t *RouterUtilTestSuite) TestProtoMsg() {
	event := &ServerEvent{AgentName: agentIds[0]}
	t.insert()
	defer t.delete()
	// protoMsg
	createRobotRsp := &web_app_sample.CreateRobotRsp{
		RobotTypeId: robotIds[0],
		RobotNumber: 10,
		State:       web_app_sample.Status_RUNNING,
	}
	t.steamEvent(createRobotRsp, event)
	t.routerMap.RouterMap(event, ProtoMsg)

	destroyRsp := &web_app_sample.DestroyRobotRsp{
		RobotRsp: &web_app_sample.RobotInfo{
			RobotTypeId: robotIds[0],
			RobotNumber: 4,
		},
	}
	t.steamEvent(destroyRsp, event)
	t.routerMap.RouterMap(event, ProtoMsg)

	getResourceRsp := &web_app_sample.GetResourceRsp{
		AutoAgent:       agentIds[0],
		TotalCpu:        16000,
		TotalMemory:     32000,
		TotalPids:       4096,
		AvailableCpu:    14000,
		AvailableMemory: 30000,
		AvailablePids:   4000,
		RobotInfo: []*web_app_sample.RobotInfo{
			{
				RobotTypeId: robotIds[0],
				RobotNumber: 3,
				State:       web_app_sample.Status_RUNNING,
			},
			{
				RobotTypeId: robotIds[0],
				RobotNumber: 7,
				State:       web_app_sample.Status_CREATING,
			},
		},
	}
	t.steamEvent(getResourceRsp, event)
	t.routerMap.RouterMap(event, ProtoMsg)

	getRobotRsp := &web_app_sample.GetRobotListRsp{
		RobotInfo: []*web_app_sample.RobotInfo{
			{
				RobotTypeId: robotIds[0],
				RobotNumber: 7,
				State:       web_app_sample.Status_RUNNING,
			},
			{
				RobotTypeId: robotIds[0],
				RobotNumber: 3,
				State:       web_app_sample.Status_CREATING,
			},
		},
	}
	t.steamEvent(getRobotRsp, event)
	t.routerMap.RouterMap(event, ProtoMsg)
}

func (t *RouterUtilTestSuite) TestQueueMsg() {
	agentStream := NewAgentStream(agentIds[0], agentIds[0])
	agentStreamMgr.Store(agentIds[0], agentStream)
	t.insert()
	defer t.delete()

	event := &ServerEvent{AgentName: agentIds[0]}

	event.Type = GetResourceEvent
	t.routerMap.RouterMap(event, MsgQueue)
	req := <-agentStream.StreamServicePassiveEventChan
	getResourceReq := &web_app_sample.GetResourceReq{}
	err := req.PEvent.UnmarshalTo(getResourceReq)
	t.Nil(err)

	event.Type = GetRobotListEvent
	t.routerMap.RouterMap(event, MsgQueue)
	req = <-agentStream.StreamServicePassiveEventChan
	getRobotsReq := &web_app_sample.GetRobotListReq{}
	err = req.PEvent.UnmarshalTo(getRobotsReq)
	t.Nil(err)

	event.Type = CreateRobotEvent
	createReq := &web_app_sample.CreateRobotReq{
		AutoAgent: agentIds[0],
		RobotReq: &web_app_sample.CreateRobotReq_RobotReq{
			RoomUrl:     "room.addr.com",
			GameId:      "game_id_1",
			DsVersion:   "v1.1.1",
			SceneId:     "scene_id_1",
			GameVersion: "v0.0.0",
			RobotTypeId: robotIds[0],
			RobotNumber: int32(10),
		},
	}
	event.Data = createReq
	t.routerMap.RouterMap(event, MsgQueue)
	req = <-agentStream.StreamServicePassiveEventChan
	createRobotsReq := &web_app_sample.CreateRobotReq{}
	err = req.PEvent.UnmarshalTo(createRobotsReq)
	t.Nil(err)

	event.Type = DestroyRobotEvent
	destroyReq := &web_app_sample.DestroyRobotReq{
		AutoAgent: agentIds[0],
		RobotReq: &web_app_sample.RobotInfo{
			RobotTypeId: robotIds[0],
			RobotNumber: 1,
		},
	}
	event.Data = destroyReq
	t.routerMap.RouterMap(event, MsgQueue)
	req = <-agentStream.StreamServicePassiveEventChan
	destroyRobotsReq := &web_app_sample.DestroyRobotReq{}
	err = req.PEvent.UnmarshalTo(destroyRobotsReq)
	t.Nil(err)
}
