package monitor

import (
	"encoding/json"
	"nodepanels-tool/command"
	"nodepanels-tool/config"
)

func SetMonitorProcessRule() {

	var processCmdList []string
	json.Unmarshal([]byte(command.GetCommandParam()), &processCmdList)

	c := config.GetConfig()
	c.Monitor.Rule.Process = processCmdList

	config.SetConfig(c)

	command.PrintSuccess()
}
