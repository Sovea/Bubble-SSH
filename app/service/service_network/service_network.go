package service_network
 
import (
 "fmt"
 "os/exec"
)
 
func NetWorkStatus() bool {
 cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
 err := cmd.Run()
 if err != nil {
 fmt.Println(err.Error())
 return false
 }
 return true
}
