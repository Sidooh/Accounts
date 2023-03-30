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
	logger := viper.GetString("LOGGER")

	if env != "TEST" {
		if logger == "GCP" {
			ClientLog.SetFormatter(NewGCEFormatter(false))
			ServerLog.SetFormatter(NewGCEFormatter(false))
		} else {
			ClientLog.SetOutput(utils.GetLogFile("client.log"))
			ServerLog.SetOutput(utils.GetLogFile("server.log"))
		}
	}
}
