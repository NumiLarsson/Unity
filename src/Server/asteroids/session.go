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
	players   []*Player
	asteroids []*asteroid
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
	world           World
	asteroids       []*asteroid // TODO remove? do we use this anymore
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
			session.world.players = session.listenerManager.getPlayers()
			session.world.asteroids = session.asteroidManager.getAsteroids()

			session.world.collisionManager()

			// Send collision ids back to asteroid manager
			deathRow, playerCollisions := session.detectCollisions()
			session.asteroidManager.updateDeathRow(deathRow)
			session.listenerManager.handleCollisions(playerCollisions)

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
				var port, newPlayer = session.listenerManager.newPlayer()
				session.currentPlayers++
				session.world.players = append(session.world.players, newPlayer)

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
func (session *session) createManagers(startPort int /*maxPlayers int, maxAsteroids*/) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, fromAsteroids := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.world.worldSize = 400 //session.worldSize
	session.world.players = make([]*Player, 1 /*maxPlayers*/)
	session.world.asteroids = make([]*asteroid, 2 /*maxAsteroids*/)

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(fromAsteroids, session.asteroids)
	go session.listenerManager.loop(fromPlayers, session.maxPlayers, startPort)

}

// =======================TODO: REBUILD IN COLLISION MANAGER? ===========================
// ======================= FIX MORE GENERIC =============================================
// detectCollisions checks each asteroid and stores all asteroids that have collided
// TODO players and use collision manager?
func (session *session) detectCollisions() ([]int, []int) {

	var asteroidCollisions []int
	var playerCollisions []int

	// First check player vs player collisions
	// second check player vs asteroid
	// last check asteroid vs asteroid
	playerCollisions = session.playerPlayerCollision(playerCollisions)
	playerCollisions, asteroidCollisions =
		session.playerAsteroidCollision(playerCollisions, asteroidCollisions)

	asteroidCollisions = session.asteroidAsteroidCollision(asteroidCollisions)

	return asteroidCollisions, playerCollisions

}

func (session *session) asteroidAsteroidCollision(asteroidCollisions []int) []int {
	for _, a1 := range session.world.asteroids {
		for _, a2 := range session.world.asteroids {

			if isCollisionAsteroidAsteroid(a1, a2) &&
				!inList(asteroidCollisions, a1.id) {
				asteroidCollisions = append(asteroidCollisions, a1.id)
			}
		}
	}
	return asteroidCollisions
}

//
func (session *session) playerAsteroidCollision(playerCollisions []int,
	asteroidCollisions []int) ([]int, []int) {
	//var playerCollisions, asteroidCollisions []int

	for _, p := range session.world.players {
		for _, a := range session.world.asteroids {
			if isCollisionPlayerAsteroid(p, a) {
				if !inList(playerCollisions, p.id) {
					playerCollisions = append(playerCollisions, p.id)
				}
				if !inList(asteroidCollisions, a.id) {
					asteroidCollisions = append(asteroidCollisions, a.id)
				}

			}
		}
	}
	return playerCollisions, asteroidCollisions
}

func (session *session) playerPlayerCollision(playerCollisions []int) []int {
	//var playerCollisions []int

	for _, p1 := range session.world.players {
		for _, p2 := range session.world.players {
			if isCollisionPlayerPlayer(p1, p2) && !inList(playerCollisions, p1.id) {
				playerCollisions = append(playerCollisions, p1.id)

			}
		}
	}
	return playerCollisions

}

// isCollisionAsteroidAsteroid checks is if two asteroids are on
// the same position causing a collision
func isCollisionAsteroidAsteroid(a1 *asteroid, a2 *asteroid) bool {

	if a1.id == a2.id {
		return false
	} else if a1.x == a2.x && a1.y == a2.y {
		return true
	}

	return false

}

// isCollisionAsteroidPlayer  TODO some sort of interface to take generic input?
func isCollisionPlayerAsteroid(p *Player, a *asteroid) bool {

	if p.x == a.x && p.y == a.y {
		return true
	}

	return false

}

// TODO some sort of interface to take generic input?
func isCollisionPlayerPlayer(p1 *Player, p2 *Player) bool {

	if p1.id == p2.id {
		return false
	} else if p1.x == p2.x && p1.y == p2.y {
		return true
	}

	return false
}

// inList checks if the item is is already in the list
func inList(list []int, item int) bool {
	for _, current := range list {
		if item == current {
			return true
		}
	}
	return false
}

// =======================TODO: REBUILD IN COLLISION MANAGER? ===========================
// ======================= FIX MORE GENERIC =============================================
