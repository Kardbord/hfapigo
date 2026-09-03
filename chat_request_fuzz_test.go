//go:build !integration

package hfgo

import (
	"encoding/json"
	"testing"
)

func FuzzChatRequestValidate(f *testing.F) {
	f.Add([]byte(`{"model":"test","messages":[{"role":"user","content":"hi"}]}`))
	f.Add([]byte(``))
	f.Add([]byte(`{}`))
	f.Add([]byte(`{"model":"test"}`))
	f.Add([]byte(`{"messages":[{"role":"user","content":"hi"}]}`))
	f.Fuzz(func(_ *testing.T, data []byte) {
		var req ChatRequest
		if err := json.Unmarshal(data, &req); err != nil {
			return
		}
		_ = req.validate()
	})
}

func FuzzChatMessageValidate(f *testing.F) {
	f.Add([]byte(`{"role":"user","content":"hi"}`))
	f.Add([]byte(``))
	f.Add([]byte(`{}`))
	f.Add([]byte(
		`{"role":"user","content":"hi","tool_calls":[` +
			`{"id":"1","type":"function","function":{"name":"test"}}]}`))
	f.Fuzz(func(_ *testing.T, data []byte) {
		var msg ChatMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			return
		}
		_ = msg.validate()
	})
}
