package limiter

import (
	"net"
	"net/http"
	"strings"
)

// access client ip from request
func getClientIP(req *http.Request) string {
	clientIP := req.Header.Get("X-Forwarded-For")
	if clientIP != "" {
		clientIP = strings.TrimSpace(strings.Split(clientIP, ",")[0])
		if clientIP == "" {
			clientIP = strings.TrimSpace(req.Header.Get("X-Real-Ip"))
		}
		if clientIP != "" {
			return clientIP
		}
	}
	if ip, _, err := net.SplitHostPort(strings.TrimSpace(req.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}
