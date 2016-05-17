package asteroids

import (
	"fmt"
	"math/rand"
)

// asteroidManager stores info regarding gameworlds boundaries, all asteroids etc.
type asteroidManager struct {
	xMax      int
	yMax      int
	nextId    int
	maxRoids  int
	treshold  int
	deathRow  []int
	input     chan (Data)
	asteroids []*Asteroid // Accessible from session.go
}

// loop …
func (manager *asteroidManager) loop(sessionConn *Connection, asteroids []*Asteroid) {

	manager.init(sessionConn, asteroids)

	for {
		//manager.print()
		select {

		case msg := <-manager.input:

			if msg.action == "session.tick" {
				//manager.print()
				manager.updateDeathRow()
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
func onDeathRow(a *Asteroid, deathRow []int) bool {
	for _, dead := range deathRow {
		if a.Id == dead {
			return true
		}
	}
	return false
}

// removeDeadAsteroids used to check if any asteroid has been in a collision
// or if it's out of bounds
func (manager *asteroidManager) removeDeadAsteroids() {

	var offset = 0

	//	fmt.Println("before",len(manager.asteroids))

	var acopy = make([]*Asteroid, len(manager.asteroids))
	copy(acopy, manager.asteroids)

	for i, asteroid := range acopy {


		// Check if inside kill list

		if !asteroid.isAlive() || !asteroid.inBounds(manager) {
			manager.removeAsteroid(i + offset)
			offset--
		}
	}
	//fmt.Println("After",len(manager.asteroids))

}

// getAsteroids return the array containing the current asteroids
func (manager *asteroidManager) getAsteroids() []*Asteroid {

	return manager.asteroids
}

// removeAsteroid removes specific asteroid from manager asteroid array
func (manager *asteroidManager) removeAsteroid(i int) {
	//fmt.Println("i:",i)

	manager.asteroids = append(manager.asteroids[:i], manager.asteroids[i+1:]...)

}

// newObject creates a new asteroid, appends it to the asteroidmanagers array
// and creates a new go-routine that ......TODO
func (manager *asteroidManager) newAsteroid() {

	asteroid := newAsteroid()
	manager.asteroids = append(manager.asteroids, asteroid)

	asteroid.init(manager.getNextId(), manager.xMax, manager.yMax)
	go asteroid.loop()

}

// newAsteroidsManager creates a new asteroid manager
func newAsteroidManager() *asteroidManager {

	fmt.Println("[AST.MAN] Created")
	return new(asteroidManager)

}

// init initiate the asteroid manager with hardcoded values TODO: input?
// and sets channels to session and
func (manager *asteroidManager) init(sessionConn *Connection, asteroids []*Asteroid) {
	// TODO fix hardcoded variables
	manager.xMax = 100
	manager.yMax = 100
	manager.asteroids = asteroids
	manager.treshold = 20
	manager.maxRoids = 20
	manager.input = sessionConn.read

}

// getNextID returns the id to be used and sets the next value
func (manager *asteroidManager) getNextId() int {
	var id = manager.nextId
	manager.nextId++
	return id
}

func (manager *asteroidManager) updateDeathRow() {

	var deathRow []int

	for _, asteroid := range manager.asteroids {
		if !asteroid.isAlive() {
			deathRow = append(deathRow, asteroid.Id)
		}
	}

	manager.deathRow = deathRow

	if len(manager.deathRow) > 0 {
		fmt.Println("[AST.MAN] Collision:", manager.deathRow)
	}
}

// ONLY FOR TEST
func (manager *asteroidManager) print() {

	var list []int

	for _, asteroid := range manager.asteroids {
		list = append(list, asteroid.Id)

	}
	fmt.Println(len(manager.asteroids))
	fmt.Println(list)

}
