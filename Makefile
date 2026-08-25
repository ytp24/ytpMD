.PHONY: all build fmt test install deb dist clean

VERSION ?= 3.1.0

all: fmt test build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ytp24 ./cmd/ytp24
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/pdf2md ./cmd/pdf2md

install: build
	mkdir -p $(HOME)/.local/bin
	cp bin/ytp24 $(HOME)/.local/bin/ytp24
	cp bin/pdf2md $(HOME)/.local/bin/pdf2md
	@echo "Installed ytp24 and pdf2md to $(HOME)/.local/bin/"

deb:
	./scripts/build-deb.sh

dist:
	./scripts/build-dist.sh

clean:
	rm -rf bin dist /tmp/ytp24*
