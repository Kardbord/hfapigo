package sdkversion

import "fmt"

// Version is the current version of the hfapigo SDK.
// This follows semantic versioning (semver.org).
const Version = "4.0.0"

// UserAgent returns the User-Agent string used for HTTP requests.
// The format is: hfapigo/<version> (Go).
func UserAgent() string {
	return fmt.Sprintf("hfapigo/%s (Go)", Version)
}
