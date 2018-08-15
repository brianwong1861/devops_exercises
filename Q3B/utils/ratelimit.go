package utils

import (
        "time"
        _ "fmt"
        "strings"
		"net/http"
		"os"
		"strconv"
)

var threshold int64

type RemoteIP struct {
	StartTime int64 `json:"start_time"`
	IPAddr string `json:"ip_addr"`
}
var remoteIP RemoteIP

func (p *RemoteIP) RateLimiting(w http.ResponseWriter, r *http.Request) bool {
	ipAddrTrimmed := strings.Split(r.RemoteAddr, ":")[0] //Split Remote IP ADDR into IP and Port
	// message = make(map[string]string)
	s1, _ := strconv.Atoi(os.Getenv("RATE_LIMITE_THRESHOLD"))
	threshold := int64(s1)

	if ( remoteIP.IPAddr == "" ) && ( remoteIP.StartTime == 0 ) {
		remoteIP.IPAddr = ipAddrTrimmed
		remoteIP.StartTime = time.Now().Unix()
	} else if ( time.Now().Unix() - remoteIP.StartTime ) <= threshold {
		w.Write([]byte("Exceeded rate limit\n"))
		remoteIP.StartTime = time.Now().Unix() 
		return false
	}
	remoteIP.StartTime = time.Now().Unix() 
	return true
}

func RateLimitingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val := remoteIP.RateLimiting(w, r)
		if val == true {
			next.ServeHTTP(w, r)
		}
        // Call the next handler, which can be another middleware in the chain, or the final handler.
    })
}
