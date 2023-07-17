package wikipedia

import (
	"context"
	"errors"

	"github.com/samber/lo"
)

// SearchOptions are the options for the Wikipedia search request.
type SearchOptions struct {
	Limit int // the max number of results returned, default 10
}

func defaultSearchOptions() *SearchOptions {
	return &SearchOptions{Limit: defaultLimit}
}

type searchRequest struct {
	Action   Action `url:"action"`
	List     string `url:"list"`
	SrProp   string `url:"srprop"`
	SrLimit  int    `url:"srlimit"`
	Limit    int    `url:"limit"`
	SrSearch string `url:"srsearch"`
}

// Search searches the Wikipedia for the given query.
func (c *Client) Search(
	ctx context.Context,
	query string,
	options ...func(searchOptions *SearchOptions),
) ([]string, error) {
	if len(query) == 0 {
		return nil, errors.New("go-wikipedia: query is empty")
	}

	o := defaultSearchOptions()
	for _, option := range options {
		option(o)
	}

	req := &searchRequest{
		Action:   ActionQuery,
		List:     "search",
		SrLimit:  o.Limit,
		Limit:    o.Limit,
		SrSearch: query,
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	return lo.Map(resp.Query.Search, func(item *searchResponse, _ int) string { return item.Title }), nil
}
