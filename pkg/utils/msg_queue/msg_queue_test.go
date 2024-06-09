package msg_queue

import (
	"github.com/sirupsen/logrus"
	"testing"
)

type MsgDa struct {
	msgT string
	data string
}

func TestQueue(t *testing.T) {
	p := GetNewMsgQueue()
	defer p.Close()

	// 订阅全部主题
	//all := p.Subscribe()
	webServerMsg := p.SubscribeTopic(func(v *MsgStruct) bool {
		if v != nil && v.MsgType == 0 {
			return true
		}
		return false
	})
	grpcServerMsg := p.SubscribeTopic(func(v *MsgStruct) bool {
		if v != nil && v.MsgType == 1 {
			return true
		}
		return false
	})

	p.Publish(&MsgStruct{0, "web-test"})
	p.Publish(&MsgStruct{0, "web-test1"})
	p.Publish(&MsgStruct{1, "grpc-test1"})
	p.Publish(&MsgStruct{1, "grpc-test2"})

	go func() {
		for v := range webServerMsg {
			msg := v.(*MsgStruct)
			data := msg.MsgData.(string)
			logrus.Infof("webServerMsg: %v", data)
		}
	}()

	go func() {
		for v := range grpcServerMsg {
			logrus.Infof("grpcServerMsg subscribe: %v", v)
		}
	}()

}
