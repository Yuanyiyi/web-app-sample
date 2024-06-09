package grpc_server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"github.com/web-app-sample/pkg/pb/web_app_sample"
	pb "github.com/web-app-sample/pkg/pb/web_app_sample"
	"github.com/web-app-sample/pkg/utils/common"
	"github.com/web-app-sample/pkg/utils/msg_queue"
	"github.com/web-app-sample/pkg/utils/startenv"
)

type Server struct {
	ctx context.Context
	web_app_sample.UnimplementedAutoAgentServer
	cb             RouterMapHandle
	autoManagerMsg chan interface{}
}
type RouterMapHandle func(*ServerEvent, EventType) error

func NewServer(ctx context.Context) *Server {
	AutoManagerMsg := publish.SubscribeTopic(func(v *msg_queue.MsgStruct) bool {
		if v != nil && v.MsgType == msg_queue.AutoManager {
			return true
		}
		return false
	})
	routerMap := NewRouterMapInterface()
	server := &Server{
		UnimplementedAutoAgentServer: web_app_sample.UnimplementedAutoAgentServer{},
		autoManagerMsg:               AutoManagerMsg,
		cb:                           routerMap.RouterMap,
		ctx:                          ctx,
	}
	go server.consumeFromMsqQueue()
	return server
}

func (s *Server) StreamService(stream web_app_sample.AutoAgent_StreamServiceServer) error {
	defer func() {
		if err1 := recover(); err1 != nil {
			logrus.WithFields(logrus.Fields{}).
				Warnf("StreamService panic error: %s", err1)
		}
	}()
	streamMD, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("peer.FromContext error")
	}
	if ok && len(streamMD.Get(common.MD_KEY_AGENTID)) == 0 {
		return errors.New("metadata name error")
	}
	agentName := streamMD.Get(common.MD_KEY_AGENTID)[0]
	logrus.WithFields(logrus.Fields{
		"agentName": agentName,
	}).Info("StreamService init ok")

	agentStream := NewAgentStream(agentName, agentName)
	agentStreamMgr.Store(agentName, agentStream)

	wg := sync.WaitGroup{}
	wg.Add(2)
	go s.RecvStreamService(wg, agentName, stream)
	go s.SendStreamService(wg, agentStream, stream)
	wg.Wait()
	return nil
}

func (s *Server) SendStreamService(wg sync.WaitGroup, agentStream *AgentStream, stream web_app_sample.AutoAgent_StreamServiceServer) {
	defer func() {
		logrus.Info("=====StreamService send close===")
		wg.Done()
	}()
	for {
		select {
		case req := <-agentStream.StreamServicePassiveEventChan:
			stream.Send(req)
		case <-s.ctx.Done():
			return
			//default:
			//	logrus.WithFields(logrus.Fields{}).Debug("SendStreamService default")
		}
	}
}

func (s *Server) RecvStreamService(wg sync.WaitGroup, agentName string, stream web_app_sample.AutoAgent_StreamServiceServer) {
	defer func() {
		logrus.Info("=====StreamService recv close===")
		wg.Done()

		logrus.WithFields(logrus.Fields{}).Infof("agentId: %s StreamService quit", agentName)
		agentStreamMgr.Delete(agentName)
		publish.Publish(&msg_queue.MsgStruct{
			MsgType: msg_queue.GrpcServerMessage,
			MsgData: &ServerEvent{AgentName: agentName, Type: AgentOffline},
		})
	}()
	logrus.WithFields(logrus.Fields{"agentName": agentName}).Info("RecvStreamService")
	for {
		recvEvent, err := stream.Recv()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"error": err,
				"data":  recvEvent,
			}).Info("StreamService error")
			break
		}
		logrus.WithFields(logrus.Fields{"msg": recvEvent}).Debug("RecvStreamService")
		m, err := recvEvent.AEvent.UnmarshalNew()
		if err != nil {
			continue
		}
		event := &ServerEvent{AgentName: agentName, Data: m}
		if s.cb != nil {
			s.cb(event, ProtoMsg)
		}
	}
}

func (s *Server) consumeFromMsqQueue() {
	for {
		select {
		case msg := <-s.autoManagerMsg:
			event := msg_queue.GetMsg(msg)
			if event == nil {
				continue
			}
			go func() {
				defer func() {
					if err := recover(); err != nil {
						logrus.WithFields(logrus.Fields{
							"error": err,
						}).Info("consumeFromMsqQueue error")
					}
				}()
				if s.cb != nil {
					s.cb(event.(*ServerEvent), MsgQueue)
				}
			}()
		case <-s.ctx.Done():
			return
		}
	}
}

func (s *Server) Close() {
	logrus.WithFields(logrus.Fields{}).Info("GrpcServer::Close")
	close(s.autoManagerMsg)
	agentStreamMgr.Range(func(key, value any) bool {
		agentStream := value.(*AgentStream)
		agentStream.Close()
		return true
	})
}

func Register(ctx context.Context) {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", startenv.GetGrpcPort()))
	if err != nil {
		logrus.Fatalf("grpc failed to listen: %s, err: %v", startenv.GetGrpcPort(), err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	newServer := NewServer(ctx)
	pb.RegisterAutoAgentServer(grpcServer, newServer)
	grpcServer.Serve(lis)
}

func GetGrpcServer(ctx context.Context) *grpc.Server {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	newServer := NewServer(ctx)
	pb.RegisterAutoAgentServer(grpcServer, newServer)
	reflection.Register(grpcServer)
	return grpcServer
}
