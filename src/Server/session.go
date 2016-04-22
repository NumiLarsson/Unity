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
		time.Sleep(1 * time.Second)
		select {
		case data := <-cServer.read:
			// Spawn listener
			fmt.Println("Read from server")
			cManager.write <- data

		case userdata := <-cManager.read:
			fmt.Printf("New data from user %d\n", userdata.result)
			//toListener <- 1
			cServer.write <- userdata

		default:
			fmt.Println("Nothing to do")
		}

		i++

	}

}
