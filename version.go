package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/version"
)

// Version is the current version of the hfapigo SDK.
// This follows semantic versioning (semver.org).
const Version = version.Version

// UserAgent returns the User-Agent string used for HTTP requests.
// The format is: hfapigo/<version> (Go).
func UserAgent() string {
	return version.UserAgent()
}
