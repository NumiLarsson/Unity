package asteroids

import "fmt"


func (world *World) collisionManager() *World {

	//Check every player to see if they collide with an asteroid
	for _, player := range world.players {
		player.checkCollision(world)

	}
	//TODO
	//enter similar ranges with every object collidable
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

//Might split up in to several functions depending on the object
func (player *Player) checkCollision(world *World) {

	for _, asteroid := range world.asteroids {
		if player.checkCoordinates(asteroid) {
			fmt.Println("Player collided with asteroid at coordinates")

			fmt.Println("(", player.XCord, player.YCord, ")")

			player.death(world)
		}
	}

}

func (player *Player) death(world *World) {
	//Make player sleep for a second or two before respawning?
	player.Lives--
	player.respawn(world)
}

//Very primitive respawn, put the dead player back at start-position (0,0)
func (player *Player) respawn(world *World) {
	player.XCord = 0
	player.YCord = 0
}
