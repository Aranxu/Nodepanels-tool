package host

import (
	"github.com/shirou/gopsutil/v3/host"
	"nodepanels-tool/util"
	"runtime"
)

func GetHostname() {
	infoStat, _ := host.Info()
	util.PrintResult(infoStat.Hostname)
}

func SetHostname() {
	if runtime.GOOS == "linux" {
		util.ExecLinuxCmd("hostnamectl set-hostname " + util.GetParam())
	} else if runtime.GOOS == "windows" {
		util.ExecWindowsCmd("WMIC computersystem where caption=\"%computername%\" rename \"" + util.GetParam() + "\"")
	}
	util.PrintSuccess()
}
