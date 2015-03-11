package watchman

type File struct {
	Name    string
	Exists  bool
	Cclock  string
	Oclock  string
	Mtime   int64
	MtimeMs int64
	MtimeUs int64
	MtimeNs int64
	MtimeF  float64
	Ctime   int64
	CtimeMs int64
	CtimeUs int64
	CtimeNs int64
	CtimeF  float64
	Size    int
	Mode    int
	Uid     int
	Gid     int
	Ino     int
	Dev     int
	Nlink   int
	New     bool
}

type base struct {
	Version string
	Err     string `json:"error"`
}

func (b *base) Error() string {
	return b.Err
}

func (c *Client) Clock(dir string) (string, error) {
	var s struct {
		base
		Clock string
	}

	if err := c.send(&s, "clock", dir); err != nil {
		return "", err
	}

	return s.Clock, nil
}

type Path struct {
	Path  string
	Depth int
}

type QueryConfig struct {
	Suffix               []string   `json:"suffix,omitempty"`
	Since                string     `json:"since,omitempty"`
	Expression           Expression `json:"expression,omitempty"`
	Fields               []string   `json:"fields,omitempty"`
	Path                 []Path     `jsonL"path,omitempty"`
	SyncTimeout          int        `json:"sync_timeout,omitempty"`
	EmptyOnFreshInstance bool       `json:"empty_on_fresh_instance,omitempty"`
}

func (c *Client) Query(dir string, conf QueryConfig) ([]File, string, error) {
	var s struct {
		base
		Clock           string
		Files           []File
		IsFreshInstance bool `json:"is_fresh_instance"`
	}

	if err := c.send(&s, "query", dir, conf); err != nil {
		return nil, "", err
	}

	return s.Files, s.Clock, nil
}

func (c *Client) Find(dir string, patterns ...string) ([]File, string, error) {
	var s struct {
		base
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

func (c *Client) Since(dir string, clock string, patterns ...string) ([]File, string, error) {
	var s struct {
		base
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

func (c *Client) GetConfig(dir string) (Config, error) {
	var s struct {
		base
		Config Config
	}

	if err := c.send(&s, "get-config", dir); err != nil {
		return Config{}, err
	}

	return s.Config, nil
}

func (c *Client) Log(level, msg string) (bool, error) {
	var s struct {
		base
		Logged bool
	}

	if err := c.send(&s, "log", level, msg); err != nil {
		return false, err
	}

	return s.Logged, nil
}

func (c *Client) Version() (string, error) {
	var v base

	if err := c.send(&v, "version"); err != nil {
		return "", err
	}

	return v.Version, nil
}

func (c *Client) GetSockname() (string, error) {
	var s struct {
		base
		Sockname string
	}

	if err := c.send(&s, "get-sockname"); err != nil {
		return "", err
	}

	return s.Sockname, nil
}

func (c *Client) ShutdownServer() (bool, error) {
	var v struct {
		base
		ShutdownServer bool `json:"shutdown-server"`
	}

	if err := c.send(&v, "shutdown-server"); err != nil {
		return false, err
	}

	return v.ShutdownServer, nil
}

func (c *Client) Watch(dir string) error {
	var v struct {
		base
		Watch string
	}

	if err := c.send(&v, "watch", dir); err != nil {
		return err
	}

	return nil
}

func (c *Client) WatchDel(dir string) error {
	var v struct {
		base
		Watch string
	}

	if err := c.send(&v, "watch-del", dir); err != nil {
		return err
	}

	return nil
}

func (c *Client) WatchLists() ([]string, error) {
	var v struct {
		base
		Roots []string
	}

	if err := c.send(&v, "watch-list"); err != nil {
		return nil, err
	}

	return v.Roots, nil
}
