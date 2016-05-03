package main

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

type World struct {
	X int
	Y int
}

func createWorld() World{
	world := World{300,300}

	fmt.Println("World created!")
	return world
}

func startLocation() World{
	position := World{1,1}
	fmt.Println("Start location set:", position.X, position.Y)
	return position
}

func location(world World){
	fmt.Println("Current location:", world.X, world.Y)
}

func getMove(world World, position World, direction string) World{
	switch direction{
	case "up":
		if position.Y == world.Y {
			fmt.Println("Cannot move that way")
		} else{
			position.Y += 1
			return position
		}

	case "down":
		if position.Y == world.Y || position.Y == 0 {
			fmt.Println("Cannot move that way")
		} else{
			position.Y -= 1
			return position
		}


	case "right":
		if position.X == world.X {
			fmt.Println("Cannot move that way")
		} else{
			position.X += 1
			return position
		}


	case "left":
		if position.X == world.X || position.X == 0{
			fmt.Println("Cannot move that way")
		} else{
			position.X -= 1
			return position
		}
	}
	return position

}
func move(world World, position World) World{
	fmt.Println("Which direction would you like to go?")
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter command: ")
	direction, _ := reader.ReadString('\n')

	switch strings.TrimSpace(direction){
	case "up":
		return getMove(world, position, "up")
	case "down":
		return getMove(world, position, "down")
	case "right":
		return getMove(world, position, "right")
	case "left":
		return getMove(world, position, "left")
	}
	return position
}

func clearScreen(){
	for i := 0; i < 10; i++{
		fmt.Println()
	}
}

func main() {
	var world *World = new(World)
	var position *World = new(World)

	loop := true
	*position = startLocation()
	*world = createWorld()

	for loop{
		fmt.Println()
		location(*position)
		fmt.Println("location/Move direction/quit")

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter command: ")
		text, _ := reader.ReadString('\n')
		switch strings.TrimSpace(text) {

		case "move":
			*position = move(*world, *position)
			break

		case "quit":
			loop = false
			break

			default:
			fmt.Println("Bad input, try again")
		
		}
	}
}
