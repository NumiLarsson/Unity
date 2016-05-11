package main

import "./asteroids"

func main() {

	var server = asteroids.CreateServer()
	server.Listen(server.CreateFakeUser())

}
