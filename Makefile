all: bin/hello

bin/hello: src/github.com/gphat/hello/hello.go
	go install github.com/gphat/hello

clean:
	rm -rf bin/*
