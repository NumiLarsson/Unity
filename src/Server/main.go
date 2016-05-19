package main

import (
	"fmt"
	"os"

	"./asteroids"
)

func main() {

	argsWithoutProg := os.Args[1:]
	var inDebugMode = true
	var inFakeMode = false

	fmt.Println("\nWelcome to the back-end of Asteroids 1.0")
	fmt.Println("Server starting…")
	fmt.Print("Parameters: ")

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-s" {
		inDebugMode = false
		fmt.Print("Default, silent\n\n")
	}

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-f" {
		inFakeMode = true
		fmt.Print("Fake users, verbose\n\n")
	}

	var server = asteroids.CreateServer(inDebugMode)

	if inFakeMode == true {
		server.Listen(server.CreateFakeUser())
	} else {
		server.Listen(make(chan asteroids.Data))
	}

}
