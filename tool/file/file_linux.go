package file

import (
	"encoding/json"
	"github.com/gookit/goutil/fsutil"
	"io/ioutil"
	"nodepanels-tool/util"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

func List() {
	if fsutil.PathExists(util.GetParam()) {
		var fileInfoList []FileInfo
		fileList, _ := ioutil.ReadDir(util.GetParam())
		for i := range fileList {
			fileStat, _ := os.Stat(filepath.Join(util.GetParam(), fileList[i].Name()))
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
		util.PrintResult(string(fileInfoJson))
	} else {
		util.PrintResult("NOTEXIST")
	}
}
