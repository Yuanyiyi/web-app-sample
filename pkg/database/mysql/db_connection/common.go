package db_connection

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/web-app-sample/pkg/utils/startenv"
)

var (
	// 定义全局对象db
	gormDb *gorm.DB
	db     *sql.DB
)

func InitSqlServer() {
	logrus.Infof("===init sqlserver param===")
	db0, err := sql.Open("mysql", getDataSourceName(""))
	if err != nil {
		panic(err)
	}
	defer db0.Close()

	// check database
	rows, err := db0.Query(fmt.Sprintf("SHOW DATABASES LIKE '%s'", startenv.GetSqlDb()))
	if err != nil {
		logrus.Infof("mysql: %s", getDataSourceName(""))
		logrus.Fatalf("check database: %s, error: %s", startenv.GetSqlDb(), err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		logrus.Infof("Database：%s already exists", startenv.GetSqlDb())
		return
	}
	// 创建数据库
	_, err = db0.Exec(fmt.Sprintf("CREATE DATABASE %s", startenv.GetSqlDb()))
	if err != nil {
		logrus.Fatalf("Database created failed, err: %s", err.Error())
	}
	logrus.Infof("Database: %s created successfully", startenv.GetSqlDb())
}

func DropDB() {
	logrus.Infof("===drop database===")
	db0, err := sql.Open("mysql", getDataSourceName(""))
	if err != nil {
		panic(err)
	}
	defer db0.Close()

	// check database
	rows, err := db0.Query(fmt.Sprintf("DROP DATABASE %s", startenv.GetSqlDb()))
	if err != nil {
		logrus.Errorf("drop database： %s, error: %s", startenv.GetSqlDb(), err.Error())
		return
	}
	logrus.Infof("Database: %s already drop", startenv.GetSqlDb())
	defer rows.Close()
}

func getDataSourceName(dbname string) string {
	strs := []string{startenv.GetMySqlUser(), ":", startenv.GetMysqlPassword(),
		"@tcp(", startenv.GetMySqlAddr(), ":", startenv.GetMySqlPort(), ")/"}
	if len(dbname) > 0 {
		strs = append(strs, dbname)
		strs = append(strs, "?charset=utf8&parseTime=True&loc=Local")
	}
	return strings.Join(strs, "")
}
