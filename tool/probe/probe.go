package probe

import (
	"io/ioutil"
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func ProbeUpgrade() {

	url := command.GetCommandParam()

	util.Download(url, filepath.Join(util.Exepath(), "nodepanels-probe.temp"))

	os.Remove(filepath.Join(util.Exepath(), "config.json"))

	ioutil.WriteFile(filepath.Join(util.Exepath(), "config.json"), []byte("{\"serverId\":\""+command.GetCommand().ServerId+"\"}"), 0666)

	command.DelCommandTempFile()

	if runtime.GOOS == "windows" {
		os.Rename(filepath.Join(util.Exepath(), "nodepanels-probe.temp"), filepath.Join(util.Exepath(), "nodepanels-probe.exe"))

		command.PrintSuccess()
		exec.Command("cmd", "/C", "net stop Nodepanels-probe").Output()
	}
	if runtime.GOOS == "linux" {
		os.Chmod(util.Exepath()+"/nodepanels-probe.temp", 0777)
		os.Rename(util.Exepath()+"/nodepanels-probe.temp", filepath.Join(util.Exepath(), "nodepanels-probe"))

		command.PrintSuccess()
		exec.Command("sh", "-c", "service nodepanels restart").Output()
	}

}
