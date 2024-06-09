package models

import (
	"fmt"
	"time"
)

type UpdateStructToMap struct {
	updateMap map[string]interface{}
}

func newUpdateStructToMap() UpdateStructToMap {
	return UpdateStructToMap{updateMap: map[string]interface{}{}}
}

func (u UpdateStructToMap) add(key string, val int32) {
	if val == 0 {
		u.updateMap[key] = val
	}
}

func (u UpdateStructToMap) len() int {
	return len(u.updateMap)
}

func (u UpdateStructToMap) get() map[string]interface{} {
	return u.updateMap
}

func getTableNameByTime(tablePrefix string, t2 time.Time) string {
	year, month, day := t2.Date()
	return fmt.Sprintf("%s_%d%d%d", tablePrefix, year, month, day)
}
