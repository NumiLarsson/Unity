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
	nextPort   int
	endPort    int
	maxPlayers int
	servConn   Connection
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
			if s.nextPort < s.endPort {
				s.write.players <- Data{"Create new player", s.nextPort}
				s.nextPort++
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

func Session(server *Connection, nextPort int, players int) {

	s := new(session)
	s.endPort = nextPort + players
	s.nextPort = nextPort
	s.maxPlayers = players

	s.write.server = server.write
	s.read.server = server.read

	// CREATE MANAGERS
	// TODO: Loopify

	createManagers(s)

	// RESPOND TO SERVER
	//

	s.write.server <- Data{"Session created", 0}
	loop(s)

}

// Setup managers and their respective connections to/from session
func createManagers(s *session) {

	toPlayers, fromPlayers := makeConnection()
	s.write.players = toPlayers.write
	s.read.players = toPlayers.read

	toAsteroids, _ := makeConnection()
	s.write.asteroids = toAsteroids.write
	s.read.asteroids = toAsteroids.read

	go createListenerManager(fromPlayers)
	// go createAsteroidManager(fromAsteroids)

}
