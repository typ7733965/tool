package env

import (
	"github.com/typ7733965/tool/comfunc/tool"
	"github.com/typ7733965/tool/config"
	"os"
)

type AppInfo struct {
	HostName string // 从系统获取的主机名
	Ip       string // 从系统获取的 pod ip，我们的运行环境是 k8s，所以这里是 pod ip
	Stage    string // 运行平台，从配置获取
	Name     string // 应用名，从配置获取
	Dev      bool   // 是否是开发环境
}

func InitEnv(app *config.App) *AppInfo {
	hostName, _ := os.Hostname()
	appInfo := &AppInfo{
		Name:     app.Name,
		Stage:    app.Stage,
		Dev:      app.Dev,
		Ip:       tool.GetIp(),
		HostName: hostName,
	}
	return appInfo
}
