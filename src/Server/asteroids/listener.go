package asteroids

import (
	//TEMP
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"strconv"
	"time" //TEMP
)

//Player is used to represent the players in the game world
type Player struct {
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

	listener.idleListener()
}

func (listener *Listener) idleListener() {
	for {
		if listener.writeBuffer[0] != nil {
			listener.conn.Write(listener.writeBuffer[0])
			listener.writeBuffer = listener.writeBuffer[1:]
		} else {
			time.Sleep(time.Second)
		}
	}
}

func (listener *Listener) Write(world World) {
	jsonWorld, err := json.Marshal(world)
	if err != nil {
		panic(err)
	}

	listener.writeBuffer = append(listener.writeBuffer, jsonWorld)
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
