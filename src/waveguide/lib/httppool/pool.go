package httppool

import (
	"errors"
	"net/http"
)

var ErrClosed = errors.New("http pool closed")

type Client struct {
	reqChan chan *request
	http    http.Client
	closing bool
}

type response struct {
	resp *http.Response
	err  error
}

type request struct {
	reply   chan response
	request *http.Request
}

func (c *Client) runWorker() {
	for {
		req := <-c.reqChan
		if req == nil {
			return
		}
		resp, err := c.http.Do(req.request)
		req.reply <- response{resp: resp, err: err}
	}
}

func (c *Client) Close() error {
	c.closing = true
	c.reqChan <- nil
	return nil
}

func (c *Client) Get(url string) (resp *http.Response, err error) {
	var req *http.Request
	req, err = http.NewRequest("GET", url, nil)
	if err == nil {
		resp, err = c.Do(req)
	}
	return
}

func (c *Client) Do(req *http.Request) (resp *http.Response, err error) {
	if c.closing {
		err = ErrClosed
		return
	}
	reply := make(chan response)
	c.reqChan <- &request{request: req, reply: reply}
	r := <-reply
	resp, err = r.resp, r.err
	return
}

func New(workers int) *Client {
	cl := &Client{
		reqChan: make(chan *request, 128),
	}
	if workers <= 0 {
		workers = 1
	}
	for workers > 0 {
		go cl.runWorker()
		workers--
	}
	return cl
}
