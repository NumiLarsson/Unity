all:
	go run main.go

doc:
	go doc
	(sleep 1; open http://localhost:6060/) &
	godoc -http=:6060

test:
	go test server.go session.go listenerManager.go asteroidManager.go asteroid.go asteroid_test.go
