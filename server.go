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
	defer db.Close()
	//TODO: Ensure in production this doesn't mess up db
	// TODO: Add a script file that accepts fresh migrate args from cmd
	//_ = db.Connection().AutoMigrate(
	//	user.Model{},
	//	account.ModelWithUser{},
	//	invite.ModelWithAccountAndInvite{},
	//	security_question.Model{},
	//	security_question_answer.ModelWithAccountAndQuestion{},
	//)

	echoServer, port, s := server.Setup()

	// TODO: Review using H2C - cleartext server
	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
