package asteroids

import (
	"fmt"
	"time"
)

// World is a placeholder for the gameboard
//type World int
// TODO: CHANGE THIS
type World struct {
	worldSize  int
	Players    []*Player
	Asteroids  []*Asteroid
	Collisions []*Collision
}

type Collision struct {
	X int
	Y int
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
	worldSize      int
	currentPlayers int
	maxPlayers     int
	world          World
	//asteroids       []*Asteroid // TODO remove? do we use this anymore
	asteroidManager *asteroidManager
	listenerManager *ListenerManager
	// For external communication
	write channels
	read  channels
}

// Session …TODO rename to init?

func Session(serverConn *Connection, startPort int, players int, worldSize int) {

	session := new(session)
	session.maxPlayers = players

	session.worldSize = worldSize
	session.write.server = serverConn.write
	session.read.server = serverConn.read
	//session.asteroids = make([]*asteroid, 0)

	session.write.server <- Data{"session_created", 200}
	session.createManagers(startPort)

	session.loop()

}

// loop is the sessions ....TODO
func (session *session) loop() {

	for {

		tick := time.After(16 * time.Millisecond)

		select {
		case <-tick:

			// Collect player and asteroid positions
			session.world.Players = session.listenerManager.getPlayers()
			session.world.Asteroids = session.asteroidManager.getAsteroids()

			//session.world.collisionManager()
			//session.detectCollisions()

			session.world.collisionManager()
			session.listenerManager.sendToClient(session.world)
			// Send collision ids back to asteroid manager

			//session.asteroidManager.updateDeathRow(deathRow)
			//session.listenerManager.handleCollisions(playerCollisions)

			// Broadcast collisions to managers

			//TEMP BROADCAST TO CLIENTS
			//session.listenerManager.sendToClient(session.world)
			//TEMP BROADCAST TO CLIENTS

			session.write.asteroids <- Data{"session.tick", 200}
			session.write.players <- Data{"session.tick", 200}

		case data := <-session.read.server:

			if data.action == "server.poke" {
				// Check if theres room inside the session
				if session.currentPlayers < session.maxPlayers {
					session.write.server <- Data{"session.has_room", 200}
				} else {
					session.write.server <- Data{"session.no_room", -1}
				}

			} else {

				// Spawn a new player
				var port, player = session.listenerManager.newPlayer()
				session.currentPlayers++
				session.world.Players = append(session.world.Players, player)

				session.write.server <- Data{"session.player_created", port}
			}

		// Send response back to server
		case userdata := <-session.read.players:

			debugPrint(fmt.Sprintf("Session: Read from manager %s\n", userdata.action))
			session.write.server <- userdata

		}

	}

}

// createManagers sets up managers and their respective connections to/from session
func (session *session) createManagers(startPort int /*maxPlayers int, maxAsteroids*/) {

	toPlayers := MakeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids := MakeConnection()

	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.world.worldSize = 400                             //session.worldSize
	session.world.Players = make([]*Player, 0 /*maxPlayers*/) // HÄR GJORDES ÄNDRING
	session.world.Asteroids = make([]*Asteroid, 0 /*maxAsteroids*/)

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(toAsteroids.FlipConnection(), session.world.Asteroids)
	go session.listenerManager.loop(toPlayers.FlipConnection(),
		session.maxPlayers, startPort, session.world.Players)

	// Wait for managers to signal that they are ready
	<-toAsteroids.read
	<-toPlayers.read
}
