package main

import (
	"accounts.sidooh/server"
	"fmt"
	"github.com/spf13/viper"
)

//var echoServer = new(echo.Echo)

func setupConfig() {
	viper.SetConfigName(".env.example")
	viper.SetConfigType("env")
	viper.AddConfigPath(".")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}

	//viper.OnConfigChange(func(e fsnotify.Event) {
	//	fmt.Println("Config file changed:", e.Name)
	//
	//	//TODO: On config change restart server
	//	//fmt.Println("Shutdown server:", e.Name)
	//	//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	//	//defer cancel()
	//	//if err := echoServer.Shutdown(ctx); err != nil {
	//	//	echoServer.Logger.Fatal(err)
	//	//} else {
	//	//	fmt.Println("Restarting server:", e.Name)
	//	//
	//	//	echoServer, port, s := server.Setup()
	//	//	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
	//	//}
	//})
	//viper.WatchConfig()
}

func main() {
	setupConfig()

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	echoServer, port, s := server.Setup()

	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
