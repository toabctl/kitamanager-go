package version

// These variables are set at build time via -ldflags.
var (
	// GitVersion is the version string from git describe (e.g. "v0.0.1" or "v0.0.1-3-g1a2b3c4").
	GitVersion = "dev"

	// GitCommit is the git commit hash.
	GitCommit = "unknown"

	// BuildTime is the time the binary was built.
	BuildTime = "unknown"
)

// Version returns the version string from git describe,
// falling back to the short commit hash if not set.
func Version() string {
	if GitVersion != "" && GitVersion != "dev" {
		return GitVersion
	}
	if len(GitCommit) >= 7 {
		return GitCommit[:7]
	}
	return GitCommit
}
