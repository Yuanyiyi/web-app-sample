package db_connection

import (
	"database/sql"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/web-app-sample/pkg/utils/startenv"
)

// sqlserver连接
func InitDB() {
	// 检查db是否存在并创建db
	InitSqlServer()

	var err error
	db, err = sql.Open("mysql", getDataSourceName(startenv.GetSqlDb()))
	if err != nil {
		logrus.Fatalf("mysql: %s connect fail: %s", getDataSourceName(startenv.GetSqlDb()), err.Error())
	}

	err = db.Ping()
	if err != nil {
		logrus.Fatalf("Connect SQL-Server: %s Fail, err: %s", getDataSourceName(startenv.GetSqlDb()), err.Error())
	}
	logrus.Info("Connect SQL-Server Success")

	// see "important settings" section.
	// db.SetConnMaxLifetime is required to ensure connections are closed by the driver safely before connection is closed by
	// MySQL manager, OS, or other middlewares.
	db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns ishighly recommended to limit the number of connection used by the application. There is no recommended
	// limit number because it depends on application ans MySQL serve.
	db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns is recommended to be set same to db.SetMaxOpenConns. When iti ts smaller than SetMaxOpenConns
	// connections can be opened and closed much more frequently than you expect.
	db.SetMaxIdleConns(10)
}

// 返回全局唯一的db
func GetSqlDB() *sql.DB {
	if db == nil {
		InitDB()
	}
	return db
}

func DBClose() {
	if db == nil {
		return
	}
	db.Close()
	db = nil
}

func DBCloseAndDropDB() {
	if db == nil {
		return
	}
	db.Close()
	db = nil
	DropDB()
}
