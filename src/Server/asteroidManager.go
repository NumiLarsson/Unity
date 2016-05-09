package main

import (
	"fmt"
	"math/rand"
)

type asteroids struct {
	x     int
	y     int
	size  int
	phase int
	input chan (Data)
}

type asteroidManager struct {
	xMax      int
	yMax      int
	input     chan (Data)
	asteroids []*asteroids // Accessible from session.go
	// spawn frequency
	// max asteroids

}

func AsteroidLoop(aManager *asteroidManager, sessionConn *Connection) {

	checkBoard(aManager)

	//sessionConn.write <- Data{"shared data",1}

	select {

	case msg := <-sessionConn.read:
		fmt.Println("Collision!! \n ", msg.action)
		// TODO: remove asteroids who has a collision or hit

	default:
		fmt.Println("default")
		break
	}

	// 1. Iterate over each of the asteroids
	// 		- Check if it's inside the board, otherwise destroy it [remove from "shared" list]
	// 2. Session reads shared Data
	// 3. Session sends back any collisions/hits and whom it affects [possibly useful to store the asteroids in a map?]
	// 4. asteroidManager broadcasts to those affected and tells them to "die"
	// 5. asteroidManager removes the affected asteroids from the "shared" list
	// 6. Depending on the outcome and parameters asteroidManager may spawn additional asteroids
	// 7. REPEAT

	print(aManager)

}

func checkBoard(aManager *asteroidManager) {

	var i = 0
	for _, a := range aManager.asteroids {

		if a.x >= 0 && a.y >= 0 && a.x < aManager.xMax && a.y < aManager.yMax {

		} else {
			removeAsteroid(i, aManager)
		}
		i++
	}

}

func removeAsteroid(i int, aManager *asteroidManager) {
	aManager.asteroids = append(aManager.asteroids[:i], aManager.asteroids[i+1:]...)
}

func createAsteroid(aManager *asteroidManager) {

	asteroid := new(asteroids)
	asteroid.x = rand.Intn(aManager.xMax)
	asteroid.y = rand.Intn(aManager.yMax)

	aManager.asteroids = append(aManager.asteroids, asteroid)
}

func createAsteroidManager(sessionConn *Connection, asteroids []*asteroids) {

	// Create relevant structs

	fmt.Println("AsteroidManager created")
	game := new(asteroidManager)
	game.xMax = 10
	game.yMax = 10

	// ONLY FOR TEST
	createAsteroid(game)
	createAsteroid(game)
	createAsteroid(game)

	// Send confirmation back to Session
	//sessionConn.write <- Data{"connect new manager",1}

	AsteroidLoop(game, sessionConn)
}

// ONLY FOR TEST
func print(aManager *asteroidManager) {
	for _, a := range aManager.asteroids {
		fmt.Println("(", a.x, ",", a.y, ")")
	}
}

// temp just for test
/*func main (){
	fmt.Println("Test MAIN started")
	ast := new(Connection)
	ast.read = make (chan Data)
	ast.write = make (chan Data)

	go createAsteroidManager(ast)

	time.Sleep(5 * time.Second)


}*/
