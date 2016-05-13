package asteroids

import (
	"encoding/json"
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

// gameSession is a local copy of a session, to be used for handling where to connect new players
type gameSession struct {
	id int
	*Connection
}

// server struct used to store information about sessions, players etc.
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

// Only used to get some kind of input from a "user"
func (server *server) CreateFakeUser() chan Data {

	fakeUser := make(chan Data)
	nextPort := server.nextPort // Prevents data race

	go func() {

		fakeUser <- Data{"server.new_user", nextPort}
		<-fakeUser

		fakeUser <- Data{"server.new_user", nextPort}
		<-fakeUser

	}()

	return fakeUser
}

// Listen is a loop that server uses to listen for new user that want to connect
// Sends correct port to use in return
func (server *server) Listen(external chan Data) {

	// TEMPORARY
	// ===
	// Kill the server after 5 seconds of inactivity
	timeout := time.After(5 * time.Second)

	go acceptNewPlayers(external)

	for {
		select {
		// TODO change external to correct input channel/port used by external comm.
		case <-external:
			fmt.Println("[SERVER] New user wants to connect")
			// TODO: Possibly in a go-routine based on performance
			port := server.addPlayer()
			external <- Data{"port", port}
			// Port to use should be sent to the user
			fmt.Println("[SERVER] Port set up for new user", port)

		case <-timeout:
			fmt.Println("\n========\n[SERVER] Terminated due to 5 seconds of inactivity\n========")
			return
		}
	}

}

func acceptNewPlayers(conn chan Data) {
	socket, err := CreateSocket(9000)
	if err != nil {
		panic(err)
	}

	for {
		tcpConn, err := socket.Accept()
		if err != nil {
			panic(err)
		}

		conn <- Data{"NewUser", 0}
		portData := <-conn

		jsonPort, err := json.Marshal(&portData.result)
		if err != nil {
			panic(err)
		}
		tcpConn.Write(jsonPort)

		tcpConn.Close()
	}
}

// addPlayer adds a player to first available session that has capacity for a new player
// If no session has capacity or available, creates a new session and player
func (server *server) addPlayer() int {

	for _, session := range server.sessions {

		// Ask a session whether there is enough room for a new player
		session.write <- Data{"server.poke", 1}
		port := <-session.read

		if port.result > -1 {
			return server.createPlayer(session)
		}
	}

	return server.createSession()
}

// CreateServer creates a new server struct
func CreateServer() *server {

	server := new(server)
	server.totalPlayers = 0
	server.nextPort = 9001 //9001 becuase it allows us to count acceptNewPlayers
	//based on the ports, and also gives us 9000 as static server port
	server.maxPlayers = 8

	return server
}

// createSession creates a local copy of the session in the server
func (server *server) createSession() int {

	serverSide, sessionSide := makeConnection()

	// Start a session and wait for it to send confirmation

	nextPort := server.getNextPort() // Prevents data race

	//Hardcoded size of the world for now
	go Session(sessionSide, nextPort, server.maxPlayers, 400)
	<-serverSide.read

	fmt.Println("[SERVER] Session created")

	// Create a local copy
	session := new(gameSession)

	session.Connection = serverSide
	session.id = server.nextSession
	server.nextSession++
	server.sessions = append(server.sessions, session)

	return server.createPlayer(session)
}

// createPlayer sends request to session to create a new player and wait for confirmation
func (server *server) createPlayer(gs *gameSession) int {

	gs.write <- Data{"server.connect_player", 1}
	data := <-gs.read
	fmt.Println("[SERVER] Player connected")
	server.totalPlayers++
	return data.result
}

// getNextPort calculates the next start port to be used by a session
func (server *server) getNextPort() int {
	var port = server.nextPort
	server.nextPort += server.maxPlayers
	return port
}
