package logger

import (
	"accounts.sidooh/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var ClientLog = &log.Logger{
	Out: nil,
}

var ServerLog = &log.Logger{
	Out: nil,
}

func Init() {
	ClientLog = log.New()
	ServerLog = log.New()

	env := viper.GetString("APP_ENV")

	if env != "TEST" {
		ClientLog.SetOutput(utils.GetLogFile("client.log"))
		ServerLog.SetOutput(utils.GetLogFile("server.log"))
	}
}
