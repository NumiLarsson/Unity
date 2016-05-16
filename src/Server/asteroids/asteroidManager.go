package asteroids

import (
	"fmt"
	//"math/rand"
)

// asteroidManager stores info regarding gameworlds boundaries, all asteroids etc.
type asteroidManager struct {
	xMax      int
	yMax      int
	nextID    int
	maxRoids  int
	treshold  int
	deathRow  []int
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
				//TODO spawn on correct x/y
				manager.spawnAsteroid()
				manager.resumeAsteroids()

			} else {
				fmt.Println("[AST.MAN] Collision!! \n ", msg.action)
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
	/*
	r := rand.Intn(101)

	if len(manager.asteroids) < manager.maxRoids && r >= manager.treshold {
		//fmt.Println("SPAWNED NEW DROID")
		manager.newAsteroid()
	}
	
	Infinite spawning new asteroids every tick?
	*/
	if len(manager.asteroids) < manager.maxRoids {
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
func onDeathRow(a *asteroid, deathRow []int) bool {
	for _, dead := range deathRow {
		if a.ID == dead {
			return true
		}
	}
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

	fmt.Println("[AST.MAN] Created")
	return new(asteroidManager)

}

// init initiate the asteroid manager with hardcoded values TODO: input?
// and sets channels to session and
func (manager *asteroidManager) init(sessionConn *Connection, asteroids []*asteroid) {

	manager.xMax = 100
	manager.yMax = 100
	manager.asteroids = asteroids
	manager.newAsteroid()
	manager.newAsteroid()
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

func (manager *asteroidManager) updateDeathRow(deathRow []int) {
	manager.deathRow = deathRow

	if len(manager.deathRow) > 0 {
		fmt.Println("[AST.MAN] Collision:", manager.deathRow)
	}

}

// ONLY FOR TEST
func (manager *asteroidManager) print() {

	for _, asteroid := range manager.asteroids {
		fmt.Println("(", asteroid.ID, ",", asteroid.X, ",", asteroid.Y, ")")
	}
	/*
		fmt.Print("\033[2J\033[;H")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
		fmt.Println(". . . . . . . . . .")
	*/
}
