run:
	go run src/main.go
	
build:
	go build -o bin/app src/main.go

test:
	go test ./...

dev:
	air -c .air.toml

swag:
	swag init -g src/main.go