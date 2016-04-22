package main

import (
	"fmt"
	"net"
	"os"
)

type Listener struct {
	id   string
	port int
	Connection
}

/*
func ListenerFunc(listener *Listener) {

	for {
		message := Data{action: "listenToMe", result: 123}
		listener.write <- message
	}

}
*/

type ListenerManager struct {
	currentPort    int
	listenerAmount []Listener
	connection     Connection
}

func createManager() *ListenerManager {
	manager := new(ListenerManager)
	manager.currentPort = 9000

	return manager
}

func createListener(port string) net.Listener {

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Listener created!")
	}

	return connection
}

func Manager(cSession *Connection) {
	manager := createManager()
	//go createListener(strconv.Itoa(manager.currentPort))

	for {
		select {
		case <-cSession.read:

			fmt.Println("Manager: Read from session")
			tempShit := Data{"port", manager.currentPort}
			cSession.write <- tempShit
			manager.currentPort++
		}

	}
	//New connection
	//NewCOnnection read = old.write
	//Newcon write = old.read

}
