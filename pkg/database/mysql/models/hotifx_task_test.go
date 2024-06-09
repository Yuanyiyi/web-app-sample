package models

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

func getHotfixTaskTableName(x int) string {
	year, month, day := time.Now().Date()
	return fmt.Sprintf("hotfix_task_%d%d%d", year, month, day+x)
}

type HotfixTaskTestSuite struct {
	ModelsTestSuite
}

func TestHotfixTaskTestSuite(t *testing.T) {
	suite.Run(t, new(HotfixTaskTestSuite))
}

func (t *HotfixTaskTestSuite) TestHotfixTaskProcess() {
	table1 := getHotfixTaskTableName(0)
	t.HotfixTaskTable.CreateTable(table1)
	err := t.HotfixTaskTable.Create(table1, &HotfixTask{
		TaskId: "task_0001",
	})
	t.Nil(err)

	table2 := getHotfixTaskTableName(-1)
	t.HotfixTaskTable.CreateTable(table2)
	err = t.HotfixTaskTable.Create(table2, &HotfixTask{
		TaskId: "task_0002",
	})
	t.Nil(err)

	err = t.HotfixTaskTable.CreateIfNotFound(table2, &HotfixTask{
		TaskId: "task_0002",
	})
	t.Nil(err)

	err = t.HotfixTaskTable.UpdateDsType(table1, "task_0001", "ds_type_0001")
	t.Nil(err)
	res, err := t.HotfixTaskTable.FindByTaskId(table1, "task_0001")
	t.Nil(err)
	t.Equal("ds_type_0001", res.DsType)

	err = t.HotfixTaskTable.UpdateToExecuting(table1, "task_0001", "ds_id_0001")
	t.Nil(err)
	res, err = t.HotfixTaskTable.FindByTaskId(table1, "task_0001")
	t.Nil(err)
	t.Equal(Executing, res.State)

	tasks, err := t.HotfixTaskTable.FindByDsType(table1, "ds_type_0001", Executing)
	t.Nil(err)
	logrus.Infof("tasks: %v", tasks[0])

	tasks, err = t.HotfixTaskTable.FindByDsID(table1, "ds_id_0001", Executing)
	t.Nil(err)
	logrus.Infof("tasks: %v", tasks[0])

	ress, err := t.HotfixTaskTable.GetTasksByState(table1, Executing)
	t.Nil(err)
	t.Equal(1, len(ress))

	ress, err = t.HotfixTaskTable.GetTasks(table1)
	t.Nil(err)
	t.Equal(1, len(ress))

	err = t.HotfixTaskTable.UpdateState(table1, "task_0001", Failed, "")
	t.Nil(err)
	res, err = t.HotfixTaskTable.FindByTaskId(table1, "task_0001")
	t.Nil(err)
	t.Nil(res)
}

func (t *HotfixTaskTestSuite) TestHotfixTaskProcess1() {
	table1 := getHotfixTaskTableName(0)
	t.HotfixTaskTable.CreateTable(table1)

	err := t.HotfixTaskTable.Create(table1, &HotfixTask{
		TaskId: "task_00011",
		State:  Executing,
		DsId:   "ds_id_0000",
	})
	t.Nil(err)
	err = t.HotfixTaskTable.UpdateState(table1, "task_00011", Executing, "")
	t.Nil(err)

	err = t.HotfixTaskTable.Create(table1, &HotfixTask{
		TaskId: "task_00012",
		State:  Completed,
		DsId:   "ds_id_0000",
	})
	t.Nil(err)
	err = t.HotfixTaskTable.UpdateState(table1, "task_00012", Completed, "")
	t.Nil(err)

	tasks, err := t.HotfixTaskTable.GetTasks(table1)
	t.Nil(err)
	t.Equal(1, len(tasks))
	logrus.Infof("======task: %v", tasks[0])

	err = t.HotfixTaskTable.UpdateBatchDsId(table1, "ds_id_0000")
	t.Nil(err)

	tasks, err = t.HotfixTaskTable.GetTasks(table1)
	t.Nil(err)
	t.Equal(1, len(tasks))
	logrus.Infof("======task: %v", tasks[0])
}

func (t *HotfixTaskTestSuite) TestHotfixTaskProcess2() {
	table1 := getHotfixTaskTableName(0)
	t.HotfixTaskTable.CreateTable(table1)
	err := t.HotfixTaskTable.Create(table1, &HotfixTask{
		TaskId: "task_0003",
		DsType: "ds_typ1",
	})
	t.Nil(err)

	err = t.HotfixTaskTable.UpdateStateByDsType(table1, "ds_typ1", "test")
	t.Nil(err)

	res, err := t.HotfixTaskTable.FindByTaskId(table1, "task_0003")
	t.Nil(err)
	t.Nil(res)
}
