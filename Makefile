## simple makefile to log workflow
.PHONY: all test clean build install

SHELL := /bin/bash

GOFLAGS ?= $(GOFLAGS:)

export KIPROTECT_TEST = yes

KIPROTECT_TEST_CONFIG ?= "$(shell pwd)/config"
KIPROTECT_TEST_SETTINGS ?= "$(shell pwd)/testing/settings"

all: dep install

build:
	@go build $(GOFLAGS) ./...

dep:
	@go get ./...

install:
	@go install $(GOFLAGS) ./...

copyright:
	python .scripts/make_copyright_headers.py

init:
	RABBITMQ_VHOST=kiprotect_test RABBITMQ_USER=kiprotect RABBITMQ_PASSWORD=kiprotect .scripts/init_rabbitmq.sh
	RABBITMQ_VHOST=kiprotect_development RABBITMQ_USER=kiprotect RABBITMQ_PASSWORD=kiprotect .scripts/init_rabbitmq.sh

# Currently we run all tests with "-p 1" to ensure that database operations do not interfere with each other

plugins: plugins/writers/example/example.so

plugins/writers/example/example.so: plugins/writers/example/example.go
	@cd plugins/writers/example; make

test: dep plugins
	@KIPROTECT_CONFIG=$(KIPROTECT_TEST_CONFIG) KIPROTECT_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test $(testargs) -p 1 -count=1 `go list ./...`

test-races: dep plugins
	@KIPROTECT_CONFIG=$(KIPROTECT_TEST_CONFIG) KIPROTECT_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test -race $(testargs) -p 1 -count=1 `go list ./...`


bench: dep
	@KIPROTECT_CONFIG=$(KIPROTECT_TEST_CONFIG) KIPROTECT_SETTINGS=$(KIPROTECT_TEST_SETTINGS) go test -p 1 -run=NONE -bench=. $(GOFLAGS) `go list ./... | grep -v api/`

clean:
	@go clean $(GOFLAGS) -i ./...
