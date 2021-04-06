.PHONY: register

BIN = ./bin
MODULE = github.com/dimchat/demo-go
BUILD = env GO111MODULE=on go build
CLEAN = env GO111MODULE=on go clean

register:
	$(BUILD) -o $(BIN)/register $(MODULE)/register

all: register

clean:
	$(CLEAN) -cache
	rm -rf ./bin
