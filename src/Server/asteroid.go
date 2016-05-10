package main

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

func newAsteroid() *asteroid {

	return new(asteroid)

}

func (asteroid *asteroid) init(id int, xMax int, yMax int) {

	asteroid.x = rand.Intn(xMax)
	asteroid.y = rand.Intn(yMax)
	asteroid.id = id
	asteroid.xStep = 1
	asteroid.yStep = 0
	asteroid.input = make(chan Data)

}
