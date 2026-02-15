.PHONY: run build build-optimized clean deploy

run:
	go run cmd/server/main.go

build:
	go build -o bin/deathmatch cmd/server/main.go

build-optimized:
	@echo "Building optimized binary for production..."
	CGO_ENABLED=0 go build \
		-ldflags="-s -w" \
		-trimpath \
		-o bin/deathmatch \
		cmd/server/main.go
	@echo "Binary size:"
	@ls -lh bin/deathmatch
	@echo "Build complete!"

deploy: build-optimized
	@echo "Deploying to server..."
	@if [ -z "$(SERVER)" ]; then \
		echo "Error: SERVER not set. Usage: make deploy SERVER=user@host"; \
		exit 1; \
	fi
	scp bin/deathmatch $(SERVER):/root/deathmatch/bin/
	ssh $(SERVER) "sudo systemctl restart deathmatch && sudo systemctl status deathmatch --no-pager"
	@echo "Deployment complete!"

clean:
	rm -rf bin/
