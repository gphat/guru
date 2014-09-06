export GOPATH:=$(shell pwd)

CANONICAL=$(GOPATH)/src/github.com/gphat/guru

all: bin/guru

bin/guru: src/github.com/gphat/guru/main.go
	go install github.com/gphat/guru

clean:
	rm -rf bin/*
