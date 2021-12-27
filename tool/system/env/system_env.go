package env

import "nodepanels-tool/util"

func GetEnv() {
	util.PrintResult(util.ExecLinuxCmd("env"))
}
