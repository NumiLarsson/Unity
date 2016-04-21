package main

import (
	"fmt"
	"net"
	"time"
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
		
		conn.Write([]byte("Hello World"))
		time.Sleep(time.Second)
	}
}
