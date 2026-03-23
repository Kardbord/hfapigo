package hfgo

import (
	"encoding/json"
	"errors"
)

// ChatFunctionCall represents a tool function call with arguments.
type ChatFunctionCall struct {
	// Required.
	Name string `json:"name"`
	// Required.
	Arguments   string  `json:"arguments"`
	Description *string `json:"description,omitempty"`
}

// MarshalJSON enforces required fields on ChatFunctionCall.
func (f ChatFunctionCall) MarshalJSON() ([]byte, error) {
	if err := f.validate(); err != nil {
		return nil, err
	}
	type alias ChatFunctionCall

	return json.Marshal(alias(f))
}

// UnmarshalJSON enforces required fields on ChatFunctionCall.
func (f *ChatFunctionCall) UnmarshalJSON(data []byte) error {
	if f == nil {
		return errors.New("chat function call: nil receiver")
	}
	type alias ChatFunctionCall
	var tmp alias
	if err := json.Unmarshal(data, &tmp); err != nil {
		return err
	}
	out := ChatFunctionCall(tmp)
	if err := out.validate(); err != nil {
		return err
	}
	*f = out

	return nil
}

func (f ChatFunctionCall) validate() error {
	if f.Name == "" {
		return &SDKError{
			Kind:    SDKErrorKindValidation,
			Message: "chat function call: name must be set",
			Err:     nil,
		}
	}
	if f.Arguments == "" {
		return &SDKError{
			Kind:    SDKErrorKindValidation,
			Message: "chat function call: arguments must be set",
			Err:     nil,
		}
	}

	return nil
}
