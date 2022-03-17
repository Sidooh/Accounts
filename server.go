package main

import (
	"accounts.sidooh/server"
	"github.com/spf13/viper"
	"log"
)

//var echoServer = new(echo.Echo)

func setupConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("Fatal error config file: ", err)
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
	setupConfig(".")

	jwtKey := viper.GetString("JWT_KEY")
	if len(jwtKey) == 0 {
		panic("JWT_KEY is not set")
	}

	echoServer, port, s := server.Setup()

	echoServer.Logger.Fatal(echoServer.StartH2CServer(":"+port, s))
}
