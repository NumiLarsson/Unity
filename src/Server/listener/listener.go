package server

import (
    "net"
)

//Listener is responsible for a client each
//Contains a tcp socket, with the specified port at creation
type Listener struct {
    socket  net.Listener
    id      string
    port    string
    //Connection
}
/*Connection is temp until we establish the proper interface.
type Connection struct {
    read    chan string
    write   chan string
}
*/

func createSocket(port string) (net.Listener, error) {
    connection, err := net.Listen("tcp", ":" + port)
    if err != nil {
        return nil, err
    }
    
    return connection, nil
}

func newListener(port string/*, conn *Connection*/) *Listener {
    listener := new(Listener)
    
    var err error
    listener.socket, err = createSocket(port)
    if err != nil {
        panic(err)
    }
    
    
    listener.port = port
    
    //listener.write = conn.read //Fan in to manager
    //listener.read = conn.write //Fan out from manager
    
    return listener
}