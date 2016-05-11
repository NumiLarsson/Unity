package main

import "fmt"

// asteroidManager stores info regarding gameworlds boundaries, all asteroids etc.
type asteroidManager struct {
	xMax      int
	yMax      int
	nextID    int
	input     chan (Data)
	asteroids []*asteroid // Accessible from session.go
	// spawn frequency
	// max asteroids

}

// loop …
func (manager *asteroidManager) loop(sessionConn *Connection, asteroids []*asteroid) {

	manager.init(sessionConn, asteroids)

	for {

		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				manager.checkBoard()
				manager.print()
				manager.resumeAsteroids()

			} else {
				fmt.Println("Collision!! \n ", msg.action)
				// TODO: remove asteroids who has a collision or hit
			}
		}

		// 1. Iterate over each of the asteroids
		// 		- Check if it's inside the board, otherwise destroy it [remove from "shared" list]
		// 2. Session reads shared Data
		// 3. Session sends back any collisions/hits and whom it affects [possibly useful to store the asteroids in a map?]
		// 4. asteroidManager broadcasts to those affected and tells them to "die"
		// 5. asteroidManager removes the affected asteroids from the "shared" list
		// 6. Depending on the outcome and parameters asteroidManager may spawn additional asteroids
		// 7. REPEAT

	}

}

// resumeAsteroids used to send "tick" to all asteroids
func (manager *asteroidManager) resumeAsteroids() {

	for _, asteroid := range manager.asteroids {
		asteroid.input <- Data{"a_manager.ok", 0}
	}

}

// checkBoard used to check if any asteroid is out of bounds
func (manager *asteroidManager) checkBoard() {

	var offset = 0
	for i, asteroid := range manager.asteroids {

		if !asteroid.inBounds(manager) {
			fmt.Println("Asteroid out of bounds. Die!")
			manager.removeAsteroid(i + offset)
			offset -= 1
		}

	}

}

func (manager *asteroidManager) getObjects() []*asteroid {

	return manager.asteroids
}

// removeAsteroid removes specific asteroid from manager asteroid array
func (manager *asteroidManager) removeAsteroid(i int) {
	manager.asteroids = append(manager.asteroids[:i], manager.asteroids[i+1:]...)
}

// newObject creates a new asteroid, appends it to the asteroidmanagers array
// and creates a new go-routine that ......TODO
func (manager *asteroidManager) newObject() {

	asteroid := newAsteroid()
	manager.asteroids = append(manager.asteroids, asteroid)

	asteroid.init(manager.getNextID(), manager.xMax, manager.yMax)
	go asteroid.loop()

}

// newAsteroidsManager creates a new asteroid manager
func newAsteroidManager() *asteroidManager {

	fmt.Println("AsteroidManager created")
	return new(asteroidManager)

}

func (manager *asteroidManager) init(sessionConn *Connection, asteroids []*asteroid) {

	manager.xMax = 10
	manager.yMax = 10
	manager.asteroids = asteroids
	manager.input = sessionConn.read

	manager.newObject()
	manager.newObject()
	manager.newObject()
	manager.newObject()

	//fmt.Printf("%d\n", (len(manager.asteroids)))

	// Send confirmation back to Session
	//sessionConn.write <- Data{"connect new manager",1}

	//manager.loop(sessionConn)
}

func (manager *asteroidManager) getNextID() int {
	var id = manager.nextID
	manager.nextID++
	return id
}

// ONLY FOR TEST
func (manager *asteroidManager) print() {
	for _, asteroid := range manager.asteroids {
		fmt.Println("(", asteroid.id, ",", asteroid.x, ",", asteroid.y, ")")
	}
}
