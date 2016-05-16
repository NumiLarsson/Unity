package asteroids

import "math/rand"

type Asteroid struct {
	Id    int
	X     int
	Y     int
	Phase int
	Alive bool
	xStep int
	yStep int
	size  int

	input chan (Data)
}

func (asteroid *Asteroid) loop() { //loop(id int, xMax int, yMax int) {

	for {

		select {
		case msg := <-asteroid.input:

			if msg.action == "kill" {
				return
			}

			asteroid.move()

		}

	}

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

	asteroid.Id = id
	asteroid.Alive = true
	asteroid.randowSpawn(xMax, yMax)

	//	asteroid.checkSizeToWorld(xMax, yMax)

	asteroid.input = make(chan Data)

}

// randomSpawn sets the location at which a asteroid is spawned
func (asteroid *Asteroid) randowSpawn(xMax int, yMax int) {

	randomDir := rand.Intn(4)

	switch randomDir {
	case 0:
		asteroid.X = rand.Intn(xMax)
		asteroid.Y = 0 - asteroid.size
		asteroid.xStep = 0
		asteroid.yStep = 1

	case 1:
		asteroid.X = xMax
		asteroid.Y = rand.Intn(yMax)
		asteroid.xStep = -1
		asteroid.yStep = 0

	case 2:
		asteroid.X = rand.Intn(xMax)
		asteroid.Y = yMax
		asteroid.xStep = 0
		asteroid.yStep = -1

	case 3:
		asteroid.X = 0 - asteroid.size
		asteroid.Y = rand.Intn(yMax)
		asteroid.xStep = 1
		asteroid.yStep = 0
	}

}
