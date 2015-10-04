package kovacs

import (
	"fmt"
	"os"
	"testing"
)

func TestClock(t *testing.T) {
	c := mustGetConnectedClient(t)
	c.Clock(testDir)
}

func TestWatch(t *testing.T) {
	c := mustGetConnectedClient(t)

	err := c.Watch(testDir)
	assert(t, err == nil, "unexpected watch error: %s", err)
}

func TestGetConfig(t *testing.T) {
	c := mustGetConnectedClient(t)

	conf, err := c.GetConfig(testDir)

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

	err = c.LogLevel("off")
	assert(t, err == nil, "unexpected config err: %s", err)
}

func TestVersion(t *testing.T) {
	c := mustGetConnectedClient(t)

	version, err := c.Version()
	assert(t, err == nil, "version err: %s", err)
	assert(t, version == "3.8.0", "incorrect version. expected 3.8.0, found %s", version)
}

func TestFind(t *testing.T) {
	c := mustGetConnectedClient(t)

	files, _, err := c.Find(testDir, "*.go")

	assert(t, err == nil, "find err: %s", err)
	assert(t, len(files) == numFiles, "expected %d files, found %d", numFiles, len(files))
}

func TestQuery(t *testing.T) {
	c := mustGetConnectedClient(t)
	wd, _ := os.Getwd()

	fmt.Println(c.Query(wd, QueryOptions{
		Suffix: []string{"go"},
	}))
}

func TestListCapabilities(t *testing.T) {
	c := mustGetConnectedClient(t)

	caps, err := c.ListCapabilities()

	assert(t, err == nil, "list capablities err: %s", err)
	fmt.Printf("caps = %+v\n", caps)
}
