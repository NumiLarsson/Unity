package main


import (
	"net"
	"fmt"
	"os"
	"strconv"
)

type Data struct{
	actions string
	result int
}

type Connection struct{
	write chan Data
	read chan Data
}

type Listener struct{
	id string
	port int
	connection Connection
}

type ListenerManager struct{
	currentPort int
	listenerAmount []Listener
	connection Connection
}

func createManager() *ListenerManager{
	manager := new(ListenerManager)
	manager.currentPort = 8083

	return manager
}

func createListener(port string) net.Listener{

	fmt.Println("Creating listener...")
	connection, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	} else{
		fmt.Println("Listener created!")
	}

	return connection
}

func manager(conn *Connection){
	manager := createManager()
	createListener(strconv.Itoa(manager.currentPort))

	conn.write <- manager.currentPort

	manager.currentPort++
}

func main() {
	manager := createManager()
	createListener(strconv.Itoa(manager.currentPort))

}

//Skapa listener med random PORT
//Return port till session

//Struct för ny port, listenermanagare
//Skapa listener
//Returnera port

//Behöver channeö och returnerar int