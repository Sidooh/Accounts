package logger

import (
	"accounts.sidooh/util"
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
		ClientLog.SetOutput(util.GetLogFile("client.log"))
		ServerLog.SetOutput(util.GetLogFile("server.log"))
	}
}
