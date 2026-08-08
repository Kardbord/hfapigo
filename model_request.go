package hfgo

import (
	"net/http"

	"github.com/Kardbord/hfgo/v4/internal/request"
)

// doModelInference posts req to the inference model endpoint for the model
// configured in opts, returning the typed response. It returns a configuration
// SDK error when no model is set. task names the inference task in that error
// message, e.g. "fill mask" or "summarization".
func doModelInference[Req, Resp any](opts request.Options, task string, req Req) (Resp, error) {
	var zero Resp

	if opts.Model == "" {
		return zero, &SDKError{
			Kind:    SDKErrorKindConfiguration,
			Message: "the model option must be set for " + task + " to succeed",
			Err:     nil,
		}
	}

	return request.DoJSON[Req, Resp](
		opts,
		http.MethodPost,
		"hf-inference/models/"+opts.Model,
		req,
	)
}
