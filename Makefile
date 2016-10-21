include version.mk

.PHONY: all clean build vendor
all: build 

PACKAGE=github.com/gertd/git-list

build:
	mkdir -p bin/linux-amd64
	cp README.md bin/linux-amd64
	cp LICENSE bin/linux-amd64
	cd bin/linux-amd64 && GOOS=linux GOARCH=amd64 go build ../..

clean:
	rm -rf bin work 

vendor:
	godep save $(shell go list ./... | grep -v /vendor)
