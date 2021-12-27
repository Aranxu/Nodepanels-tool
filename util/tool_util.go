package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	ServerId string             `json:"serverId"`
	Warning  Warning            `json:"warning"`
	Monitor  Monitor            `json:"monitor"`
	Command  map[string]Command `json:"command"`
}

type Warning struct {
	Switch int         `json:"switch"`
	Rule   WarningRule `json:"rule"`
}

type WarningRule struct {
	Cpu WarningRuleCpu `json:"cpu"`
	Mem WarningRuleMem `json:"mem"`
}

type WarningRuleCpu struct {
	Switch   int `json:"switch"`
	Value    int `json:"value"`
	Duration int `json:"duration"`
	Count    int `json:"count"`
}

type WarningRuleMem struct {
	Switch   int `json:"switch"`
	Value    int `json:"value"`
	Duration int `json:"duration"`
	Count    int `json:"count"`
}

type Monitor struct {
	Rule MonitorRule `json:"rule"`
}

type MonitorRule struct {
	Process []string `json:"process"`
}

type Command struct {
	Timeout int  `json:"timeout"`
	Stop    bool `json:"stop"`
}

func GetToolType() string {
	return strings.Replace(os.Args[1], "-", "", 1)
}

func GetHostId() string {
	f, _ := ioutil.ReadFile(Exepath() + "/config")
	config := Config{}
	json.Unmarshal(f, &config)

	return strings.Split(config.ServerId, "\n")[0]
}

func GetConfig() Config {
	f, err := ioutil.ReadFile(Exepath() + "/config")
	if err != nil {
		return Config{}
	}

	c := Config{}
	err = json.Unmarshal(f, &c)
	if err != nil {
		return Config{}
	}
	return c
}

func GetParam() string {
	_, tempFileExist := os.Stat(filepath.Join(Exepath(), GetToolType()+"-"+os.Args[2]+".temp"))
	if tempFileExist == nil {
		paramByte, _ := ioutil.ReadFile(filepath.Join(Exepath(), GetToolType()+"-"+os.Args[2]+".temp"))
		return string(paramByte)
	} else {
		return ""
	}
}

func DelParam() {
	os.Remove(filepath.Join(Exepath(), GetToolType()+"-"+os.Args[2]+".temp"))
}

func PrintResult(msg string) {
	msg = strings.ReplaceAll(msg, "\\", "\\\\")
	msg = strings.ReplaceAll(msg, "\n", "\\n")
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	fmt.Println("{\"toolType\":\"" + GetToolType() + "\",\"serverId\":\"" + GetHostId() + "\",\"msg\":\"" + msg + "\"}")
}

func PrintSuccess() {
	fmt.Println("{\"toolType\":\"" + GetToolType() + "\",\"serverId\":\"" + GetHostId() + "\",\"msg\":\"SUCCESS\"}")
}

func PrintError(msg string) {
	msg = strings.ReplaceAll(msg, "\\", "\\\\")
	msg = strings.ReplaceAll(msg, "\n", "\\n")
	msg = strings.ReplaceAll(msg, "\"", "\\\"")
	fmt.Println("{\"toolType\":\"" + GetToolType() + "\",\"serverId\":\"" + GetHostId() + "\",\"msg\":\"ERROR:" + msg + "\"}")
	fmt.Println("{\"toolType\":\"" + GetToolType() + "\",\"serverId\":\"" + GetHostId() + "\",\"msg\":\"ERROR\"}")
}

func PrintEnd() {
	fmt.Println("{\"toolType\":\"" + GetToolType() + "\",\"serverId\":\"" + GetHostId() + "\",\"msg\":\"END\"}")
}

func CheckCompleteness() bool {
	_, e := os.Stat(filepath.Join(Exepath(), "config"))
	if e != nil {
		return false
	}

	f, _ := ioutil.ReadFile(Exepath() + "/config")
	if !strings.Contains(string(f), "\"serverId\": \"") {
		return false
	}

	return true
}
