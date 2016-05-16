package asteroids

import "fmt"


func (world *World) collisionManager() *World {
	var deadPlayerIDs []int
	var deadAsteroidIDs []int

	deadPlayerIDs, deadAsteroidIDs = checkCollision(world)

	//Used to make it compile
	if len(deadPlayerIDs) > 0 || len(deadAsteroidIDs) > 0 {
		fmt.Println("[COL.MAN] Collisions, Players:", deadPlayerIDs, "Asteroids:", deadAsteroidIDs)
	}

	//TODO
	//enter similar ranges with every object destructible
	return world
}

//Check the objects coordinates to see if collision occurs
//COuld be made more generic using overriding
func (player *Player) checkCoordinates(asteroid *asteroid) bool {

	//Size of the objects will alter the collision hitbox
	//For now every object is only a dot


	if player.XCord == asteroid.x && player.YCord == asteroid.y {

		return true
	}

	return false
}

//Checks the collisions during the tick and returns two arrays
//Of the player and asteroid IDs which were destroyed
//Could be made to act as a hub for every collision at once
//Thus becoming the real collisionManager (Consider changing name)
func checkCollision(world *World) ([]int, []int) {
	var deadPlayerIDs []int
	var deadAsteroidIDs []int

	for _, player := range world.players {
		for _, asteroid := range world.asteroids {
			if player.checkCoordinates(asteroid) {
				fmt.Println("[COL.MAN] Player collided with asteroid at coordinates (", player.x, player.y, ")")

				//Player collision with an asteroid will
				//Kill the player and the asteroid
				//It only makes sense... Right?
				deadPlayerIDs = append(deadPlayerIDs, player.id)
				deadAsteroidIDs = append(deadAsteroidIDs, asteroid.id)
				
			}
		}
		
		
	for _, asteroid := range world.asteroids {
		if player.checkCoordinates(asteroid) {
			fmt.Println("Player collided with asteroid at coordinates")

			fmt.Println("(", player.XCord, player.YCord, ")")
			
			player.death(world)
			
		}
	}
	}

	return deadPlayerIDs, deadAsteroidIDs
}

func (player *Player) death(world *World) {
	//Make player sleep for a second or two before respawning?
	player.lives--
	player.respawn(world)
}

//Very primitive respawn, put the dead player back at start-position (0,0)
func (player *Player) respawn(world *World) {
	player.x = 0
	player.y = 0
}
