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
    
    $ PORT=38800 ./waveguided worker



## project structure


### lib

javascript frontend code

### src

go webapp code

### static

static web files

### contrib

configs and stuff

### templates

golang templates for webapp
