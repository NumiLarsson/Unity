package asteroids

import "fmt"

// ======================= FIX MORE GENERIC =============================================
// detectCollisions checks each asteroid and stores all asteroids that have collided
// TODO players and use collision manager?
/*func (session *session) detectCollisions() ([]int, []int) {

	var asteroidCollisions []int
	var playerCollisions []int

	// First check player vs player collisions

}*/

func (world *World) asteroidAsteroidCollision() {

	for _, a1 := range world.asteroids {
		for _, a2 := range world.asteroids {
			if isCollision(a1.X, a1.Y, a2.X, a2.Y) && a1.Id != a2.Id {
				a1.Alive = false
			}
		}
	}

}

//
func (world *World) playerAsteroidCollision() {

	for _, p := range world.players {
		for _, a := range world.asteroids {
			if isCollision(p.X, p.Y, a.X, p.Y) {
				p.Alive = false
				a.Alive = false
			}
		}
	}
}

func (world *World) playerPlayerCollision() {

	for _, p1 := range world.players {
		for _, p2 := range world.players {
			if isCollision(p1.X, p1.Y, p2.X, p2.Y) && p1.Id != p2.Id {
				p1.Alive = false
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

// inList checks if the item is is already in the list
func inList(list []int, item int) bool {
	for _, current := range list {
		if item == current {
			return true
		}
	}
	return false
}

// ======================= FIX MORE GENERIC =============================================
// detectCollisions checks each asteroid and stores all asteroids that have collided

func (world *World) collisionManager() {

	// First check player vs player

	world.playerPlayerCollision()
	// second check player vs asteroid
	world.playerAsteroidCollision()
	// last check asteroid vs asteroid
	world.asteroidAsteroidCollision()

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

	//Used to make it compile
	if len(deadPlayerIDs) > 0 || len(deadAsteroidIDs) > 0 {
		fmt.Println("[COL.MAN] Collisions, Players:", deadPlayerIDs,
			"Asteroids:", deadAsteroidIDs)
	}

}

//Check the objects coordinates to see if collision occurs
//COuld be made more generic using overriding
func (player *Player) checkCoordinates(asteroid *Asteroid) bool {

	//Size of the objects will alter the collision hitbox
	//For now every object is only a dot

	if player.X == asteroid.X && player.Y == asteroid.Y {

		return true
	}

	return false
}

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
