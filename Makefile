APP_NAME=forge

.PHONY: build test lint run

build:
	go build -o $(APP_NAME) .

run: build
	./$(APP_NAME)

test:
	go test ./...

lint:
	@echo "(placeholder lint)"