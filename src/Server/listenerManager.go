package main

import (
    "github.com/numilarsson/ospp-2016-group-08/src/Server/listener"
)

//ListenerManager is used as a struct to basically emulate an object
type ListenerManager struct {
    MaxPlayers      int
    CurrentPlayers  int
    CurrentPort     int
    listenerList    []listener.Listener
    //WorldArray      World
}

//NewListenerManager does exactly what it says, with a cap on maxPlayers 
//connected and maxPlayers numbers of ports in a row from firstPort
func NewListenerManager(maxPlayers int, firstPort int) *ListenerManager {
    lisManager := new(ListenerManager)
    lisManager.MaxPlayers = maxPlayers
    lisManager.CurrentPlayers = 0
    lisManager.CurrentPort = firstPort
    lisManager.listenerList = make([]listener.Listener, maxPlayers)
    
    return lisManager
}
//IncrementPort increments the current port in the manager.
func (lisManager *ListenerManager) IncrementPort() {
    lisManager.CurrentPort++
}

//NewObject creates a new listener for the listener manager, used to connect to a new player.
func (lisManager *ListenerManager) NewObject() int {
    defer lisManager.IncrementPort()
    
    tempListener := listener.NewListener(lisManager.CurrentPort)
    //fmt.Println(tempListener.Port)
    //used only to be able to compile
    
    //Fix this
    lisManager.listenerList[0] = *tempListener 
    //Need to be proper later
    return lisManager.CurrentPort
}

//Write is the function that fans out all the data to be sent to clients
//to the listeners that are responsible for doing so!
func (lisManager *ListenerManager) Write(/*world *listener.World*/) {
    for _, value := range lisManager.listenerList {
        //Must check if value exists before we do this function call
        //but if value != nil doesn't work!
        go value.Write(/*world*/)
    }
}

func main() {
    ok := make(chan bool)
    lisManager := NewListenerManager(1, 9000)
    lisManager.NewObject()
    lisManager.Write()
    <- ok
}
