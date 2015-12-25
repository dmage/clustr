package service

import (
	"fmt"
	"os"
)

func getProcExe(pid int) (string, error) {
	filename := fmt.Sprintf("/proc/%d/exe", pid)
	return os.Readlink(filename)
}
