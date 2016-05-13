package main

//import "./asteroids"
import "github.com/numilarsson/ospp-2016-group-08/src/server/asteroids"

func main() {

	var server = asteroids.CreateServer()
	server.Listen(make(chan asteroids.Data)/*server.CreateFakeUser()*/)

}
