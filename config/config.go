package config

import (
	"github.com/gookit/goutil/jsonutil"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

var BinPath = ExePath()

type Config struct {
	ServerId string  `json:"serverId"`
	Monitor  Monitor `json:"monitor"`
	Usage    int64   `json:"usage"` //循环获取实时使用率的end时间，防止无限重复调用
}

type Monitor struct {
	Rule MonitorRule `json:"rule"`
}

type MonitorRule struct {
	Process []string `json:"process"`
}

// GetSid 获取服务器id
func GetSid() string {
	return GetConfig().ServerId
}

// GetConfig 获取配置文件
func GetConfig() Config {
	c := Config{}
	f, _ := ioutil.ReadFile(filepath.Join(BinPath, "config.json"))
	jsonutil.Decode(f, &c)
	return c
}

// SetConfig 设置配置文件
func SetConfig(c Config) {
	json, _ := jsonutil.EncodePretty(c)
	ioutil.WriteFile(filepath.Join(BinPath, "config.json"), json, 0666)
}

func ExePath() string {
	file, err := exec.LookPath(os.Args[0])
	if err != nil {
		return ""
	}
	path, err := filepath.Abs(file)
	if err != nil {
		return ""
	}
	i := strings.LastIndex(path, "/")
	if i < 0 {
		i = strings.LastIndex(path, "\\")
	}
	if i < 0 {
		return ""
	}
	return path[0 : i+1]
}
