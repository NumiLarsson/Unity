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

	fakeUser := make(chan Data)

	//Only used to get som kind of input from a "user"
	go func() {

		time.Sleep(2 * time.Second)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

		time.Sleep(1 * time.Second)
		//fakeUser <- 1

	}()

	listen(fakeUser)

}

func listen(external chan Data) {

	// Create channels to session
	cInternal, cExternal := makeConnection()

	createSession(cExternal)
	//fmt.Printf("server: cInternal write %d\n", cInternal.write)
	fmt.Printf("server: cInternal read  %d\n", cInternal.read)
	fmt.Printf("server: cExternal write %d\n", cExternal.write)
	//fmt.Printf("server: cExternal read %d\n", cExternal.read)
	for {

		select {
		case message := <-external:
			//TODO Should store new user to be able to send back Port once
			//the listener is created for specific user
			fmt.Printf("server: message from new user \n", message.action)
			// Probably a go-routine in the future to prevent blocking
			// cInternal should be changed to
			port := connectToSession(cInternal, server.nextPort)
			fmt.Printf("Port %d\n", port)

			server.totalPlayers++
			server.nextPort++

		// Receive confirmation that listener is created
		// TO FIX: Currently only working when sending through both channels?
		case message := <-cInternal.read:
			fmt.Printf("Server: message from session %d\n", message.action)

		case message := <-cInternal.write:
			fmt.Printf("Server: message from session write!! %d\n", message.action)

			//TODO: Send response back to user
			/*default:
			time.Sleep(500 * time.Millisecond)
			fmt.Println("Server: ")
			*/
		}
	}

	//Next step to implement creating more sessions
	//Be able to store channels to different sessions
	//Read input from new user.
	//Create new channels to new session
	//store sessioninfo + channels
	//...

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
