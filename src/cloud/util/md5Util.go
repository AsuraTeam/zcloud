package util

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"time"
	"math/rand"
)

func Md5String(str string) string {
	m := md5.New()
	m.Write([]byte(str))
	return hex.EncodeToString(m.Sum(nil))
}

// 获取一个 MD5 穿
func Md5Uuid() string {
	str := strconv.FormatInt(time.Now().Unix(), 10) + strconv.Itoa(time.Now().Nanosecond()) + strconv.Itoa(rand.Intn((10000000)))
	return Md5String(str)
}
