package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	omitTaskCols = []string{"scene_id", "scene_version", "template_scene_id", "template_version", "is_deleted", "index"}
)

// 热更任务表
type HotfixTask struct {
	gorm.Model
	TaskId    string `json:"task_id" gorm:"column:task_id;size:64;not null;primary_key"`
	Version   string `json:"ds_version" gorm:"column:ds_version;size:64;not null"`
	State     State  `json:"state" gorm:"column:state;size:16;not null"`
	DsType    string `json:"ds_type" gorm:"column:ds_type;size:64;not null"`
	DsId      string `json:"ds_id" gorm:"column:ds_id;size:64;not null"`
	IsDeleted string `json:"is_deleted" gorm:"column:is_deleted;size:4;not null"`
	Reason    string `json:"reason" gorm:"column:reason;size:64;not null"`
}

func (h HotfixTask) TableNamePrefix() string {
	return "hotfix_task"
}

func (h HotfixTask) TableName() string {
	// 默认表名
	year, month, day := time.Now().Date()
	return fmt.Sprintf("%s_%d%d%d", h.TableNamePrefix(), year, month, day)
}

type HotfixTaskGormDB struct {
	table  HotfixTask
	gormDb *gorm.DB
}

func NewHotfixTaskGormDB(db *gorm.DB) *HotfixTaskGormDB {
	table := &HotfixTaskGormDB{
		table:  HotfixTask{},
		gormDb: db,
	}
	if !table.ExistTable("") {
		table.init("")
	}

	return table
}

func (c *HotfixTaskGormDB) TableNamePrefix() string {
	return c.table.TableNamePrefix()
}

func (c *HotfixTaskGormDB) CreateTable(tableName string) {
	if !c.ExistTable(tableName) {
		c.init(tableName)
	}
}

func (c *HotfixTaskGormDB) ExistTable(tableName string) bool {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	return c.gormDb.Migrator().HasTable(tableName)
}

func (c *HotfixTaskGormDB) DropTable(tableName string) error {
	if tableName == "" {
		return nil
	}
	return c.gormDb.Table(tableName).Migrator().DropTable(&HotfixTask{})
}

// 没有表时执行创建表
func (c *HotfixTaskGormDB) init(tableName string) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	logrus.Infof("create table: %s", tableName)
	c.gormDb.Table(tableName).AutoMigrate(&HotfixTask{})
}

func (c *HotfixTaskGormDB) Create(tableName string, req *HotfixTask) error {
	if tableName == "" {
		tableName = getTableNameByTime(c.TableNamePrefix(), time.Now())
	}
	req.State = NoStarted
	req.IsDeleted = "no"
	err := c.gormDb.Table(tableName).Create(req).Error
	return err
}

func (c *HotfixTaskGormDB) CreateIfNotFound(tableName string, req *HotfixTask) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	var record *HotfixTask
	err := c.gormDb.Table(tableName).
		Select("task_id").
		Where(&HotfixTask{TaskId: req.TaskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Take(&record).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Create("", req)
		}
		return err
	}
	return err
}

func (c *HotfixTaskGormDB) UpdateState(tableName, taskId string, state State, reason string) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{State: state, Reason: reason}

	switch state {
	case Failed, Stop:
		data.IsDeleted = "yes"
	}
	err := c.gormDb.Table(tableName).
		Where(&HotfixTask{TaskId: taskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Updates(&data).Error
	return err
}

func (c *HotfixTaskGormDB) UpdateBatchState(tableName string, taskIds []string, state State) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{State: state}
	switch state {
	case Failed, Stop:
		data.IsDeleted = "yes"
	}
	err := c.gormDb.Table(tableName).
		Where("task_id in (?) and state != 'completed' and is_deleted='no'", taskIds).
		Updates(&data).Error
	return err
}

func (c *HotfixTaskGormDB) UpdateBatchDsId(tableName, dsId string) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	err := c.gormDb.Table(tableName).
		Where("ds_id=? and state != 'completed' and is_deleted='no'", dsId).
		Updates(map[string]interface{}{
			"ds_id": "",
			"state": "no-started",
		}).Error
	return err
}

