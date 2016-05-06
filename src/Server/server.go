package main

import (
	"fmt"
	"time"
)

// Data struct to be sent in channels
type Data struct {
	action string
	result int
}

type portRequest int

// Connection struct, containing one write and one read channel
type Connection struct {
	write chan Data
	read  chan Data
}

// Local copy of a session, to be used for handling where to connect new players
type gameSession struct {
	id      int
	players int
	*Connection
}

// TODO: type struct?
// Used to store information about sessions, players etc.
var server struct {
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

	server.totalPlayers = 0
	server.nextPort = 9000
	server.maxPlayers = 8

	listen(createFakeUser())

}

// Only used to get som kind of input from a "user"
func createFakeUser() chan Data {

	fakeUser := make(chan Data)

	go func() {

		time.Sleep(2 * time.Second)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

		time.Sleep(1 * time.Second)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

	}()

	return fakeUser
}

// Server listen for new user that want to connect
// Sends correct port to use in return
func listen(external chan Data) {

	for {
		select {
		// TODO change external to correct input channel/port used by external comm.
		case message := <-external:
			fmt.Println("Server: New user wants to connect \n", message.action)
			var port = addPlayer()

			// Port to use should be sent to the user
			fmt.Println("Server: Port set up for new user", port)
		}
	}

}

// Add player to first available session that has capacity for a new player
// If no session has capacity or available, creates a new session and player
func addPlayer() int {

	for _, s := range server.sessions {
		s.write <- Data{"Connect", 1}
		port <- s.read
		if port.result > -1 {

		}
		if s.players < server.maxPlayers {
			return createPlayer(s)
		} else {
			return createSession()
		}
	}

	return createSession()
}

// Used to create a local copy of the session in the server
func createSession() int {

	cInternal, cExternal := makeConnection()
	session := new(gameSession)

	session.Connection = cInternal
	session.id = server.nextSession
	server.nextSession++
	server.sessions = append(server.sessions, session)

	// Start a session and wait for it to send confirmation
	go Session(cExternal, nextPort(), server.maxPlayers)
	<-cInternal.read

	return createPlayer(session)
}

// Send request to session to create a new player and wait for confirmation
func createPlayer(gs *gameSession) int {

	gs.write <- Data{"Connect new player", 1}
	data := <-gs.read
	fmt.Println("Player connected")
	server.totalPlayers++
	return data.result
}

// Calculate the next start port to be used by a session
func nextPort() int {
	var port = server.nextPort
	server.nextPort += server.maxPlayers
	return port
}
