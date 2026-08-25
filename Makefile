.PHONY: all build fmt test install deb dist clean

VERSION ?= 3.2.0

all: fmt test build

fmt:
	go fmt ./...

test:
	go test -v ./...

build:
	mkdir -p bin
	CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/ytpmd ./cmd/ytpmd
	cp -f bin/ytpmd bin/ytpMD
	cp -f bin/ytpmd bin/ytp24
	cp -f bin/ytpmd bin/pdf2md

install: build
	mkdir -p $(HOME)/.local/bin
	rm -f $(HOME)/.local/bin/ytpmd $(HOME)/.local/bin/ytpMD $(HOME)/.local/bin/ytp24 $(HOME)/.local/bin/pdf2md
	cp -f bin/ytpmd $(HOME)/.local/bin/ytpmd
	cp -f bin/ytpmd $(HOME)/.local/bin/ytpMD
	cp -f bin/ytpmd $(HOME)/.local/bin/ytp24
	cp -f bin/ytpmd $(HOME)/.local/bin/pdf2md
	@echo "Installed ytpmd (aliases: ytpMD, ytp24, pdf2md) to $(HOME)/.local/bin/"

deb:
	./scripts/build-deb.sh

dist:
	./scripts/build-dist.sh

clean:
	rm -rf bin dist /tmp/ytp*
