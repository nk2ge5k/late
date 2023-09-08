package late

import "fmt"

// BuildInfo represents the information about Ruler build.
type BuildInfo struct {
	// Version is the current git tag with v prefix stripped
	Version string `json:"version"`
	// Commit is the current git commit SHA
	Commit string `json:"commit"`
	// Date is the build date in RFC3339
	Date string `json:"date"`
}

// String returns string representation of the build info.
func (bi BuildInfo) String() string {
	return fmt.Sprintf("%s (sha: %s) %s", bi.Version, bi.Commit, bi.Date)
}

// Global variable to access build info of the binary.
var buildInfo BuildInfo

// SetBuildInfo sets the build information for the binary.
func SetBuildInfo(version, commit, date string) {
	buildInfo.Version = version
	buildInfo.Commit = commit
	buildInfo.Date = date
}

// GetBuildInfo returns the current build information for the binary.
func GetBuildInfo() BuildInfo {
	return buildInfo
}
