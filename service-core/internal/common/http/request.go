package apphttp

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"
)

func DecodeJSON(r *http.Request, v any) error {
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(v)
}

func ClientIP(r *http.Request) string {
	var ipVal string
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		parts := strings.Split(ip, ",")
		ipVal = strings.TrimSpace(parts[0])
	} else if ip := r.Header.Get("X-Real-IP"); ip != "" {
		ipVal = ip
	} else {
		ipVal = r.RemoteAddr
	}

	if host, _, err := net.SplitHostPort(ipVal); err == nil {
		return host
	}
	return ipVal
}
