## simple makefile to log workflow
.PHONY: all test clean build install

SHELL := /bin/bash

GOFLAGS ?= -ldflags="-extldflags=-static" $(GOFLAGS:)

export KIPROTECT_TEST = yes

KIPROTECT_TEST_CONFIG ?= "$(shell pwd)/config"
KIPROTECT_TEST_SETTINGS ?= "$(shell pwd)/testing/settings"
KIPROTECT_TEST_API_SETTINGS ?= "$(shell pwd)/testing/api/settings"

all: dep build install

build:
	@CGO_ENABLED=0 go build $(GOFLAGS) ./...

dep:
	@go get ./...

install:
	@CGO_ENABLED=0 go install $(GOFLAGS) ./...

copyright:
	python3 .scripts/make_copyright_headers.py

init:
	RABBITMQ_VHOST=kiprotect_test RABBITMQ_USER=kiprotect RABBITMQ_PASSWORD=kiprotect .scripts/init_rabbitmq.sh
	RABBITMQ_VHOST=kiprotect_development RABBITMQ_USER=kiprotect RABBITMQ_PASSWORD=kiprotect .scripts/init_rabbitmq.sh

# Currently we run all tests with "-p 1" to ensure that database operations do not interfere with each other

plugins: plugins/writers/example/example.so

plugins/writers/example/example.so: plugins/writers/example/example.go
	@cd plugins/writers/example; make

test: dep
	@KODEX_CONFIG=$(KIPROTECT_TEST_CONFIG) KODEX_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test $(GOTAGS) $(testargs) -p 1 -count=1 `go list ./... | grep -v api/`

test-api: dep
	@KODEX_CONFIG=$(KIPROTECT_TEST_CONFIG) KODEX_SETTINGS=$(KIPROTECT_TEST_API_SETTINGS) go test $(GOTAGS) $(testargs) -p 1 -count=1 `go list ./... | grep api/`

test-races: dep
	@KODEX_CONFIG=$(KIPROTECT_TEST_CONFIG) KODEX_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test -race $(testargs) -p 1 -count=1 `go list ./...`


bench: dep
	@KODEX_CONFIG=$(KIPROTECT_TEST_CONFIG) KODEX_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test -p 1 -run=NONE -bench=. $(GOFLAGS) `go list ./... | grep -v api/`

clean:
	@go clean $(GOFLAGS) -i ./...
	find plugins -name "*.so" -exec rm {} \;
