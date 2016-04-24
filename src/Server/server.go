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

	//Only used to get som kind of input from a "user"
	go func() {

		time.Sleep(2 * time.Second)
		fakeUser <- 1

		time.Sleep(1 * time.Second)
		//fakeUser <- 1

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
			//TODO Should store new user to be able to send back Port once
			//the listener is created for specific user

			// Probably a go-routine in the future to prevent blocking
			//cInternal should be changed to
			port := connectToSession(cInternal, server.nextPort)
			fmt.Printf("Port %d\n", port)

			server.totalPlayers++
			server.nextPort++

			//TODO: Send response back to user

		}

		//Next step to implement creating more sessions
		//Be able to store channels to different sessions
		//Read input from new user.
		//Create new channels to new session
		//store sessioninfo + channels
		//...

	}

}

func createSession(conn *Connection) {

	// Swap read and write
	// connSwap := new(Connection)
	//connSwap.write = conn.read
	// connSwap.read = conn.write

	go Session(conn)

}

//Used to create a new listener for a new user
func connectToSession(conn *Connection, port int) int {

	//Should send request to a session below used only for testing
	fmt.Printf("connect to session %d\n", port)
	conn.write <- Data{"connect", port}
	response := <-conn.read

	//If go-routine should it loop and wait for response that listener created
	//and then send response/port back to user

	return response.result

}
