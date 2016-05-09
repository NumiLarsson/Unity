package main

import (
	"fmt"
)

type asteroid struct {
	x     int
	y     int
	id    int
	size  int
	phase int
	input chan (Data)
}

func move(asteroid *asteroid) {

	fmt.Println("Moving")

	for {

		select {
		case msg := <-asteroid.input:

			if msg.action == "kill" {
				return
			}

			//time.Sleep(500 * time.Millisecond)
			asteroid.x++

		}

	}

}

func spawnAsteroid(asteroid *asteroid) {

	move(asteroid)

}
