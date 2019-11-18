package logger

import (
	"github.com/cihub/seelog"
	"io/ioutil"
	"log"
	"strings"
)

func LogDebug(v ...interface{}) {
	seelog.Debug(v)
}
func LogInfo(v ...interface{}) {
	seelog.Info(v)
}
func LogError(v ...interface{}) {
	seelog.Error(v)
}

func InitLogger(currentDir string) {

	config, err := ioutil.ReadFile(currentDir + "seelog.xml")
	if err != nil {
		log.Fatalln(err)
	}
	configStr := string(config)
	newConfigStr := strings.Replace(configStr, "filename=\"./", "filename=\""+currentDir, -1)
	Logger, err := seelog.LoggerFromConfigAsString(newConfigStr)
	if err != nil {
		log.Fatalln(err)
	}
	seelog.ReplaceLogger(Logger)
}
