package main

import (
	"fmt"
	"net"
	"os"
	"time"
	"encoding/binary"
	"strconv"
	"bufio"
)

type Data struct {
	action string
	result int
}

type Connection struct {
	write chan Data
	read  chan Data
}

type Listener struct {
	id   string
	port int
	Connection
}

/*
func sendUInt16(intString int) {
	portUInt := *(*uint16)(unsafe.Pointer(&intString))
	fmt.Println("portUInt is of type %T\n", portUInt)
}
*/

func main() {

	conn := new(Connection)
	conn.write = make(chan Data)
	conn.read = make(chan Data)

	createSession(conn)
	port := connectToSession(conn)

	//Send Port to client?

	fmt.Printf("Port %d\n", port)

	listenConn := new(Connection)
	listenConn.write = make(chan Data)
	listenConn.read = make(chan Data)
	
	go listen(listenConn)
	
	portData := <- conn.read
	listenConn.write <- portData
	
}

func listen(conn *Connection) {
	socket := createListener("9000")
	fmt.Println("Waiting for client to connect")
	connection, err := socket.Accept()
	if err != nil {
		panic(err)
	}
	//portData := <- conn.write
	//Rewrite this to use portData.result
	//Convert portData.result in to a uint16
	var portUInt uint16 = 9001
	
	//This is production code again
	port := make([]byte, 2)
	binary.LittleEndian.PutUint16(port, portUInt)
	connection.Write(port)
	connection.Close()
}


func createSession(conn *Connection) {

	// Swap read and write
	connSwap := new(Connection)
	connSwap.write = conn.read
	connSwap.read = conn.write

	go Session(connSwap)

}

func connectToSession(conn *Connection) int {

	conn.write <- Data{"new", 0}

	response := <-conn.read
	return response.result

}

func Session(conn *Connection) {

	//fromListener := make(chan int)
	//toListener := make(chan int
	connList := new(Connection)
	connList.write = make(chan Data)
	connList.read = make(chan Data)

	i := 0
	for i < 100 {
		time.Sleep(time.Second)
		select {
		case <-conn.read:
			go manager(connList)

		case userdata := <-connList.write:
			fmt.Printf("New data from user %d\n", userdata.result)
			//toListener <- 1
			conn.write <- userdata

		default:
		}

		i++

	}

}

/*
func ListenerFunc(listener *Listener) {

	for {
		message := Data{action: "listenToMe", result: 123}
		listener.write <- message
	}

}
*/

type ListenerManager struct {
	currentPort    int
	listenerList []Listener
	Connection
}

func createManager() *ListenerManager {
	manager := new(ListenerManager)
	manager.currentPort = 9001

	return manager
}

func createListener(port string) net.Listener {

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else {
		fmt.Println("Listener created!")
	}

	return connection
}

func listener(conn *Connection, port string) {
	listenerConn := new(Listener)
	listenerConn.write = conn.read
	listenerConn.read = conn.write
	
	listener := createListener(port)
	fmt.Println("New Listener created", port)
	if listener == nil {
		panic("Listener creation failed")
	}
	portInt, _ := strconv.Atoi(port)
	portData := Data{"port", portInt}
	conn.write <- portData
	
	connection, err := listener.Accept()
	if err != nil {
		panic(err)
	}
	fmt.Println("Client connected to the new listener!")
	
	for {
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

func manager(conn *Connection) {
	manager := createManager()
	
	listenerConn := new(Connection)
	listenerConn.write = make(chan Data)
	listenerConn.read = make(chan Data)
	go listener(listenerConn, strconv.Itoa(manager.currentPort))

	//New connection
	//NewCOnnection read = old.write
	//Newcon write = old.read
	readMessage := <- listenerConn.write

	conn.write <- readMessage

	manager.currentPort++
}
