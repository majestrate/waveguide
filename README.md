# waveguide

video cdn server that offloads via webtorrent


## building

Requirements:

* go 1.9

* GNU make

Building:

    $ make
    
## Running (demo)


copy `contrib/waveguided.ini` to `/etc/waveguided.ini`

Run frontend server:

    $ ./waveguided frontend
    
Run worker server:
    
    $ ./waveguided worker

