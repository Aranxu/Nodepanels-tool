package speedtest

import (
	"bufio"
	"encoding/json"
	"github.com/gookit/goutil/fsutil"
	"io"
	"nodepanels-tool/util"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

func SpeedTest() {

	speedtestFileName := "np-speedtest-cli"
	speedtestDownloadUrl := "https://nodepanels-file-1256221051.cos.accelerate.myqcloud.com/speedtest/speedtest-linux-" + runtime.GOARCH

	if !fsutil.PathExists(filepath.Join(util.Exepath(), speedtestFileName)) {
		util.Download(speedtestDownloadUrl, filepath.Join(util.Exepath(), speedtestFileName))
		os.Chmod(filepath.Join(util.Exepath(), speedtestFileName), 0777)
	}

	var cmd *exec.Cmd

	nodeIds := strings.Split(util.GetParam(), " ")

	for _, value := range nodeIds {

		cmd = exec.Command("sh", "-c", "ExecStart=env HOME=/tmp/ "+filepath.Join(util.Exepath(), speedtestFileName)+" --accept-license -s "+value)

		stdout, _ := cmd.StdoutPipe()
		stderr, _ := cmd.StderrPipe()

		if err := cmd.Start(); err != nil {
		}

		go asyncLog(value, stdout)
		go asyncLog(value, stderr)

		if err := cmd.Wait(); err != nil {
		}
	}

	os.Remove(filepath.Join(util.Exepath(), speedtestFileName))
}

func asyncLog(nodeId string, std io.ReadCloser) error {
	reader := bufio.NewReader(std)
	for {
		resultMap := ResultMap{}
		resultMap.NodeId = nodeId

		readString, err := reader.ReadBytes('\n')

		if err != nil || err == io.EOF {
			return err
		}
		if strings.Contains(string(readString), "error") {
			resultMap.Latency = "-1"
			SpeedtestSendBack(resultMap)
		} else if strings.Contains(string(readString), "ms") {
			latency := strings.Split(strings.Split(strings.ReplaceAll(strings.TrimSpace(string(readString)), " ", ""), ":")[1], "ms")[0]
			resultMap.Latency = latency
			SpeedtestSendBack(resultMap)
			if "performance-net-speedtest-ping" == util.GetToolType() {
				SpeedtestKill()
			}
		} else if strings.Contains(string(readString), "Download") {
			download := strings.Split(strings.Split(strings.ReplaceAll(strings.TrimSpace(string(readString)), " ", ""), "Download:")[1], "Mbps")[0]
			resultMap.Download = download
			SpeedtestSendBack(resultMap)
		} else if strings.Contains(string(readString), "Upload") {
			upload := strings.Split(strings.Split(strings.ReplaceAll(strings.TrimSpace(string(readString)), " ", ""), "Upload:")[1], "Mbps")[0]
			resultMap.Upload = upload
			SpeedtestSendBack(resultMap)
		}

	}
}

func SpeedtestSendBack(resultMap ResultMap) {
	resultMsg, _ := json.Marshal(resultMap)
	util.PrintResult(string(resultMsg))
}
