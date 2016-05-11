package main

import (
    "net"
    "time"
)

func main() {
    connection, err1 := net.Listen("tcp", ":9000")
    if err1 != nil {
        panic(err1)
    }
    
    var conn net.Conn
    conn, err2 := connection.Accept()
    if err2 != nil {
        panic(err2)
    }
    
    for {
        time.Sleep(time.Second)
        bytes := make([]byte, 1024)
        time := time.Now()
        tempString := time.Format("Mon Jan _2 15:04:05 2006")
        bytes = []byte(tempString)
        conn.Write(bytes)
    }
}