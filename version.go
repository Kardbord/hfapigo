package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/sdkversion"
)

// Version is the current version of the hfgo SDK.
// This follows semantic versioning (semver.org).
const Version = sdkversion.Version

// UserAgent returns the User-Agent string used for HTTP requests.
// The format is: hfgo/<version> (Go).
func UserAgent() string {
	return sdkversion.UserAgent()
}
