package asteroids

import "fmt"

// asteroidCollision is used to check if two asteroids have collided
func (world *World) asteroidCollision() {

	for _, a1 := range world.asteroids {
		for _, a2 := range world.asteroids {
			if isCollision(a1.X, a1.Y, a2.X, a2.Y) && a1.Id != a2.Id {
				a1.Alive = false
			}
		}
	}

}

// playerCollision is used to check if a player has collided with another player or asteroid
func (world *World) playerCollision() {

	for _, p := range world.players {
		for _, a := range world.asteroids {
			if isCollision(p.X, p.Y, a.X, p.Y) {
				p.Alive = false
				a.Alive = false
			}
		}

		for _, p2 := range world.players {
			if isCollision(p.X, p.Y, p2.X, p2.Y) && p.Id != p2.Id {
				p.Alive = false
			}

		}
	}
}

// isCollision checks if two objects are located at the same position
func isCollision(x1 int, y1 int, x2 int, y2 int) bool {

	if x1 == x2 && y1 == y2 {
		return true
	}
	return false
}

// collisionManager used to get all collision that have occured during the last tick
func (world *World) collisionManager() {

	// First check player vs player and asteroid
	world.playerCollision()

	// Second check asteroid vs asteroid
	world.asteroidCollision()

	//////////////////////////////////////////////////////////
	// Below used to have the same prints as before///////////
	//////////////////////////////////////////////////////////
	var deadPlayerIDs []int
	var deadAsteroidIDs []int

	for _, player := range world.players {
		if player.Alive == false {
			deadPlayerIDs = append(deadPlayerIDs, player.Id)
		}
	}

	for _, asteroid := range world.asteroids {
		if asteroid.Alive == false {
			deadAsteroidIDs = append(deadAsteroidIDs, asteroid.Id)
		}
	}

	if len(deadPlayerIDs) > 0 || len(deadAsteroidIDs) > 0 {
		debugPrint(fmt.Sprintln("[COL.MAN] Collisions, Players:", deadPlayerIDs,
			"Asteroids:", deadAsteroidIDs))
	}

}

/////////////////////////////////////////////////////
/////////////////// below to be removed
/////////////////////////////////////////////////////
func (player *Player) death(world *World) {
	//Make player sleep for a second or two before respawning?
	player.Lives--
	player.respawn(world)
}

//Very primitive respawn, put the dead player back at start-position (0,0)
func (player *Player) respawn(world *World) {
	player.X = 0
	player.Y = 0
}
