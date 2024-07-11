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
	shouldWarn := hfapigo.APIKey() == ""
	if shouldWarn {
		fmt.Printf("%s not found in env, tests may fail due to rate limiting.\n", HuggingFaceTokenEnv)
	}
	m.Run()
	if shouldWarn {
		fmt.Printf("%s not found in env, tests may fail due to rate limiting.\n", HuggingFaceTokenEnv)
	}
}
