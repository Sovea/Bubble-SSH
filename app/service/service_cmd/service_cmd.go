package service_cmd

import (
	"fmt"
	"os/exec"
)

func CMD_Exec(name string, options ...string) (bool,error) {
	cmd := exec.Command(name, options...)
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		return false,err
	}
	return true,nil
}
