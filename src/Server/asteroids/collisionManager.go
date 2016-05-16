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
		if a1.Alive {
			for _, a2 := range world.asteroids {

				if isCollisionAsteroidAsteroid(a1, a2) {
					a1.Alive = false
					a2.Alive = false
				}
			}
		}
	}

}

//
func (world *World) playerAsteroidCollision() {

	for _, p := range world.players {
		if p.Alive {
			for _, a := range world.asteroids {
				if isCollisionPlayerAsteroid(p, a) {
					p.Alive = false
					a.Alive = false
				}
			}
		}
	}
}

func (world *World) playerPlayerCollision() {

	for _, p1 := range world.players {
		if p1.Alive {
			for _, p2 := range world.players {
				if isCollisionPlayerPlayer(p1, p2) {
					p1.Alive = false
					p2.Alive = false
				}
			}
		}
	}
}

// isCollisionAsteroidAsteroid checks is if two asteroids are on
// the same position causing a collision
func isCollisionAsteroidAsteroid(a1 *Asteroid, a2 *Asteroid) bool {

	if a1.Id == a2.Id {
		return false
	} else if a1.X == a2.X && a1.Y == a2.Y {
		return true
	}

	return false

}

// isCollisionAsteroidPlayer  TODO some sort of interface to take generic input?
func isCollisionPlayerAsteroid(p *Player, a *Asteroid) bool {

	if p.X == a.X && p.Y == a.Y {
		return true
	}

	return false

}

// TODO some sort of interface to take generic input?
func isCollisionPlayerPlayer(p1 *Player, p2 *Player) bool {

	if p1.Id == p2.Id {
		return false
	} else if p1.X == p2.X && p1.Y == p2.Y {
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
