package db_connection

import (
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/web-app-sample/pkg/utils/startenv"
)

func InitGormDB() {
	// 检查db是否存在并创建db
	InitSqlServer()

	logLevel := logger.Warn
	if startenv.GetLogLevel() == "DEBUG" {
		logLevel = logger.Info
	}
	// 添加链路追踪与扩展打印日志内容
	newLogger := logger.New(
		logrus.New(), // io writer
		logger.Config{
			SlowThreshold:             time.Second, //Slow SQL threshold
			LogLevel:                  logLevel,    //Log level
			IgnoreRecordNotFoundError: true,        //Ignore ErrRecordNotFound error for logger
			Colorful:                  false,       //Disable color
		},
	)

	var err error
	// 预编码：全局模式
	gormDb, err = gorm.Open(mysql.Open(getDataSourceName(startenv.GetSqlDb())),
		&gorm.Config{
			Logger:                 newLogger,
			PrepareStmt:            true,
			SkipDefaultTransaction: true, // 禁用事务
		},
	)

	// 设置默认引擎: InnoDB
	gormDb.InstanceSet("gorm:table_options", "ENGINE=InnoDB")

	sqlDb, err := gormDb.DB()
	if err != nil {
		logrus.Fatalf("sqlDb error: %s", err.Error())
	}

	err = sqlDb.Ping()
	if err != nil {
		logrus.Fatalf("Connect SQL-Server: %s Fail, err: %s", getDataSourceName(startenv.GetSqlDb()), err.Error())
		return
	}
	logrus.Info("Connect SQL-Server Success")

	gormDb.Logger.LogMode(0)
	//gormDb = gormDb.Session(&gorm.Session{PrepareStmt: true})

	// see "important settings" section.
	// db.SetConnMaxLifetime is required to ensure connections are closed by the driver safely before connection is closed by
	// MySQL manager, OS, or other middlewares.
	sqlDb.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns ishighly recommended to limit the number of connection used by the application. There is no recommended
	// limit number because it depends on application ans MySQL serve.
	sqlDb.SetMaxOpenConns(100)
	// db.SetMaxIdleConns is recommended to be set same to db.SetMaxOpenConns. When iti ts smaller than SetMaxOpenConns
	// connections can be opened and closed much more frequently than you expect.
	sqlDb.SetMaxIdleConns(10)
}

func GetGormDB() *gorm.DB {
	if gormDb == nil {
		InitGormDB()
	}
	return gormDb
}

func NewGormDB() *gorm.DB {
	return gormDb
}

func GormDBCloseAndDropDB() {
	if gormDb == nil {
		return
	}
	db0, err := gormDb.DB()
	if err == nil {
		db0.Close()
	}
	gormDb = nil
	DropDB()
}

func GormDBClose() {
	if gormDb == nil {
		return
	}
	db0, err := gormDb.DB()
	if err == nil {
		db0.Close()
	}
	gormDb = nil
}
