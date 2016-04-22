package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

type Data struct {
	action string
	result int
}

type Connection struct {
	write chan Data
	read  chan Data
}

type Listener struct {
	id   string
	port int
	Connection
}

func main() {

	conn := new(Connection)
	conn.write = make(chan Data)
	conn.read = make(chan Data)

	createSession(conn)
	port := connectToSession(conn)

	//Send Port to client?

	fmt.Printf("Port %d\n", port)

	listen()

}

func listen() {

}

func createSession(conn *Connection) {

	// Swap read and write
	connSwap := new(Connection)
	connSwap.write = conn.read
	connSwap.read = conn.write

	go Session(connSwap)

}

func connectToSession(conn *Connection) int {

	conn.write <- Data{"new", 0}

	response := <-conn.read
	return response.result

}

func Session(conn *Connection) {

	//fromListener := make(chan int)
	//toListener := make(chan int
	connList := new(Connection)
	connList.write = make(chan Data)
	connList.read = make(chan Data)

	i := 0
	for i < 100 {
		time.Sleep(time.Second)
		select {
		case <-conn.read:
			go manager(connList)

		case userdata := <-connList.write:
			fmt.Printf("New data from user %d\n", userdata.result)
			//toListener <- 1
			conn.write <- userdata

		default:
			fmt.Println("Nothing to do")
		}

		i++

	}

}

func CreateListener(conn *Connection) *Listener {

	listener := new(Listener)
	listener.write = conn.read
	listener.read = conn.write

	return listener

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

func manager(conn *Connection) {
	manager := createManager()
	go createListener(strconv.Itoa(manager.currentPort))

	//New connection
	//NewCOnnection read = old.write
	//Newcon write = old.read

	tempShit := Data{"port", manager.currentPort}

	conn.write <- tempShit

	manager.currentPort++
}
