package main

import (
	"fmt"
	"net"
	"time"
	"encoding/binary"
)

func main() {
	fmt.Println("Launching server...")
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		panic(err)
	}
	conn, err := ln.Accept()
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	fmt.Print("Connection Accepted")

	for {
		/*
			message, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				panic(err)
			}

			fmt.Print("Message Received:", string(message))

			newmessage := strings.ToUpper(message)

			conn.Write([]byte(newmessage + "\n"))
		*/
		//This is the echo function we used earlier
		var portUInt uint16 = 9001
		port := make([]byte, 2)
		binary.LittleEndian.PutUint16(port, portUInt);
		conn.Write(port)
		time.Sleep(time.Second)
	}
}
