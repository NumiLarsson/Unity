package main

import (
	"fmt"
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

	listener := CreateListener(connList)

	go ListenerFunc(listener)

	i := 0
	for i < 100 {

		select {
		case <-conn.read:
			fmt.Printf("Spawn user\n")

			message := Data{action: "port", result: 9000}
			conn.write <- message

		case userdata := <-connList.read:
			fmt.Printf("New data from user %d\n", userdata.result)
			//toListener <- 1

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

func ListenerFunc(listener *Listener) {

	for {
		message := Data{action: "listenToMe", result: 123}
		listener.write <- message
	}

}
