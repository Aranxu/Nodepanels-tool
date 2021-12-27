package util

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func Exepath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return ""
	}
	return path[0 : i+1]
}

func ExecLinuxCmd(cmd string) string {
	output, _ := exec.Command("sh", "-c", cmd).Output()
	return strings.TrimRight(string(output), "\n")
}

func ExecWindowsCmd(cmd string) string {
	output, _ := exec.Command("cmd", "/C", cmd).Output()
	return string(output)
}
