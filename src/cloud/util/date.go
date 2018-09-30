package util

import (
	"strings"
	"time"
	"strconv"
	"sort"
	"encoding/json"
	"bytes"
)

// 获取时间
func GetDate() string {
	return strings.Split(time.Now().Local().String(), ".")[0]
}

// 2018-01-26 15:11
// 生成镜像tag按时间戳
func MakeImageTag() string {
	t := time.Now().Local().String()
	t = ReplaceTime(t)
	t = strings.Replace(t, "-", "", -1)
	t = strings.Replace(t, ":", "", -1)
	t = strings.Replace(t, " ", "-", -1)
	return t
}

// 初始化时间
// 2018-01-15 13:49
var timeMap = Lock{}
func initTimeMap() {
	if len(timeMap.GetData()) > 0 {
		return
	}
	timeMap.Put("60", "1分钟前")
	timeMap.Put("3600", "1小时前")
	timeMap.Put("86400", "1天前")
	timeMap.Put("604800", "1周前")
	timeMap.Put("2419200", "1月前")
	timeMap.Put("29030400", "1年前")
	timeMap.Put("290304000", "N年前")
}

// 2018-01-15 13:30
// 获取时间简单格式
//util.GetMinTime("2018-01-15 14:40:49")
func GetMinTime(ctime string) string {
	if ctime == "" {
		return "未知"
	}
	initTimeMap()
	stamp := TimeToStamp(ctime)
	now := time.Now().Unix()
	sortt := make([]int, 0)
	interval := now - stamp
	for k := range timeMap.GetData() {
		kv, _ := strconv.ParseInt(k, 10, 64)
		if interval < kv {
			sortt = append(sortt, int(kv))
		}
	}
	sort.Ints(sortt)
	if len(sortt) == 0 {
		return "未知"
	}
	tr := timeMap.GetV(strconv.Itoa(sortt[0]))
	var r string
	max := 0
	switch tr {
	case "1年前":
		max = int(interval / 2419200)
		r = strconv.Itoa(max) + "月前"
		break
	case "1月前":
		max = int(interval / 604800)
		r = strconv.Itoa(max) + "周前"
		break
	case "1周前":
		max = int(interval / 86400)
		r = strconv.Itoa(max) + "天前"
		break
	case "1天前":
		max = int(interval / 3600)
		r = strconv.Itoa(max) + "小时前"
		break
	case "1小时前":
		max = int(interval / 60)
		r = strconv.Itoa(max) + "分钟前"
		break
	case "1分钟前":
		max = int(interval / 1)
		if max < 1 {
			max = 30
		}
		r = strconv.Itoa(max) + "秒前"
		break
	default:
		max = int(interval / 29030400)
		r = strconv.Itoa(max) + "年前"
		break
	}
	return r
}

// 2018-01-20 9:30
// 将对象转成字符串
func ObjToString(v interface{}) string  {
	t, _ := json.Marshal(v)
	return string(t)
}

func StringsToJSON(str string) string {
	var jsons bytes.Buffer
	for _, r := range str {
		rint := int(r)
		if rint < 128 {
			jsons.WriteRune(r)
		} else {
			jsons.WriteString("\\u")
			jsons.WriteString(strconv.FormatInt(int64(rint), 16))
		}
	}
	return jsons.String()
}

// 时间转成时间戳
// 2018-01-15 13:40
func TimeToStamp(ctime string) int64 {
	ctime = strings.TrimSpace(ctime)
	//获取本地location 	//待转化为时间戳的字符串 注意 这里的小时和分钟还要秒必须写 因为是跟着模板走的 修改模板的话也可以不写
	timeLayout := "2006-01-02 15:04:05"                          //转化所需模板
	loc, _ := time.LoadLocation("Local")                         //重要：获取时区
	theTime, err := time.ParseInLocation(timeLayout, ctime, loc) //使用模板在对应时区转化为time.time类型
	if err != nil {
		return time.Now().Unix()
	}
	sr := theTime.Unix() //转化为时间戳 类型是int64
	return sr
}

// 替换时间的T和Z
func ReplaceTime(t string) string {
	t = strings.Replace(t, "T", " ", -1)
	t = strings.Replace(t, "Z", "", -1)
	t = strings.Replace(t, "+0800 CS", "", -1)
	ts := strings.Split(t,".")
	return ts[0]
}
