package asteroids

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
)

// asteroidManager stores info regarding gameworlds boundaries, all asteroids etc.
type asteroidManager struct {
	xMax      int
	yMax      int
	nextID    int
	maxRoids  int
	treshold  int
	deathRow  *[]int
	input     chan (Data)
	asteroids []*asteroid // Accessible from session.go
}

// loop …
func (manager *asteroidManager) loop(sessionConn *Connection, asteroids []*asteroid) {

	manager.init(sessionConn, asteroids)

	for {
		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				manager.removeDeadAsteroids()
				//manager.print()
				manager.spawnAsteroid()
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

// spawnAsteroid spawns a new asteroid if current asteroids in world below maxValue and
// if the randomized int that is set has a higher value than the worlds threshold
func (manager *asteroidManager) spawnAsteroid() {

	r := rand.Intn(101)

	if len(manager.asteroids) < manager.maxRoids && r >= manager.treshold {
		//fmt.Println("SPAWNED NEW DROID")
		manager.newAsteroid()
	}

}

// resumeAsteroids used to send "tick" to all asteroids
func (manager *asteroidManager) resumeAsteroids() {

	for _, asteroid := range manager.asteroids {
		asteroid.input <- Data{"a_manager.ok", 0}
	}

}

// onDeathRow TODO: implement! should check if current asteroid is on deathrow and can be removed
func onDeathRow(a *asteroid, deathRow *[]int) bool {
	return false
}

// removeDeadAsteroids used to check if any asteroid has been in a collision
// or if it's out of bounds
func (manager *asteroidManager) removeDeadAsteroids() {

	var offset = 0
	for i, asteroid := range manager.asteroids {

		// Check if inside kill list

		if onDeathRow(asteroid, manager.deathRow) || !asteroid.inBounds(manager) {
			manager.removeAsteroid(i + offset)
			offset--
		}

	}

}

// getAsteroids return the array containing the current asteroids
func (manager *asteroidManager) getAsteroids() []*asteroid {

	return manager.asteroids
}

// removeAsteroid removes specific asteroid from manager asteroid array
func (manager *asteroidManager) removeAsteroid(i int) {
	manager.asteroids = append(manager.asteroids[:i], manager.asteroids[i+1:]...)
}

// newObject creates a new asteroid, appends it to the asteroidmanagers array
// and creates a new go-routine that ......TODO
func (manager *asteroidManager) newAsteroid() {

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

// init initiate the asteroid manager with hardcoded values TODO: input?
// and sets channels to session and
func (manager *asteroidManager) init(sessionConn *Connection, asteroids []*asteroid) {

	manager.xMax = 100
	manager.yMax = 100
	manager.asteroids = asteroids
	manager.treshold = 20
	manager.maxRoids = 20
	manager.input = sessionConn.read

}

// getNextID returns the id to be used and sets the next value
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

// ONLY FOR TEST
func (manager *asteroidManager) printWorld(){


	for y := 0; y < manager.yMax ; y++ {
		fmt.Println("")
		for x := 0; x < manager.xMax ; x++{
			fmt.Print("* ")
			
		}
	}
	
	c := exec.Command("clear")
	c.Stdout = os.Stdout
	c.Run()
	
}


