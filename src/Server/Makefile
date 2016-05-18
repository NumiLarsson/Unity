all:
	go run main.go -s

verbose:
	go run main.go 

race:
	go run -race main.go

doc:
	@go doc
	@echo "Documentation available at: http://localhost:6060/"
	@godoc -http=:6060

test:
	@cd asteroids && \
	go test *.go -v -coverprofile=coverage.out

openatom:
	@cd ../../.. && \
	atom ospp-2016-group-08
