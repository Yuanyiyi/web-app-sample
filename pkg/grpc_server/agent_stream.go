package grpc_server

import "github.com/web-app-sample/pkg/pb/web_app_sample"

type AgentStream struct {
	StreamServicePassiveEventChan chan *web_app_sample.StreamPassiveEvent
	Addr                          string
	Name                          string
}

func NewAgentStream(addr string, name string) *AgentStream {
	return &AgentStream{
		StreamServicePassiveEventChan: make(chan *web_app_sample.StreamPassiveEvent, 1024),
		Addr:                          addr,
		Name:                          name,
	}
}

func (agentStream *AgentStream) SendStreamPassiveEvent(streamEvent *web_app_sample.StreamPassiveEvent) {
	agentStream.StreamServicePassiveEventChan <- streamEvent
}

func (agentStream *AgentStream) Close() {
	agentStream.StreamServicePassiveEventChan <- nil
}
