package main

import (
	//"container/list"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Listener struct {
	id   string
	port int
	*Connection
	cManager Connection
}

type ListenerManager struct {
	currentPort    int
	startPort      int
	listenerAmount []*Listener
	*Connection
}

// Create ListenerManager and set comm. channels to session
func createListenerManager(cSession *Connection, startPort int) {
	listenerManager := new(ListenerManager)
	listenerManager.Connection = cSession

	Manager(listenerManager)
}

func createListener(lManager *ListenerManager, port int) /*net.Listener*/ {

	fmt.Println("Creating listener...")

	listener := new(Listener)
	cListener, cListenerExt := makeConnection()
	listener.port = port

	// Set channels to new listener
	// Add to array/linked list
	listener.Connection = cListenerExt
	listener.cManager.write = cListener.write
	listener.cManager.read = cListener.read

	// TODO set up next available arrayindex
	// Should we use slices?
	lManager.listenerAmount = append(lManager.listenerAmount, listener)

	userConnection, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Listener created!")
	}

	Listenerfunc(listener, userConnection)

	//return connection
}

func Listenerfunc(l *Listener, userConnection net.Listener) {
	fmt.Println("Now in listener function!")

	// Send response to listenermanager
	l.write <- Data{"Listener set up", 1}
}

func Manager(lManager *ListenerManager) {

	for {

		select {
		case data := <-lManager.read:
			// Read from session, meaning create new listener
			fmt.Println("Manager: Read from session")
			fmt.Println("Manager: rcvd", data.action)
			go createListener(lManager, lManager.currentPort)
			time.Sleep(1 * time.Second)
			fmt.Println("Size:", len(lManager.listenerAmount))
			lManager.write <- Data{"Listener set up", 1}
		}

		// TODO: Set up how to listen to all listeners?

		//Needed to set sleep so listener would have time to be created
		/*time.Sleep(1 * time.Second)
		if lManager.listenerAmount[0].cManager.read != nil {
			for i, _ := range lManager.listenerAmount {
				select {
				case response := <-lManager.listenerAmount[i].cManager.read:
					fmt.Println("Manager: Read from listener")
					fmt.Println("Manager: sending response back to session")
					lManager.write <- response
				}
			}
		}*/
	}
}
