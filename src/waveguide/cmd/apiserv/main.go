package apiserv

import (
	"net"
	"os"
	"os/signal"
	"syscall"
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
		err = server.Configure(conf)
		if err != nil {
			log.Fatalf("failed to set up api server: %s", err.Error())
		}
		var listener net.Listener
		listener, err = net.Listen("tcp", conf.ApiServer.Addr)
		if err != nil {
			log.Fatal(err.Error())
		}
		go func() {
			sigchnl := make(chan os.Signal)
			signal.Notify(sigchnl, syscall.SIGHUP, os.Interrupt)
			for sig := range sigchnl {
				switch sig {
				case syscall.SIGHUP:
					log.Info("SIGHUP")
					e := conf.Load(configFname)
					if e == nil {
						log.Infof("reloaded config file %s", configFname)
						e = server.Configure(conf)
						if e == nil {
							log.Info("reloaded config and internal state is good")
						} else {
							log.Errorf("reconfigure failed: %s", e.Error())
						}
					} else {
						log.Errorf("failed to reload config: %s", e.Error())
					}
				case os.Interrupt:
					log.Info("closing api listener")
					listener.Close()
				}
			}
		}()

		log.Infof("api serving at http://%s/", conf.ApiServer.Addr)
		err = server.Serve(listener)
		if err != nil {
			log.Errorf("server error: %s", err.Error())
		}
	} else {
		log.Info("api server is disabled in config")
	}
}
