package startup

import (
	"fmt"
	"nodepanels-tool/config"
)

func GetStartup() {
	fmt.Println("{\"toolType\":\"system-startup-get\",\"serverId\":\"" + config.GetSid() + "\",\"msg\":\"END\"}")
}
