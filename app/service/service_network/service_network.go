package service_network
 
import (
 "fmt"
 "os/exec"
//  "time"
)
 
func NetWorkStatus() bool {
 cmd := exec.Command("ping", "baidu.com", "-c", "1", "-W", "5")
//  fmt.Println("NetWorkStatus Start:", time.Now().Unix())
 err := cmd.Run()
//  fmt.Println("NetWorkStatus End :", time.Now().Unix())
 if err != nil {
 fmt.Println(err.Error())
 return false
 } else {
//  fmt.Println("Net Status , OK")
 }
 return true
}
