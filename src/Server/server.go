package main

import (
	"fmt"
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

var server struct {
	totalPlayers int
	nextPort     int
}

func makeConnection() (c1, c2 *Connection) {

	c1 = new(Connection)
	c2 = new(Connection)

	c1.read = make(chan Data)
	c1.write = make(chan Data)

	c2.read = c1.write
	c2.write = c1.read

	return
}

func main() {

	server.totalPlayers = 0
	server.nextPort = 0

	fakeUser := make(chan int)

	//Send Port to client?
	//fmt.Printf("Port %d\n", port)

	go func() {

		time.Sleep(2 * time.Second)
		fakeUser <- 1

		time.Sleep(1 * time.Second)
		fakeUser <- 1

	}()

	listen(fakeUser)

}

func listen(external chan int) {

	//conn := new(Connection)
	//conn.write = make(chan Data)
	//conn.read = make(chan Data)

	cInternal, cExternal := makeConnection()

	createSession(cExternal)

	for {
		select {
		case <-external:
			// Probably a go-routine in the future to prevent blocking
			port := connectToSession(cInternal, server.nextPort)
			fmt.Printf("Port %d\n", port)

			server.totalPlayers++
			server.nextPort++

		}
	}

}

func createSession(conn *Connection) {

	// Swap read and write
	// connSwap := new(Connection)
	//connSwap.write = conn.read
	// connSwap.read = conn.write

	go Session(conn)

}

func connectToSession(conn *Connection, port int) int {

	conn.write <- Data{"connect", port}

	response := <-conn.read
	return response.result

}
