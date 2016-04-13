package main

import (
	"fmt"
)
func main() {
	output_channel := make(chan string)
	go say_channel(output_channel) 
	output_channel <- "Hello, World, Channel"
	output_channel <- "This will also be printed"
	output_channel <- "I can do this all day"
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

