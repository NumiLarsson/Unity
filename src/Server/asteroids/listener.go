package asteroids

import (
	//TEMP
	"encoding/json"
	"net"
	"strconv"
	"time" //TEMP
)

//Player is used to represent the players in the game world
type Player struct {
	Id    int
	X     int
	Y     int
	Lives int
	Alive bool
}

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
	socket      net.Listener
	ID          string
	port        int
	player      *Player
	conn        net.Conn
	writeBuffer [][]byte
	//Connection
}

//CreateSocket creates a tcp listener at the specified port
func CreateSocket(port int) (net.Listener, error) {
	connection, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	return connection, nil
}

//NewListener creates a new socket then runs this socket as a go routine
func newListener() *Listener {
	return new(Listener)
}

func (listener *Listener) init(port int) {

	var err error
	listener.socket, err = CreateSocket(port)
	if err != nil {
		panic(err)
	}

	listener.port = port

}

func newPlayer() *Player {
	return new(Player)
}

func (player *Player) init(id int) {
	player.Id = id
	player.X = 0
	player.Y = 0
	player.Lives = 3
	player.Alive = true
}

func (listen *Listener) startUpListener() {
	var err error
	listen.conn, err = listen.socket.Accept()
	if err != nil {
		panic(err)
	}

	listen.idleListener()
}

func (listen *Listener) idleListener() {
	for {
		if listen.writeBuffer[0] != nil {
			listen.conn.Write(listen.writeBuffer[0])
			listen.writeBuffer = listen.writeBuffer[1:]
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (listen *Listener) Write(world *World) {
	jsonWorld, err := json.Marshal(world)
	if err != nil {
		panic(err)
	}

	listen.writeBuffer = append(listen.writeBuffer, jsonWorld)
}

func (listen *Listener) getPlayer() *Player {
	return listen.player
}
