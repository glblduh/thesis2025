CC=go
BIN_NAME=thesis2025-backend

build_web:
	cd web && npm run build

build: build_web
	$(CC) build -o bin/$(BIN_NAME)

run: build
	./bin/$(BIN_NAME)
