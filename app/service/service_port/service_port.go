package service_port
import (
	"net"
	"time"
	// "os/exec"
	// "fmt"
)
func CheckPort(port string) (bool) {
	// var err error
	
	tcpAddress, _ := net.ResolveTCPAddr("tcp4", "127.0.0.1:" + port)

	
	for i := 0; i < 3; i++ {
		listener, err := net.ListenTCP("tcp", tcpAddress)
		if err != nil {
			// time.Sleep(time.Duration(100) * time.Millisecond)
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
            // fmt.Println("Connecting error:", err)
			return false
        }
        if conn != nil {
            defer conn.Close()
			return true
            // fmt.Println("Opened", net.JoinHostPort(host, port))
        }
    }
	return false
}
// func CheckPort(port int) bool {
//     checkStatement := fmt.Sprintf("netstat -tunlp | grep %d", port)
//     output, err := exec.Command("sh", "-c", checkStatement).CombinedOutput()
//     if err != nil {
//         return false   // err != nil 说明端口没被占用
//     }

//     if len(output) > 0 {
//         return true
//     }

//     return false
// }