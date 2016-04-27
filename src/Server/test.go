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

func listen(conn *Connection) {
	socket := createSocket("9000")
	
	fmt.Println("Waiting for client to connect")
	connection, err := socket.Accept()
	defer connection.Close()
	if err != nil {	
		panic(err)
	}
	conn.write <- Data{"NewUser", ""} //Tell Session that there's a new client
	
	port := make([]byte, 1024) 
	portData := <- conn.write
	
	port = []byte(portData.result) //connection.Write has to be in bytes, as it's pure network.
	connection.Write(port)
	fmt.Println("Port has been sent to client")
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

	listenerManager := newManager(connList)

	for {
		select {
		case managerData := <- connList.write:
			fmt.Println("New data from manager:", managerData.result)
		case <- conn.read:
			fmt.Println("Creating new listener")
			port := listenerManager.NewObject()
			fmt.Println("Listener created with port:", port)
			conn.write <- Data{"port", port}
			fmt.Println("Listener port sent to server")
			
		}
	}

}

//Listener and managers
//manager is an interface to allow session to store managers as a list of managers for itteration later.
type manager interface {
	NewObject ()
}

//Listenermanager is the specific manager for listeners
type ListenerManager struct {
	currentPort    string
	listenerList []Listener
	ListenerConnection Connection
	Connection
}

//Create a new manager, doesn't do much more than that.
func createManager() *ListenerManager {
	manager := new(ListenerManager)
	manager.currentPort = "9001"

	return manager
}

//NewObject creates a new listener to the listenermanager, it returns the port for the new listener
//and creates a new goroutine for the listener, to allow the user to connect.
func (manager ListenerManager) NewObject() string {
	
	listenerConn := new(Connection)
	listenerConn.read = make(chan Data)
	listenerConn.write = make(chan Data)
	
	fmt.Println("Current port: ", manager.currentPort)
	
	listener := new(Listener)
	listener.write = manager.ListenerConnection.read
	listener.read = manager.ListenerConnection.write
	listener.port = manager.currentPort
	
	defer manager.IncrementPort()
	
	go listener.StartUpListener(manager)
	
	return manager.currentPort
}

//StartUpListener is the function that actually creates the socket, it waits for an ID from the client, then enters echo mode.
func (listener *Listener) StartUpListener(manager ListenerManager) {
	listener.socket = createSocket(manager.currentPort)
	
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
	
	fmt.Println(listener.id, "has joined and its listener is now in echo mode")
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

	fmt.Println("Creating listener: ", port)
	connection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Printf("Listener %s created!\n", port)
	}

	return connection
}

//IncrementPort is only used to update the port for the manager, it's a seperate function so that we can use go IncrementPort()
func (manager ListenerManager) IncrementPort() {
	portInt, _ := strconv.Atoi(manager.currentPort)
	manager.currentPort = strconv.Itoa((portInt + 1))
	fmt.Println(manager.currentPort)
}

//newManager creates a new listener manager
func newManager(conn *Connection) ListenerManager {
	manager := createManager()
	
	return *manager
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
