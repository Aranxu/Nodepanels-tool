package host

import (
	"github.com/shirou/gopsutil/v3/host"
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"runtime"
)

func GetHostname() {
	infoStat, _ := host.Info()
	command.PrintResult(infoStat.Hostname)
}

func SetHostname() {
	if runtime.GOOS == "linux" {
		util.ExecLinuxCmd("hostnamectl set-hostname " + command.GetCommandParam())
	} else if runtime.GOOS == "windows" {
		util.ExecWindowsCmd("WMIC computersystem where caption=\"%computername%\" rename \"" + command.GetCommandParam() + "\"")
	}
	command.PrintSuccess()
}
