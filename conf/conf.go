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
	MysqlConfig
	RedisConfig
	MongoConfig
	KafkaConfig
	EsConfig
	ConsulConfig
	ConsulCheckConfig
}

type BaseConfig struct {
	Port int //服务启动端口
}
type MysqlConfig struct {
	MysqlUrl          string //mysql连接
	MysqlTbPrefix     string //Mysql表前缀
	MysqlMaxIdleConns int    //Mysql最大空闲连接
	MysqlMaxOpenConns int    //Mysql最大连接
}
type RedisConfig struct {
	RedisHost         string //redis ip:port
	RedisPoolSize     int    //redis连接池大小
	RedisReadTimeout  int    //redis 读数据超时
	RedisWriteTimeout int    //redis 写数据超时
	RedisIdleTimeout  int    //空闲连接时长
}
type MongoConfig struct {
	MongoHost string //mongodb连接
}
type KafkaConfig struct {
	KafkaHosts string //kafka连接
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
		MysqlConfig: MysqlConfig{
			MysqlUrl:          "root:123456@tcp(127.0.0.1:3306)/comprehensive?charset=utf8",
			MysqlTbPrefix:     "",
			MysqlMaxIdleConns: 150,
			MysqlMaxOpenConns: 250,
		},
		RedisConfig: RedisConfig{
			RedisHost:         "127.0.0.1:6379",
			RedisPoolSize:     1000,
			RedisReadTimeout:  100,
			RedisWriteTimeout: 100,
			RedisIdleTimeout:  60,
		},
		MongoConfig: MongoConfig{
			MongoHost: "127.0.0.1:27017",
		},

		KafkaConfig: KafkaConfig{
			KafkaHosts: "127.0.0.1:9092",
		},
		EsConfig: EsConfig{
			EsHost: "http://127.0.0.1:9200",
		},
		ConsulConfig: ConsulConfig{
			ConsulAddress:         "127.0.0.1:8500",
			ConsulRegisterId:      "serverNode_1",
			ConsulRegisterName:    "serverNode",
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

	//初始化mysql,创建表
	//InitMysqlDb()
	InitTable()

	//初始化redis
	//InitRedis()

	//初始化Kafka
	InitKafka()
}
