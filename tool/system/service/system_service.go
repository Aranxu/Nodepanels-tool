package service

import (
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"strings"
)

func GetService() {
	var startup = util.ExecLinuxCmd("systemctl | grep '.service' | awk '{print $1,$3,$4,$5}'")
	if strings.Index(startup, "command not found") < 0 {
		command.PrintResult(startup)
	}
}
