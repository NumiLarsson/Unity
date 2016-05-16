package asteroids

import (
	"encoding/json"
	"fmt"
	"time"
)

// Data is a semi-generic struct to be sent in channels
type Data struct {
	action string
	result int
}

// Connection struct, containing one write and one read channel
type Connection struct {
	write chan Data
	read  chan Data
}

// gameSession is a local copy of a session, to be used for handling where to connect new players
type gameSession struct {
	id int
	*Connection
}

// Server struct used to store information about sessions, players etc.
type Server struct {
	totalPlayers int
	nextPort     int
	maxPlayers   int
	nextSession  int
	sessions     []*gameSession
}

var inDebugMode = true

func debugPrint(str string) {

	if inDebugMode {
		fmt.Print("DEBUG: ", str)
	}

}

// makeTwoWayConnection creates two connections where the outlets of the second is
// the reverse of the first
func makeTwoWayConnection() (c1, c2 *Connection) {

	c1 = new(Connection)
	c2 = new(Connection)

	c1.read = make(chan Data)
	c1.write = make(chan Data)

	c2.read = c1.write
	c2.write = c1.read

	return
}

// CreateFakeUser is a temporary/debug function to create a mock user
func (server *Server) CreateFakeUser() chan Data {

	fakeUser := make(chan Data)
	nextPort := server.nextPort

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
func (server *Server) Listen(external chan Data) {

	// TEMPORARY
	// ===
	// Kill the server after 5 seconds of inactivity
	timeout := time.After(60 * time.Second)

	go acceptNewPlayers(external)

	for {
		select {
		// TODO change external to correct input channel/port used by external comm.
		case <-external:
			// TODO: Possibly in a go-routine based on performance
			port := server.addPlayer()
			external <- Data{"port", port}
			// Port to use should be sent to the user
			//fmt.Println("[SERVER] Port set up for new user", port)
			debugPrint(fmt.Sprintln("[SERVER] Port set up for new user", port))

		case <-timeout:
			fmt.Println("\n========\n[SERVER] Terminated due to 60 seconds of inactivity\n========")
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
func (server *Server) addPlayer() int {

	for _, session := range server.sessions {

		// Ask a session whether there is enough room for a new player
		session.write <- Data{"server.poke", 1}
		port := <-session.read

		if port.result > -1 {
			return server.createPlayer(session)
		}
	}

	// If no session has capacity or available, creates a new session and player
	return server.createSession()
}

// CreateServer creates a new server struct
func CreateServer(debugMode bool) *Server {

	server := new(Server)
	server.totalPlayers = 0
	server.nextPort = 9001 //9001 becuase it allows us to count acceptNewPlayers
	//based on the ports, and also gives us 9000 as static server port
	server.maxPlayers = 8

	inDebugMode = debugMode

	return server
}

// createSession creates a local copy of the session in the server
func (server *Server) createSession() int {

	serverSide, sessionSide := makeTwoWayConnection()

	// Start a session and wait for it to send confirmation

	nextPort := server.getNextPort() // Prevents data race

	//Hardcoded size of the world for now
	go Session(sessionSide, nextPort, server.maxPlayers, 400)
	<-serverSide.read

	debugPrint(fmt.Sprintln("[SERVER] Session created"))

	// Create a local copy
	session := new(gameSession)

	session.Connection = serverSide
	session.id = server.nextSession
	server.nextSession++
	server.sessions = append(server.sessions, session)

	return server.createPlayer(session)
}

// createPlayer sends request to session to create a new player and wait for confirmation
func (server *Server) createPlayer(gs *gameSession) int {

	gs.write <- Data{"server.connect_player", 1}
	data := <-gs.read
	debugPrint(fmt.Sprintln("[SERVER] Player connected"))
	server.totalPlayers++
	return data.result
}

// getNextPort calculates the next start port to be used by a session
func (server *Server) getNextPort() int {
	var port = server.nextPort
	server.nextPort += server.maxPlayers
	return port
}
