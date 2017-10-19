# waveguide

video cdn server that offloads via webtorrent


## building

Requirements:

* go 1.9

* yarn

* GNU make

* postgresql

* rabbitmq

Building:

    $ make
    
## Running (demo)

copy `contrib/waveguide.ini` to `waveguide.ini`

Run frontend server:

    $ PORT=28800 ./waveguided frontend
    
Run a worker daemon or as many instances as desired:
    
    $ ./waveguided worker

Run CDN Server, make sure to firewall it.

    $ PORT=48800 ./waveguided cdn


## project structure

### src/

all golang code

### src/waveguide/lib

all core libraries for all daemons

### src/vendor

vendored dependencies

### src/waveguide/cmd/worker

worker daemon

### src/waveguide/cmd/frontend

frontend webapp

### src/waveguide/cmd/cdn

simple demo cdn server

### js

javascript frontend code


### static

static web files

### contrib

configs and stuff

### templates

golang templates for webapp


## makefile targets

### all

default target, runs `clean` and `build`

### js

only build js frontend

### go

only build webapp daemon

### clean

clean all files

### distclean

clear all js files and run `clean` target

### clean-js

clean js frontend files only

### clean-go

clean go files only

