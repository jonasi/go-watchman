package watchman

import "testing"

func TestVersion(t *testing.T) {
	c := mustGetConnectedClient(t)

	version, err := c.Version()
	assert(t, err == nil, "version err: %s", err)
	assert(t, version == "3.0.0", "incorrect version. expected 3.0.0, found %s", version)
}
