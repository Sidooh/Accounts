package utils

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	Environment string
	JwtKey      string
	Port        int
}

func SetupConfig(path string) {
	// Set the path to look for the configurations file
	viper.AddConfigPath(path)

	// Set the file name of the configurations file
	viper.SetConfigName(".env")
	viper.SetConfigType("env")

	// Enable VIPER to read Environment Variables
	viper.AutomaticEnv()
	//
	viper.SetDefault("INVITE_EXPIRY", 48)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
		} else {
			// Config file was found but another error was produced
			log.Fatal("Fatal error: ", err)
		}
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
