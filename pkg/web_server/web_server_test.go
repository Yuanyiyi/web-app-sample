package web_server

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"

	"github.com/web-app-sample/pkg/controller"
	httpUtil "github.com/web-app-sample/pkg/utils/http"
	"github.com/web-app-sample/tests"
)

type WebServerTestSuite struct {
	tests.BaseTestSuite
}

func TestWebServerTestSuite(t *testing.T) {
	suite.Run(t, new(WebServerTestSuite))
}

func (t *WebServerTestSuite) TestWebServer() {
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan interface{})
	defer func() {
		c <- struct{}{}
	}()
	go func() {
		<-c
		cancel()
	}()

	managerServer := controller.NewController("auto_test")
	webServer := NewWebServer(ctx, managerServer)

	ctl := gomock.NewController(&testing.T{})
	defer ctl.Finish()

	var httpRsp *Response
	logrus.Info("=========web-server add-task ==========")
	// createHotfixTask
	req1 := &CreateHotfixTaskReq{
		DsVersion:               "v0.31.0.1.0000",
		SceneId:                 "scene_id_001",
		SceneVersion:            "scene_version_001",
		TemplateSceneId:         "templated_id_001",
		TemplateVersion:         "templated_version_001",
		OriginalTemplateVersion: "templated_version_000",
	}
	success, resp := httpUtil.CreatePostReqCtx(req1, webServer.createHotfixTask)
	t.True(true, success)
	err := json.Unmarshal([]byte(resp), &httpRsp)
	t.Nil(err)
	t.Equal(http.StatusOK, httpRsp.Code)
}
