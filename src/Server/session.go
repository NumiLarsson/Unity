package main

import (
	"fmt"
)

// Change state by shifting x bits
type World int

type channels struct {
	server    chan (Data)
	players   chan (Data)
	asteroids chan (Data)
}

type session struct {
	players    int
	maxPlayers int
	world      World

	// For external communication
	write channels
	read  channels
}

func loop(s *session) {

	for {

		select {
		case data := <-s.read.server:

			// Receive info to spawn new listener
			fmt.Println("Session: Read from server: ", data.action)

			// Should we double check if maxplayer reached?
			if s.players < s.maxPlayers {
				s.write.players <- Data{"Create new player", 100}

				s.players++
			} else {
				s.write.server <- Data{"Session full", -1}
			}

		// Send response back to server
		case userdata := <-s.read.players:

			fmt.Printf("Session: Read from manager %d\n", userdata.action)
			s.write.server <- userdata
		}

	}

}

//
func Session(serverConn *Connection, startPort int, players int) {

	s := new(session)
	s.maxPlayers = players

	s.write.server = serverConn.write
	s.read.server = serverConn.read

	// CREATE MANAGERS
	// TODO: Loopify

	createManagers(s, startPort)

	// RESPOND TO SERVER
	//

	s.write.server <- Data{"Session created", 0}
	loop(s)

}

// Setup managers and their respective connections to/from session
func createManagers(s *session, startPort int) {

	toPlayers, fromPlayers := makeConnection()
	s.write.players = toPlayers.write
	s.read.players = toPlayers.read

	toAsteroids, _ := makeConnection()
	s.write.asteroids = toAsteroids.write
	s.read.asteroids = toAsteroids.read

	go createListenerManager(fromPlayers, startPort)
	// go createAsteroidManager(fromAsteroids)

}
