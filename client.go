// https://facebook.github.io/watchman/docs/socket-interface.html

package kovacs

import (
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"sync/atomic"
)

func socketLoc() (string, error) {
	if addr := os.Getenv("WATCHMAN_SOCK"); addr != "" {
		return addr, nil
	}

	var loc struct {
		Version  string
		Sockname string
	}

	b, err := exec.Command("watchman", "get-sockname").Output()

	if err != nil {
		return "", err
	}

	if err := json.Unmarshal(b, &loc); err != nil {
		return "", err
	}

	return loc.Sockname, nil
}

type req struct {
	dest   interface{}
	respCh chan error
	cmd    []interface{}
}

func NewClient(logHandler func(string)) *Client {
	return &Client{
		conn:        nil,
		reqCh:       make(chan *req),
		closeCh:     make(chan chan error),
		logHandler:  logHandler,
		subHandlers: map[string]func(*SubscriptionEvent){},
	}
}

type Client struct {
	conn        *net.UnixConn
	reqCh       chan *req
	closeCh     chan chan error
	logHandler  func(string)
	subHandlers map[string]func(*SubscriptionEvent)
}

func (c *Client) Connect() error {
	logf("locating watchman socket")
	addr, err := socketLoc()

	if err != nil {
		return err
	}

	logf("connecting to %s", addr)
	conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: addr, Net: "unix"})

	if err != nil {
		return err
	}

	c.conn = conn

	go c.listen()

	return nil
}

func (c *Client) listen() {
	var (
		enc    = json.NewEncoder(logWriter("request", c.conn))
		dec    = json.NewDecoder(logReader("response", c.conn))
		respCh = make(chan json.RawMessage)
		closed int32
	)

	go func() {
		for {
			var ev struct {
				Log          *string `json:"log"`
				Subscription *string `json:"subscription"`
			}

			if atomic.LoadInt32(&closed) != 0 {
				return
			}

			var m json.RawMessage

			if err := dec.Decode(&m); err != nil {
				panic("JSON decoding error " + err.Error())
			}

			if err := json.Unmarshal(m, &ev); err != nil {
				panic("JSON decoding error " + err.Error())
			}

			if ev.Log != nil {
				if c.logHandler != nil {
					c.logHandler(*ev.Log)
				}

				continue
			}

			if ev.Subscription != nil {
				var ev SubscriptionEvent

				if err := json.Unmarshal(m, &ev); err != nil {
					panic("JSON decoding error " + err.Error())
				}

				c.subHandlers[ev.Subscription](&ev)
				continue
			}

			respCh <- m
		}
	}()

	var curReq *req

OUTER:
	for {
		select {
		case resp := <-respCh:
			if curReq == nil {
				panic("Got a response without a request: " + string(resp))
			}

			req := curReq
			curReq = nil

			if err := json.Unmarshal(resp, req.dest); err != nil {
				req.respCh <- err
				continue
			}

			if err, ok := req.dest.(error); ok && err.Error() != "" {
				req.respCh <- err
				continue
			}

			req.respCh <- nil
		case req := <-c.reqCh:
			if err := enc.Encode(req.cmd); err != nil {
				req.respCh <- err
				continue
			}

			curReq = req
		case closeCh := <-c.closeCh:
			atomic.StoreInt32(&closed, 1)
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

func (c *Client) send(dest interface{}, args ...interface{}) error {
	req := req{
		dest:   dest,
		respCh: make(chan error),
		cmd:    args,
	}

	c.reqCh <- &req

	if err := <-req.respCh; err != nil {
		return err
	}

	return nil
}
