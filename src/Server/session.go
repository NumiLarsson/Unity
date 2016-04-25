package main

import (
	"fmt"
	"time"
)

func Session(cServer *Connection) {

	// Create channels between session and listenermanager
	cManager, cManagerExt := makeConnection()

	go createListenerManager(cManagerExt)

	i := 0
	for i < 100 {
		time.Sleep(500 * time.Millisecond)
		select {
		case data := <-cServer.read:

			// Receive info to spawn new listener
			// Should this be a go-routine?
			// data.action should contain the port that the new listener should use
			fmt.Println("Session: Read from server: ", data.action)
			cManager.write <- data

		// Send response back to server
		case userdata := <-cManager.read:
			fmt.Printf("Session: New data from manager %d\n", userdata.action)
			cServer.write <- userdata

		default:
			fmt.Println("Session: Nothing to do")

		}

		i++

	}

}
