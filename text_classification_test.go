//go:build !integration

package hfgo

import (
	"testing"

	"github.com/Kardbord/hfgo/v4/internal/testutils"
	"github.com/stretchr/testify/require"
)

func TestTextClassificationRequestClone_Deep(t *testing.T) {
	t.Parallel()

	req := TextClassificationRequest{
		Input: "I feel great today!",
		Parameters: &TextClassificationParameters{
			FunctionToApply: testutils.Ptr(TextClassificationFuncSoftmax),
			TopK:            testutils.Ptr(2),
		},
	}

	cloned := req.Clone()

	*cloned.Parameters.FunctionToApply = TextClassificationFuncSigmoid
	*cloned.Parameters.TopK = 5
	cloned.Input = "changed"

	require.Equal(t, "I feel great today!", req.Input)
	require.Equal(t, TextClassificationFuncSoftmax, *req.Parameters.FunctionToApply)
	require.Equal(t, 2, *req.Parameters.TopK)
	require.Equal(t, TextClassificationFuncSigmoid, *cloned.Parameters.FunctionToApply)
	require.Equal(t, 5, *cloned.Parameters.TopK)
}

func TestTextClassificationBatchRequestClone_Deep(t *testing.T) {
	t.Parallel()

	req := TextClassificationBatchRequest{
		Inputs: []string{"text1", "text2"},
		Parameters: &TextClassificationParameters{
			FunctionToApply: testutils.Ptr(TextClassificationFuncNone),
		},
	}

	cloned := req.Clone()

	cloned.Inputs[0] = "changed"
	*cloned.Parameters.FunctionToApply = TextClassificationFuncSoftmax

	require.Equal(t, []string{"text1", "text2"}, req.Inputs)
	require.Equal(t, TextClassificationFuncNone, *req.Parameters.FunctionToApply)
	require.Equal(t, "changed", cloned.Inputs[0])
	require.Equal(t, TextClassificationFuncSoftmax, *cloned.Parameters.FunctionToApply)
}

func TestTextClassificationParametersClone_Deep(t *testing.T) {
	t.Parallel()

	params := &TextClassificationParameters{
		FunctionToApply: testutils.Ptr(TextClassificationFuncSoftmax),
		TopK:            testutils.Ptr(2),
	}

	cloned := params.Clone()

	*cloned.FunctionToApply = TextClassificationFuncNone
	*cloned.TopK = 3

	require.Equal(t, TextClassificationFuncSoftmax, *params.FunctionToApply)
	require.Equal(t, 2, *params.TopK)
	require.Equal(t, TextClassificationFuncNone, *cloned.FunctionToApply)
	require.Equal(t, 3, *cloned.TopK)
}

func TestTextClassificationClone_Nil(t *testing.T) {
	t.Parallel()

	var r *TextClassificationRequest
	require.Empty(t, r.Clone())

	var b *TextClassificationBatchRequest
	require.Empty(t, b.Clone())

	var p *TextClassificationParameters
	require.Empty(t, p.Clone())
}
