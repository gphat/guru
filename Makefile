export GOPATH:=$(shell pwd)/_work

GURU_VERSION=$(shell git rev-parse HEAD)
CANONICAL=$(GOPATH)/src/github.com/gphat/guru
BUILDTAGS=debug
GURU=github.com/gphat/guru

all: prebuild
		go get -tags '$(BUILDTAGS)' -t -v $(GURU)
		go install -tags '$(BUILDTAGS)' -v $(GURU)

$(CANONICAL):
	mkdir -p $(GOPATH)/src/github.com/gphat
	ln -s $(shell pwd) $(GOPATH)/src/$(GURU)

prebuild: fmt $(CANONICAL)

fmt:
	go fmt ./...

clean:
	go clean -i -r $(GURU)

dist:
	GOOS=linux GOARCH=amd64 make release
