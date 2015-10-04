package kovacs

type baseResponse struct {
	Version string
	Err     string `json:"error"`
}

func (b *baseResponse) Error() string {
	return b.Err
}

// https://facebook.github.io/watchman/docs/cmd/clock.html
func (c *Client) Clock(dir string) (string, error) {
	var s struct {
		baseResponse
		Clock string
	}

	if err := c.send(&s, "clock", dir); err != nil {
		return "", err
	}

	return s.Clock, nil
}

// https://facebook.github.io/watchman/docs/cmd/find.html
func (c *Client) Find(dir string, patterns ...string) ([]File, string, error) {
	var s struct {
		baseResponse
		Clock string
		Files []File
	}

	params := []interface{}{"find", dir}
	for _, p := range patterns {
		params = append(params, p)
	}

	if err := c.send(&s, params...); err != nil {
		return nil, "", err
	}

	return s.Files, s.Clock, nil
}

// https://facebook.github.io/watchman/docs/cmd/get-config.html
func (c *Client) GetConfig(dir string) (*Config, error) {
	var s struct {
		baseResponse
		Config Config
	}

	if err := c.send(&s, "get-config", dir); err != nil {
		return nil, err
	}

	return &s.Config, nil
}

// https://facebook.github.io/watchman/docs/cmd/get-sockname.html
func (c *Client) GetSockname() (string, error) {
	var s struct {
		baseResponse
		Sockname string
	}

	if err := c.send(&s, "get-sockname"); err != nil {
		return "", err
	}

	return s.Sockname, nil
}

// https://facebook.github.io/watchman/docs/cmd/list-capabilities.html
func (c *Client) ListCapabilities() ([]string, error) {
	var s struct {
		baseResponse
		Capabilities []string `json:"capabilities"`
	}

	if err := c.send(&s, "list-capabilities"); err != nil {
		return nil, err
	}

	return s.Capabilities, nil
}

// https://facebook.github.io/watchman/docs/cmd/log.html
func (c *Client) Log(level, msg string) (bool, error) {
	var s struct {
		baseResponse
		Logged bool
	}

	if err := c.send(&s, "log", level, msg); err != nil {
		return false, err
	}

	return s.Logged, nil
}

// https://facebook.github.io/watchman/docs/cmd/log-level.html
func (c *Client) LogLevel(level string) error {
	var s struct {
		baseResponse
	}

	if err := c.send(&s, "log-level", level); err != nil {
		return err
	}

	return nil
}

// https://facebook.github.io/watchman/docs/cmd/query.html
func (c *Client) Query(dir string, conf QueryConfig) ([]File, string, error) {
	var s struct {
		baseResponse
		Clock           string
		Files           []File
		IsFreshInstance bool `json:"is_fresh_instance"`
	}

	if err := c.send(&s, "query", dir, conf); err != nil {
		return nil, "", err
	}

	return s.Files, s.Clock, nil
}

// https://facebook.github.io/watchman/docs/cmd/shutdown-server.html
func (c *Client) ShutdownServer() (bool, error) {
	var v struct {
		baseResponse
		ShutdownServer bool `json:"shutdown-server"`
	}

	if err := c.send(&v, "shutdown-server"); err != nil {
		return false, err
	}

	return v.ShutdownServer, nil
}

// https://facebook.github.io/watchman/docs/cmd/since.html
func (c *Client) Since(dir string, clock string, patterns ...string) ([]File, string, error) {
	var s struct {
		baseResponse
		Clock string
		Files []File
	}

	params := []interface{}{"since", clock, dir}
	for _, p := range patterns {
		params = append(params, p)
	}

	if err := c.send(&s, params...); err != nil {
		return nil, "", err
	}

	return s.Files, s.Clock, nil
}

// https://facebook.github.io/watchman/docs/cmd/subscribe.html
func (c *Client) Subscribe() {

}

// https://facebook.github.io/watchman/docs/cmd/trigger.html
func (c *Client) Trigger() {

}

// https://facebook.github.io/watchman/docs/cmd/trigger-del.html
func (c *Client) TriggerDel() {

}

// https://facebook.github.io/watchman/docs/cmd/trigger-list.html
func (c *Client) TriggerList() {

}

// https://facebook.github.io/watchman/docs/cmd/unsubscribe.html
func (c *Client) Unsubscribe() {

}

// https://facebook.github.io/watchman/docs/cmd/version.html
func (c *Client) Version() (string, error) {
	var v baseResponse

	if err := c.send(&v, "version"); err != nil {
		return "", err
	}

	return v.Version, nil
}

// https://facebook.github.io/watchman/docs/cmd/watch.html
func (c *Client) Watch(dir string) error {
	var v struct {
		baseResponse
		Watch string
	}

	if err := c.send(&v, "watch", dir); err != nil {
		return err
	}

	return nil
}

// https://facebook.github.io/watchman/docs/cmd/watch-del.html
func (c *Client) WatchDel(dir string) error {
	var v struct {
		baseResponse
		Watch string
	}

	if err := c.send(&v, "watch-del", dir); err != nil {
		return err
	}

	return nil
}

// https://facebook.github.io/watchman/docs/cmd/watch-del-all.html
func (c *Client) WatchDelAll() {
}

// https://facebook.github.io/watchman/docs/cmd/watch-list.html
func (c *Client) WatchList() ([]string, error) {
	var v struct {
		baseResponse
		Roots []string
	}

	if err := c.send(&v, "watch-list"); err != nil {
		return nil, err
	}

	return v.Roots, nil
}

// https://facebook.github.io/watchman/docs/cmd/watch-project.html
func (c *Client) WatchProject() {

}
