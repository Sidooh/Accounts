package main

import (
	"accounts.sidooh/server"
	"os"
)

func main() {
	_, ok := os.LookupEnv("JWT_KEY")
	if !ok {
		panic("JWT_KEY is not set")
	}

	e, port, s := server.Setup()

	e.Logger.Fatal(e.StartH2CServer(":"+port, s))
}
