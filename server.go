package main

import (
	"accounts.sidooh/server"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func setupConfig() {
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
	viper.WatchConfig()
}

func main() {
	setupConfig()

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	e, port, s := server.Setup()

	e.Logger.Fatal(e.StartH2CServer(":"+port, s))
}
