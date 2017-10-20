package api

import (
	"github.com/streadway/amqp"
	"waveguide/lib/config"
)

func newConn(c *config.MQConfig) *mqConn {
	return &mqConn{
		mqConf: c,
	}
}

type mqConn struct {
	conn   *amqp.Connection
	chnl   *amqp.Channel
	mqConf *config.MQConfig
	queue  amqp.Queue
}

func (c *mqConn) Close() error {
	if c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

func (c *mqConn) ensureConnection(visit func(*amqp.Connection) error) (err error) {
	if c.conn == nil {
		c.conn, err = amqp.Dial(c.mqConf.URL)
	}
	if err == nil {
		err = visit(c.conn)
	}
	return err
}

func (c *mqConn) ensureQueue(visit func(amqp.Queue, *amqp.Channel) error) error {
	return c.ensureConnection(func(conn *amqp.Connection) error {
		var err error
		if c.chnl == nil {
			c.chnl, err = conn.Channel()
			if err == nil {
				c.queue, err = c.chnl.QueueDeclare(
					mqQueueName, // name
					false,       // durable
					false,       // delete when unused
					false,       // exclusive
					false,       // no-wait
					nil)
			}
		}
		if err == nil {
			err = visit(c.queue, c.chnl)
		} else {
			conn.Close()
			c.conn = nil
			c.chnl = nil
		}
		return err
	})
}
