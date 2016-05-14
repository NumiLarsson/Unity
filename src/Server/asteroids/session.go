package asteroids

import (
	"fmt"
	"time"
	//"encoding/json"
)

// World is a placeholder for the gameboard
//type World int
// TODO: CHANGE THIS
type World struct {
	worldSize int
	players   []*Player
	asteroids []*asteroid
}

//TEMPSTRUCT since I can't send asteroids
type WorldClient struct {
	players []*Player
	asteroids []*AsteroidClient
}
//TEMPSTRUCT since I can't send asteroids

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

		tick := time.After(160 * time.Millisecond)

		select {
		case <-tick:
			// Collect player and asteroid positions
			session.world.players = session.listenerManager.getPlayers()
			session.world.asteroids = session.asteroidManager.getAsteroids()

			session.world.collisionManager()

			// Send collision ids back to asteroid manager
			deathRow := session.detectCollisions()
			session.asteroidManager.updateDeathRow(deathRow)

			// Broadcast collisions to managers
			
			//TEMP BROADCAST TO CLIENTS
			tempWorld := new(WorldClient)
			tempWorld.players = make([]*Player, 1)
			tempWorld.asteroids = make([]*AsteroidClient, 2)
			
			tempPlayer := new(Player)
			tempPlayer.XCord = 1;
			tempPlayer.YCord = 2;
			tempPlayer.Lives = 3;
			
			tempAsteroid1 := AsteroidClient{1,1,1}
			tempAsteroid2 := AsteroidClient{2,2,2}
			tempWorld.players[0] = tempPlayer
			tempWorld.asteroids[0] = &tempAsteroid1
			tempWorld.asteroids[1] = &tempAsteroid2
			/*jsonArray, err := json.Marshal()
			if err != nil {
				panic(err)
			}
			fmt.Println(string(jsonArray))*/
			session.listenerManager.sendToClient(tempWorld)
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
func (session *session) createManagers(startPort int, /*maxPlayers int, maxAsteroids*/) {

	toPlayers, fromPlayers := makeConnection()
	session.write.players = toPlayers.write
	session.read.players = toPlayers.read

	toAsteroids, fromAsteroids := makeConnection()
	session.write.asteroids = toAsteroids.write
	session.read.asteroids = toAsteroids.read

	session.world = new(World)
	session.world.worldSize = session.worldSize
	session.world.players = make([]*Player, 0/*maxPlayers*/)
	session.world.asteroids = make([]*asteroid, 0/*maxAsteroids*/)

	session.asteroidManager = newAsteroidManager()
	session.listenerManager = newListenerManager()

	go session.asteroidManager.loop(fromAsteroids, session.asteroids)
	go session.listenerManager.loop(fromPlayers, session.maxPlayers, startPort)

}

func (session *session) detectCollisions() []int {

	var collisions []int

	for _, a1 := range session.world.asteroids {
		for _, a2 := range session.world.asteroids {

			if isCollision(a1, a2) && !inList(collisions, a1.id) {
				collisions = append(collisions, a1.id)
			}
		}
	}

	return collisions

}

func isCollision(a1 *asteroid, a2 *asteroid) bool {

	if a1.id == a2.id {
		return false
	} else if a1.x == a2.x && a1.y == a2.y {
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
