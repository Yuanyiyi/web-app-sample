package controller

import "github.com/web-app-sample/pkg/database/mysql/models"

type Controller struct{}

func NewController(env string) *Controller {
	return &Controller{}
}

func (c *Controller) CreateHotfixTask(req *models.HotfixTask) error {
	return nil
}

func (c *Controller) DeleteHotfixTask(taskId string) error {
	return nil
}
