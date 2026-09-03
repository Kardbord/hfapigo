//go:build !integration

package request

import (
	"testing"
)

func FuzzParseSSEField(f *testing.F) {
	f.Add("data: hello")
	f.Fuzz(func(_ *testing.T, line string) {
		_, _ = parseSSEField(line)
	})
}
