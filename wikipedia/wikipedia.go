package wikipedia

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/google/go-querystring/query"
)

// Action is the action of the Wikipedia API request.
type Action string

const (
	ActionQuery Action = "query" // query action
)

const (
	urlWithPlaceholder = "https://%s.wikipedia.org/w/api.php"
	defaultLimit       = 10
)

// Client is a client for the Wikipedia API requests.
// Wikipedia API main page: https://www.mediawiki.org/wiki/API:Main_page
// Wikipedia API docs: https://en.wikipedia.org/api/rest_v1/
type Client struct {
	c *http.Client
	o *options

	url string
}

// NewClient returns a new instance of the Wikipedia client.
func NewClient(opts ...Option) (*Client, error) {
	o := defaultOptions()
	for _, opt := range opts {
		opt.apply(o)
	}

	return &Client{
		c:   new(http.Client),
		o:   o,
		url: fmt.Sprintf(urlWithPlaceholder, o.language),
	}, nil
}

type innerQuery struct {
	Format string `url:"format"`
	V      any
}

func newInnerQuery(v any) *innerQuery {
	return &innerQuery{
		Format: "json",
		V:      v,
	}
}

type requestError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

type responseQuery struct {
	Search []*searchResponse `json:"search"`
}

type searchResponse struct {
	Ns        int    `json:"ns"`
	Title     string `json:"title"`
	PageID    int    `json:"pageid"`
	Size      int    `json:"size"`
	WordCount int    `json:"wordcount"`
	Snippet   string `json:"snippet"`
	Timestamp string `json:"timestamp"`
}

type apiResult struct {
	Error    requestError   `json:"error"`
	Query    responseQuery  `json:"query"`
	Continue map[string]any `json:"continue"`
	Parse    map[string]any `json:"parse"`
}

func (c *Client) do(ctx context.Context, v any) (*apiResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: failed to create http request: %w", err)
	}

	req.Header.Set("User-Agent", c.o.userAgent)

	q, err := query.Values(newInnerQuery(v))
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: encode query parameters: %w", err)
	}

	uq := req.URL.Query()
	for k, v := range q {
		for _, vv := range v {
			uq.Add(k, vv)
		}
	}

	req.URL.RawQuery = uq.Encode()

	resp, err := c.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: http request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: read response body: %w", err)
	}

	res := new(apiResult)
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("go-wikipedia: unmarshal response body: %w", err)
	}

	return res, nil
}
