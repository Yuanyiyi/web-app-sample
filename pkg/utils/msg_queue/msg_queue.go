package msg_queue

import (
	"sync"
	"time"
)

const (
	WebServerMessage MsgType = iota
	GrpcServerMessage
	AutoManager
)

var (
	msgQue         *MessageQueue
	autoManagerMsg chan interface{}
	webServerMsq   chan interface{}
	grpcServerMsg  chan interface{}
)

func SetNewMsgQueue(size int, t1 time.Duration) *MessageQueue {
	if msgQue == nil {
		msgQue = NewMessageQueue(size, t1)
	}
	return msgQue
}

func GetNewMsgQueue() *MessageQueue {
	if msgQue == nil {
		msgQue = NewMessageQueue(1024, 5*time.Second)
	}
	return msgQue
}

type (
	MsgType    int
	MsgData    any
	subscriber chan interface{}        // 订阅者为一个通道
	topicFunc  func(v *MsgStruct) bool // 主题为一个过滤器
	MsgStruct  struct {
		MsgType
		MsgData
	}
)

type MessageQueue struct {
	m           sync.RWMutex             // 读写锁
	buffer      int                      // 订阅队列的缓存长度
	timeout     time.Duration            // 生产者生产消息的超时时间
	subscribers map[subscriber]topicFunc // 订阅者消息
}

// MessageQueue 构建发送者对象
func NewMessageQueue(buffer int, timeout time.Duration) *MessageQueue {
	return &MessageQueue{
		m:           sync.RWMutex{},
		buffer:      buffer,
		timeout:     timeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 订阅全部主题
func (m *MessageQueue) Subscribe() chan interface{} {
	return m.SubscribeTopic(nil)
}

func (m *MessageQueue) SubscribeTopic(topic topicFunc) chan interface{} {
	ch := make(chan interface{}, m.buffer)
	m.m.Lock()
	defer m.m.Unlock()
	m.subscribers[ch] = topic
	return ch
}

// 退出主题
func (m *MessageQueue) Evict(sub chan interface{}) {
	m.m.Lock()
	defer m.m.Unlock()
	delete(m.subscribers, sub)
	close(sub)
}

// close chan of all
func (m *MessageQueue) Close() {
	m.m.Lock()
	defer m.m.Unlock()

	for sub := range m.subscribers {
		delete(m.subscribers, sub)
		close(sub)
	}
}

// publish 向所有满足条件的主题发送消息
func (m *MessageQueue) Publish(v *MsgStruct) {
	if v == nil {
		return
	}
	m.m.Lock()
	defer m.m.Unlock()

	var wg sync.WaitGroup
	for sub, topic := range m.subscribers {
		wg.Add(1)
		go m.SendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// sendTopic 向某一主题发送消息
func (m *MessageQueue) SendTopic(sub subscriber, topic topicFunc, v *MsgStruct, wg *sync.WaitGroup) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}
	select {
	case sub <- v:
	case <-time.After(m.timeout):
	}
}

func GetMsg(v interface{}) any {
	if v == nil {
		return nil
	}
	val, ok := v.(*MsgStruct)
	if !ok {
		return nil
	}
	return val.MsgData
}

func GetAutoManagerChan() chan interface{} {
	if autoManagerMsg == nil {
		msg := GetNewMsgQueue()
		autoManagerMsg = msg.SubscribeTopic(func(v *MsgStruct) bool {
			if v != nil && v.MsgType == AutoManager {
				return true
			}
			return false
		})
	}
	return autoManagerMsg
}

func GetWebServerMsqChan() chan interface{} {
	if webServerMsq == nil {
		msg := GetNewMsgQueue()
		webServerMsq = msg.SubscribeTopic(func(v *MsgStruct) bool {
			if v != nil && v.MsgType == WebServerMessage {
				return true
			}
			return false
		})
	}
	return webServerMsq
}

func GetGrpcServerMsgChan() chan interface{} {
	if grpcServerMsg == nil {
		msg := GetNewMsgQueue()
		grpcServerMsg = msg.SubscribeTopic(func(v *MsgStruct) bool {
			if v != nil && v.MsgType == GrpcServerMessage {
				return true
			}
			return false
		})
	}
	return grpcServerMsg
}
