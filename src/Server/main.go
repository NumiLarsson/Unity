package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":9000")
	defer ln.Close()
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		defer conn.Close()
		if err != nil {
			panic(err)
		}

		io.WriteString(conn, fmt.Sprint("Hello World\n", time.Now(), "\n"))
	}
}

func say(print_this string, sync_channel chan bool) {
	fmt.Println(print_this)
	sync_channel <- true
}

func say_channel(input_chan chan string) {
	print_this := <-input_chan
	fmt.Println(print_this)
	say_channel(input_chan)
}
