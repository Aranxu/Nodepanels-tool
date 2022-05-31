package util

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"nodepanels-tool/command"
	"os"
	"strings"
	"unsafe"
)

func PostJson(url string, jsonParam []byte) string {

	request, err := http.NewRequest("POST", url, bytes.NewReader(jsonParam))
	if err != nil {
		command.PrintError(err.Error())
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		command.PrintError(err.Error())
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		command.PrintError(err.Error())
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func PutFile(url string, filePath string) string {
	client := &http.Client{}
	fileContents, _ := ioutil.ReadFile(filePath)
	request, _ := http.NewRequest("PUT", url, strings.NewReader(string(fileContents)))
	request.ContentLength = int64(len(string(fileContents)))

	resp, err := client.Do(request)
	if err != nil {
		return "ERROR:-1"
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "ERROR:-1"
	}
	str := (*string)(unsafe.Pointer(&respBytes))
	return *str
}

func Download(url string, target string) {
	res, _ := http.Get(url)
	newFile, _ := os.Create(target)
	io.Copy(newFile, res.Body)
	defer res.Body.Close()
	defer newFile.Close()
}
