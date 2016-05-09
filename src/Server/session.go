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
	asteroids  []*asteroid

	// For external communication
	write channels
	read  channels
}

func (session *session) loop() {

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

			fmt.Printf("Session: Read from manager %s\n", userdata.action)
			session.write.server <- userdata

		default:
			// Nothing
		}

	}

}

// Session …
func Session(serverConn *Connection, startPort int, players int) {

	session := new(session)
	session.maxPlayers = players

	session.write.server = serverConn.write
	session.read.server = serverConn.read
	session.asteroids = make([]*asteroid, 0)

	// CREATE MANAGERS
	// TODO: Loopify

	session.createManagers(startPort)

	// RESPOND TO SERVER
	//

	session.write.server <- Data{"Session created", 0}
	session.loop()

}

// Setup managers and their respective connections to/from session
func (session *session) createManagers(startPort int) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, _ := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	go createListenerManager(fromPlayers, startPort)
	go createAsteroidManager(toAsteroids, session.asteroids)

}
