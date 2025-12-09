CC=go
BIN_NAME=thesis2025-backend

build:
	$(CC) build -o bin/$(BIN_NAME)

run:
	$(CC) build -o bin/$(BIN_NAME)
	./bin/$(BIN_NAME)
