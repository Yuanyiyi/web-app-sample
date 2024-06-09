package web_server

const (
	SystemError = 500
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type CreateHotfixTaskReq struct {
	DsVersion               string `json:"ds_version" binding:"required"`
	SceneId                 string `json:"scene_id" binding:"required"`
	SceneVersion            string `json:"scene_version" binding:"required"`
	TemplateSceneId         string `json:"template_scene_id" binding:"required"`
	TemplateVersion         string `json:"template_version" binding:"required"`
	OriginalTemplateVersion string `json:"original_template_version" binding:"required"`
}

type DeleteHotfixTaskReq struct {
	TaskId string `json:"task_id" binding:"required"`
}

type CreateHotfixTaskRsp struct {
	TaskId string `json:"task_id" binding:"required"`
	State  string `json:"state" binding:"required"`
	Reason string `json:"reason"`
}

type CreateDsResp struct {
	GameServerId  string `json:"gameServerId"`
	DsId          string `json:"dsId" binding:"required"`
	InternetIp    string `json:"internetIp" binding:"required"`
	IntranetIp    string `json:"intranetIp" binding:"required"`
	DSPort        int32  `json:"dsPort"`
	LocalPort     int32  `json:"localPort" binding:"required"`
	Reason        string `json:"error"`
	LastPointTime int64  `json:"lastPointTime"`
}

type CreateDsFailResp struct {
	DsId         string `json:"DsId" binding:"required"`
	GameServerId string `json:"GameServerId"`
	InternetIp   string `json:"InternetIp"`
	IntranetIp   string `json:"IntranetIp"`
	Reason       string `json:"Reason"`
}

type DeleteDsResp struct {
	DsId string `json:"dsId" binding:"required"`
	//success or fail
	Status string `json:"status"`
	Reason string `json:"error"`
}

type ReportHotfixTaskReq struct {
	DsId  string       `json:"ds_id" binding:"required"`
	Tasks []HotfixTask `json:"tasks"`
}

type HotfixTask struct {
	TaskId string `json:"task_id" binding:"required"`
	State  string `json:"state" binding:"required"`
}

type ReportHotfixTaskRsp struct{}
