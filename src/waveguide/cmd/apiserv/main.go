package apiserv

import (
	"net"
	"waveguide/lib/api"
	"waveguide/lib/config"
	"waveguide/lib/log"
)

func Run() {
	var server api.Server
	var conf config.Config
	const configFname = "waveguide.ini"

	err := conf.Load(configFname)
	if err != nil {
		log.Fatalf("failed to load config: %s", err.Error())
	}

	if conf.ApiServer.Enabled {
		var listener net.Listener
		listener, err = net.Listen("tcp", conf.ApiServer.Addr)
		if err != nil {
			log.Fatal(err.Error())
		}
		// TODO: sighup
		log.Infof("api serving on %s", conf.ApiServer.Addr)
		err = server.Serve(listener)
		if err != nil {
			log.Errorf("api server: %s", err.Error())
		}
	} else {
		log.Info("api server is disabled in config")
	}
}
