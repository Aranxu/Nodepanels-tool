package trace

import (
	"bufio"
	"io"
	"nodepanels-tool/util"
	"os/exec"
)

func Trace() {

	//check traceroute
	util.ExecLinuxCmd("traceroute")
	if util.ExecLinuxCmd("echo $?") != "0" {
		util.PrintError("ERROR")
		return
	}

	cmd := exec.Command("sh", "-c", "traceroute -n -I -w 1 "+util.GetParam())

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
		util.PrintResult(string(readString))
	}
}
