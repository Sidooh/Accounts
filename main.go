package main

import (
	"accounts.sidooh/clients"
	"accounts.sidooh/db"
	"accounts.sidooh/models/account"
	"accounts.sidooh/models/invite"
	"accounts.sidooh/models/security_question"
	"accounts.sidooh/models/security_question_answer"
	"accounts.sidooh/models/user"
	"accounts.sidooh/server"
	"accounts.sidooh/util"
	"accounts.sidooh/util/cache"
	"accounts.sidooh/util/logger"
	"github.com/spf13/viper"
)

func main() {
	util.SetupConfig(".")

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	if viper.GetInt("INVITE_LEVEL_LIMIT") < 1 {
		panic("INVITE_LEVEL_LIMIT is not set")
	}

	logger.Init()
	db.Init()
	defer db.Close()
	//TODO: Ensure in production this doesn't mess up db
	// TODO: Add a script file that accepts fresh migrate args from cmd
	if viper.GetBool("MIGRATE_DB") {
		err := db.Connection().AutoMigrate(
			user.Model{},
			account.ModelWithUser{},
			invite.ModelWithAccountAndInvite{},
			security_question.Model{},
			security_question_answer.ModelWithAccountAndQuestion{},
		)

		if err != nil {
			panic("failed to auto-migrate")
		}
	}

	cache.Init()
	clients.Init()

	echoServer, port, s := server.Setup()

	// TODO: Review using H2C - cleartext server
	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
