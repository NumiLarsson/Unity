package asteroids

import (
	"fmt"
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
func (manager *ListenerManager) loop(sessionConn *Connection,
	maxPlayers int, startPort int, players []*Player) {

	manager.init(sessionConn, maxPlayers, startPort, players)
	sessionConn.write <- Data{"l.manager_ready", 200}

	for {

		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				// Read current values
				// TODO: Where should input from user be checked
				manager.handleCollisions()
				manager.players = manager.collectPlayerPositions()

				// Send update + world to players

			} else {
				debugPrint(fmt.Sprintln("[LIST.MAN] Collision!! \n ", msg.action))
				// TODO: remove asteroids who has a collision or hit
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
	maxPlayers int, firstPort int, players []*Player) {

	// TODO fix hardcoded variables
	manager.xMax = 100
	manager.yMax = 100

	manager.maxPlayers = maxPlayers
	manager.currentPlayers = 0
	manager.nextID = 1
	manager.currentPort = firstPort
	manager.input = sessionConn.read
	manager.players = players

	manager.listeners = make([]*Listener, 0)

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

//NewPlayer creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) newPlayer() (int, *Player) {

	debugPrint(fmt.Sprintln("[LIST.MAN] Creating new object in listener manager"))

	//Creation of the listener and listener-player
	listener := newListener()
	listener.init(manager.currentPort)

	player := newPlayer()
	player.init(manager.getNextID(), manager.xMax, manager.yMax)
	listener.player = player

	//Insert the created listener to listenerList
	//Insert the created player to Players
	manager.listeners = append(manager.listeners, listener)
	manager.players = append(manager.players, player)

	manager.incrementCurrentPlayers()

	go listener.startUpListener()

	return manager.getNextPort(), player
}

// getNextID returns the id to be used and sets the next value
func (manager *ListenerManager) getNextID() int {
	defer func() { manager.nextID++ }()
	return manager.nextID
}

// collectPlayerPositions collect all player positions and return an array of them
func (manager *ListenerManager) collectPlayerPositions() []*Player {
	//playerList := make([]Player, 0)
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
func (manager *ListenerManager) sendToClient(world World) {
	for _, listener := range manager.listeners {
		if listener.ID != "" {
			go listener.Write(world)
		}
	}
}

// handleCollisons handles collisons with a player
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
