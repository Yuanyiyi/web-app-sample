package models

import (
	"math/rand"
	"os"
	"strconv"

	"github.com/stretchr/testify/suite"

	"github.com/web-app-sample/pkg/database/mysql/db_connection"
	"github.com/web-app-sample/tests"
)

type ModelsTestSuite struct {
	tests.BaseTestSuite
	robotIds        []string
	agentIds        []string
	HotfixDsTable   *HotfixDSGormDB
	HotfixTaskTable *HotfixTaskGormDB
}

func (tb *ModelsTestSuite) SetupSuite() {
	dbName := "hotfix_manager_" + strconv.Itoa(rand.Intn(100))
	os.Setenv("MYSQL_DB", dbName)
	os.Setenv("MYSQL_ADDR", "172.16.32.3")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "!QAZ2wsx3eDC")
	os.Setenv("ROOM_ADDR", "test.room.addr.com")

	db_connection.InitGormDB()

	gormDb := db_connection.GetGormDB()
	tb.HotfixDsTable = NewHotfixDSGormDB(gormDb)
	tb.HotfixTaskTable = NewHotfixTaskGormDB(gormDb)
}

func (tb *ModelsTestSuite) TearDownSuite() {
	db_connection.GormDBCloseAndDropDB()
}

type NoDBTestSuite struct {
	suite.Suite
}
