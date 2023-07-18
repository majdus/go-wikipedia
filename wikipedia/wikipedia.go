package wikipedia

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/samber/lo"

	"github.com/google/go-querystring/query"
)

// Action is the action of the Wikipedia API request.
type Action string

const (
	ActionQuery Action = "query" // query action
)

const (
	urlWithPlaceholder = "http://%s.wikipedia.org/w/api.php"
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

type requestError struct {
	Code string `json:"code"`
	Info string `json:"info"`
}

type warnings struct {
	Main map[string]any `json:"main"`
}

type searchInfo struct {
	TotalHits int `json:"totalhits"`
}

type revision struct {
	RevID    int    `json:"revid"`
	ParentID int    `json:"parentid"`
	Star     string `json:"*"`
}

type innerPage struct {
	Ns                  int                 `json:"ns"`
	Title               string              `json:"title"`
	PageID              int                 `json:"pageid"`
	ContentModel        string              `json:"contentmodel"`
	PageLanguage        string              `json:"pagelanguage"`
	PageLanguageTmlCode string              `json:"pagelanguagetmlcode"`
	PageLanguageDir     string              `json:"pagelanguagedir"`
	Touched             string              `json:"touched"`
	LastRevid           int                 `json:"lastrevid"`
	Length              int                 `json:"length"`
	FullURL             string              `json:"fullurl"`
	EditURL             string              `json:"editurl"`
	CanonicalURL        string              `json:"canonicalurl"`
	PageProps           map[string]string   `json:"pageprops"`
	Missing             string              `json:"missing"`
	Extract             string              `json:"extract"`
	Revisions           []revision          `json:"revisions"`
	Extlink             []map[string]string `json:"extlinks"`
	Link                []map[string]any    `json:"links"`
	Category            []map[string]any    `json:"categories"`
	ImageInfo           []map[string]string `json:"imageinfo"`
	Coordinate          []map[string]any    `json:"coordinates"`
}

type normalize struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type responseQuery struct {
	Search     []*searchResponse    `json:"search"`
	SearchInfo searchInfo           `json:"searchinfo"`
	Pages      map[string]innerPage `json:"pages"`
	Redirect   []normalize          `json:"redirects"`
	Normalize  []normalize          `json:"normalized"`
}

func (rq *responseQuery) oneOfPage() (innerPage, error) {
	pageIDs := lo.Keys(rq.Pages)
	if len(pageIDs) == 0 {
		return innerPage{}, errors.New("go-wikipedia: no pages found")
	}

	return rq.Pages[pageIDs[0]], nil
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

type responseContinue struct {
	Sroffset int    `json:"sroffset"`
	Continue string `json:"continue"`
}

type apiResult struct {
	Error         requestError     `json:"error"`
	Warnings      warnings         `json:"warnings"`
	BatchComplete string           `json:"batchcomplete"`
	Continue      responseContinue `json:"continue"`
	Query         responseQuery    `json:"query"`
}

func (c *Client) do(ctx context.Context, v any) (*apiResult, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.url, nil)
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: failed to create http request: %w", err)
	}

	req.Header.Set("User-Agent", c.o.userAgent)

	q, err := query.Values(v)
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

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("go-wikipedia: http request: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("go-wikipedia: read response body: %w", err)
	}

	res := new(apiResult)
	if err := json.Unmarshal(body, res); err != nil {
		return nil, fmt.Errorf("go-wikipedia: unmarshal response body: %w", err)
	}

	if len(res.Error.Code) > 0 {
		return nil, fmt.Errorf("go-wikipedia: returns error, code: %s, info: %s", res.Error.Code, res.Error.Info)
	}

	return res, nil
}
