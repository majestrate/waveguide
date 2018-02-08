# waveguide

video sharing server that offloads via webtorrent

## building

Requirements:

* go 1.9.x

* yarn

* GNU make

* postgresql

* rabbitmq

* ffmpeg

* nginx-rtmp

Building:

    $ make
    
## Running (demo)

copy `contrib/waveguide.ini` to `waveguide.ini`

**make sure to set your domains correctly in the config**

Run it:

    $ ./waveguided
    
and throw nginx in front:

    # cp contrib/waveguide.nginx /etc/nginx/sites-enabled/waveguide


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

