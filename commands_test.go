package watchman

import (
	"fmt"
	"os"
	"testing"
)

var wd string

func init() {
	wd, _ = os.Getwd()
}

func TestWatch(t *testing.T) {
	c := mustGetConnectedClient(t)

	err := c.Watch(wd)
	assert(t, err == nil, "unexpected watch error: %s", err)
}

func TestGetConfig(t *testing.T) {
	c := mustGetConnectedClient(t)

	conf, err := c.GetConfig(wd)

	assert(t, err == nil, "unexpected config err: %s", err)
	assert(t, conf != nil, "unexpected nil config")
}

func TestGetSockname(t *testing.T) {
	c := mustGetConnectedClient(t)

	sock, err := c.GetSockname()

	assert(t, err == nil, "unexpected config err: %s", err)
	assert(t, sock != "", "unexpected empty sockname")
}

func TestLog(t *testing.T) {
	c := mustGetConnectedClient(t)

	ok, err := c.Log("debug", "HELLO")

	assert(t, err == nil, "unexpected config err: %s", err)
	assert(t, ok, "unexpected log response")
}

func TestLog_InvalidLevel(t *testing.T) {
	c := mustGetConnectedClient(t)

	ok, err := c.Log("asf", "HELLO")

	assert(t, err != nil, "unexpected nil error")
	assert(t, !ok, "unexpected log response")
}

func TestLogLevel(t *testing.T) {
	c := mustGetConnectedClient(t)
	err := c.LogLevel("debug")

	assert(t, err == nil, "unexpected config err: %s", err)

	ok, _ := c.Log("error", "HELLO")
	fmt.Printf("ok = %+v\n", ok)
}

func TestVersion(t *testing.T) {
	c := mustGetConnectedClient(t)

	version, err := c.Version()
	assert(t, err == nil, "version err: %s", err)
	assert(t, version == "3.8.0", "incorrect version. expected 3.8.0, found %s", version)
}

func TestFind(t *testing.T) {
	c := mustGetConnectedClient(t)
	wd, _ := os.Getwd()

	fmt.Println(c.Find(wd, "*.go"))
}

func TestQuery(t *testing.T) {
	c := mustGetConnectedClient(t)
	wd, _ := os.Getwd()

	fmt.Println(c.Query(wd, QueryConfig{
		Suffix: []string{"go"},
	}))
}
