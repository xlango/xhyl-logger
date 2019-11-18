package conf

import (
	"encoding/json"
	"io/ioutil"
	"logconnection/logger"
	"logconnection/utils"
	"os"
	"runtime"
)

var GlobalConfig *ConfigInfo

type ConfigInfo struct {
	BaseConfig
	EsConfig
	ConsulConfig
	ConsulCheckConfig
}

type BaseConfig struct {
	Port int //服务启动端口
}
type EsConfig struct {
	EsHost string //es连接
}
type ConsulConfig struct {
	ConsulAddress         string //consul ip:port
	ConsulRegisterId      string //consul 服务节点的名称
	ConsulRegisterName    string // 服务名称
	ConsulRegisterPort    int    // 服务端口
	ConsulRegisterTags    string // tag，可以为空,服务注册版本
	ConsulRegisterAddress string // 当前服务 IP
}
type ConsulCheckConfig struct {
	ConsulCheckPort                int    // 健康检查 Port
	ConsulCheckTimeout             string // 健康检查 请求超时时间 : "3s"
	ConsulCheckInterval            string // 健康检查 间隔 : "3s"
	DeregisterCriticalServiceAfter string // 健康检查 check失败后30秒删除本服务，注销时间，相当于过期时间
}

func InitConfig() {
	//解决windows无法debug启动的问题
	currentDir := utils.GetCurrentExeDir()
	if len(currentDir) > 0 {
		currentDir = currentDir + string(os.PathSeparator)
	}

	//获取主机可用cpu数，配置程序使用cpu核数
	cpuNumber := runtime.NumCPU()
	runtime.GOMAXPROCS(cpuNumber)

	//如果配置文件未配置，则使用默认配置
	GlobalConfig = &ConfigInfo{
		EsConfig: EsConfig{
			EsHost: "http://127.0.0.1:9200",
		},
		ConsulConfig: ConsulConfig{
			ConsulAddress:         "127.0.0.1:8500",
			ConsulRegisterId:      "logserver",
			ConsulRegisterName:    "logserver",
			ConsulRegisterTags:    "v0001",
			ConsulRegisterAddress: "127.0.0.1",
		},
		ConsulCheckConfig: ConsulCheckConfig{
			ConsulCheckTimeout:             "3s",
			ConsulCheckInterval:            "5s",
			DeregisterCriticalServiceAfter: "30s",
		},
	}

	//初始化日志
	logger.InitLogger(currentDir)

	//读取config.json配置文件相关配置
	configPath := currentDir + "config.json"
	logger.LogInfo("Config path: " + configPath)
	configContent, err := ioutil.ReadFile(configPath)
	if err != nil {
		logger.LogError(err)
	}
	err = json.Unmarshal(configContent, GlobalConfig)
	if err != nil {
		logger.LogError(err)
	}
}
