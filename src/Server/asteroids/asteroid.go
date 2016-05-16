package asteroids

import "math/rand"

type asteroid struct {
	x     int
	y     int
	xStep int
	yStep int
	id    int
	size  int
	phase int
	input chan (Data)
}

func (asteroid *asteroid) loop() { //loop(id int, xMax int, yMax int) {

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
func (asteroid *asteroid) move() {

	asteroid.x += asteroid.xStep
	asteroid.y += asteroid.yStep
}

// inBounds checks if a given asteroid a is inside the bounds
func (asteroid *asteroid) inBounds(manager *asteroidManager) bool {

	return asteroid.x >= 0 &&
		asteroid.y >= 0 &&
		asteroid.x <= manager.xMax &&
		asteroid.y <= manager.yMax

}

// newAsteroid allocates a new astroid
func newAsteroid() *asteroid {

	return new(asteroid)

}

// init sets the asteroids values, id,channel and spawn location
func (asteroid *asteroid) init(id int, xMax int, yMax int) {

	asteroid.id = id
	asteroid.randowSpawn(xMax, yMax)

	//	asteroid.checkSizeToWorld(xMax, yMax)

	asteroid.input = make(chan Data)

}

// randomSpawn sets the location at which a asteroid is spawned
func (asteroid *asteroid) randowSpawn(xMax int, yMax int) {

	randomDir := rand.Intn(4)

	switch randomDir {
	case 0:
		asteroid.x = rand.Intn(xMax)
		asteroid.y = 0 - asteroid.size
		asteroid.xStep = 0
		asteroid.yStep = 1

	case 1:
		asteroid.x = xMax
		asteroid.y = rand.Intn(yMax)
		asteroid.xStep = -1
		asteroid.yStep = 0

	case 2:
		asteroid.x = rand.Intn(xMax)
		asteroid.y = yMax
		asteroid.xStep = 0
		asteroid.yStep = -1

	case 3:
		asteroid.x = 0 - asteroid.size
		asteroid.y = rand.Intn(yMax)
		asteroid.xStep = 1
		asteroid.yStep = 0
	}

}
