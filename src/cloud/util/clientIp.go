package util

import (
	"net/http"
	"strings"
)

// 获取客户端IP地址
func GetClientIp(r *http.Request) string {
	var ip string

	xfor := r.Header.Get("X-Forwarded-For")
	if len(xfor) > 0 {
		ip = strings.Split(strings.Replace(xfor, "\"", "", -1), ",")[0]
	}

	if ip == "" {
		ip = r.Header.Get("Remote_addr")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	return strings.Split(ip, ":")[0]
}