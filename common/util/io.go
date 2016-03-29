package util

import (
	// "fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

/*获取当前文件执行的路径*/
func GetCurrDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	splitstring := strings.Split(path, "/")
	size := len(splitstring)
	return strings.Join(splitstring[0:size-1], "/")
}

//文件夹或文件是否存在
func IsExistFileOrDir(filename string) bool {
	var exist = true
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

// 获得可执行程序所在目录
func ExecutableDir() (string, error) {
	pathAbs, err := filepath.Abs(os.Args[0])
	if err != nil {
		return "", err
	}
	return filepath.Dir(pathAbs), nil
}
