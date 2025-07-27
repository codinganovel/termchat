.PHONY: build run test clean install

BINARY_NAME=termchat
BINARY_PATH=cmd/termchat
BUILD_DIR=build

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(BINARY_PATH)

run: build
	./$(BUILD_DIR)/$(BINARY_NAME)

test:
	go test -v ./...

clean:
	rm -rf $(BUILD_DIR)
	go clean

install:
	go install ./$(BINARY_PATH)

lint:
	golangci-lint run

fmt:
	go fmt ./...

deps:
	go mod download
	go mod tidy