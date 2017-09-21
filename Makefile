REPO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ifdef GOROOT
	GO = $(GOROOT)/bin/go
endif

GOPATH := $(REPO)

GO ?= $(shell which go)

WAVED := $(REPO)/waveguided

JS = $(REPO)/static/waveguide.min.js

all: build

build: $(WAVED) js

$(WAVED):
	GOPATH=$(GOPATH) $(GO) build -v -ldflags "-X waveguide/lib/version.Git=-$(shell git rev-parse --short HEAD)" -o $(WAVED)

js: $(JS)

$(JS): 
	yarn install
	yarn dist

test:
	GOPATH=$(GOPATH) $(GO) test waveguide/...

clean: clean-js clean-go

clean-go:
	GOPATH=$(GOPATH) $(GO) clean -v
	rm -f $(WAVED)

clean-js:
	rm -f $(JS)

distclean: clean
	rm -fr $(REPO)/node_modules
