.PHONY: all build test clean

all: test build

build:
	mkdir -p bin
	go build -o bin/pdf2md ./cmd/pdf2md

test:
	go test -v ./...

clean:
	rm -rf bin
