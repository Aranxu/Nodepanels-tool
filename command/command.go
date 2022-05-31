package command

import (
	"fmt"
	"github.com/gookit/goutil/jsonutil"
	"io/ioutil"
	"nodepanels-tool/config"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var BinPath = ExePath()

type Command struct {
	ServerId string      `json:"serverId"`
	Page     string      `json:"page"`
	Tool     CommandTool `json:"tool"`
}

type CommandTool struct {
	Version string `json:"version"`
	Type    string `json:"type"`
	Param   string `json:"param"`
}

func GetCommand() Command {
	command := Command{}
	temp, _ := ioutil.ReadFile(filepath.Join(BinPath, "temp", os.Args[1]+".temp"))
	jsonutil.Decode(temp, &command)
	return command
}

func GetCommandType() string {
	return GetCommand().Tool.Type
}

func GetCommandParam() string {
	return GetCommand().Tool.Param
}

func PrintResult(msg string) {
	msg = strings.ReplaceAll(msg, "\\", "\\\\")
	msg = strings.ReplaceAll(msg, "\n", "\\n")
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	fmt.Println("{\"pid\":\"" + GetCommand().Page + "\"," + "\"toolType\":\"" + GetCommand().Tool.Type + "\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"" + msg + "\"}")
}

func PrintSuccess() {
	fmt.Println("{\"pid\":\"" + GetCommand().Page + "\"," + "\"toolType\":\"" + GetCommand().Tool.Type + "\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"SUCCESS\"}")
}

func PrintError(msg string) {
	msg = strings.ReplaceAll(msg, "\\", "\\\\")
	msg = strings.ReplaceAll(msg, "\n", "\\n")
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	fmt.Println("{\"toolType\":\"" + GetCommand().Tool.Type + "\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"ERROR:" + msg + "\"}")
	fmt.Println("{\"toolType\":\"" + GetCommand().Tool.Type + "\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"ERROR\"}")
}

func PrintEnd() {
	fmt.Println("{\"pid\":\"" + GetCommand().Page + "\"," + "\"toolType\":\"" + GetCommand().Tool.Type + "\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"END\"}")
}

func ExePath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return ""
	}
	return path[0 : i+1]
}
