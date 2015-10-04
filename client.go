// https://facebook.github.io/watchman/docs/socket-interface.html

package kovacs

import (
	"encoding/json"
	"errors"
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

// a mix of base fields (version, error) and other fields
// that indicate that is either an event or an error has occurred
type rawResp struct {
	Version      string  `json:"version"`      // base field
	Error        *string `json:"error"`        // an error response
	Warning      *string `json:"warning"`      // a supplemental warning (not handled right now)
	Log          *string `json:"log"`          // log event
	Subscription *string `json:"subscription"` // subscription event
}

type resp struct {
	rawResp
	Raw json.RawMessage `json:"-"` // store the raw json
}

func (r *resp) UnmarshalJSON(b []byte) error {
	if err := json.Unmarshal(b, &r.rawResp); err != nil {
		return err
	}

	r.Raw = make(json.RawMessage, len(b))
	copy(r.Raw, b)

	return nil
}

// NewClient returns a new Client.  Connect must be called before any other
// client methods are used.  An optional logHandler may be passed in to
// handle an watchman generated log messages.
func NewClient(logHandler func(string)) *Client {
	return &Client{
		conn:        nil,
		reqCh:       make(chan *req),
		closeCh:     make(chan chan error),
		logHandler:  logHandler,
		subHandlers: map[string]func(*SubscriptionEvent){},
	}
}

// A Client manages a single socket connection to the watchman server. Connect
// must be called before any other method is used.
type Client struct {
	conn        *net.UnixConn
	reqCh       chan *req
	closeCh     chan chan error
	logHandler  func(string)
	subHandlers map[string]func(*SubscriptionEvent)
}

// Connect initializes the connection the watchman server.  It assumes that
// watchman server is running locally and attempts a unix socket connection
// addr is the path to the watchman server socket location. If an empty string
// is provided, Client will attempty to infer the location from the env var
// WATCHMAN_SOCK and, if that fails, by shelling out a `watchman get-sockname` call
func (c *Client) Connect(addr string) error {
	if addr == "" {
		logf("locating watchman socket")

		var err error
		addr, err = socketLoc()

		if err != nil {
			return err
		}
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
		respCh = make(chan *resp)
		closed int32
	)

	// read json values off the socket and send back to main
	// select loop
	go func() {
		for {
			if atomic.LoadInt32(&closed) != 0 {
				return
			}

			var r resp

			if err := dec.Decode(&r); err != nil {
				panic("JSON decoding error " + err.Error())
			}

			respCh <- &r
		}
	}()

	var curReq *req

OUTER:
	for {
		select {
		case resp := <-respCh:
			if resp.Log != nil {
				if c.logHandler != nil {
					c.logHandler(*resp.Log)
				}

				continue
			}

			if resp.Subscription != nil {
				var ev SubscriptionEvent

				if err := json.Unmarshal(resp.Raw, &ev); err != nil {
					panic("JSON decoding error " + err.Error())
				}

				c.subHandlers[ev.Subscription](&ev)
				continue
			}

			if curReq == nil {
				panic("Got a response without a request: " + string(resp.Raw))
			}

			req := curReq
			curReq = nil

			if resp.Error != nil {
				req.respCh <- errors.New(*resp.Error)
				continue
			}

			if req.dest != nil {
				if err := json.Unmarshal(resp.Raw, req.dest); err != nil {
					req.respCh <- err
					continue
				}
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

// Close shuts down the connection to the server
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
