package util

import "net/http"

// 获取客户端IP地址
func GetClientIp(r *http.Request) string {
	ip := r.Header.Get("Remote_addr")
	if ip == "" {
		ip = r.RemoteAddr
	}
	return ip
}


