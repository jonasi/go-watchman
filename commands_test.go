package watchman

import (
	"fmt"
	"os"
	"testing"
)

func TestVersion(t *testing.T) {
	c := mustGetConnectedClient(t)

	version, err := c.Version()
	assert(t, err == nil, "version err: %s", err)
	assert(t, version == "3.0.0", "incorrect version. expected 3.0.0, found %s", version)
}

func TestFind(t *testing.T) {
	c := mustGetConnectedClient(t)
	wd, _ := os.Getwd()

	fmt.Println(c.Find(wd, "*.go"))
}
