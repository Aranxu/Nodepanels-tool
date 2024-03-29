package yum

import (
	"github.com/shirou/gopsutil/v3/host"
	"io"
	"io/ioutil"
	"nodepanels-tool/command"
	"nodepanels-tool/util"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func GetYum() {
	var url string
	var path string
	var file []byte

	platform, _, _, _ := host.PlatformInformation()

	if platform == "centos" {
		path = "/etc/yum.repos.d/CentOS-Base.repo"
		file, _ = ioutil.ReadFile("/etc/yum.repos.d/CentOS-Base.repo")
		r, _ := regexp.Compile("(?m).*(baseurl=http).*")
		url = r.FindAllString(string(file), 1)[0]
		url = strings.Split(url, "baseurl=")[1]
		url = strings.Split(url, "/centos")[0]
	} else if platform == "ubuntu" {
		path = "/etc/apt/sources.list"
		file, _ = ioutil.ReadFile("/etc/apt/sources.list")
		r, _ := regexp.Compile("(?m)^(deb http).*")
		url = r.FindAllString(string(file), 1)[0]
		url = strings.Split(url, "deb ")[1]
		url = strings.Split(url, "/ubuntu")[0]
	} else if platform == "debian" {
		path = "/etc/apt/sources.list"
		file, _ = ioutil.ReadFile("/etc/apt/sources.list")
		r, _ := regexp.Compile("(?m)^(deb http).*")
		url = r.FindAllString(string(file), 1)[0]
		url = strings.Split(url, "deb ")[1]
		url = strings.Split(url, "/debian")[0]
	}
	content := strings.ReplaceAll(strings.ReplaceAll(string(file), "\n", "\\n"), "\"", "\\\"")
	result := "{\"path\":\"" + path + "\",\"file\":\"" + content + "\",\"url\":\"" + url + "\"}"
	command.PrintResult(result)
}

func SetYum() {
	param := command.GetCommandParam()

	platform, _, platformVersion, _ := host.PlatformInformation()

	if platform == "centos" {
		platformVersion = strings.Split(platformVersion, ".")[0]
		util.Download("https://file.nodepanels.com/yum/"+param+"/centos/"+platformVersion+"/CentOS-Base.repo", "/etc/yum.repos.d/CentOS-Base.repo")
		os.Chmod("/etc/yum.repos.d/CentOS-Base.repo", 0644)
		command.PrintResult("CLEAN")
		util.ExecLinuxCmd("yum clean all")
		command.PrintResult("MAKECACHE")
		util.ExecLinuxCmd("yum makecache")
	} else if platform == "ubuntu" {
		util.Download("https://file.nodepanels.com/yum/"+param+"/ubuntu/"+platformVersion+"/sources.list", "/etc/apt/sources.list")
		os.Chmod("/etc/apt/sources.list", 0644)
	} else if platform == "debian" {
		platformVersion = strings.Split(platformVersion, ".")[0]
		util.Download("https://file.nodepanels.com/yum/"+param+"/debian/"+platformVersion+"/sources.list", "/etc/apt/sources.list")
		os.Chmod("/etc/apt/sources.list", 0644)
	}

	command.PrintSuccess()
}

func SetYumFile() {
	param := command.GetCommandParam()

	platform, _, _, _ := host.PlatformInformation()

	if platform == "centos" {
		os.WriteFile("/etc/yum.repos.d/CentOS-Base.repo", []byte(param), 0644)
		command.PrintResult("CLEAN")
		util.ExecLinuxCmd("yum clean all")
		command.PrintResult("MAKECACHE")
		util.ExecLinuxCmd("yum makecache")
	} else if platform == "ubuntu" {
		os.WriteFile("/etc/apt/sources.list", []byte(param), 0644)
	} else if platform == "debian" {
		os.WriteFile("/etc/apt/sources.list", []byte(param), 0644)
	}

	command.PrintSuccess()
}

func BackupYum() {
	var srcPath string
	var dstPath string

	platform, _, _, _ := host.PlatformInformation()

	if platform == "centos" {
		srcPath = "/etc/yum.repos.d/CentOS-Base.repo"
		dstPath = filepath.Join(util.Exepath(), "backup", "yum", "CentOS-Base.repo")
	} else if platform == "ubuntu" {
		srcPath = "/etc/apt/sources.list"
		dstPath = filepath.Join(util.Exepath(), "backup", "yum", "sources.list")
	} else if platform == "debian" {
		srcPath = "/etc/apt/sources.list"
		dstPath = filepath.Join(util.Exepath(), "backup", "yum", "sources.list")
	}

	source, _ := os.Open(srcPath)
	defer source.Close()

	if _, err := os.Stat(filepath.Join(util.Exepath(), "backup", "yum")); os.IsNotExist(err) {
		os.MkdirAll(filepath.Join(util.Exepath(), "backup", "yum"), 0777)
	}
	destination, _ := os.Create(dstPath)
	defer destination.Close()

	io.Copy(destination, source)

	command.PrintSuccess()
}

func RestoreYum() {
	var srcPath string
	var dstPath string
	var filename string

	platform, _, _, _ := host.PlatformInformation()

	if platform == "centos" {
		srcPath = filepath.Join(util.Exepath(), "backup", "yum", "CentOS-Base.repo")
		dstPath = "/etc/yum.repos.d/CentOS-Base.repo"
		filename = "CentOS-Base.repo"
	} else if platform == "ubuntu" {
		srcPath = filepath.Join(util.Exepath(), "backup", "yum", "sources.list")
		dstPath = "/etc/apt/sources.list"
		filename = "sources.list"
	} else if platform == "debian" {
		srcPath = filepath.Join(util.Exepath(), "backup", "yum", "sources.list")
		dstPath = "/etc/apt/sources.list"
		filename = "sources.list"
	}

	if _, err := os.Stat(filepath.Join(util.Exepath(), "backup", "yum", filename)); os.IsNotExist(err) {
		command.PrintError("")
	} else {
		source, _ := os.Open(srcPath)
		defer source.Close()

		destination, _ := os.Create(dstPath)
		defer destination.Close()

		io.Copy(destination, source)

		if platform == "centos" {
			command.PrintResult("CLEAN")
			util.ExecLinuxCmd("yum clean all")
			command.PrintResult("MAKECACHE")
			util.ExecLinuxCmd("yum makecache")
		}

		command.PrintSuccess()
	}
}
