package service_port
import (
	"net"
	"time"
)
func CheckPort(port string) bool {
	var err error
	
	tcpAddress, err := net.ResolveTCPAddr("tcp4", ":" + port)
	if err != nil {
		return true
	}
	
	for i := 0; i < 3; i++ {
		listener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			time.Sleep(time.Duration(100) * time.Millisecond)
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