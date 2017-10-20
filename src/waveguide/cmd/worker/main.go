package worker

import (
	"os"
	"os/signal"
	"syscall"
	"time"
	"waveguide/lib/api"
	"waveguide/lib/config"
	"waveguide/lib/log"
)

func Run() {
	const configFname = "waveguide.ini"
	log.SetLevel("debug")
	var app Worker
	var conf config.Config
	err := conf.Load(configFname)
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	err = app.Configure(&conf)
	if err != nil {
		log.Fatalf("failed to configure worker: %s", err)
	}
	sigchnl := make(chan os.Signal)
	signal.Notify(sigchnl, syscall.SIGHUP, os.Interrupt)
	w := api.NewWorker(&conf.MQ, app.FindWorker)
	runWorker := func(worker *api.Worker) {
		log.Info("running worker")
		for app.Continue() {
			err = worker.Run()
			if err != nil {
				log.Errorf("worker died: %s", err)
			}
			time.Sleep(time.Second)
		}
	}
	go runWorker(w)
	for sig := range sigchnl {
		switch sig {
		case os.Interrupt:
			log.Info("Interrupted")
			app.Stop()
			w.Close()
			log.Info("stopped worker")
			return
		case syscall.SIGHUP:
			err = conf.Load(configFname)
			if err != nil {
				log.Errorf("failed to reload config: %s", err)
				continue
			}
			log.Infof("reloaded %s", configFname)
			err = app.Reconfigure(&conf)
			if err != nil {
				log.Fatalf("failed to reconfigure worker: %s", err)
			}
			log.Info("restarting worker")
			w.Close()
			w = api.NewWorker(&conf.MQ, app.FindWorker)
			go runWorker(w)
		}
	}
}
