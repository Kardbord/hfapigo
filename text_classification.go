package hfgo

// TextClassificationInput represents either a single string to run text
// classification on, or a batch of strings to run classification on.
type TextClassificationInput []string

// NewTextClassificationInput is a helper function for creating single or
// batched TextClassificationInput types.
func NewTextClassificationInput(inputs ...string) TextClassificationInput {
	return inputs
}

// TextClassificationRequest represents a text classification
// inference request to the API.
type TextClassificationRequest struct {
	// The text to classify.
	// Required.
	Inputs TextClassificationInput `json:"inputs"`

	// Additional inference parameters for text classification
	Parameters *TextClassificationParameters `json:"parameters,omitempty"`
}

// TextClassificationParameters specify additional inference
// parameters for text classification.
type TextClassificationParameters struct {
	// Possible values: sigmoid, softmax, none.
	FunctionToApply *string `json:"function_to_apply,omitempty"`

	// When specified, limits the output to the top K most probable classes.
	TopK *int `json:"top_k,omitempty"`
}

const (
	// TextClassificationFuncSigmoid is useful for binary classification tasks.
	TextClassificationFuncSigmoid = "sigmoid"
	// TextClassificationFuncSoftmax is useful for multi-class classification tasks.
	TextClassificationFuncSoftmax = "softmax"
	// TextClassificationFuncNone is useful for returning raw scores without transformation.
	TextClassificationFuncNone = "none"
)

// TextClassification represents a text classification output.
type TextClassification struct {
	// The predicted class label.
	Label string `json:"label"`

	// The corresponding probability.
	Score float64 `json:"score"`
}
