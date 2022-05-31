package process

import (
	"encoding/json"
	"github.com/shirou/gopsutil/v3/process"
	"nodepanels-tool/command"
	"sort"
	"strconv"
	"strings"
)

type ProcessInfo struct {
	Name                      string `json:"name"`
	Cmd                       string `json:"cmd"`
	Cwd                       string `json:"cwd"`
	Exe                       string `json:"exe"`
	CreateTime                int64  `json:"createTime"`
	Foreground                bool   `json:"foreground"`
	Nice                      int32  `json:"nice"`
	NumCtxSwitchesVoluntary   int64  `json:"voluntary"`
	NumCtxSwitchesInvoluntary int64  `json:"involuntary"`
	NumThreads                int32  `json:"numThreads"`
	OpenFiles                 int    `json:"openFiles"`
	Status                    string `json:"status"`
	Username                  string `json:"username"`
}

type ProcessUsage struct {
	Name       string  `json:"name"`
	CpuPercent float64 `json:"cpu"`
	MemPercent float32 `json:"mem"`
	DiskWrite  uint64  `json:"write"`
	DiskRead   uint64  `json:"read"`
	Cmd        string  `json:"cmd"`
	Pid        int32   `json:"pid"`
}

type ProcessUsageSlice []ProcessUsage

func (p ProcessUsageSlice) Len() int {
	return len(p)
}

func (p ProcessUsageSlice) Less(i, j int) bool {
	if p[i].CpuPercent == p[j].CpuPercent {
		return p[i].MemPercent > p[j].MemPercent
	} else {
		return p[i].CpuPercent > p[j].CpuPercent
	}
}

func (p ProcessUsageSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func GetProcessesList() {

	var list []ProcessUsage

	processes, _ := process.Processes()
	for _, val := range processes {
		name, _ := val.Name()
		if strings.Index(strings.ToLower(name), "nodepanels-tool") < 0 {
			cmd, _ := val.Cmdline()
			pid := val.Pid
			cpuPercent, _ := val.CPUPercent()
			memPercent, _ := val.MemoryPercent()

			processUsage := ProcessUsage{}
			processUsage.Name = name
			processUsage.Cmd = cmd
			if strings.Index(strconv.FormatFloat(cpuPercent, 'f', -1, 64), "0.0") == 0 {
				processUsage.CpuPercent = float64(0)
			} else {
				processUsage.CpuPercent = cpuPercent
			}
			processUsage.MemPercent = memPercent
			processUsage.Pid = pid

			list = append(list, processUsage)
		}
	}
	sort.Sort(ProcessUsageSlice(list))
	if len(list) >= 30 {
		list = list[0:30]
	}

	result, _ := json.Marshal(list)

	command.PrintResult(string(result))
}

func GetProcessInfo() {

	processId, _ := strconv.Atoi(command.GetCommandParam())

	process, _ := process.NewProcess(int32(processId))

	if process != nil {

		processInfo := ProcessInfo{}
		cmd, _ := process.Cmdline()
		name, _ := process.Name()
		cwd, _ := process.Cwd()
		exe, _ := process.Exe()
		createTime, _ := process.CreateTime()
		foreground, _ := process.Foreground()
		nice, _ := process.Nice()
		numCtxSwitches, _ := process.NumCtxSwitches()
		numThreads, _ := process.NumThreads()
		openFiles, _ := process.OpenFiles()
		status, _ := process.Status()
		username, _ := process.Username()

		processInfo.Cmd = cmd
		processInfo.Name = name
		processInfo.Cwd = cwd
		processInfo.Exe = exe
		processInfo.CreateTime = createTime
		processInfo.Foreground = foreground
		processInfo.Nice = nice
		if numCtxSwitches != nil {
			processInfo.NumCtxSwitchesVoluntary = numCtxSwitches.Voluntary
			processInfo.NumCtxSwitchesInvoluntary = numCtxSwitches.Involuntary
		}
		processInfo.NumThreads = numThreads
		processInfo.OpenFiles = len(openFiles)
		processInfo.Status = status[0]
		processInfo.Username = username

		msg, _ := json.Marshal(processInfo)

		command.PrintResult(string(msg))
	}
}
