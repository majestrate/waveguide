REPO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

ifdef GOROOT
	GO = $(GOROOT)/bin/go
endif

GOPATH := $(REPO)

GO ?= $(shell which go)

WAVED := $(REPO)/waveguided

JS = $(REPO)/static/waveguide.min.js
GOPHERJS = $(GOPATH)/bin/gopherjs

all: build

build: $(WAVED) $(JS)



$(GOPHERJS):
	GOPATH=$(GOPATH) $(GO) get -v github.com/gopherjs/gopherjs

$(WAVED): $(GOPHERJS)
	GOPATH=$(GOPATH) $(GO) build -v -ldflags "-X waveguide/lib/version.Git=-$(shell git rev-parse --short HEAD)" -o $(WAVED)

$(JS): $(GOPHERJS)
	GOPATH=$(GOPATH) $(GOPHERJS) build waveguide/js/waveguide -v -m -o $(JS)

test:
	GOPATH=$(GOPATH) $(GO) test waveguide/...

clean:
	rm -fr $(GOPATH)/pkg
	rm -f $(WAVED)
	rm -f $(JS)

distclean: clean
	rm -f $(GOPHERJS)
	rm -fr $(GOPATH)/src/github.com
	rm -fr $(GOPATH)/src/golang.org
	rm -fr $(REPO)/node_modules
	rm -fr $(GOPATH)/pkg
