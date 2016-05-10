package main

import (
	"fmt"
	"net" //Used for the unimportable Listener

	"./listener"
	//"strconv"
	"time"
)

//export GOPATH=$HOME/Golang2

//type World struct{}

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
	MaxPlayers     int
	CurrentPlayers int
	CurrentPort    int
	input          chan (Data)
	listeners      []*listener.Listener
	//World           [][]int
}

//Does not import correctly for some reason
type Listener struct {
	socket net.Listener
	ID     string
	Port   string
	//Connection
}

func (manager *ListenerManager) loop(sessionConn *Connection, maxPlayers int, startPort int) {
	manager.init(sessionConn, maxPlayers, startPort)

	for {
		time.Sleep(250 * time.Millisecond)
		//fmt.Println(manager.CurrentPlayers)
		//manager.ReceiveFromSession()
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
	//lisManager := new(ListenerManager)
	manager.MaxPlayers = maxPlayers
	manager.CurrentPlayers = 0
	manager.CurrentPort = firstPort
	//manager.NewObject()
	manager.listeners = make([]*listener.Listener, 0)

}

// getNextPort calculates the next start port to be used by a session
func (manager *ListenerManager) getNextPort() int {
	var port = manager.CurrentPort
	manager.CurrentPort++
	return port
}

func (manager *ListenerManager) incrementCurrentPlayers() {
	manager.CurrentPlayers++
}

//NewObject creates a new listener for the listener manager, used to connect to a new player.
func (manager *ListenerManager) NewObject() int {

	fmt.Println("Creating new object in listener manager")
	//Creation of the listener
	newListener := listener.NewListener(manager.CurrentPort)
	//Insert the created listener to listenerList
	manager.listeners = append(manager.listeners, newListener)

	manager.incrementCurrentPlayers()

	return manager.getNextPort()
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

func (listener *Listener) sendToClient(world *World) {
	fmt.Println("Todo" + listener.ID)
}

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
