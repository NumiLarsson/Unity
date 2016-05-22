package asteroids

import (
	"fmt"
	"time"
)

// World is a placeholder for the gameboard
//type World int
// TODO: CHANGE THIS
type World struct {
	WorldSize  int
	Players    []*Player
	Asteroids  []*Asteroid
	Collisions []*Collision
}

// Collision holds the coordinates of a collision
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
	worldSize       int
	currentPlayers  int
	maxPlayers      int
	world           *World
	asteroidManager *asteroidManager
	listenerManager *ListenerManager
	// For external communication
	write channels
	read  channels
}

// Session initiates the sessions values and creates a new go-routine for the session
func Session(serverConn *Connection, startPort int, players int, worldSize int) {

	session := new(session)
	session.maxPlayers = players

	session.worldSize = worldSize
	session.write.server = serverConn.write
	session.read.server = serverConn.read

	session.write.server <- Data{"session_created", 200}

	session.createManagers(startPort)

	go session.loop()

}

// loop is the sessions game loop, sending ticks to it's managers and collecting data and
// calculating collisions and distribute the new world
func (session *session) loop() {

	for {

		tick := time.After(2000 * time.Millisecond)
		//TEMP, tick should be 16 * millisecond

		select {
		case <-tick:
			// Collect player and asteroid positions
			session.world.Players = session.listenerManager.getPlayers()
			session.world.Asteroids = session.asteroidManager.getAsteroids()
			session.world.Players[0].fakeMovePlayer();
			// Calculate collisions
			session.world.collisionManager()
			
			// BROADCAST TO CLIENTS
			session.listenerManager.sendToClient(session.world)

			//session.world.Players[0].fakeMovePlayer()
			//Faking player movement so that I have something to draw

			//Empty world {}, something is going wrong.
			//session.world.players jsons fine, but world just doesn't

			session.write.asteroids <- Data{"session.tick", 200}
			session.write.players <- Data{"session.tick", 200}

		case input := <-session.read.server:

			if input.action == "server.poke" {
				// Check if theres room inside the session
				fmt.Println("POKE")
				if session.currentPlayers < session.maxPlayers {
					session.write.server <- Data{"session.has_room", 200}
				} else {
					session.write.server <- Data{"session.no_room", -1}
				}

			} else {
				// Spawn a new player
				fmt.Println("SPAWN")
				var port, player = session.listenerManager.newPlayer()
				session.currentPlayers++
				session.world.Players = append(session.world.Players, player)

				session.write.server <- Data{"session.player_created", port}
			}
		}
	}
}

// createManagers sets up managers and their respective connections to/from session
func (session *session) createManagers(startPort int /*maxPlayers int, maxAsteroids int*/) {

	toPlayers := MakeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids := MakeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.world = new(World)
	session.world.WorldSize = 400
	//session.world.Players = make([]*Player, 0)
	//session.world.Asteroids = make([]*Asteroid, 0)
	//session.world.Collisions = make([]*Collision, 0)

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(toAsteroids.FlipConnection() /*,session.world.Asteroids*/)
	go session.listenerManager.loop(toPlayers.FlipConnection(),
		session.maxPlayers, startPort)

	// Wait for managers to signal that they are ready
	<-toAsteroids.read
	<-toPlayers.read

}
