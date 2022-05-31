package env

import (
	"nodepanels-tool/command"
	"nodepanels-tool/util"
)

func GetEnv() {
	command.PrintResult(util.ExecLinuxCmd("env"))
}
