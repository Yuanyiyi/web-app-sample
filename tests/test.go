package tests

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"github.com/web-app-sample/pkg/utils/runtime"
)

type BaseTestSuite struct {
	suite.Suite
}

func (tb *BaseTestSuite) SetupSuite() {
	runtime.SetLevel(logrus.DebugLevel)
}

func (tb *BaseTestSuite) TearDownSuite() {
	//TODO: drop table
}

type NoDBTestSuite struct {
	suite.Suite
}
