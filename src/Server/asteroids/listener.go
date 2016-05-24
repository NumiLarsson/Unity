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
	worldX	int
	worldY	int
	Name  string
	ID    int
	X     int
	Y     int
	Lives int
	Alive bool
	step  int
}

type playerMessage struct {
	Action string
	Value  string
}

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
	worldX 		int
	worldY		int
	socket      net.Listener
	ID          string
	port        int
	player      *Player
	conn        net.Conn
	writeBuffer chan []byte
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

// init initiates the listeners values
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
	//1 second worth of writes
}

//newPlayer returns a new player
func newPlayer() *Player {
	return new(Player)
}

//init initiates a new players values
func (player *Player) init(id int, xMax int, yMax int) {
	player.ID = id
	seed := time.Now().UnixNano()
	rand.Seed(seed)

	fmt.Println(seed)
	player.Name = strconv.Itoa(id);
	player.worldX = xMax
	player.worldY = yMax
	player.step = 1;
	player.randomSpawn(player.worldX, player.worldY)
	player.Lives = 3 // updated
	player.Alive = true
}

// startUpListener
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
	clientChan := make(chan *playerMessage)
	go listener.readFromClient(clientChan)
	
	for {
		select {
		case jsonWorld := <-listener.writeBuffer:
			listener.conn.Write(jsonWorld)
			//fmt.Println(string(jsonWorld))
		case message := <- clientChan :
			if ( !listener.player.newInput(message) ) {
				//fmt.Println("Input from player was invalid")
			} else {
				//fmt.Println("Look at me:", listener.ID, listener.player.X, listener.player.Y);
			}
		}
	}
}

func (listener *Listener) readFromClient(clientChan chan *playerMessage) {
	defer listener.panicCatcher(clientChan)
	for {
		bytes := make([]byte, 1024)
		bytesRead, err := listener.conn.Read(bytes)
		if err != nil {
			panic(err)
		}
		//fmt.Println("CLIENT SENT A MESSAGE!", string(bytes[:bytesRead]))
		message := new(playerMessage)
		err = json.Unmarshal(bytes[:bytesRead], &message)
		if err != nil {
			panic(err)
		}
		clientChan <- message
	}
}

func (listener *Listener) panicCatcher(clientChan chan *playerMessage) {
	fmt.Println(recover())
	err := listener.conn.Close()
	if (err != nil) {
		panic(err)
	}
	err = listener.socket.Close()
	if (err != nil) {
		panic(err)
	}
	listener.socket, err = CreateSocket(listener.port)
	if err != nil {
		panic(err)
	}
	listener.conn, err = listener.socket.Accept()
	if err != nil {
		panic(err)
	}
	listener.readFromClient(clientChan)
}

//write writes game world to clients
func (listener *Listener) Write(world *World) {
	jsonWorld, err := json.Marshal(world)
	if err != nil {
		panic(err)
	}

	listener.writeBuffer <- jsonWorld
	//listener.writeBuffer = append(listener.writeBuffer, jsonWorld)
}

// getPlayer returns a listeners player
func (listener *Listener) getPlayer() *Player {
	return listener.player
}

// randomSpawn spawn a player on a random location
func (player *Player) randomSpawn(xMax int, yMax int) {

	player.X = rand.Intn(xMax)
	player.Y = rand.Intn(yMax)

	fmt.Println(player.ID, player.X, player.Y)
}

//isAlive return if the player is alive or not
func (player *Player) isAlive() bool {
	return player.Alive
}

//getLives returns the amount of lives the player has left
func (player *Player) getLives() int {
	return player.Lives
}

//setAlive sets the Alive state to true
func (player *Player) setAlive() {
	player.Alive = true
}

func (player *Player) tryMove(value string) bool {
	switch (value) {
	case "North": //North
		if (player.Y + 1 > player.worldY) {
			return false
		} 
		//Else
		player.Y += player.step
		return true
		
	case "East": //East
		if (player.X + 1 > player.worldX) {
			return false
		} 
		//Else
		player.X += player.step
		return true
		
	case "South": //South
		if (player.Y < 0) {
			return false
		} 
		//Else
		player.Y -= player.step
		return true
		
	case "West": //West
		if (player.X - 1 < 0) {
			return false
		} 
		//Else
		player.X -= player.step
		return true
	}
	return false;
}

//newInput returns true if the input was valid.
func (player *Player) newInput(playMessage *playerMessage) bool {
	switch playMessage.Action {
	case "Move":
		return player.tryMove(playMessage.Value)
	case "Name":
		player.Name = playMessage.Value
		return true;
	}
	return false;
}