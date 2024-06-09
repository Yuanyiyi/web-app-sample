package startenv

import (
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// 需要重启才会生效的配置
const (
	AutoTest              = "AUTO_TEST"
	MySqlAddr             = "MYSQL_ADDR"
	MySqlPort             = "MYSQL_PORT"
	MySqlDB               = "MYSQL_DB"
	MySqlUser             = "MYSQL_USER"
	MysqlPassword         = "MYSQL_PASSWORD"
	RedisAddr             = "REDIS_ADDR"
	RedisPassword         = "REDIS_PASSWORD"
	RedisDB               = "REDIS_DB"
	RedisDeployment       = "REDIS_DEPLOYMENT"
	RedisPoolSize         = "REDIS_POOL_SIZE"
	NacosIpAddr           = "NACOS_IP_ADDR"
	NacosPort             = "NACOS_PORT"
	NacosNamespaceId      = "NACOS_NAMESPACE_ID"
	NacosScheme           = "NACOS_SCHEME"
	NacosLogDir           = "NACOS_LOGDIR"
	NacosCache            = "NACOS_CACHE"
	NacosDataId           = "NACOS_DATAID"
	NacosGroup            = "NACOS_GROUP"
	LogFileName           = "LOG_FILE_NAME"
	Environment           = "ENVIRONMENT"
	OtelCollector         = "OTEL_COLLECTOR"
	LogSelect             = "LOG_SELECT"
	LogLevel              = "LOG_LEVEL"
	WebPort               = "HTTP_PORT"
	Env                   = "ENV"
	HotfixManagerAddr     = "HOTFIX_MANAGER_ADDR"
	CreateHotfixTaskQps   = "CREATE_HOTFIX_TASK_QPS"
	CreateTaskForLocalQps = "CREATE_HOTFIX_FOR_LOCAL_QPS"
	NuwagameQps           = "NUWA_GAME_QPS"
)

// 动态调整的配置
const (
	NuwaGameAddr       = "NUWA_GAME_ADDR"
	TStationAccessAddr = "T_STATION_ADDR"
	LocalManagerAddr   = "LOCAL_MANAGER_ADDR"
	LarkWebhook        = "LARK_WEBHOOK"
	LarkSecret         = "LARK_SECRET"
	AutoScaling        = "AUTO_SCALING"
	TaskLimit          = "TASK_LIMIT"
	ResourceId         = "RESOURCE_ID"
	StressTest         = "STRESS_TEST"
	ExpiredDate        = "EXPIRED_DATE"
	HotfixSwitch       = "HOTFIX_SWITCH"
)

func GetSentryDSN() string {
	return os.Getenv("AUTOTEST_SENTRY_DSN")
}

func GetLogsFileName() string {
	file := os.Getenv(LogFileName)
	file = filepath.Join(file, "manager-"+time.Now().Format("20060102_150405")+".log")
	return file
}

func GetEnvironment() string {
	name := os.Getenv(Environment)
	if name == "" {
		name = "dev"
	}
	return name
}

func GetLogLevel() string {
	if os.Getenv(LogLevel) == "" {
		return "info"
	}
	return os.Getenv(LogLevel)
}

func GetOtlpCollector() string {
	collector := os.Getenv(OtelCollector)
	if collector == "" {
		collector = "http://localhost:30080"
	}
	return collector
}

func GetLogSet() string {
	return os.Getenv(LogSelect)
}

func GetMySqlAddr() (addr string) {
	addr = os.Getenv(MySqlAddr)
	return
}

func GetMySqlUser() (user string) {
	user = os.Getenv(MySqlUser)
	if user == "" {
		user = "root"
	}
	return
}

func GetMysqlPassword() (password string) {
	password = os.Getenv(MysqlPassword)
	if password == "" {
		password = "!QAZ2wsx3eDC"
	}
	return
}

func GetMySqlPort() (port string) {
	port = os.Getenv(MySqlPort)
	if port == "" {
		port = "3306"
	}
	return
}

func GetSqlDb() (dbName string) {
	dbName = os.Getenv(MySqlDB)
	if dbName == "" {
		dbName = "auto_test"
	}
	return
}

func GetNuwaGameAddr() string {
	return os.Getenv(NuwaGameAddr)
}

func GetTStationAddr() string {
	return os.Getenv(TStationAccessAddr)
}

func GetLocalManagerAddr() string {
	return os.Getenv(LocalManagerAddr)
}

func GetHotfixManagerAddr() string {
	return os.Getenv(HotfixManagerAddr)
}

func GetHttpPort() string {
	if os.Getenv(WebPort) == "" {
		return "8449"
	}
	return os.Getenv(WebPort)
}

func GetRedisAddr() string {
	return os.Getenv(RedisAddr)
}

func GetRedisPassword() string {
	return os.Getenv(RedisPassword)
}

func GetRedisDB() string {
	return os.Getenv(RedisDB)
}

func GetEnv() string {
	return os.Getenv(Env)
}

func GetRedisDeployment() string {
	return os.Getenv(RedisDeployment)
}

func GetRedisPoolSize() string {
	return os.Getenv(RedisPoolSize)
}

func GetTaskLimit() int {
	limit, err := strconv.Atoi(os.Getenv(TaskLimit))
	if err != nil {
		return 20
	}
	return limit
}

func GetResourceId() string {
	return os.Getenv(ResourceId)
}

func GetLarkRobot() string {
	return os.Getenv(LarkWebhook)
}

func GetLarkSecret() string {
	return os.Getenv(LarkSecret)
}

func GetAutoTest() bool {
	return os.Getenv(AutoTest) == "true"
}

func GetAutoScaling() bool {
	return os.Getenv(AutoScaling) == "true"
}

func GetHotfixSwitch() bool {
	return os.Getenv(HotfixSwitch) == "" || os.Getenv(HotfixSwitch) == "true"
}

func GetStressTest() bool {
	return os.Getenv(StressTest) == "true"
}

func GetNacosIpAddr() string {
	return os.Getenv(NacosIpAddr)
}

func GetNacosPort() int {
	i, err := strconv.Atoi(os.Getenv(NacosPort))
	if err != nil {
		return 80
	}
	return i
}

func GetNacosId() string {
	return os.Getenv(NacosNamespaceId)
}

func GetNacosScheme() string {
	if os.Getenv(NacosScheme) == "" {
		return "http"
	}
	return os.Getenv(NacosScheme)
}

func GetNacosLogDir() string {
	return os.Getenv(NacosLogDir)
}

func GetNacosCache() string {
	return os.Getenv(NacosCache)
}

func GetNacosDataId() string {
	return os.Getenv(NacosDataId)
}

func GetNacosGroup() string {
	return os.Getenv(NacosGroup)
}

func GetCreateHotfixTaskQps() int {
	i, err := strconv.Atoi(os.Getenv(CreateHotfixTaskQps))
	if err != nil {
		return 50
	}
	return i
}

func GetCreateTaskForLocalQps() int {
	i, err := strconv.Atoi(os.Getenv(CreateTaskForLocalQps))
	if err != nil {
		return 100
	}
	return i
}

func GetNuwagameQps() int {
	i, err := strconv.Atoi(os.Getenv(NuwagameQps))
	if err != nil {
		return 20
	}
	return i
}

func GetExpiredDate() int {
	i, err := strconv.Atoi(os.Getenv(ExpiredDate))
	if err != nil {
		return 15
	}
	return i
}

func GetGrpcPort() string {
	return "3344"
}
