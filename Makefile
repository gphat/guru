all: bin/guru

bin/guru: src/github.com/gphat/guru/guru.go
	go install github.com/gphat/guru

clean:
	rm -rf bin/*
