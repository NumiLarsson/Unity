package asteroids

import (
	"fmt"
	//"net" //Used for the unimportable Listener
	//"strconv"
	//"time"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
	MaxPlayers     int
	CurrentPlayers int
	CurrentPort    int
	input          chan (Data)
	listeners      []*Listener
	Players        []Player
	//World           [][]int
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

	manager.listeners = make([]*Listener, 0)

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

//NewObject creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) NewObject() int {

	fmt.Println("Creating new object in listener manager")
	//Creation of the listener
	newListener := NewListener(manager.CurrentPort)
	//Insert the created listener to listenerList
	manager.listeners = append(manager.listeners, newListener)

	manager.incrementCurrentPlayers()

	return manager.getNextPort()
}

// collectPlayerPositions collect all player positions and return an array of them
func (manager *ListenerManager) collectPlayerPositions() []Player {
	playerList := make([]Player, 0)

	for _, listener := range manager.listeners {

		var player = listener.getPlayer()
		playerList = append(playerList, player)
	}

	return playerList
}

// getObjects returns an array of playerpositions
// TODO CHANGE listener.Player....
func (manager *ListenerManager) getObjects() []Player {

	return manager.Players
}

// ReceiveFromSession handles data from session
func (manager *ListenerManager) ReceiveFromSession( /*world *World*/ ) {

	//Range over the array of listeners, sending the info from session
	//To each of the listeners via SendToClient-function

	for _, listener := range manager.listeners {
		if listener.ID != "" {
			//			listener.SendToClient(world)

			fmt.Println("todo "+listener.ID, listener.Port)
		}
	}

}

// SendToClient broadcasts world-info to every listener
func (manager *ListenerManager) sendToClient(world *World) {
	for _, listener := range manager.listeners {
		if listener.ID != "" {
			//Function call where the world-info
			//Is sent to each listener in the list
		}
	}
}

/*func (listener *listener.Listener) sendToClient(world *World) {
	fmt.Println("Todo" + listener.ID)
}*/

/*func main() {
	//Main used for testing, will be removed upon final product
	fmt.Println()

	world := new(World)
	lisManager := NewListenerManager(10, 9000)

	port := lisManager.NewObject()
	fmt.Println()

	port2 := lisManager.NewObject()

	fmt.Println()
	fmt.Println(port, lisManager.listenerList[0].ID)
	fmt.Println(port2, lisManager.listenerList[1].ID)
	fmt.Println(len(lisManager.listenerList))
	fmt.Println()
	lisManager.ReceiveFromSession(world)
}*/
