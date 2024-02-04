.PHONY: build test clean

GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GOCLEAN=$(GOCMD) clean
SRC_PATH=cmd/bluecoins-to-splitwise-go
# BUILD_PATH=build
BINARY_NAME=bluecoins-to-splitwise

all: test build

build: 
	cd $(SRC_PATH) && $(GOBUILD) -o ../../$(BINARY_NAME) -v

test: 
	$(GOTEST) -v ./...

clean: 
	$(GOCLEAN)
	rm -f $(BINARY_NAME)