func (c *HotfixTaskGormDB) UpdateToExecuting(tableName, taskId, dsId string) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{State: Executing, DsId: dsId}
	err := c.gormDb.Table(tableName).
		Where(&HotfixTask{TaskId: taskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Updates(&data).Error
	return err
}

func (c *HotfixTaskGormDB) UpdateBatchToExecuting(tableName, dsId string, taskIds []string) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{State: Executing, DsId: dsId}
	err := c.gormDb.Table(tableName).
		Where("task_id in (?) and is_deleted='no'", taskIds).
		Updates(&data).Error
	return err
}

func (c *HotfixTaskGormDB) UpdateDsType(tableName, taskId, dsType string) (err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{DsType: dsType}
	err = c.gormDb.Table(tableName).
		Where(&HotfixTask{TaskId: taskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Updates(&data).Error
	return
}

func (c *HotfixTaskGormDB) Deleted(tableName, taskId string) (err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	data := &HotfixTask{IsDeleted: "yes"}
	err = c.gormDb.Table(tableName).
		Where(&HotfixTask{TaskId: taskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Updates(&data).Error
	return
}

func (c *HotfixTaskGormDB) UpdateStateByDsType(tableName, dsType, reason string) error {
	if tableName == "" {
		tableName = c.table.TableName()
	}

	err := c.gormDb.Table(tableName).
		Where("ds_type=? and state='no-started' and is_deleted='no'", dsType).
		Updates(&HotfixTask{State: Failed, IsDeleted: "yes", Reason: reason}).Error
	return err
}

func (c *HotfixTaskGormDB) FindByTaskId(tableName, taskId string) (res *HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	err = c.gormDb.Table(tableName).
		Omit("is_deleted", "reason").
		Where(&HotfixTask{TaskId: taskId, IsDeleted: "no"}, "task_id", "is_deleted").
		Take(&res).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return
	}
	return
}

func (c *HotfixTaskGormDB) FindByDsType(tableName, dsType string, state State) (res []*HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	omitCols1 := append(omitTaskCols, "ds_type", "ds_version")
	db := c.gormDb.Table(tableName).Omit(omitCols1...).
		Where(&HotfixTask{DsType: dsType, IsDeleted: "no"}, "ds_type", "is_deleted")
	if state != "" {
		omitCols1 = append(omitCols1, "state")
		db = c.gormDb.Table(tableName).Omit(omitCols1...).
			Where(&HotfixTask{DsType: dsType, State: state, IsDeleted: "no"}, "ds_type", "state", "is_deleted")
	}

	err = db.Find(&res).Error
	if err != nil {
		return
	}

	return
}

func (c *HotfixTaskGormDB) FindByDsID(tableName, dsId string, state State) (res []*HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	omitCols1 := append(omitTaskCols, "ds_type", "ds_version")
	db := c.gormDb.Table(tableName).Omit(omitCols1...).
		Where(&HotfixTask{DsId: dsId, IsDeleted: "no"}, "ds_id", "is_deleted")
	if state != "" {
		omitCols1 = append(omitCols1, "state")
		db = c.gormDb.Table(tableName).Omit(omitCols1...).
			Where(&HotfixTask{DsId: dsId, State: state, IsDeleted: "no"}, "ds_id", "state", "is_deleted")
	}

	err = db.Find(&res).Error
	if err != nil {
		return
	}

	return
}

func (c *HotfixTaskGormDB) FindByDsIDByNotExecuting(tableName, dsId string) (res []*HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	omitCols1 := append(omitTaskCols, "ds_type", "ds_version")
	err = c.gormDb.Table(tableName).Omit(omitCols1...).
		Where("ds_id=? and is_deleted='no' and state != 'executing'", dsId).
		Find(&res).Error

	return
}

func (c *HotfixTaskGormDB) GetTasksByState(tableName string, state State) (taskData []*HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	omitCols1 := append(omitTaskCols, "ds_id", "ds_version")
	err = c.gormDb.Table(tableName).Omit(omitCols1...).
		Where(&HotfixTask{State: state, IsDeleted: "no"},
			"state", "is_deleted").
		Find(&taskData).Error
	return
}

func (c *HotfixTaskGormDB) GetTasks(tableName string) (taskData []*HotfixTask, err error) {
	if tableName == "" {
		tableName = c.table.TableName()
	}
	omitCols1 := append(omitTaskCols, "ds_version")
	err = c.gormDb.Table(tableName).Omit(omitCols1...).
		Where("state != 'completed' and is_deleted='no'").
		Find(&taskData).Error

	return
}
