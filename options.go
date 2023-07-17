package wikipedia

// Option is a functional option for the Client.
type Option interface {
	apply(*options)
}

type options struct {
	language  string
	userAgent string
}

func defaultOptions() *options {
	return &options{
		language:  "en",
		userAgent: "wikipedia (https://github.com/scottzhlin/go-wikipedia/)",
	}
}

type functionalOption struct {
	f func(*options)
}

func newFunctionalOption(f func(*options)) *functionalOption {
	return &functionalOption{
		f: f,
	}
}

func (fo functionalOption) apply(o *options) {
	fo.f(o)
}

// WithLanguage sets the language code of the API being requested.
// It sets url domain prefix to one of the two letter prefixes found on the [All wikipedias ordered by number of articles](http://meta.wikimedia.org/wiki/List_of_Wikipedias).
// For example, `en` for English (https://en.wikipedia.org), `zh` for Chinese (https://zh.wikipedia.org).
func WithLanguage(language string) Option {
	return newFunctionalOption(func(o *options) { o.language = language })
}

// WithUserAgent sets the HTTP header User-Agent for all requests to access wikipedia API.
func WithUserAgent(userAgent string) Option {
	return newFunctionalOption(func(o *options) { o.userAgent = userAgent })
}
