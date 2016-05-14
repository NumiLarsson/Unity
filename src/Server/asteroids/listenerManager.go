package asteroids

import (
	"fmt"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
	MaxPlayers     int
	CurrentPlayers int
	CurrentPort    int
	input          chan (Data)
	listeners      []*Listener
	Players        []*Player
}

// loop TODO
func (manager *ListenerManager) loop(sessionConn *Connection, maxPlayers int, startPort int) {
	manager.init(sessionConn, maxPlayers, startPort)

	for {

		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				// Read current values
				// TODO: Where should input from user be checked
				manager.Players = manager.collectPlayerPositions()

				// Send update + world to players

			} else {
				fmt.Println("Collision!! \n ", msg.action)
				// TODO: remove asteroids who has a collision or hit
			}
		}
	}
}

// newAsteroidsManager creates a new asteroid manager
func newListenerManager() *ListenerManager {

	fmt.Println("ListenerManager created")
	return new(ListenerManager)

}

//NewListenerManager does exactly what it says, with a cap on maxPlayers
//connected and maxPlayers numbers of ports in a row from firstPort
func (manager *ListenerManager) init(sessionConn *Connection, maxPlayers int, firstPort int) {

	manager.MaxPlayers = maxPlayers
	manager.CurrentPlayers = 0
	manager.CurrentPort = firstPort
	manager.input = sessionConn.read

	manager.listeners = make([]*Listener, 1)

}

// getNextPort calculates the next start port to be used by a session
func (manager *ListenerManager) getNextPort() int {
	var port = manager.CurrentPort
	manager.CurrentPort++
	return port
}

// incrementCurrentPlayers increments currentPlayers
func (manager *ListenerManager) incrementCurrentPlayers() {
	manager.CurrentPlayers++
}

//NewPlayer creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) newPlayer() (int, *Player) {

	fmt.Println("Creating new object in listener manager")
	//Creation of the listener and listener-player
	newListener := NewListener(manager.CurrentPort)

	newPlayer := newListener.player
	//Insert the created listener to listenerList
	//Insert the created player to Players
	//manager.listeners = append(manager.listeners, newListener)
	manager.listeners[0] = newListener;
	//APPEND IS NOT WORKING, STOP USING IT PLEASE
	
	manager.Players = append(manager.Players, newPlayer)

	manager.incrementCurrentPlayers()

	return manager.getNextPort(), newPlayer
}

// collectPlayerPositions collect all player positions and return an array of them
func (manager *ListenerManager) collectPlayerPositions() []*Player {
	//playerList := make([]Player, 0)
	var playerList []*Player
	for _, listener := range manager.listeners {

		var player = listener.getPlayer()
		playerList = append(playerList, player)
		fmt.Println(player.XCord)
	}

	return playerList
}

// getPlayers returns an array of playerpositions
func (manager *ListenerManager) getPlayers() []*Player {
	return manager.Players
}

// SendToClient broadcasts world-info to every listener
func (manager *ListenerManager) sendToClient(world *World) {
	for _, listener := range manager.listeners {
		if listener.ID != "" {
			go listener.Write(world)
		}
	}
}
