package main

type Config struct {
	Settle               int
	RootRestrictFiles    []string `json:"root_restrict_files"`
	IllegalFSTypes       []string `json:"illegal_fstypes"`
	IllegalFSTypesAdvice string   `json:"illegal_fstypes_advice"`
	IgnoreVCS            []string `json:"ignore_vcs"`
	IgnoreDirs           []string `json:"ignore_dirs"`
	GCAgeSeconds         int      `json:"gc_age_seconds"`
	GCIntervalSeconds    int      `json:"gc_interval_seconds"`
}
