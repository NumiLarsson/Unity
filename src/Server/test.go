package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"strconv"
	"bufio"
)

type Data struct {
	action string
	result string //Temporary rewrite to string because Data is what our channels use
}

type Connection struct {
	write chan Data
	read  chan Data
}

type Listener struct {
	id   	string
	port 	string /*Changed this to be a string instead of an int, 
	as it's easier to deal with native strings here and convert in C#*/
	socket 	net.Listener
	Connection
} 

func main() {

	conn := new(Connection)
	conn.write = make(chan Data)
	conn.read = make(chan Data)

	createSession(conn)

	//Create a channel to speak to listen()
	listenConn := new(Connection)
	listenConn.write = make(chan Data)
	listenConn.read = make(chan Data)
	
	//Run listen() concurrently
	go listen(listenConn)
	
	for {
		select {
			case portData := <- conn.read:
				//Session wrote to me
				listenConn.write <- portData
				//Write to listen() so that it can send it to the client
			case <- listenConn.write:
				//Listen() wrote to me
				conn.write <- Data{"NewUser", "NewUser"}
				//Write to session to tell it to spawn a new listener
			default:
		}	
	}	
}

//Specific listener function for the server, to allow new users to connect.
func listen(conn *Connection) {
	
	for {
		socket := createSocket("9000")
		fmt.Println("Waiting for new client to connect")
		connection, err := socket.Accept()
		if err != nil {	
			panic(err)
		}
		conn.write <- Data{"NewUser", ""} //Tell Session that there's a new client
		
		port := make([]byte, 1024) 
		portData := <- conn.write
		fmt.Println("Sending port:", portData.result)
		port = []byte(portData.result) //connection.Write has to be in bytes, as it's pure network.
		connection.Write(port)
		fmt.Println("Port has been sent to client")
		
		connection.Close()
		socket.Close()
	}
}


func createSession(conn *Connection) {

	// Swap read and write
	connSwap := new(Connection)
	connSwap.write = conn.read
	connSwap.read = conn.write

	go Session(connSwap)

}

func Session(conn *Connection) {

	connList := new(Connection)
	connList.write = make(chan Data)
	connList.read = make(chan Data)

	//Create a new manager with startingPort set to 9001 (as server uses 9000 to talk to clients)
	startingPort := "9001"
	listenerManager := createManager(connList, startingPort)

	for {
		select {
		case managerData := <- connList.write:
			fmt.Println("New data from manager:", managerData.result)
		case <- conn.read:
			//Create new listener, which returns the port as a string
			port := listenerManager.NewObject()
			//Tell the server about the port.
			conn.write <- Data{"port", port}
			
			//Temp code, NYI.
			for i, listener := range listenerManager.listenerList {
				if listener != nil {
					fmt.Println(i, listener.id)
				}
			} //Not yet fully implemented, should preferably be a linked list or something like that.
			
		}
	}

}

//Listener and managers
//manager is an interface to allow session to store managers as a list of managers for itteration later.
type manager interface {
	NewObject ()
}

//ListenerManager is the specific manager for listeners
type ListenerManager struct {
	currentPort    string
	//Temp code, NYI.
	listenerList []*Listener //NYI.
	//Temp code, NYI.
	ListenerConnection Connection
	Connection
}

//createManager, does what it says, incrementing users from startingPort.
func createManager(conn *Connection, startingPort string) *ListenerManager {
	manager := new(ListenerManager)
	//Temp code, NYI.
	manager.listenerList = make([]*Listener, 500) //Magic number because this is NYI.
	//Temp code, NYI.
	manager.currentPort = startingPort
	manager.read = conn.write
	manager.write = conn.read

	return manager
}

//NewObject creates a new listener to the listenermanager, it returns the port for the new listener
//and creates a new goroutine for the listener, to allow the user to connect.
func (manager *ListenerManager) NewObject() string {
	
	listenerConn := new(Connection)
	listenerConn.read = make(chan Data)
	listenerConn.write = make(chan Data)
	
	fmt.Println("Current port: ", manager.currentPort)
	
	listener := new(Listener)
	listener.write = manager.ListenerConnection.read
	listener.read = manager.ListenerConnection.write
	listener.port = manager.currentPort
	
	//Temp code, NYI. Proof of concept
	for x, value := range manager.listenerList {
		if value == nil {
			manager.listenerList[x] = listener
			break
		}
	} //Temp code, NYI.
	
	go listener.StartUpListener(manager)
	
	return manager.currentPort
}

//StartUpListener is the function that actually creates the socket, it waits for an ID from the client, then enters echo mode.
func (listener *Listener) StartUpListener(manager *ListenerManager) {
	listener.socket = createSocket(manager.currentPort)
	manager.IncrementPort()
	connection, err := listener.socket.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Println("Client connected to listener: ", manager.currentPort)
	
	message, err := bufio.NewReader(connection).ReadString('\n')
	if err != nil {
		panic(err)
	}
	listener.id = string(message)
	
	go listener.IdleListener(connection)
}

//IdleListener is the "standard state" for listener, while it's not actively doing anything.
func (listener *Listener) IdleListener(connection net.Conn) {
	for {
		defer connection.Close()
		defer listener.socket.Close()
		message, err := bufio.NewReader(connection).ReadString('\n')
		if err != nil {
			panic(err)
		}

		fmt.Print("Received: ", string(message))

		bytes := make([]byte, 1024)
		bytes = []byte(message)
		connection.Write(bytes)
	}
}

//createSocket accepts a new tcp connection using the supplied port.
func createSocket(port string) net.Listener {

	connection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return connection
}

//IncrementPort is only used to update the port for the manager, it's a seperate function so that we can use go IncrementPort()
func (manager *ListenerManager) IncrementPort() {
	portInt, _ := strconv.Atoi(manager.currentPort)
	manager.currentPort = strconv.Itoa((portInt + 1))
	fmt.Println("Manager port is now:", manager.currentPort)
}	

//IdleManager is the idle state for the manager function, it's the resting state of the function.
func (manager ListenerManager) IdleManager (conn *Connection) {
	for {
		select {
			case readMessage := <- manager.ListenerConnection.read:
				conn.write <- readMessage
			default: 
				time.Sleep(time.Second)
		}
	}
}
