package startup

import (
	"fmt"
	"nodepanels-tool/util"
)

func GetStartup() {
	fmt.Println("{\"toolType\":\"system-startup-get\",\"serverId\":\"" + util.GetHostId() + "\",\"msg\":\"END\"}")
}
