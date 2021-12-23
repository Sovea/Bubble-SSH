package service_file
import (
	"os"
)
func Exists(path string) bool {
	_, err := os.Stat(path)    //os.Stat get detail of the file
	if err != nil {
		return !os.IsNotExist(err)
	}
	return true
}
