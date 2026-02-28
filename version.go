package hfapigo

import (
	"github.com/Kardbord/hfapigo/v4/internal/sdkversion"
)

// Version is the current version of the hfapigo SDK.
// This follows semantic versioning (semver.org).
const Version = sdkversion.Version

// UserAgent returns the User-Agent string used for HTTP requests.
// The format is: hfapigo/<version> (Go).
func UserAgent() string {
	return sdkversion.UserAgent()
}
