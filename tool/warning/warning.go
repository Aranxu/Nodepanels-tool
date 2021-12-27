package warning

import (
	"encoding/json"
	"io/ioutil"
	"nodepanels-tool/util"
	"strconv"
	"strings"
)

func SetWarningRule() {
	param := util.GetParam()

	totalSwitch := strings.Split(param, "|")[0]

	cpuSwitch := strings.Split(strings.Split(param, "|")[1], ";")[0]
	cpuValue := strings.Split(strings.Split(param, "|")[1], ";")[1]
	cpuDuration := strings.Split(strings.Split(param, "|")[1], ";")[2]

	c := util.GetConfig()

	c.Warning.Switch, _ = strconv.Atoi(totalSwitch)
	c.Warning.Rule.Cpu.Switch, _ = strconv.Atoi(cpuSwitch)
	c.Warning.Rule.Cpu.Value, _ = strconv.Atoi(cpuValue)
	c.Warning.Rule.Cpu.Duration, _ = strconv.Atoi(cpuDuration)
	c.Warning.Rule.Cpu.Count = 0

	memSwitch := strings.Split(strings.Split(param, "|")[2], ";")[0]
	memValue := strings.Split(strings.Split(param, "|")[2], ";")[1]
	memDuration := strings.Split(strings.Split(param, "|")[2], ";")[2]

	c.Warning.Rule.Mem.Switch, _ = strconv.Atoi(memSwitch)
	c.Warning.Rule.Mem.Value, _ = strconv.Atoi(memValue)
	c.Warning.Rule.Mem.Duration, _ = strconv.Atoi(memDuration)
	c.Warning.Rule.Mem.Count = 0

	data, _ := json.MarshalIndent(c, "", "\t")
	ioutil.WriteFile(util.Exepath()+"/config", data, 0666)

	util.PrintSuccess()
}
