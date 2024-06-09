package web_server

import (
	sentinel2 "github.com/web-app-sample/pkg/utils/sentinel"
	"github.com/web-app-sample/pkg/utils/startenv"
)

const (
	createHotfixTask = "create_hotfix_task"
)

var (
	qps = startenv.GetCreateHotfixTaskQps()
)

func init() {
	sentinel2.LimitingPolicy(createHotfixTask, float64(qps))
}
