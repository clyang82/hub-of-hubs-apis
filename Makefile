SHELL :=/bin/bash

all: build
.PHONY: all

build:
	go build -o hub-of-hubs-apis cmd/main.go
