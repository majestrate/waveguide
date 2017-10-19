# waveguide

video cdn server that offloads via webtorrent


## building

Requirements:

* go 1.9

* yarn

* GNU make

* postgresql (for demo)

Building:

    $ make
    
## Running (demo)

copy `contrib/waveguide.ini` to `waveguide.ini`

Run frontend server:

    $ PORT=28800 ./waveguided frontend
    
Run worker server:
    
    $ ./waveguided worker

Run CDN Server, make sure to firewall it.

    $ PORT=48800 ./waveguided cdn


## project structure

### js

javascript frontend code

### src

go webapp code

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

