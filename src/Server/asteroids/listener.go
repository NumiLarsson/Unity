package asteroids

import (
	"fmt" //TEMP
	"net"
	"strconv"
	//"time" //TEMP
	"encoding/json"
)

//Player is used to represent the players in the game world
type Player struct {
	XCord int
	YCord int
	Lives int
}

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
	socket      net.Listener
	ID          string
	Port        int
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
func NewListener(port int /*, conn *Connection*/) (*Listener) {

	listener := new(Listener)

	var err error
	listener.socket, err = CreateSocket(port)
	if err != nil {
		panic(err)
	}

	listener.Port = port
	listener.player = new(Player)
	listener.player.XCord = 0
	listener.player.YCord = 0
	listener.player.Lives = 3
	
	listener.writeBuffer = make(chan []byte, 60)

	//listener.write = conn.read //Fan in to manager
	//listener.read = conn.write //Fan out from manager

	go listener.startUpListener()

	return listener //Listener has player in it!
}

func (listen *Listener) startUpListener() {
	var err error
	listen.conn, err = listen.socket.Accept()
	if err != nil {
		panic(err)
	}
	listen.ID = "Hello World"

	listen.idleListener()
}

func (listen *Listener) idleListener() {
	for {
		select{
			case jsonWorld := <- listen.writeBuffer:
				listen.conn.Write(jsonWorld)
				fmt.Println("Sent info to client", string(jsonWorld))
			default:
		}		 
	}
}

func (listen *Listener) Write(world /**World*/*WorldClient) {
	jsonWorld, err := json.Marshal(world)
	if err != nil {
		panic(err)
	}
	listen.writeBuffer <- jsonWorld
}

func (listen *Listener) getPlayer() *Player {
	return listen.player
}
