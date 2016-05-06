package main

import (
	"fmt"
	"time"
)

type session struct {
	nextPort   int
	endPort    int
	maxPlayers int
	//conn
}

func Session(cServer *Connection, nextPort int, players int) {

	var endPort = nextPort + players

	// Create channels between session and listenermanager
	cManager, cManagerExt := makeConnection()

	go createListenerManager(cManagerExt)
	cServer.write <- Data{"Session created", 0}

	for {
		time.Sleep(500 * time.Millisecond)
		select {
		case data := <-cServer.read:

			// Receive info to spawn new listener
			fmt.Println("Session: Read from server: ", data.action)

			// Should we double check if maxplayer reached?
			if nextPort < endPort {
				cManager.write <- Data{"Create new player", nextPort}
				nextPort++
			} else {
				cServer.write <- Data{"Session full", -1}
			}

		// Send response back to server
		case userdata := <-cManager.read:
			fmt.Printf("Session: Read from manager %d\n", userdata.action)
			cServer.write <- userdata

		}
	}

}
