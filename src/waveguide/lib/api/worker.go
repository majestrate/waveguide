package api

import (
	"encoding/json"
	"errors"
	"github.com/streadway/amqp"
	"waveguide/lib/config"
	"waveguide/lib/log"
)

var ErrWorkerNotReady = errors.New("worker not ready")

type WorkerFunc func(*Request) error
type WorkFinder func(string) (WorkerFunc, bool)

type Worker struct {
	findWorker WorkFinder
	conn       *mqConn
}

func NewWorker(c *config.MQConfig, f WorkFinder) *Worker {
	return &Worker{
		conn:       newConn(c),
		findWorker: f,
	}
}

func (w *Worker) Run() error {
	if w.findWorker == nil {
		return ErrWorkerNotReady
	}
	return w.conn.ensureQueue(func(q amqp.Queue, ch *amqp.Channel) error {
		chnl, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait,
			nil)    // args
		if err != nil {
			return err
		}
		for delivery := range chnl {
			req := new(Request)
			err = json.Unmarshal(delivery.Body, req)
			if err == nil {
				worker, ok := w.findWorker(req.Method)
				if ok {
					err = worker(req)
				} else {
					err = ErrNoSuchMethod
				}
			}
			if err != nil {
				log.Errorf("worker failed to dispatch job: %s", err)
			}
		}
		return nil
	})
}
