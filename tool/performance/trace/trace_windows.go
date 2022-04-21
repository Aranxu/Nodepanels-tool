package trace

import (
	"bufio"
	"io"
	"nodepanels-tool/util"
	"os/exec"
	"strings"
)

func Trace() {

	//check traceroute
	util.ExecWindowsCmd("tracert")
	if util.ExecWindowsCmd("echo %errorlevel%") != "0" {
		util.PrintError("ERROR")
		return
	}

	cmd := exec.Command("cmd", "/C", "tracert -d -h 30 -w 1 "+util.GetParam())

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
	}

	go asyncLog(stdout)
	go asyncLog(stderr)

	if err := cmd.Wait(); err != nil {
	}

	util.PrintEnd()

}

func asyncLog(std io.ReadCloser) error {
	reader := bufio.NewReader(std)
	for {
		readString, err := reader.ReadBytes('\n')

		if err != nil || err == io.EOF {
			return err
		}
		util.PrintResult(strings.TrimRight(strings.TrimRight(string(readString), "\n"), "\r"))
	}
}
