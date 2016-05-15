package main

import "./asteroids"

func main() {

	var server = asteroids.CreateServer()
	//server.Listen(make(chan asteroids.Data) /*server.CreateFakeUser()*/)
	server.Listen(server.CreateFakeUser())

}
