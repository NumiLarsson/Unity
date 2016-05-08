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
	asteroids []*asteroids // Accessible from session.go
	// spawn frequency
	// max asteroids

}

func loop() {

		// 1. Iterate over each of the asteroids
		// 		- Check if it's inside the board, otherwise destroy it [remove from "shared" list]
		// 2. Session reads shared Data
		// 3. Session sends back any collisions/hits and whom it affects [possibly useful to store the asteroids in a map?]
		// 4. asteroidManager broadcasts to those affected and tells them to "die"
		// 5. asteroidManager removes the affected asteroids from the "shared" list
		// 6. Depending on the outcome and parameters asteroidManager may spawn additional asteroids
		// 7. REPEAT

}

func createAsteroidManager(sessionConn *Connection) {

		// Create relevant structs

		// Send confirmation back to Session

		loop()
}
