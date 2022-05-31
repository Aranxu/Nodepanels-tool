package file

import (
	"bufio"
	"crypto/md5"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/gookit/goutil/fsutil"
	"io"
	"io/ioutil"
	"net/http"
	"nodepanels-tool/command"
	"nodepanels-tool/config"
	"nodepanels-tool/util"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type FileInfo struct {
	Name       string `json:"n"`
	ModifyTime int64  `json:"t"`
	Size       int64  `json:"s"`
	Nlink      int    `json:"l"`
	Mode       string `json:"m"`
	User       string `json:"u"`
	Group      string `json:"g"`
}

func FileNewFile() {
	fileName := strings.Split(command.GetCommandParam(), ",")[0]
	filePath := strings.Split(command.GetCommandParam(), ",")[1]
	os.Create(filepath.Join(filePath, fileName))
	command.PrintSuccess()
}

func FileNewFolder() {
	fileName := strings.Split(command.GetCommandParam(), ",")[0]
	filePath := strings.Split(command.GetCommandParam(), ",")[1]
	os.Mkdir(filepath.Join(filePath, fileName), 0666)
	command.PrintSuccess()
}

func FileAttr() {
	param := command.GetCommandParam()
	result := ""
	statStr := util.ExecLinuxCmd("stat \"" + param + "\"")
	fileStr := util.ExecLinuxCmd("file -b \"" + param + "\"")
	lsattrStr := util.ExecLinuxCmd("lsattr \"" + param + "\"")
	/*fileStat, _ := os.Stat(GetParam())

	linuxFileAttr := fileStat.Sys().(*syscall.Stat_t)
	time.Unix(linuxFileAttr.Atim.Sec, 0).Unix()
	time.Unix(linuxFileAttr.Mtim.Sec, 0).Unix()
	time.Unix(linuxFileAttr.Ctim.Sec, 0).Unix()*/

	if string(fileStr) != "directory" && lsattrStr != "" {
		result = "{\"stat\":\"" + statStr + "\",\"file\":\"" + fileStr + "\",\"lsattr\":\"" + lsattrStr + "\"}"
	} else {
		result = "{\"stat\":\"" + statStr + "\",\"file\":\"" + fileStr + "\"}"
	}

	command.PrintResult(result)
}

func FileMd5() {
	data, err := ioutil.ReadFile(command.GetCommandParam())
	if err != nil {
		command.PrintError("")
	}
	command.PrintResult(fmt.Sprintf("%x", md5.Sum(data)))
}

func FileSha1() {
	data, err := ioutil.ReadFile(command.GetCommandParam())
	if err != nil {
		command.PrintError("")
	}
	command.PrintResult(fmt.Sprintf("%x", sha1.Sum(data)))
}

func FilePermission() {
	userList := strings.ReplaceAll(util.ExecLinuxCmd("cat /etc/passwd | awk -F':' '{ print $1}'"), "\n", ",")
	groupList := strings.ReplaceAll(util.ExecLinuxCmd("cat /etc/group | cut -d : -f 1"), "\n", ",")
	filePermission := strings.Split(strings.Split(util.ExecLinuxCmd("stat \""+command.GetCommandParam()+"\""), "Access: (")[1], ")")[0]
	command.PrintResult("{\"user\":\"" + userList + "\",\"group\":\"" + groupList + "\",\"permission\":\"" + filePermission + "\"}")
}

func FilePermissionSet() {
	param := strings.Split(command.GetCommandParam(), ",")
	filePath := param[0]
	user := param[1]
	group := param[2]
	permission := param[3]
	containChild := param[4]

	if containChild == "true" {
		util.ExecLinuxCmd("chmod " + permission + " -R \"" + filePath + "\"")
		util.ExecLinuxCmd("chown " + user + ":" + group + " -R \"" + filePath + "\"")
	} else {
		util.ExecLinuxCmd("chmod " + permission + " \"" + filePath + "\"")
		util.ExecLinuxCmd("chown " + user + ":" + group + " \"" + filePath + "\"")
	}
	command.PrintSuccess()
}

func FileDelete() {
	filePaths := command.GetCommandParam()
	for _, filePath := range strings.Split(filePaths, ",") {
		os.RemoveAll(filePath)
	}
	command.PrintSuccess()
}

func FileCopy() {
	tasks := strings.Split(command.GetCommandParam(), ";")
	for _, val := range tasks {
		resFile := strings.Split(val, ",")[0]
		desFile := strings.Split(val, ",")[1]
		copy(resFile, desFile)
	}
	command.PrintSuccess()
}

func FileMove() {
	tasks := strings.Split(command.GetCommandParam(), ";")
	for _, val := range tasks {
		resFile := strings.Split(val, ",")[0]
		desFile := strings.Split(val, ",")[1]
		if !fsutil.PathExists(desFile) {
			os.MkdirAll(desFile, 0666)
		}
		copy(resFile, desFile)
		os.RemoveAll(resFile)
	}
	command.PrintSuccess()
}

func FileRename() {
	resFile := strings.Split(command.GetCommandParam(), ",")[0]
	desFile := strings.Split(command.GetCommandParam(), ",")[1]
	os.Rename(resFile, desFile)
	command.PrintSuccess()
}

func FileEdit() {
	fileEditDto := FileEditDto{}
	json.Unmarshal([]byte(command.GetCommandParam()), &fileEditDto)

	os.WriteFile(fileEditDto.FilePath, []byte(fileEditDto.Content), 0644)
	command.PrintSuccess()
}

type FileEditDto struct {
	FilePath string `json:"filePath"`
	Content  string `json:"content"`
}

func FileSize() {
	command.PrintResult(strconv.FormatInt(GetDirSize(command.GetCommandParam()), 10))
}

func FileTrashAdd() {
	trashFileDto := TrashFileDto{}
	json.Unmarshal([]byte(command.GetCommandParam()), &trashFileDto)

	trashPath := filepath.Join(trashFileDto.RecyclePath, "recycle")
	trashSize := trashFileDto.RecycleSize
	fileList := trashFileDto.TrashFileList

	if _, err := os.Stat(trashPath); os.IsNotExist(err) {
		os.MkdirAll(trashPath, 0666)
	}

	if _, err := os.Stat(filepath.Join(trashPath, "index.json")); os.IsNotExist(err) {
		ioutil.WriteFile(filepath.Join(trashPath, "index.json"), []byte("[]"), 0666)
	}

	var tempSize int64
	for _, file := range fileList {
		tempSize += GetDirSize(file.Path)
	}

	if trashSize <= GetDirSize(trashPath)+tempSize {
		command.PrintResult("LIMIT")
		return
	}

	var trashFileSlice []TrashFile
	indexJson, _ := ioutil.ReadFile(filepath.Join(trashPath, "index.json"))
	json.Unmarshal(indexJson, &trashFileSlice)

	for _, file := range fileList {
		trashFile := TrashFile{
			Name:   filepath.Base(file.Path),
			Path:   file.Path,
			Time:   strconv.FormatInt(time.Now().UnixNano(), 10),
			Size:   GetDirSize(file.Path),
			IsPath: fsutil.IsDir(file.Path),
		}
		trashFileSlice = append(trashFileSlice, trashFile)

		os.MkdirAll(filepath.Join(trashPath, trashFile.Time), 0666)

		util.ExecLinuxCmd("mv \"" + file.Path + "\" \"" + filepath.Join(trashPath, trashFile.Time) + "\"")
	}

	indexJson, _ = json.MarshalIndent(trashFileSlice, "", "\t")
	ioutil.WriteFile(filepath.Join(trashPath, "index.json"), indexJson, 0666)

	command.PrintSuccess()
}

func FileTrashRecover() {
	success := true

	trashFileDto := TrashFileDto{}
	json.Unmarshal([]byte(command.GetCommandParam()), &trashFileDto)

	trashPath := filepath.Join(trashFileDto.RecyclePath, "recycle")
	trashFileList := trashFileDto.TrashFileList

	for _, trashFile := range trashFileList {
		filePath := trashFile.Path
		fileName := trashFile.Name
		fileTime := trashFile.Time
		fileIsPath := trashFile.IsPath

		if _, err := os.Stat(filepath.Join(trashPath, fileTime)); os.IsNotExist(err) {
			success = false
		} else {
			if fileIsPath {
				util.ExecLinuxCmd("cp -rp \"" + filepath.Join(trashPath, fileTime, fileName) + "\" \"" + filePath + "\"")
			} else {
				util.ExecLinuxCmd("cp -p \"" + filepath.Join(trashPath, fileTime, fileName) + "\" \"" + filePath + "\"")
			}
		}
	}

	if !success {
		command.PrintError("")
	}

	FileTrashDelete()
}

func FileTrashDelete() {
	trashFileDto := TrashFileDto{}
	json.Unmarshal([]byte(command.GetCommandParam()), &trashFileDto)

	trashPath := filepath.Join(trashFileDto.RecyclePath, "recycle")
	deleteTrashFileList := trashFileDto.TrashFileList

	if fsutil.PathExists(filepath.Join(trashPath, "index.json")) {

		var trashFileList []TrashFile
		indexJson, _ := ioutil.ReadFile(filepath.Join(trashPath, "index.json"))
		json.Unmarshal(indexJson, &trashFileList)

		var tempTrashFileList []TrashFile

		for _, t := range trashFileList {
			var exist = false
			for _, trashFile := range deleteTrashFileList {
				if t.Time == trashFile.Time {
					exist = true
				}
			}
			if exist {
				os.RemoveAll(filepath.Join(trashPath, t.Time))
			} else {
				tempTrashFileList = append(tempTrashFileList, t)
			}
		}

		indexJson, _ = json.MarshalIndent(tempTrashFileList, "", "\t")
		if string(indexJson) == "null" {
			ioutil.WriteFile(filepath.Join(trashPath, "index.json"), []byte("[]"), 0666)
		} else {
			ioutil.WriteFile(filepath.Join(trashPath, "index.json"), indexJson, 0666)
		}
	}
	command.PrintSuccess()
}

func FileTrashList() {
	trashPath := command.GetCommandParam()
	if fsutil.PathExists(filepath.Join(trashPath, "index.json")) {
		indexJson, _ := ioutil.ReadFile(filepath.Join(trashPath, "index.json"))
		command.PrintResult(string(indexJson))
	}
}

type TrashFileDto struct {
	RecyclePath   string      `json:"recyclePath"`
	RecycleSize   int64       `json:"recycleSize"`
	TrashFileList []TrashFile `json:"fileList"`
}

type TrashFile struct {
	Name   string `json:"name"`
	Path   string `json:"path"`
	Time   string `json:"time"`
	Size   int64  `json:"size"`
	IsPath bool   `json:"isPath"`
}

func FileUpload() {
	var uploadFileDto UploadFileDto
	json.Unmarshal([]byte(command.GetCommandParam()), &uploadFileDto)

	fileInfo, _ := os.Stat(uploadFileDto.FilePath)
	if fileInfo.Size() > 2*1024*1024 {
		command.PrintResult("LIMIT")
		return
	}

	url := util.PostJson(uploadFileDto.AgentUrl+"/cos/url/get", []byte("{\"serverId\":\""+config.GetSid()+"\",\"tempFile\":\""+uploadFileDto.TempFile+"\"}"))

	res, _ := http.Get(url)
	newFile, _ := os.Create(filepath.Join(uploadFileDto.FilePath, uploadFileDto.FileName))
	io.Copy(newFile, res.Body)
	defer res.Body.Close()
	defer newFile.Close()

	command.PrintSuccess()
}

type UploadFileDto struct {
	AgentUrl string `json:"agentUrl"`
	TempFile string `json:"tempFile"`
	FileName string `json:"fileName"`
	FilePath string `json:"filePath"`
}

func FileDownload() {
	var uploadFileDto UploadFileDto
	json.Unmarshal([]byte(command.GetCommandParam()), &uploadFileDto)

	filePath := uploadFileDto.FilePath
	tempPath := strings.ReplaceAll(uuid.New().String(), "-", "")
	fileInfo, _ := os.Stat(filePath)
	if fileInfo.Size() > 2*1024*1024 {
		command.PrintResult("LIMIT")
		return
	}
	url := util.PostJson(uploadFileDto.AgentUrl+"/cos/url/put", []byte("{\"serverId\":\""+config.GetSid()+"\",\"tempFile\":\""+tempPath+"\"}"))
	putResult := util.PutFile(url, filePath)
	if putResult == "ERROR:-1" {
		command.PrintResult(putResult)
	} else {
		command.PrintResult(tempPath)
	}
}

func GetDirSize(path string) int64 {
	var size int64
	filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			size += info.Size()
		}
		return err
	})
	return size
}

func copy(res, des string) {
	if fsutil.IsDir(res) {
		if list, e := ioutil.ReadDir(res); e == nil {
			if len(list) == 0 {
				os.MkdirAll(des, 0777)
			}
			for _, item := range list {
				copy(filepath.Join(res, item.Name()), filepath.Join(des, item.Name()))
			}
		}
	} else {
		if !fsutil.PathExists(filepath.Dir(des)) {
			os.MkdirAll(filepath.Dir(des), 0777)
		}
		file, _ := os.Open(res)
		defer file.Close()
		bufReader := bufio.NewReader(file)
		out, _ := os.Create(des)
		defer out.Close()
		io.Copy(out, bufReader)
	}
}
