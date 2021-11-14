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
