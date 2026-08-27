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
func (s tableQuestionAnsweringService) answer(
	req TableQuestionAnsweringRequest,
	opts ...Option,
) ([]TableQuestionAnswering, error) {
	return doModelInference[TableQuestionAnsweringRequest, []TableQuestionAnswering](
		s.opts.With(opts...),
		"table question answering",
		req,
	)
}
