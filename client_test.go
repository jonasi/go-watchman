package kovacs

import "testing"

func assert(t *testing.T, test bool, fmt string, args ...interface{}) {
	if test {
		return
	}

	t.Fatalf(fmt, args...)
}

func mustGetConnectedClient(t *testing.T) *Client {
	c := NewClient(nil)
	err := c.Connect()

	assert(t, err == nil, "connect error %s", err)

	return c
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	assert(t, c != nil, "nil client")
}

func TestClientConnect(t *testing.T) {
	c := NewClient(nil)
	err := c.Connect()

	assert(t, err == nil, "connect error %s", err)
}

func TestClientClose(t *testing.T) {
	c := mustGetConnectedClient(t)

	err := c.Close()
	assert(t, err == nil, "connect error %s", err)
}
