package util

import (
	"sync"
	"encoding/json"
	"os"
)

type Param struct {
	length int64
	start int64
	draw int64
}

type MapLock struct {
	Data map[string]interface{}
	Lock sync.RWMutex
}

func (m *MapLock) Set(k string, v interface{}){
	m.Lock.Lock()
	if len(m.Data) < 1 {
		m.Data = make(map[string]interface{})
	}
	m.Data[k] = v
	defer m.Lock.Unlock()
}

func returnMap(data interface{})string  {
	v, err := json.Marshal(data)
	if err == nil{
		return string(v)
	}
	return "{}"
}

// 为表格提供的数据
// return
func ResponseMap(listResult interface{}, total interface{}, draw interface{}) map[string]interface{}{
	maps := MapLock{}
	maps.Set("data", listResult)
	maps.Set("recordsTotal", total)
	maps.Set("recordsFiltered", total)
	maps.Set("draw", draw)
	return maps.Data
}

// 2018-01-15
//  获取表的行数
func GetTableRows(table string)  {

}

// 响应错误信息
func ResponseMapError(err string)map[string]interface{} {
	maps := MapLock{}
	maps.Set("data", err)
	maps.Set("recordsTotal", 0)
	maps.Set("recordsFiltered", 0)
	maps.Set("draw", 1)
	return maps.Data
}

// API响应信息
func ApiResponse(status bool, info interface{}) map[string]interface{} {
	maps := MapLock{}
	maps.Set("data", info)
	maps.Set("status", status)
	maps.Set("date",GetDate())
	hostname,_ := os.Hostname()
	maps.Set("server", hostname)
	maps.Set("code", 0)
	if !status{
		maps.Set("code", -1)
	}
	return maps.Data
}

