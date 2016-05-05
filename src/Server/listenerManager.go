package main

import (
    "fmt"
    "github.com/numilarsson/ospp-2016-group-08/src/Server/listener"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
    maxPlayers      int
    currentPlayers  int
    currentPort     int
    listeners       []listener.Listener
    World           [][]int
}

type World struct {
}

//NewListenerManager does exactly what it says, with a cap on maxPlayers 
//connected and maxPlayers numbers of ports in a row from firstPort
func NewListenerManager(maxPlayers int, firstPort int) *ListenerManager {
    lisManager := new(ListenerManager)
    lisManager.maxPlayers = maxPlayers
    lisManager.currentPlayers = 0
    lisManager.currentPort = firstPort
    lisManager.listeners = make([]listener.Listener, maxPlayers)
    
    return lisManager
}

func (lisManager *ListenerManager) IncrementPort() {
    lisManager.currentPort++
}

//NewObject creates a new listener for the listener manager, used to connect to a new player.
func (lisManager *ListenerManager) NewObject() int {
    defer lisManager.IncrementPort()
    
    tempListener := listener.NewListener(lisManager.currentPort)
    fmt.Println(tempListener.Port)
    return lisManager.currentPort
}

//Write is the function that fans out all the data to be sent to clients
//to the listeners that are responsible for doing so!
func (lisManager *ListenerManager) Write(world *World) {
    for key, value := range lisManager.listeners {
        if key != nil {
            go key.Write(world)
        }
    }
}

func main() {
    listener := listener.NewListener(9000)
    fmt.Println(listener.Port)
    fmt.Println("Listener created")
    
}
