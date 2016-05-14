package asteroids

import "math/rand"

type asteroid struct {
	x     int
	y     int
	xStep int //These do not belong in the asteroid we send to C#
	yStep int //These do not belong in the asteroid we send to C#
	id    int //These do not belong in the asteroid we send to C#
	size  int //These do not belong in the asteroid we send to C#
	phase int 
	input chan (Data) //These do not belong in the asteroid we send to C#
}

//AsteroidClient is the public data that we send to clients
type AsteroidClient struct {
	x		int
	y 		int
	stage 	int
}

func (asteroid *asteroid) getClientData() *AsteroidClient{
	asteroidClient := new(AsteroidClient)
	asteroidClient.x = asteroid.x
	asteroidClient.y = asteroid.y
	asteroidClient.stage = asteroid.phase
	return asteroidClient
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
	asteroid.xStep = rand.Intn(3) - 1
	asteroid.yStep = rand.Intn(3) - 1

//	asteroid.checkSizeToWorld(xMax, yMax)

	asteroid.input = make(chan Data)
}
