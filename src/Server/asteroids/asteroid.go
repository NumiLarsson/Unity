package asteroids

import "math/rand"

// Asteroid represents a single asteroid
type Asteroid struct {
	ID    int
	X     int
	Y     int
	Phase int
	Alive bool
	xStep int
	yStep int
	size  int
	input chan (Data)
}

// loop runs an asteroid until an kill message is sent
func (asteroid *Asteroid) loop() {

	for {
		select {
		case msg := <-asteroid.input:

			if msg.action == terminate {
				asteroid.input <- Data{terminated, ok}
				return
			}

			asteroid.move()
		}
	}

}

// IsAlive checks if an asteroid still is alive
func (asteroid *Asteroid) isAlive() bool {
	return asteroid.Alive
}

// move updates the asteroids location with each tick
func (asteroid *Asteroid) move() {
	asteroid.X += asteroid.xStep
	asteroid.Y += asteroid.yStep
}

// inBounds checks if a given asteroid a is inside the bounds
func (asteroid *Asteroid) inBounds(manager *asteroidManager) bool {

	return asteroid.X >= 0 &&
		asteroid.Y >= 0 &&
		asteroid.X <= manager.xMax &&
		asteroid.Y <= manager.yMax

}

// newAsteroid allocates a new astroid
func newAsteroid() *Asteroid {
	return new(Asteroid)
}

// init sets the asteroids values, id,channel and spawn location
func (asteroid *Asteroid) init(id int, xMax int, yMax int) {

	asteroid.ID = id
	asteroid.size = 10
	asteroid.Alive = true

	asteroid.randowSpawn(xMax, yMax)
	//asteroid.checkSizeToWorld(xMax, yMax)
	asteroid.input = make(chan Data)
}

// randomSpawn sets the location at which a asteroid is spawned
func (asteroid *Asteroid) randowSpawn(xMax int, yMax int) {

	randomSpot := rand.Intn(4)
	randomDir := rand.Intn(3) - 1

	switch randomSpot {
	case 0:
		// Spawn top edge
		asteroid.X = rand.Intn(xMax)
		asteroid.Y = 0 - asteroid.size
		asteroid.xStep = randomDir
		asteroid.yStep = 1

	case 1:
		// Spawn right edge
		asteroid.X = xMax
		asteroid.Y = rand.Intn(yMax)
		asteroid.xStep = -1
		asteroid.yStep = randomDir

	case 2:
		// Spawn bottom edge
		asteroid.X = rand.Intn(xMax)
		asteroid.Y = yMax
		asteroid.xStep = randomDir
		asteroid.yStep = -1

	case 3:
		// Spawn left edge
		asteroid.X = 0
		asteroid.Y = rand.Intn(yMax)
		asteroid.xStep = 1
		asteroid.yStep = randomDir
	}

}
