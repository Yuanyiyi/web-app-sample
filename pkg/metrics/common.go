package metrics

import (
	"sync"
	"time"
)

type CollectorMap struct {
	Record sync.Map
}

func (c *CollectorMap) GetAndResetValue(key string) int {
	res := c.GetValue(key)
	c.SetValue(key, 0)
	return res
}

func (c *CollectorMap) GetValue(key string) int {
	value, ok := c.Record.Load(key)
	if !ok {
		return 0
	}
	res, ok := value.(int)
	if !ok {
		return 0
	}
	return res
}

func (c *CollectorMap) SetValue(key string, value interface{}) {
	c.Record.Store(key, value)
}

func (c *CollectorMap) Inc(key string) {
	res := c.GetValue(key) + 1
	c.SetValue(key, res)
}

type CrdStateRecord struct {
	lock   sync.RWMutex
	record map[string]int64
}

func newCrdStateRecord() *CrdStateRecord {
	return &CrdStateRecord{lock: sync.RWMutex{}, record: make(map[string]int64)}
}

func (d *CrdStateRecord) set(key string, val int64) {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.record[key] = val
}

func (d *CrdStateRecord) get(key string) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	val, ok := d.record[key]
	if !ok {
		val = time.Now().UnixMilli()
	}
	return val
}

func (d *CrdStateRecord) getAndReset(key string) int64 {
	val := d.get(key)
	delete(d.record, key)
	return val
}

type StatisticsInfo struct {
	Number         int
	CumulativeTime int64
}

func newStatisticsInfo() *StatisticsInfo {
	return &StatisticsInfo{}
}

func (d *StatisticsInfo) Inc(duration int64) {
	d.Number++
	d.CumulativeTime += duration
}

func (d *StatisticsInfo) average() int64 {
	defer d.reset()
	if d.Number > 0 {
		return d.CumulativeTime / int64(d.Number)
	}
	return 0
}

func (d *StatisticsInfo) reset() {
	d.Number = 0
	d.CumulativeTime = 0
}
