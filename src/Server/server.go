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

type gameSession struct {
	id      int
	players int
	*Connection
}

var server struct {
	totalPlayers int
	nextPort     int
	maxPlayers   int
	nextSession  int
	sessions     []*gameSession
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
	server.nextPort = 9000
	server.maxPlayers = 8

	fakeUser := make(chan Data)

	//Only used to get som kind of input from a "user"
	go func() {

		time.Sleep(2 * time.Second)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

		time.Sleep(1 * time.Second)
		fakeUser <- Data{"New user wants to connect", server.nextPort}

	}()

	listen(fakeUser)

}

func listen(external chan Data) {

	for {
		select {
		// TODO change external to correct input channel/port used by external comm.
		case message := <-external:
			fmt.Println("server: message from new user \n", message.action)
			var port = addPlayer()
			// Port to use should be sent to the user
			fmt.Println("Set up new port", port)
		}
	}

}

func addPlayer() int {

	if server.sessions == nil {
		return createSession()
	}

	for _, s := range server.sessions {

		if s.players <= server.maxPlayers {
			return createPlayer(s)
		} else {
			return createSession()
		}
	}
	return -1
}

// Used to create a local copy of the session in the server
func createSession() int {

	cInternal, cExternal := makeConnection()
	session := new(gameSession)

	session.Connection = cInternal
	session.id = server.nextSession
	server.nextSession++
	server.sessions = append(server.sessions, session)

	go Session(cExternal, nextPort(), server.maxPlayers)
	<-cInternal.read

	return createPlayer(session)
}

func createPlayer(session *gameSession) int {

	session.write <- Data{"connect", 0}
	data := <-session.read
	fmt.Println("Player connected")
	server.totalPlayers++
	return data.result
}

func nextPort() int {
	var port = server.nextPort
	server.nextPort += server.maxPlayers
	return port
}
