package request

type RequestOptions struct {
	BaseURL   string
	Token     string
	Model     string
	Provider  string
	Transport Transport
}
