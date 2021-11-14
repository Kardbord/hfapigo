package hfapigo

import (
	"encoding/json"
	"fmt"
	"strings"
)

type APIError struct {
	Errors []string `json:"error"`
}

func (e APIError) Error() string {
	bytes, err := json.Marshal(e)
	if err != nil {
		return fmt.Sprintf(`{error=["%v"]}`, strings.Join(e.Errors, `", "`))
	}
	return string(bytes)
}

// Attempts to unmarshal the response body into an APIError, and return it.
// If the unmarshal fails, returns the orig error.
func respBodyToAPIError(respBody []byte, orig error) error {
	apiErr := APIError{}
	err := json.Unmarshal(respBody, &apiErr)
	if err != nil {
		return orig
	}
	return apiErr
}
