# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=goflipdot-server
BINARY_UNIX=$(BINARY_NAME)_unix
BINARY_RPI=$(BINARY_NAME)_rpi

all: test build

build:
	$(GOBUILD) -o $(BINARY_NAME) -v

test:
	$(GOTEST) -v ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -f $(BINARY_RPI)

run:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...
	./$(BINARY_NAME)

deps:
	$(GOGET) github.com/harperreed/goflipdot/pkg/goflipdot
	$(GOGET) github.com/spf13/viper

tidy:
	$(GOMOD) tidy

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Build for Raspberry Pi
build-rpi:
	GOARCH=arm GOARM=7 GOOS=linux $(GOBUILD) -o $(BINARY_RPI) -v

docker-build:
	docker build -t $(BINARY_NAME):latest .

.PHONY: all build test clean run deps tidy build-linux build-rpi docker-build
