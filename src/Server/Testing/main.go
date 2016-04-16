package main

import (
	"fmt"
	"net"
	"io"
	"time"
)
func main() {
	addrs, err := net.Listen("tcp", ":9000")
	if (err != nil){
		panic(err)
	}

	io.Copy(conn, conn)

	conn.Close()
}

func say(print_this string, sync_channel chan bool){
	fmt.Println(print_this)
	sync_channel <- true
}

func say_channel(input_chan chan string) {
	print_this := <- input_chan
	fmt.Println(print_this)
	say_channel(input_chan)
}

