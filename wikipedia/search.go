package wikipedia

import (
	"context"
	"errors"
)

// SearchOptions are the options for the Wikipedia search request.
type SearchOptions struct {
	SrLimit int // the max number of results returned, default 10
	Limit   int // the max number of results returned, default 10
}

func defaultSearchOptions() *SearchOptions {
	return &SearchOptions{
		SrLimit: defaultLimit,
		Limit:   defaultLimit,
	}
}

type SearchRequest struct {
	Action   Action `url:"action"`
	List     string `url:"list"`
	SrProp   string `url:"srprop"`
	SrLimit  int    `url:"srlimit"`
	Limit    int    `url:"limit"`
	SrSearch string `url:"srsearch"`
	Format   string `url:"format"`
}

// Search searches the Wikipedia for the given query.
func (c *Client) Search(
	ctx context.Context,
	query string,
	searchOptions *SearchOptions,
) ([]*SearchResponse, error) {
	if len(query) == 0 {
		return nil, errors.New("go-wikipedia: query is empty")
	}

	if searchOptions == nil {
		searchOptions = defaultSearchOptions()
	}

	req := &SearchRequest{
		Action:   ActionQuery,
		List:     "search",
		SrLimit:  searchOptions.SrLimit,
		Limit:    searchOptions.Limit,
		SrSearch: query,
		Format:   "json",
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Query.Search, nil
}
