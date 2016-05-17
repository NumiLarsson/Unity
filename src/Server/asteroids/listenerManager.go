package asteroids

import (
	"fmt"
	"math/rand"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
	xMax           int
	yMax           int
	maxPlayers     int
	currentPlayers int
	currentPort    int
	nextID         int
	input          chan (Data)
	listeners      []*Listener
	players        []*Player
}

// loop TODO
func (manager *ListenerManager) loop(sessionConn *Connection, maxPlayers int, startPort int) {

	manager.init(sessionConn, maxPlayers, startPort)

	for {
		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				//manager.print()
				manager.handleCollisions()
				manager.players = manager.collectPlayerPositions()

				// Send update + world to players
			}
		}
	}
}

// newAsteroidsManager creates a new asteroid manager
func newListenerManager() *ListenerManager {

	debugPrint(fmt.Sprintln("[LIST.MAN] Created"))
	return new(ListenerManager)

}

//NewListenerManager does exactly what it says, with a cap on maxPlayers
//connected and maxPlayers numbers of ports in a row from firstPort
func (manager *ListenerManager) init(sessionConn *Connection,
	maxPlayers int, firstPort int) {

	// TODO fix hardcoded variables
	manager.xMax = 100
	manager.yMax = 100
	manager.maxPlayers = maxPlayers
	manager.nextID = 1
	manager.currentPort = firstPort
	manager.input = sessionConn.read

	//manager.listeners = make([]*Listener, 0)

	sessionConn.write <- Data{"l.manager_ready", 200}
}

//NewPlayer creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) newPlayer() (int, *Player) {

	debugPrint(fmt.Sprintln("[LIST.MAN] Creating new object in listener manager"))

	//Creation of the listener and listener-player
	listener := newListener()
	listener.init(manager.currentPort)

	listener.player = newPlayer()
	listener.player.init(manager.getNextID(), manager.xMax, manager.yMax)

	//Insert the created listener to listenerList
	//Insert the created player to Players
	manager.listeners = append(manager.listeners, listener)
	manager.players = append(manager.players, listener.player)

	manager.incrementCurrentPlayers()

	go listener.startUpListener()

	return manager.getNextPort(), listener.player
}

// getNextID returns the id to be used and sets the next value
func (manager *ListenerManager) getNextPort() int {
	defer func() { manager.currentPort++ }()
	return manager.currentPort
}

// incrementCurrentPlayers increments currentPlayers
func (manager *ListenerManager) incrementCurrentPlayers() {
	manager.currentPlayers++
}

// getNextID returns the id to be used and sets the next value
func (manager *ListenerManager) getNextID() int {
	defer func() { manager.nextID++ }()
	return manager.nextID
}

//TEMP FAKE func
func (player *Player) fakeMovePlayer() {
	player.X = rand.Intn(5)
	player.Y = rand.Intn(5)
}

// collectPlayerPositions collect all player positions and return an array of them
func (manager *ListenerManager) collectPlayerPositions() []*Player {

	var playerList []*Player
	for _, listener := range manager.listeners {

		var player = listener.getPlayer()
		playerList = append(playerList, player)
	}

	return playerList
}

// getPlayers returns an array of playerpositions
func (manager *ListenerManager) getPlayers() []*Player {
	return manager.players
}

// SendToClient broadcasts world-info to every listener
func (manager *ListenerManager) sendToClient(world *World) {
	for _, listener := range manager.listeners {
		if listener.ID != "" {
			go listener.Write(world)
		}
	}
}

// handleCollisions handles collisons with a player
func (manager *ListenerManager) handleCollisions() {

	for _, player := range manager.players {
		if !player.isAlive() {
			if player.getLives() > 1 {
				player.setAlive()
				player.Lives--
			} else {
				// TODO : remove player
			}
		}
	}
}

//
func (manager *ListenerManager) print() {

	var list []int

	for _, player := range manager.players {
		if !player.isAlive() {
			list = append(list, player.ID)
		}
	}
	if len(list) > 0 {
		debugPrint(fmt.Sprintln("[LIST.MAN] Collision:", list))
	}
}
