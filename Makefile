REPO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ifdef GOROOT
	GO = $(GOROOT)/bin/go
endif

GO ?= $(shell which go)

WAVED := $(REPO)/waveguided

GOPATH := $(REPO)

all: clean build test

build: $(WAVED)

$(WAVED):
	GOPATH=$(GOPATH) $(GO) build -v -ldflags "-X waveguide/lib/version.Git=-$(shell git rev-parse --short HEAD)" -o $(WAVED)

test:
	GOPATH=$(GOPATH) $(GO) test waveguide/...

clean:
	GOPATH=$(GOPATH) $(GO) clean
	rm -f $(WAVED)
