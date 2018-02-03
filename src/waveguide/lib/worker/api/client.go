package api

import (
	"encoding/json"
	"github.com/streadway/amqp"
	"waveguide/lib/config"
)

type Client struct {
	conn *mqConn
}

func (cl *Client) Close() error {
	return cl.conn.Close()
}

func (cl *Client) Do(r *Request) error {
	body, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return cl.conn.ensureQueue(func(q amqp.Queue, ch *amqp.Channel) error {
		return ch.Publish(
			"",     // exchnage
			q.Name, // routing key
			false,  // manditory
			false,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  mqContentType,
				Body:         body,
			})
	})
}

func NewClient(mq *config.MQConfig) (*Client, error) {
	cl := &Client{
		conn: newConn(mq),
	}
	err := cl.conn.ensureConnection(func(_ *amqp.Connection) error {
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cl, nil
}
