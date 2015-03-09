package main

import (
	"encoding/json"
	"net"
)

type req struct {
	ptr    interface{}
	respCh chan error
	cmd    []string
}

func NewClient() *Client {
	cl := &Client{
		conn:    nil,
		reqCh:   make(chan req),
		closeCh: make(chan chan error),
	}

	return cl
}

type Client struct {
	conn    *net.UnixConn
	reqCh   chan req
	closeCh chan chan error
}

func (c *Client) Connect() error {
	logf("locating watchman socket")
	addr, err := socketLoc()

	if err != nil {
		return err
	}

	logf("connecting to %s", addr)
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{addr, "unix"})

	if err != nil {
		return err
	}

	c.conn = conn

	go c.listen()

	return nil
}

func (c *Client) listen() {
	var (
		enc = json.NewEncoder(logWriter("request", c.conn))
		dec = json.NewDecoder(logReader("response", c.conn))
	)

OUTER:
	for {
		select {
		case req := <-c.reqCh:
			if err := enc.Encode(req.cmd); err != nil {
				req.respCh <- err
				continue
			}

			if err := dec.Decode(req.ptr); err != nil {
				req.respCh <- err
				continue
			}

			if err, ok := req.ptr.(error); ok && err.Error() != "" {
				req.respCh <- err
				continue
			}

			req.respCh <- nil
		case closeCh := <-c.closeCh:
			closeCh <- c.conn.Close()
			break OUTER
		}
	}
}

func (c *Client) Close() error {
	ch := make(chan error)
	c.closeCh <- ch

	return <-ch
}

func (c *Client) send(ptr interface{}, args ...string) error {
	req := req{
		ptr:    ptr,
		respCh: make(chan error),
		cmd:    args,
	}

	c.reqCh <- req

	if err := <-req.respCh; err != nil {
		return err
	}

	return nil
}
