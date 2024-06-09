package common

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
)

const (
	Cpu            int64 = 500 // 0.5c
	Memory         int64 = 700 // 700m
	MD_KEY_AGENTID       = "agentID"
)

func RemoveElemFromList(lst []string, removeStr string) []string {
	for idx, str := range lst {
		if str == removeStr {
			return append(lst[:idx], lst[idx+1:]...)
		}
	}
	return lst
}

/*
gameId+version+dsVersion
*/
func GenerateMd5Id(strs []string) string {
	label := strings.Join(strs, "_")
	res := md5.Sum([]byte(label))
	return hex.EncodeToString(res[:])
}

func Min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func Max(a, b int64) int64 {
	if a < b {
		return b
	}
	return a
}

func GetRobotNumber(avaCpu, avaMemory, n int64) int64 {
	if avaCpu-Cpu*n <= 0 || avaMemory-Memory*n <= 0 {
		return 0
	}
	cpuCnt := (avaCpu - Cpu*n) / Cpu
	memCnt := (avaMemory - Memory*n) / Memory
	return Max(Min(cpuCnt, memCnt), 0)
}
