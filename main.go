package main

import (
	"fmt"
	"nodepanels-tool/command"
	"nodepanels-tool/tool/file"
	"nodepanels-tool/tool/monitor"
	"nodepanels-tool/tool/performance/speedtest"
	"nodepanels-tool/tool/performance/trace"
	"nodepanels-tool/tool/probe"
	"nodepanels-tool/tool/process"
	"nodepanels-tool/tool/system/dns"
	"nodepanels-tool/tool/system/env"
	"nodepanels-tool/tool/system/host"
	"nodepanels-tool/tool/system/service"
	"nodepanels-tool/tool/system/startup"
	"nodepanels-tool/tool/system/time"
	"nodepanels-tool/tool/system/yum"
	"os"
)

//go:generate goversioninfo -arm -o=resource_windows.syso -icon=favicon.ico

func main() {
	defer func() {
		err := recover()
		if err != nil {
			command.PrintError(fmt.Sprintf("%s", err))
			command.PrintEnd()
		}
	}()

	version := "v2.0.2"

	if len(os.Args) > 1 {

		if os.Args[1] == "-version" {
			fmt.Print(version)
			return
		}

		switch command.GetCommandType() {
		case "performance-net-speedtest-ping":
			speedtest.SpeedTest()
		case "performance-net-speedtest-all":
			speedtest.SpeedTest()
		case "probe-upgrade":
			probe.ProbeUpgrade()
		case "process-list":
			process.GetProcessesList()
		case "process-info":
			process.GetProcessInfo()
		case "process-monitor-switch":
			monitor.SetMonitorProcessRule()
		case "system-hostname-get":
			host.GetHostname()
		case "system-hostname-set":
			host.SetHostname()
		case "system-dns-get":
			dns.GetDns()
		case "system-dns-set":
			dns.SetDns()
		case "system-dns-backup":
			dns.BackupDns()
		case "system-dns-restore":
			dns.RestoreDns()
		case "system-yum-get":
			yum.GetYum()
		case "system-yum-set":
			yum.SetYum()
		case "system-yum-file-set":
			yum.SetYumFile()
		case "system-yum-backup":
			yum.BackupYum()
		case "system-yum-restore":
			yum.RestoreYum()
		case "system-time-info-get":
			time.GetTimeInfo()
		case "system-time-zone-set":
			time.SetTimeZone()
		case "system-time-set":
			time.SetTime()
		case "system-env-get":
			env.GetEnv()
		case "system-startup-get":
			startup.GetStartup()
		case "system-service-get":
			service.GetService()
		case "file-list":
			file.List()
		case "file-new-file":
			file.FileNewFile()
		case "file-new-folder":
			file.FileNewFolder()
		case "file-attr":
			file.FileAttr()
		case "file-md5":
			file.FileMd5()
		case "file-sha1":
			file.FileSha1()
		case "file-permission":
			file.FilePermission()
		case "file-permission-set":
			file.FilePermissionSet()
		case "file-delete":
			file.FileDelete()
		case "file-copy":
			file.FileCopy()
		case "file-move":
			file.FileMove()
		case "file-rename":
			file.FileRename()
		case "file-edit":
			file.FileEdit()
		case "file-size":
			file.FileSize()
		case "file-trash-add":
			file.FileTrashAdd()
		case "file-trash-recover":
			file.FileTrashRecover()
		case "file-trash-delete":
			file.FileTrashDelete()
		case "file-trash-list":
			file.FileTrashList()
		case "file-upload":
			file.FileUpload()
		case "file-download":
			file.FileDownload()
		case "trace":
			trace.Trace()
		default:
			command.PrintError("WRONG PARAM")
		}
	} else {
		fmt.Println("Wrong parameter.")
	}

}
