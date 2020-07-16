## simple makefile to log workflow
.PHONY: all test build install dep

GOFLAGS ?= $(GOFLAGS:)

all: dep build install

build:
	@go build -buildmode=plugin $(GOFLAGS) ./...

clean:
	rm example.so

dep:
	@go get ./...

install:
	@go install $(GOFLAGS) ./...

# Currently we run all tests with "-p 1" to ensure that database operations do not interfere with each other

test: dep
	go test ./...
