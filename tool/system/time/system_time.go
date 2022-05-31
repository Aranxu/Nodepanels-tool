package time

import (
	"nodepanels-tool/command"
	"nodepanels-tool/util"
)

func GetTimeInfo() {
	var time = util.ExecLinuxCmd("date +\"%Y-%m-%d %H:%M:%S\"")
	var timestamp = util.ExecLinuxCmd("date +%s")
	var timezone = util.ExecLinuxCmd("ls -il /etc | grep localtime | awk '{print $12}' | awk -F zoneinfo/ '{print $2}'")
	var timezoneNum = util.ExecLinuxCmd("date +\"%z\"")
	command.PrintResult("{\"timezoneNum\":\"" + timezoneNum + "\",\"timezone\":\"" + timezone + "\",\"time\":\"" + time + "\",\"timestamp\":\"" + timestamp + "\"}")
}

func SetTimeZone() {
	util.ExecLinuxCmd("ln -snf /usr/share/zoneinfo/" + command.GetCommandParam() + " /etc/localtime")
	command.PrintSuccess()
}

func SetTime() {
	util.ExecLinuxCmd("date -s \"" + command.GetCommandParam() + "\"")
	util.ExecLinuxCmd("hwclock -w")
	util.ExecLinuxCmd("ln -snf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime")
	command.PrintSuccess()
}
