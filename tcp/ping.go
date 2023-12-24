package tcp

import (
	"net"
	"strconv"
	"time"
)

func PingTCP(host string, port int, timeoutMs time.Duration) (time.Duration, error) {
	// Start the timer
	start := time.Now()

	// Connect to the host on the specified port and timeout
	conn, err := net.DialTimeout("tcp", host+":"+strconv.Itoa(port), timeoutMs)
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	// Calculate the round-trip time
	rtt := time.Since(start)

	return rtt, nil
}
