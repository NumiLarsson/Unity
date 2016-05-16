package asteroids

import (
	"fmt"
	"time"
)

// World is a placeholder for the gameboard
//type World int
// TODO: CHANGE THIS
type World struct {
	worldSize int
	Players   []*Player //Json only exports exported objects,
	Asteroids []*asteroid //keep these as capital first letter
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
	players         int
	maxPlayers      int
	world           *World
	asteroids       []*asteroid
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

	//session.write.server <- Data{"session_created", 200}
	//This is not using GO so it's 100000% deadlocked.
	
	session.createManagers(startPort)

	go session.loop()

}

// loop is the sessions ....TODO
func (session *session) loop() {

	for {

		tick := time.After(16 * time.Millisecond)
		//TEMP, tick should be 16 * millisecond

		select {
		case <-tick:
			// Collect player and asteroid positions
			session.world.Players = session.listenerManager.getPlayers()
			session.world.Asteroids = session.asteroidManager.getAsteroids()
			
			session.world.Players[0].fakeMovePlayer()
			
			session.world.collisionManager()
			
			// Send collision ids back to asteroid manager
			deathRow := session.detectCollisions()
			session.asteroidManager.updateDeathRow(deathRow)

			//Empty world {}, something is going wrong.
			//session.world.players jsons fine, but world just doesn't
			
			
			//TEMP BROADCAST TO CLIENTS
			session.listenerManager.sendToClient(session.world)
			//TEMP BROADCAST TO CLIENTS
			
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
				var port, newPlayer = session.listenerManager.newPlayer()
				session.players++
				session.world.Players = append(session.world.Players, newPlayer)

				session.write.server <- Data{"session.player_created", port}
			}

		// Send response back to server
		case userdata := <-session.read.players:

			fmt.Printf("Session: Read from manager %s\n", userdata.action)
			session.write.server <- userdata

		}

	}

}

// createManagers sets up managers and their respective connections to/from session
func (session *session) createManagers(startPort int /*maxPlayers int, maxAsteroids int*/) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, fromAsteroids := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.world = new(World)
	session.world.worldSize = session.worldSize
	session.world.Players = make([]*Player, 0/*maxPlayers*/)
	session.world.Asteroids = make([]*asteroid, 1/*maxAsteroids*/)
	 
	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(fromAsteroids, session.asteroids)
	go session.listenerManager.loop(fromPlayers, session.maxPlayers, startPort)
}

func (session *session) detectCollisions() []int {

	var collisions []int

	for _, a1 := range session.world.Asteroids {
		for _, a2 := range session.world.Asteroids {

			if isCollision(a1, a2) && !inList(collisions, a1.ID) {
				collisions = append(collisions, a1.ID)
			}
		}
	}

	return collisions

}

func isCollision(a1 *asteroid, a2 *asteroid) bool {

	if a1.ID == a2.ID {
		return false
	} else if a1.X == a2.X && a1.Y == a2.Y {
		return true
	}

	return false

}

func inList(list []int, item int) bool {
	for _, current := range list {
		if item == current {
			return true
		}
	}
	return false
}
