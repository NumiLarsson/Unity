package listener

import (
    "net"
)

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
    socket  net.Listener
    ID      string
    Port    string
    //Connection
}

func createSocket(port string) (net.Listener, error) {
    connection, err := net.Listen("tcp", ":" + port)
    if err != nil {
        return nil, err
    }
    
    return connection, nil
}

//NewListener creates a new socket then runs this socket as a go routine
func NewListener(port string/*, conn *Connection*/) *Listener {
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