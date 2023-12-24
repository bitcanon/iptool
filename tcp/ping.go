package tcp

import (
	"net"
	"strconv"
	"time"

	"golang.org/x/net/ipv4"
)

func PingTCP(host string, port int, ttl int, timeoutMs time.Duration) (time.Duration, error) {
	start := time.Now()

	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), timeoutMs)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	p := ipv4.NewConn(conn)
	if err := p.SetTTL(ttl); err != nil {
		return 0, err
	}
	rtt := time.Since(start)

	return rtt, nil
}
