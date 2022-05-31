package probe

import (
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func ProbeUpgrade() {

	url := strings.Split(command.GetCommandParam(), " ")[1]

	util.Download(url, filepath.Join(util.Exepath(), "nodepanels-probe.temp"))

	if runtime.GOOS == "windows" {
		command.PrintSuccess()
		exec.Command("cmd", "/C", "net stop Nodepanels-probe").Output()
	}
	if runtime.GOOS == "linux" {
		os.Chmod(util.Exepath()+"/nodepanels-probe.temp", 0777)
		os.Rename(util.Exepath()+"/nodepanels-probe.temp", filepath.Join(util.Exepath(), "/nodepanels-probe"))

		command.PrintSuccess()
		exec.Command("sh", "-c", "service nodepanels restart").Output()
	}

}
