package request

import "context"

// NormalizeContext returns ctx when non-nil, otherwise context.Background.
func NormalizeContext(ctx context.Context) context.Context {
	if ctx != nil {
		return ctx
	}

	return context.Background()
}
