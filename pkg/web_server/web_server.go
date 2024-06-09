package web_server

import (
	"context"
	"net/http"

	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"

	"github.com/web-app-sample/pkg/controller"
	"github.com/web-app-sample/pkg/database/mysql/models"
)

type WebServer struct {
	*gin.Engine
	ctx     context.Context
	manager *controller.Controller
}

func NewWebServer(ctx context.Context, manager *controller.Controller) *WebServer {
	server := &WebServer{
		ctx:     ctx,
		Engine:  gin.Default(),
		manager: manager,
	}
	server.register()
	return server
}

func (s *WebServer) register() {
	s.GET("/health", s.health)
	s.GET("/metrics", s.PromHandler(promhttp.Handler()))
	s.POST("/hotfix_manager/create_hotfix_task", s.createHotfixTask)
	s.POST("/hotfix_manager/delete_hotfix_task", s.deleteHotfixTask)
}

func (s *WebServer) health(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func (s *WebServer) PromHandler(handler http.Handler) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}

func (s *WebServer) createHotfixTask(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("createHotfixTask panic: %v", err)
		}
	}()
	entry, err1 := sentinel.Entry(createHotfixTask)
	if err1 != nil {
		logrus.Warnf("===createHotfixTask Limiting: %s", err1)
		ctx.JSON(http.StatusOK,
			gin.H{"code": SystemError, "msg": "Limiting"})
		return
	}
	defer entry.Exit()

	req := &CreateHotfixTaskReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		logrus.Warnf("===createHotfixTask error: %s", err.Error())
		ctx.JSON(http.StatusOK,
			gin.H{"code": SystemError, "msg": "other error"})
		return
	}

	logrus.Debugf("createHotfixTask req: %v", req)
	hotfixTask := &models.HotfixTask{}
	err = s.manager.CreateHotfixTask(hotfixTask)
	if err != nil {
		logrus.Warnf("===createHotfixTask error: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": "ok"})
}

func (s *WebServer) deleteHotfixTask(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Warnf("deleteHotfixTask panic: %v", err)
		}
	}()

	req := &DeleteHotfixTaskReq{}
	err := ctx.BindJSON(req)
	if err != nil {
		logrus.Warnf("===deleteHotfixTask error: %s", err.Error())
		ctx.JSON(http.StatusOK,
			gin.H{"code": SystemError, "msg": "other error"})
		return
	}

	logrus.Debugf("deleteHotfixTask req: %v", req)
	err = s.manager.DeleteHotfixTask(req.TaskId)
	if err != nil {
		logrus.Warnf("===deleteHotfixTask error: %s", err.Error())
	}
	ctx.JSON(http.StatusOK, gin.H{"code": http.StatusOK, "status": "ok"})
}
