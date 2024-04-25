all: test build

build:
	go build -o bin/kv .

test:
	go test -v ./...

clean:
	rm -f bin/kv

run: build
	./bin/kv

fmt:
	go fmt ./...

.PHONY: all build test clean run fmt
