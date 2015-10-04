package kovacs

var (
	StdinDevNull     StdinType = stdinString("/dev/null")
	StdinNamePerLine StdinType = stdinString("NAME_PER_LINE")
)

type StdinType interface {
	stdinNoop()
}

type stdinString string

func (s stdinString) stdinNoop() {}

type StdinArray []string

func (s StdinArray) stdinNoop() {}

type Config struct {
	Settle               int        `json:"settle"`
	RootRestrictFiles    []string   `json:"root_restrict_files"`
	RootFiles            [][]string `json:"root_files"`
	EnforceRootFiles     bool       `json:"enforce_root_files"`
	IllegalFSTypes       []string   `json:"illegal_fstypes"`
	IllegalFSTypesAdvice string     `json:"illegal_fstypes_advice"`
	IgnoreVCS            []string   `json:"ignore_vcs"`
	IgnoreDirs           []string   `json:"ignore_dirs"`
	GCAgeSeconds         int        `json:"gc_age_seconds"`
	GCIntervalSeconds    int        `json:"gc_interval_seconds"`
	FSEventsLatency      float64    `json:"fsevents_latency"`
	IdleReapAgeSeconds   int        `json:"idle_reap_age_seconds"`
}

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

type Path struct {
	Path  string
	Depth int
}

type QueryOptions struct {
	Suffix               []string   `json:"suffix,omitempty"`
	Since                string     `json:"since,omitempty"`
	Expression           Expression `json:"expression,omitempty"`
	Fields               []string   `json:"fields,omitempty"`
	Path                 []Path     `json:"path,omitempty"`
	SyncTimeout          int        `json:"sync_timeout,omitempty"`
	EmptyOnFreshInstance bool       `json:"empty_on_fresh_instance,omitempty"`
	RelativeRoot         string     `json:"relative_root,omitempty"`
}

type SubscriptionEvent struct {
	Version      string   `json:"version"`
	Clock        string   `json:"clock"`
	Files        []string `json:"files"`
	Root         string   `json:"root"`
	Subscription string   `json:"subscription"`
}

type SubscriptionOptions struct {
	Since    string
	Expr     Expression
	Fields   []string
	DeferVCS bool
}

type TriggerOptions struct {
	Name          string      `json:"name"`
	Command       []string    `json:"command"`
	AppendFiles   bool        `json:"append_files,omitempty"`
	Expression    interface{} `json:"expression"`
	Stdin         StdinType   `json:"stdin"`
	Stdout        string      `json:"stdout"`
	Stderr        string      `json:"stderr"`
	MaxFilesStdin int         `json:"max_files_stdin"`
	Chdir         string      `json:"chdir"`
	RelativeRoot  string      `json:"relative_root"`
}
