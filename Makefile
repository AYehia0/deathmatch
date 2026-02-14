.PHONY: run build clean

run:
	go run cmd/server/main.go

build:
	go build -o bin/deathmatch cmd/server/main.go

clean:
	rm -rf bin/
