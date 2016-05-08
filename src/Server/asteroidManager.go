package main

import (
	"fmt"
)

//
type asteroid struct {
	x     int
	y     int
	size  int
	phase int
	input chan (Data)
}

type asteroidManager struct {
	xMax      int
	yMax      int
	input     chan (Data)
	asteroids []*asteroids
	// spawn frequency
	// max asteroids

}

func createAsteroidManager(sessionConn *Connection) {

}
