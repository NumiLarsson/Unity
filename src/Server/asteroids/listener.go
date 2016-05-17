package asteroids

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time"
)

//Player is used to represent the players in the game world
type Player struct {
	Name  string
	ID    int
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
	writeBuffer chan []byte
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

	/*
		This should be synchronized with listener?

		listener.player = new(Player)
		listener.player.Name = strconv.Itoa(port)
		listener.player.X = 0
		listener.player.Y = 0
		listener.player.Lives = 3
	*/

	listener.writeBuffer = make(chan []byte, 60)

}

func newPlayer() *Player {
	return new(Player)
}

func (player *Player) init(id int, xMax int, yMax int) {
	player.ID = id

	rand.Seed(time.Now().UnixNano())

	player.randomSpawn(xMax, yMax)
	player.Lives = 3
	player.Alive = true
}

func (listener *Listener) startUpListener() {
	var err error
	listener.conn, err = listener.socket.Accept()
	if err != nil {
		panic(err)
	}
	listener.ID = "Hello World"

	listener.idleListener()
}

func (listener *Listener) idleListener() {

	/*
		for {
			if listener.writeBuffer[0] != nil {
				listener.conn.Write(listener.writeBuffer[0])
				listener.writeBuffer = listener.writeBuffer[1:]
			} else {
				time.Sleep(time.Second)
			}
		}
	*/

	for {
		select {
		case jsonWorld := <-listener.writeBuffer:
			//sizeWorld := binary.Size(jsonWorld)
			//jsonSize, err := json.Marshal(sizeWorld)
			//if err != nil {
			//	panic(err)
			//}
			//listen.conn.Write(jsonSize)
			listener.conn.Write(jsonWorld)
		default:
		}
	}
}

//Write writes world to clients
func (listener *Listener) Write(world *World) {
	jsonWorld, err := json.Marshal(world)
	if err != nil {
		panic(err)
	}

	listener.writeBuffer <- jsonWorld
	//listener.writeBuffer = append(listener.writeBuffer, jsonWorld)
}

func (listener *Listener) getPlayer() *Player {
	return listener.player
}

func (player *Player) randomSpawn(xMax int, yMax int) {

	player.X = rand.Intn(xMax)
	player.Y = rand.Intn(yMax)

	fmt.Println(player.ID, player.X, player.Y)
}

func (player *Player) isAlive() bool {
	return player.Alive
}

func (player *Player) getLives() int {
	return player.Lives
}

func (player *Player) setAlive() {
	player.Alive = true
}
