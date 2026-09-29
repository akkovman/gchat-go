package ip

import (
	"net"
	"net/http"
	"strings"
)

/*
Functions gets real ip address from X-Forwaded-For header

In nginx.conf:

	proxy_set_header X-Real-IP $remote_addr;
	proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
	proxy_set_header X-Forwarded-Proto $scheme;
*/
func GetClientAddress(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")

	if commaIndex := strings.Index(xff, ","); commaIndex != -1 {
		return strings.TrimSpace(xff[:commaIndex])
	}

	ipStr := strings.TrimSpace(xff)
	if ip := net.ParseIP(ipStr); ip != nil {
		return ip.String()
	}

	return ""
}
