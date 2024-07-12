# Makefile for build

SERVER=punch-server
CLIENT=punch-client

PLATFORMS=darwin linux windows
ARCHITECTURES=amd64 arm64

LDFLAGS=-ldflags '-s -w -extldflags "-static"' 


all: clean build_all

build:
	go build ${LDFLAGS} -o dist/${SERVER} server/server.go
	go build ${LDFLAGS} -o dist/${CLIENT} client/client.go

build_all:
	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o dist/$(SERVER)-$(GOOS)-$(GOARCH) server/server.go)))

	$(foreach GOOS, $(PLATFORMS),\
	$(foreach GOARCH, $(ARCHITECTURES), $(shell export GOOS=$(GOOS); export GOARCH=$(GOARCH); go build $(LDFLAGS) -o dist/$(CLIENT)-$(GOOS)-$(GOARCH) client/client.go)))
	@echo ok

clean:
	@rm -rf dist/punch-*

.PHONY: all build build_all clean
