package nacos

import (
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestSendMsg(t *testing.T) {
	SendMsg("test")
	text, err := GetMsg()
	assert.Nil(t, err)
	assert.Equal(t, "test", text)
}

func TestListenMsg(t *testing.T) {
	ListenMsg(func(namespace, group, dataId, data string) {
		logrus.Infof("namespace: %s, group: %s, dataId: %s, data: %s", namespace, group, dataId, data)
	})
	SendMsg("test")
	SendMsg("test1")
	SendMsg("test2")
}
