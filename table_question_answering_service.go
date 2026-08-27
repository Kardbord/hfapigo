package hfgo

import (
	"github.com/Kardbord/hfgo/v4/internal/request"
)

// tableQuestionAnsweringService implements table question answering calls using the configured request options.
type tableQuestionAnsweringService struct {
	opts request.Options
}

// newTableQuestionAnsweringService builds a table question answering service with a snapshot of the provided options.
func newTableQuestionAnsweringService(opts request.Options) tableQuestionAnsweringService {
	return tableQuestionAnsweringService{opts: opts}
}

// answer sends a table question answering request and returns the answers.
//
// NOTE: The HuggingFace API returns a bare JSON object for table question
// answering, not an array — despite the upstream schema declaring an array
// response. We unmarshal into a single TableQuestionAnswering and wrap it
// in a slice so the caller sees a uniform []TableQuestionAnswering return.
func (s tableQuestionAnsweringService) answer(
	req TableQuestionAnsweringRequest,
	opts ...Option,
) ([]TableQuestionAnswering, error) {
	single, err := doModelInference[TableQuestionAnsweringRequest, TableQuestionAnswering](
		s.opts.With(opts...),
		"table question answering",
		req,
	)
	if err != nil {
		return nil, err
	}

	return []TableQuestionAnswering{single}, nil
}
