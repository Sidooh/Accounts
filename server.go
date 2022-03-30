package main

import (
	"accounts.sidooh/db"
	"accounts.sidooh/server"
	"accounts.sidooh/util"
	"github.com/spf13/viper"
)

func main() {
	util.SetupConfig(".")

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	db.Init()

	echoServer, port, s := server.Setup()

	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
