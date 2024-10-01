package hfapigo_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/Kardbord/hfapigo/v3"
)

const HuggingFaceTokenEnv = "HUGGING_FACE_TOKEN"

const TestFilesDir = "./test_files"

func init() {
	key := os.Getenv(HuggingFaceTokenEnv)
	if key != "" {
		hfapigo.SetAPIKey(key)
	}
}

func TestMain(m *testing.M) {
	if hfapigo.APIKey() == "" {
		fmt.Fprintf(os.Stderr, "%s not set, tests will fail due to rate limiting.", HuggingFaceTokenEnv)
		os.Exit(1)
	}
	m.Run()
}
