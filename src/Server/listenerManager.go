package main

import (
    "fmt"
    "github.com/numilarsson/ospp-2016-group-08/src/server/listener"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
    maxPlayers      int
    currentPlayers  int
    currentPort     int
    listeners       []listener.Listener
    World           [][]int
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
    
    listener := new(listener.Listener)
            
    return lisManager.currentPort
}

func main() {
    listener := listener.NewListener("9000")
    fmt.Println(listener.Port)
    fmt.Println("Listener created")
    
}