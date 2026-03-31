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
	// TextClassificationFuncSigmoid applies a sigmoid to each score independently.
	// Useful for multi-label classification tasks, where multiple classes may apply simultaneously.
	TextClassificationFuncSigmoid = "sigmoid"
	// TextClassificationFuncSoftmax normalizes scores into a probability distribution summing to 1.
	// Useful for single-label multi-class classification tasks, where exactly one class applies.
	TextClassificationFuncSoftmax = "softmax"
	// TextClassificationFuncNone returns raw scores without any transformation applied.
	TextClassificationFuncNone = "none"
)

// TextClassification represents a text classification output.
type TextClassification struct {
	// The predicted class label.
	Label string `json:"label"`

	// The corresponding probability.
	Score float64 `json:"score"`
}
