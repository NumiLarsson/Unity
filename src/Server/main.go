package main

import (
	"fmt"
	"os"

	"./asteroids"
)

func main() {

	argsWithoutProg := os.Args[1:]
	var inDebugMode = true

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-s" {
		fmt.Println("\nServer starting in silent mode\n======")
	}

	var server = asteroids.CreateServer(inDebugMode)
	//server.Listen(make(chan asteroids.Data) /*server.CreateFakeUser()*/)
	server.Listen(server.CreateFakeUser())

}
