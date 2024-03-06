package config

type TargetBuildConfig struct {
	// Enable data race detection (supported only on amd64)
	Race *bool `json:"race,omitempty"`
	// List of build tags to consider satisfied during the build
	Tags string `json:"tags,omitempty"`
	// Arguments to pass on each go tool link invocation
	LdFlags string `json:"ldflags,omitempty"`
	// Indicates which kind of object file to build
	Mode string `json:"mode,omitempty"`
	// Whether to stamp binaries with version control information
	VCS string `json:"vcs,omitempty"`
	// Remove all file system paths from the resulting executable
	TrimPath *bool `json:"trimpath,omitempty"`
	// CGO dependency configure arguments
	CrossArgs string `json:"cross_args,omitempty"`
}

type TargetConfig struct {
	TargetBuildConfig TargetBuildConfig
	OutBinPath        string
}

type XGoConfig struct {
	// Go release to use for cross compilation (flag: go)
	GoVersion string `json:"go_version,omitempty"`
	// Set a Global Proxy for Go Modules
	GoProxy string `json:"go_proxy,omitempty"`
	Verbose *bool  `json:"verbose,omitempty"`
}

type Config struct {
	XGoConfig XGoConfig
	// Root of the project (where go.mod is located)
	RootPath   string
	BinTempDir string
	// Relative path to the main package (from the root)
	MainPkgRelPath string
	// List of targets to build
	Targets map[Target]TargetConfig
}
