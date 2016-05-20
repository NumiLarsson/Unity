package main

import (
	"fmt"
	"os"

	"./asteroids"
)

func main() {

	argsWithoutProg := os.Args[1:]
	inactivityLimit := 60
	inDebugMode := true
	inFakeMode := false

	fmt.Println("\nWelcome to the back-end of Asteroids 1.0")
	fmt.Println("Server starting…")
	fmt.Print("Parameters: ")

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-s" {
		inDebugMode = false
		fmt.Print("Default, silent\n\n")
	}

	if len(argsWithoutProg) > 0 && argsWithoutProg[0] == "-f" {
		inFakeMode = true
		inactivityLimit = 5
		fmt.Print("Fake users, verbose\n\n")
	}

	var server = asteroids.CreateServer(inDebugMode, inactivityLimit)

	if inFakeMode == true {
		server.Listen(server.CreateFakeUser())
	} else {
		server.Listen(make(chan asteroids.Data))
	}

	asteroids.DebugPrint(fmt.Sprintln("[MAIN] Dead"))

}
