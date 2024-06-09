package main

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"

	"hotfix.manager.com/pkg/config/nacos"
	"hotfix.manager.com/pkg/database/mysql/db_connection"
	"hotfix.manager.com/pkg/manager"
	"hotfix.manager.com/pkg/utils/runtime"
	"hotfix.manager.com/pkg/utils/signals"
	"hotfix.manager.com/pkg/utils/startenv"
	"hotfix.manager.com/pkg/web_server"
)

func main() {
	ctx := signals.GetSigKillContext()
	logrus.WithField("logLevel", startenv.GetLogLevel()).Info("Setting LogLevel configuration")
	level, err := logrus.ParseLevel(strings.ToLower(startenv.GetLogLevel()))
	if err == nil {
		runtime.SetLevel(level)
	} else {
		logrus.WithError(err).Info("Specified wrong Logging.hotfix_manager. Setting default loglevel - Info")
		runtime.SetLevel(logrus.InfoLevel)
	}
	// 初始化配置
	nacos.InitNacosConfig()
	// 连接mysql
	db_connection.InitGormDB()
	managerServer, _ := manager.NewManager(*ctx, "")
	webServer := web_server.NewWebServer(*ctx, managerServer)
	go managerServer.Run()
	if err = webServer.Run(fmt.Sprintf(":%s", startenv.GetHttpPort())); err != nil {
		logrus.Fatalf("http-server error: %s", err.Error())
	} else {
		logrus.Infof("http-server run success, port: %s", startenv.GetHttpPort())
	}
	// close db connection
	db_connection.GormDBClose()
	logrus.Info("Shut down hotfix-manager server")
}
