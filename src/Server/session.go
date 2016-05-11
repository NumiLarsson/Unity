package main

import (
	"./listener"
	"fmt"
	"time"
)

// World is a placeholder for the gameboard
//type World int
// TODO: CHANGE THIS
type World struct {
	Players   []listener.Player
	Asteroids []*listener.Asteroid
}

// channels struct used to implement a structured way to handle multiple
// write/read channels for session
type channels struct {
	server    chan (Data)
	players   chan (Data)
	asteroids chan (Data)
}

// session struct stores info regarding players,session managers,
// read/write channels etc.
type session struct {
	players         int
	maxPlayers      int
	world           World
	asteroids       []*asteroid
	asteroidManager *asteroidManager
	listenerManager *ListenerManager
	// For external communication
	write channels
	read  channels
}

// loop is the sessions ....TODO
func (session *session) loop() {

	for {

		fakeTick := time.After(16 * time.Millisecond)

		select {
		case <-fakeTick:
			session.write.asteroids <- Data{"session.tick", 200}
			session.write.players <- Data{"session.tick", 200}
				

		case data := <-session.read.server:

			if data.action == "server.poke" {

				// Check if theres room inside the session
				if session.players < session.maxPlayers {
					session.write.server <- Data{"session.has_room", 200}
				} else {
					session.write.server <- Data{"session.no_room", -1}
				}

			} else {

				// Spawn a new player
				var port = session.listenerManager.NewObject()
				session.players++
				session.write.server <- Data{"session.player_created", port}
			}

		// Send response back to server
		case userdata := <-session.read.players:

			fmt.Printf("Session: Read from manager %s\n", userdata.action)
			session.write.server <- userdata

		}

		// Collect player and asteroid positions
		session.world.Players = session.listenerManager.getObjects()
		// TODO: implement below
	
		//session.world.Asteroids = session.asteroidManager.getObjects()
	
	}

}

// Session …
func Session(serverConn *Connection, startPort int, players int) {

	session := new(session)
	session.maxPlayers = players

	session.write.server = serverConn.write
	session.read.server = serverConn.read
	session.asteroids = make([]*asteroid, 0)

	session.write.server <- Data{"session_created", 200}
	session.createManagers(startPort)


	session.loop()	

}

// createManagers sets up managers and their respective connections to/from session
func (session *session) createManagers(startPort int) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, fromAsteroids := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(fromAsteroids, session.asteroids)
	go session.listenerManager.loop(fromPlayers, session.maxPlayers, startPort)


}


