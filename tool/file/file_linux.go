package file

import (
	"encoding/json"
	"github.com/gookit/goutil/fsutil"
	"io/ioutil"
	"nodepanels-tool/command"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func List() {
	if fsutil.PathExists(command.GetCommandParam()) {
		var fileInfoList []FileInfo
		fileList, _ := ioutil.ReadDir(command.GetCommandParam())
		for i := range fileList {
			fileStat, _ := os.Stat(filepath.Join(command.GetCommandParam(), fileList[i].Name()))
			var username string
			u, err := user.LookupId(strconv.FormatUint(uint64(fileStat.Sys().(*syscall.Stat_t).Uid), 10))
			if err == nil {
				username = u.Username
			}
			var groupname string
			g, err := user.LookupGroupId(strconv.FormatUint(uint64(fileStat.Sys().(*syscall.Stat_t).Gid), 10))
			if err == nil {
				groupname = g.Name
			}
			fileInfo := FileInfo{
				Name:       fileList[i].Name(),
				Size:       fileList[i].Size(),
				ModifyTime: fileList[i].ModTime().Unix(),
				Nlink:      int(fileStat.Sys().(*syscall.Stat_t).Nlink),
				Mode:       strings.ToLower(fileList[i].Mode().String()),
				User:       username,
				Group:      groupname,
			}
			fileInfoList = append(fileInfoList, fileInfo)
		}
		fileInfoJson, _ := json.Marshal(fileInfoList)
		command.PrintResult(string(fileInfoJson))
	} else {
		command.PrintResult("NOTEXIST")
	}
}
