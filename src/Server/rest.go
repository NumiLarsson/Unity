package main

import (
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
	listenerAmount []Listener
	*Connection
}

// Create ListenerManager and set comm. channels to session
func createListenerManager(cSession *Connection) {
	listenerManager := new(ListenerManager)
	listenerManager.Connection = cSession

	Manager(listenerManager)
}

func createListener(lManager *ListenerManager, port int) net.Listener {

	listener := new(Listener)
	cListener, cListenerExt := makeConnection()
	listener.port = port
	listener.Connection = cListenerExt

	// Set channels to new listener
	// Add to array/linked list
	// TODO: find a smart way to store the listener in the listenermanager array
	// but with the correct channels( *listener has the reversed channels)

	go Listenerfunc(listener)

	time.Sleep(500 * time.Millisecond)

	listener.cManager.write = cListener.write
	listener.cManager.read = cListener.read

	fmt.Println("createListener created!")
	fmt.Println("createListener channels", listener.cManager.write, listener.cManager.read)

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Listener created!")
	}

	return connection
}

func Listenerfunc(l *Listener) {
	fmt.Println("Listener created!")
	fmt.Println("Listener channels", l.write, l.read)
	time.Sleep(1000 * time.Millisecond)
	fmt.Println("Listener created!")
	fmt.Println("Listener channels", l.write, l.read)

}

func Manager(lManager *ListenerManager) {
	//manager := createManager()

	for {
		select {
		case <-lManager.read:
			//
			fmt.Println("Manager: Read from session")
			tempShit := Data{"port", lManager.currentPort}
			go createListener(lManager, lManager.currentPort)
			lManager.write <- tempShit
			lManager.currentPort++
		}

	}

}
