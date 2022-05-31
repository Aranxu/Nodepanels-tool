package trace

import (
	"bufio"
	"io"
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"os/exec"
)

func Trace() {

	//check traceroute
	util.ExecLinuxCmd("traceroute")
	if util.ExecLinuxCmd("echo $?") != "0" {
		command.PrintError("ERROR")
		return
	}

	cmd := exec.Command("sh", "-c", "traceroute -n -I -w 1 "+command.GetCommandParam())

	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
	}

	go asyncLog(stdout)
	go asyncLog(stderr)

	if err := cmd.Wait(); err != nil {
	}

	command.PrintEnd()

}

func asyncLog(std io.ReadCloser) error {
	reader := bufio.NewReader(std)
	for {
		readString, err := reader.ReadBytes('\n')

		if err != nil || err == io.EOF {
			return err
		}
		command.PrintResult(string(readString))
	}
}
