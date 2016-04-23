package main

import (
	"fmt"
	"time"
)

func Session(cServer *Connection) {

	//fromListener := make(chan int)
	//toListener := make(chan int
	//cManager := new(Connection)
	//cManager.write = make(chan Data)
	//cManager.read = make(chan Data)

	cManager, cManagerExt := makeConnection()

	// TODO!
	// Swap the write/read before sending
	go Manager(cManagerExt)

	i := 0
	for i < 100 {
		time.Sleep(500 * time.Millisecond)
		select {
		case data := <-cServer.read:

			// Receive info to spawn new listener
			// Should this be a go-routine?
			//
			fmt.Println("Session Read from server: ", data.action)
			cManager.write <- data

		case userdata := <-cManager.read:
			fmt.Printf("Session New data from user %d\n", userdata.result)
			//toListener <- 1
			cServer.write <- userdata

		default:
			fmt.Println("Session Nothing to do")
		}

		i++

	}

}
