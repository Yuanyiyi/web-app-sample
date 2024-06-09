package db_connection

import (
	"math/rand"
	"os"
	"strconv"
	"testing"
)

func TestDBConnection(t *testing.T) {
	dbName := "hotfix_manager_" + strconv.Itoa(rand.Intn(20))
	os.Setenv("MYSQL_ADDR", "172.16.32.3")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB", dbName)
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "!QAZ2wsx3eDC")
	os.Setenv("ROOM_ADDR", "test.room.addr.com")
	os.Setenv("GRPC_PORT", "8446")
	os.Setenv("HTTP_PORT", "8447")
	InitDB()
	DBCloseAndDropDB()
}

func TestGormDBConnection(t *testing.T) {
	dbName := "hotfix_manager_" + strconv.Itoa(rand.Intn(20))
	os.Setenv("MYSQL_ADDR", "172.16.32.3")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_DB", dbName)
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "!QAZ2wsx3eDC")
	InitGormDB()
	GormDBCloseAndDropDB()
}
