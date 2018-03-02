package util

import (
	"github.com/astaxie/beego/logs"
	"encoding/base64"
)

// base64 加密
// 2018-01-21 18:39
func Base64Encoding(str string)  string {
	input := []byte(str)
	// 演示base64编码
	encodeString := base64.StdEncoding.EncodeToString(input)
	return encodeString
}

// base64 解密
// 2018-01-21 18:40
func Base64Decoding(encodeString string) string {
	decodeBytes, err := base64.StdEncoding.DecodeString(encodeString)
	if err != nil {
		logs.Error(err)
		return ""
	}
	return string(decodeBytes)
}