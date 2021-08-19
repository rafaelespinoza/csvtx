package version

// These are pieces of version metadata that can be set through -ldflags.
var (
	BranchName string
	BuildTime  string
	CommitHash string
	GoOSArch   string
	GoVersion  string
)
