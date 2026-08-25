.PHONY: all build fmt test install clean

all: fmt test build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	mkdir -p bin
	go build -o bin/ytp24 ./cmd/ytp24
	go build -o bin/pdf2md ./cmd/pdf2md

install: build
	mkdir -p $(HOME)/.local/bin
	cp bin/ytp24 $(HOME)/.local/bin/ytp24
	cp bin/pdf2md $(HOME)/.local/bin/pdf2md
	@echo "Installed ytp24 and pdf2md to $(HOME)/.local/bin/"

clean:
	rm -rf bin
