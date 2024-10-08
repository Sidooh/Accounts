package main

import (
	"accounts.sidooh/api"
	"accounts.sidooh/pkg/clients"
	"accounts.sidooh/pkg/db"
	"accounts.sidooh/pkg/entities"
	"accounts.sidooh/utils"
	"accounts.sidooh/utils/cache"
	"accounts.sidooh/utils/logger"
	"github.com/spf13/viper"
)

func main() {
	utils.SetupConfig(".")

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
	// TODO: Ensure in production this doesn't mess up db
	// TODO: Add a script file that accepts fresh migrate args from cmd
	if viper.GetBool("MIGRATE_DB") {
		err := db.Connection().AutoMigrate(
			entities.User{},
			entities.AccountWithUser{},
			entities.InviteWithAccountAndInviter{},
			entities.Question{},
			entities.QuestionAnswerWithAccountAndQuestion{},
		)

		if err != nil {
			panic("failed to auto-migrate")
		}
	}

	cache.Init()
	clients.Init()

	echoServer, port, s := api.Setup()

	// TODO: Review using H2C - cleartext server
	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
