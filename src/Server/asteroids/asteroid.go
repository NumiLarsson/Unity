package asteroids

import "math/rand"

type asteroid struct {
	X     int
	Y     int
	ID    int
	Phase int
	xStep int
	yStep int
	size  int
	input chan (Data)
}



func (asteroid *asteroid) loop() { //loop(id int, XMax int, yMax int) {

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

func (asteroid *asteroid) move() {

	asteroid.X += asteroid.xStep
	asteroid.Y += asteroid.yStep
}

// inBounds checks if a given asteroid a is inside the bounds
func (asteroid *asteroid) inBounds(manager *asteroidManager) bool {

	return asteroid.X >= 0 &&
		asteroid.Y >= 0 &&
		asteroid.X <= manager.xMax &&
		asteroid.Y <= manager.yMax

}

func newAsteroid() *asteroid {

	return new(asteroid)

}

func (asteroid *asteroid) init(id int, xMax int, yMax int) {

	asteroid.ID = id
	asteroid.randowSpawn(xMax, yMax)

	//	asteroid.checkSizeToWorld(xMax, yMax)

	asteroid.input = make(chan Data)
}

func (asteroid *asteroid) randowSpawn(xMax int, yMax int) {

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
