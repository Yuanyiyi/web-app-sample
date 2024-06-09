package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/web-app-sample/pkg/utils/startenv"
)

func SendMsg(jsonStrData string) error {
	initNacosClient()
	//发布配置
	success, err := nacosClient.PublishConfig(vo.ConfigParam{
		DataId:  startenv.GetNacosDataId(),
		Group:   startenv.GetNacosGroup(),
		Content: jsonStrData})
	if err != nil {
		return err
	}
	if !success {
		return errors.Errorf("massege publish error")
	}

	return nil
}

func ListenMsg(f func(namespace, group, dataId, data string)) {
	initNacosClient()
	// 监听配置
	nacosClient.ListenConfig(vo.ConfigParam{
		DataId: startenv.GetNacosDataId(),
		Group:  startenv.GetNacosGroup(),
		OnChange: func(namespace, group, dataId, data string) {
			logrus.Infof("ListenConfig namespace: %s, group: %s, dataId: %s, data: %s", namespace, group, dataId, data)
			go f(namespace, group, dataId, data)
		},
	})
	return
}

func GetMsg() (string, error) {
	initNacosClient()
	// 获取配置
	content, err := nacosClient.GetConfig(vo.ConfigParam{
		DataId: startenv.GetNacosDataId(),
		Group:  startenv.GetNacosGroup()})
	if err != nil {
		logrus.Errorf("getConfig from nacos: {dataId: %s, group: %s} error: %s",
			startenv.GetNacosDataId(), startenv.GetNacosGroup(), err.Error())
	}
	return content, err
}
