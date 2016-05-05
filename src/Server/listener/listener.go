package listener

import (
    "net"
    "strconv"
    "encoding/json"
    "fmt" //TEMP
)

//Player is used to represent the players in the game world
type Player struct {
    XCord   int
    YCord   int
    Lives   int
}
//Asteroid is used to represent the asteroids in the game world
type Asteroid struct {
    XCord   int
    YCord   int
    Stage   int
}

//World is used to represent the entire gameworld to send to clients
type World struct {
    Players     []*Player
    Asteroids   []*Asteroid
}

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
    socket  net.Listener
    ID      string
    Port    int
    player  Player
    //Connection
}

func createSocket(port int) (net.Listener, error) {
    connection, err := net.Listen("tcp", ":" + strconv.Itoa(port))
    if err != nil {
        return nil, err
    }
    
    return connection, nil
}

//NewListener creates a new socket then runs this socket as a go routine
func NewListener(port int/*, conn *Connection*/) *Listener {
    listener := new(Listener)
    
    var err error
    listener.socket, err = createSocket(port)
    if err != nil {
        panic(err)
    }
    
    
    listener.Port = port
    
    //listener.write = conn.read //Fan in to manager
    //listener.read = conn.write //Fan out from manager
    
    return listener
}

func (listen *Listener) Write(/*world *World*/) {
    //TEMPCODE
    currentWorld := new(World)
    currentWorld.Players = make([]*Player, 1)
    currentWorld.Asteroids = make([]*Asteroid, 1)
    
    player := new(Player)
    asteroid := new(Asteroid)
    player.Lives = 1
    player.XCord = 1
    player.YCord = 1
    
    asteroid.Stage = 2
    asteroid.XCord = 2
    asteroid.YCord = 2
    
    currentWorld.Players[0] = player
    currentWorld.Asteroids[0] = asteroid
    
    jsonWorld, err := json.Marshal(&currentWorld)
    if err != nil {
        panic(err)
    }
    
    fmt.Println(string(jsonWorld))
}