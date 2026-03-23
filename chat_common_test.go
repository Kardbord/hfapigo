package hfgo

import (
	"encoding/json"
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/hferrors"
	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

type toolCallDecodeCase struct {
	name        string
	data        []byte
	wantErr     bool
	wantErrKind hferrors.SDKErrorKind
}

type toolCallMarshalCase struct {
	name        string
	value       any
	wantErrKind hferrors.SDKErrorKind
}

func runToolCallDecodeTests(t *testing.T, cases []toolCallDecodeCase, decode func([]byte) error) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := decode(tc.data)
			if tc.wantErr {
				require.Error(t, err)
				testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)

				return
			}
			require.NoError(t, err)
		})
	}
}

func runToolCallMarshalTests(t *testing.T, cases []toolCallMarshalCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			require.Error(t, err)
			testutils.AssertSDKErrorKind(t, err, tc.wantErrKind)
		})
	}
}

func TestChatFunctionCall_Validation(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name  string
		value ChatFunctionCall
	}{
		{
			name:  "missing name",
			value: ChatFunctionCall{Arguments: "{}"},
		},
		{
			name:  "missing arguments",
			value: ChatFunctionCall{Name: "fn"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := json.Marshal(tc.value)
			require.Error(t, err)
			testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindValidation)
		})
	}

	var got ChatFunctionCall
	err := json.Unmarshal([]byte(`{"name":"fn","arguments":""}`), &got)
	require.Error(t, err)
	testutils.AssertSDKErrorKind(t, err, hferrors.SDKErrorKindValidation)
}
