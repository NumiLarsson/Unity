package main

import (
	"fmt"
	"time"
)

// Data struct to be sent in channels
// TODO: change to data
type Data struct {
	action string
	result int
}

// Connection struct, containing one write and one read channel
// TODO: change to connection
type Connection struct {
	write chan Data
	read  chan Data
}

// Local copy of a session, to be used for handling where to connect new players
type gameSession struct {
	id int
	*Connection
}

// TODO: type struct?
// Used to store information about sessions, players etc.
type server struct {
	totalPlayers int
	nextPort     int
	maxPlayers   int
	nextSession  int
	sessions     []*gameSession
}

// Creates two mirrored connections
// TODO: Currently using Data channels, implement generic channels?
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

	var server = createServer()
	listen(server, createFakeUser(server))

}

// Only used to get some kind of input from a "user"
func createFakeUser(server *server) chan Data {

	fakeUser := make(chan Data)

	go func() {

		time.Sleep(250 * time.Millisecond)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

		time.Sleep(500 * time.Millisecond)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

	}()

	return fakeUser
}

// Server listen for new user that want to connect
// Sends correct port to use in return
func listen(server *server, external chan Data) {

	// TEMPORARY
	// ===
	// Kill the server after 5 seconds of inactivity
	timeout := time.After(5 * time.Second)

	for {
		select {
		// TODO change external to correct input channel/port used by external comm.
		case message := <-external:
			fmt.Println("Server: New user wants to connect \n", message.action)
			// TODO: Possibly in a go-routine based on performance
			var port = addPlayer(server)

			// Port to use should be sent to the user
			fmt.Println("Server: Port set up for new user", port)

		case <-timeout:
			return
		}
	}

}

// Add player to first available session that has capacity for a new player
// If no session has capacity or available, creates a new session and player
func addPlayer(server *server) int {

	for _, s := range server.sessions {

		// Ask a session whether there is enough room for a new player
		s.write <- Data{"Connect", 1}
		port := <-s.read

		if port.result > -1 {
			return createPlayer(server, s)
		}
	}

	return createSession(server)
}

// Creates a new server struct
func createServer() *server {

	s := new(server)
	s.totalPlayers = 0
	s.nextPort = 9000
	s.maxPlayers = 8

	return s
}

// Used to create a local copy of the session in the server
func createSession(server *server) int {

	cInternal, cExternal := makeConnection()
	session := new(gameSession)

	session.Connection = cInternal
	session.id = server.nextSession
	server.nextSession++
	server.sessions = append(server.sessions, session)

	// Start a session and wait for it to send confirmation
	go Session(cExternal, nextPort(server), server.maxPlayers)
	<-cInternal.read

	fmt.Println("Session created")

	return createPlayer(server, session)
}

// Send request to session to create a new player and wait for confirmation
func createPlayer(server *server, gs *gameSession) int {

	gs.write <- Data{"Connect new player", 1}
	data := <-gs.read
	fmt.Println("Player connected")
	server.totalPlayers++
	return data.result
}

// Calculate the next start port to be used by a session
func nextPort(server *server) int {
	var port = server.nextPort
	server.nextPort += server.maxPlayers
	return port
}
