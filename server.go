package main

import (
	"accounts.sidooh/db"
	"accounts.sidooh/models/account"
	"accounts.sidooh/models/invite"
	"accounts.sidooh/models/user"
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
	//TODO: Ensure in production this doesn't mess up db
	_ = db.Connection().AutoMigrate(user.Model{}, account.ModelWithUser{}, invite.ModelWithAccountAndInvite{})

	echoServer, port, s := server.Setup()

	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
