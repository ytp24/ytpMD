.PHONY: all build test install clean

all: test build

build:
	mkdir -p bin
	go build -o bin/pdf2md ./cmd/pdf2md
	cp bin/pdf2md bin/ytp24

test:
	go test -v ./...

install: build
	mkdir -p $(HOME)/.local/bin
	cp bin/pdf2md $(HOME)/.local/bin/pdf2md
	cp bin/pdf2md $(HOME)/.local/bin/ytp24
	@echo "Installed ytp24 and pdf2md to $(HOME)/.local/bin/"

clean:
	rm -rf bin
