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
	asteroids  []*asteroids

	// For external communication
	write channels
	read  channels
}

func loop(session *session) {

	for {

		select {
		case data := <-session.read.server:

			// Receive info to spawn new listener
			fmt.Println("Session: Read from server: ", data.action)

			// Should we double check if maxplayer reached?
			if session.players < session.maxPlayers {
				session.write.players <- Data{"Create new player", 100}

				session.players++
			} else {
				session.write.server <- Data{"Session full", -1}
			}

		// Send response back to server
		case userdata := <-session.read.players:

			fmt.Printf("Session: Read from manager %d\n", userdata.action)
			session.write.server <- userdata
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
func createManagers(session *session, startPort int) {

	//toPlayers, fromPlayers := makeConnection()
	//s.write.players = toPlayers.write
	//s.read.players = toPlayers.read

	toAsteroids, _ := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	//go createListenerManager(fromPlayers, startPort)
	go createAsteroidManager(toAsteroids, session.asteroids)

}
