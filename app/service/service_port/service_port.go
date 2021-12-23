package service_port
import (
	"net"
	"time"
)
func CheckPort(port string) (bool) {
	// var err error
	
	tcpAddress, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:" + port)

	
	for i := 0; i < 3; i++ {
		listener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			if i == 3 {
				return true
			}
			continue
		} else {
			listener.Close()
			break
		}
	}

	return false
}
func Raw_connect(host string, ports []string) bool {
    for _, port := range ports {
        timeout := time.Second
        conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
        if err != nil {
			return false
        }
        if conn != nil {
            defer conn.Close()
			return true
        }
    }
	return false
}
