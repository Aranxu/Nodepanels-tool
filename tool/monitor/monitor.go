package monitor

import (
	"encoding/json"
	"io/ioutil"
	"nodepanels-tool/util"
)

func SetMonitorProcessRule() {

	var processCmdList []string
	json.Unmarshal([]byte(util.GetParam()), &processCmdList)

	c := util.GetConfig()
	c.Monitor.Rule.Process = processCmdList

	data, _ := json.MarshalIndent(c, "", "\t")
	ioutil.WriteFile(util.Exepath()+"/config", data, 0666)

	util.PrintSuccess()
}
