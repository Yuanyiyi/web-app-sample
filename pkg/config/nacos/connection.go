package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"github.com/sirupsen/logrus"

	"github.com/web-app-sample/pkg/utils/startenv"
)

var (
	nacosClient config_client.IConfigClient
)

func getNacosConfigClient() config_client.IConfigClient {
	//Another way of create clientConfig
	clientConfig := *constant.NewClientConfig(
		constant.WithNamespaceId(startenv.GetNacosId()), //When namespace is public, fill in the blank string here.
		constant.WithTimeoutMs(5000),
		constant.WithNotLoadCacheAtStart(true),
		constant.WithLogLevel(startenv.GetLogLevel()),
	)

	if startenv.GetNacosLogDir() != "" {
		clientConfig.LogDir = startenv.GetNacosLogDir()
	}

	if startenv.GetNacosCache() != "" {
		clientConfig.CacheDir = startenv.GetNacosCache()
	}

	//Another way of create serverConfigs
	serverConfigs := []constant.ServerConfig{
		*constant.NewServerConfig(
			startenv.GetNacosIpAddr(),
			uint64(startenv.GetNacosPort()),
			constant.WithScheme(startenv.GetNacosScheme()),
			constant.WithContextPath("/nacos"),
		),
	}

	// Another way of create config client for dynamic configuration (recommend)
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	if err != nil {
		logrus.Fatalf("connect nacos: {addr: %s, port: %d, namespaceId: %s } error: %s",
			startenv.GetNacosIpAddr(), startenv.GetNacosPort(), startenv.GetNacosId(), err.Error())
	} else {
		logrus.Infof("connect nacos: {addr: %s, port: %d, namespaceId: %s } success",
			startenv.GetNacosIpAddr(), startenv.GetNacosPort(), startenv.GetNacosId())
	}
	return configClient
}

func initNacosClient() {
	if &nacosClient == nil {
		nacosClient = getNacosConfigClient()
	}
}
