TARGET = ""

.PHONY : all test cover install

default: install

all : test cover install

test:
	go test -v -race ./...
cover:
	go test -v -coverprofile=coverage.txt -covermode=atomic ./...
install:
	go install
