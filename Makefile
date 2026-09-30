CC=go
BIN_NAME=thesis2025-backend

build_arm: export GOOS := linux
build_arm: export GOARCH := arm
build_arm: export GOARM := 7

build_web:
	cd web && npm run build

build: build_web
	$(CC) build -o bin/$(BIN_NAME)

run_web:
	cd web && npm run dev

run: build
	./bin/$(BIN_NAME)

build_arm: build_web
	$(CC) build -o bin/$(BIN_NAME)-armv7
