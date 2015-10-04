package kovacs

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"
	"time"
)

var (
	testDir  string
	numFiles int
)

func assert(t *testing.T, test bool, fmt string, args ...interface{}) {
	if test {
		return
	}

	t.Fatalf(fmt, args...)
}

func mustGetConnectedClient(t *testing.T) *Client {
	c := NewClient(nil)
	err := c.Connect(path.Join(testDir, "sock"))

	assert(t, err == nil, "connect error %s", err)

	return c
}

func TestMain(m *testing.M) {
	wd, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	testDir = path.Join(wd, "test")

	files, err := ioutil.ReadDir(wd)

	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".go") {
			numFiles++
		}
	}

	setup()
	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestNewClient(t *testing.T) {
	c := NewClient(nil)
	assert(t, c != nil, "nil client")
}

func TestClientConnect(t *testing.T) {
	c := NewClient(nil)
	err := c.Connect(path.Join(testDir, "sock"))

	assert(t, err == nil, "connect error %s", err)
}

func TestClientClose(t *testing.T) {
	c := mustGetConnectedClient(t)

	err := c.Close()
	assert(t, err == nil, "connect error %s", err)
}

var cmd *exec.Cmd

func setup() {
	os.Remove(path.Join(testDir, "sock"))
	os.Remove(path.Join(testDir, "log"))
	os.Remove(path.Join(testDir, "state"))

	ioutil.WriteFile(path.Join(testDir, "state"), []byte(`{}`), 0644)

	cmd = exec.Command(
		"watchman",
		"--foreground",
		"--statefile="+path.Join(testDir, "state"),
		"--logfile="+path.Join(testDir, "log"),
		"--log-level=2",
		"--sockname="+path.Join(testDir, "sock"),
	)

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	done := make(chan struct{})

	go func() {
		for {
			err := exec.Command("watchman", "--sockname="+path.Join(testDir, "sock"), "version").Run()

			if err == nil {
				done <- struct{}{}
			}

			time.Sleep(time.Millisecond)
		}
	}()

	select {
	case <-time.After(30 * time.Second):
		panic("watchman server was not initialized after 30 seconds")
	case <-done:
	}
}

func teardown() {
	if err := cmd.Process.Kill(); err != nil {
		panic(err)
	}
}
