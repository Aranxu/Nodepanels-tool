package speedtest

import (
	"github.com/shirou/gopsutil/v3/process"
	"strings"
)

type ResultMap struct {
	NodeId   string `json:"nodeId"`
	Latency  string `json:"latency"`
	Download string `json:"download"`
	Upload   string `json:"upload"`
}

func SpeedtestKill() {
	processes, _ := process.Processes()
	for _, val := range processes {
		cmd, _ := val.Cmdline()
		if strings.Contains(cmd, "np-speedtest-cli") {
			val.Kill()
		}
	}
}
