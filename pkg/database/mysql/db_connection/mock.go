package db_connection

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	mockDb     *sql.DB
	mockGormDb *gorm.DB
	mockSql    sqlmock.Sqlmock
)

func MockDB() {
	var err error
	// 创建一个模拟的MySQL数据库连接和查询
	mockDb, mockSql, err = sqlmock.New()
	if err != nil {
		logrus.Fatalf("Failed to create mock database connection: %v", err.Error())
	}

	// 使用模拟的数据库连接创建GORM实例
	mockGormDb, err = gorm.Open(
		mysql.New(mysql.Config{
			Conn:                      mockDb,
			SkipInitializeWithVersion: true,
		}),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
	)
	if err != nil {
		logrus.Fatalf("Failed to creagte GORM instance: %s", err.Error())
	}

	// 设置GORM实例的连接
	// 这里可以使用应用程序中的真实连接设置代码
	// 示例中使用了模拟的数据库连接
	//mockGormDb.DB() = mockDb
}

func GetMockDB() *gorm.DB {
	if mockGormDb == nil {
		MockDB()
	}
	return mockGormDb
}

func MockDBClose() {
	mockDb.Close()
	mockDb = nil
	mockGormDb = nil
}
