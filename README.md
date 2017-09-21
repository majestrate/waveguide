# waveguide

video cdn server that offloads via webtorrent


## building

Requirements:

* go 1.9

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

