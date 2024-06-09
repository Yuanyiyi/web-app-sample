package nacos

type Config struct {
	NuwaGameAddr       string `json:"nuwa_game_addr"`
	TStationAccessAddr string `json:"t_station_addr"`
	LocalManagerAddr   string `json:"local_manager_addr"`
	LarkWebhook        string `json:"lark_webhook"`
	LarkSecret         string `json:"lark_secret"`
	AutoScaling        string `json:"auto_scaling"`
	TaskLimit          string `json:"task_limit"`
	ResourceId         string `json:"dsversion_resource_id"`
	StressTest         string `json:"stress_test"`
	ExpiredDate        string `json:"expired_date"`
	HotfixSwitch       string `json:"hotfix_switch"`
}
