package wikipedia

import (
	"context"
	"fmt"
	"strconv"

	"github.com/anaskhan96/soup"
	"github.com/samber/lo"
)

// GetPageOption is the option func for the Wikipedia page request.
type GetPageOption func(*GetPageOptions)

// GetPageOptions are the options for the Wikipedia page request.
type GetPageOptions struct {
	Redirects bool
}

// WithGetPageRedirects sets the redirects option for the Wikipedia page request.
func WithGetPageRedirects(r bool) GetPageOption {
	return func(o *GetPageOptions) {
		o.Redirects = r
	}
}

func defaultGetPageOptions() *GetPageOptions {
	return &GetPageOptions{Redirects: false}
}

type pageRequest struct {
	Action  Action   `url:"action" json:"action"`
	PageIDs []int    `url:"pageids" del:"|" json:"pageids"`
	Titles  []string `url:"titles" del:"|" json:"titles"`
	Props   []string `url:"prop" del:"|" json:"prop"`
	InProp  string   `url:"inprop" json:"inprop"`
	PpProp  string   `url:"ppprop" json:"ppprop"`
	Format  string   `url:"format"`
}

// GetPage returns a wikipedia page from the wikipedia API endpoint by given page id.
func (c *Client) GetPage(ctx context.Context, id int, opts ...GetPageOption) (*Page, error) {
	r := &pageRequest{
		Action:  ActionQuery,
		PageIDs: []int{id},
		Props:   []string{"info", "pageprops"},
		InProp:  "url",
		PpProp:  "disambiguation",
		Format:  "json",
	}
	return c.page(ctx, r, opts...)
}

// GetPageByTitle returns a wikipedia page from the wikipedia API endpoint by given page title.
func (c *Client) GetPageByTitle(ctx context.Context, title string, opts ...GetPageOption) (*Page, error) {
	r := &pageRequest{
		Action: ActionQuery,
		Titles: []string{title},
		Props:  []string{"info", "pageprops"},
		InProp: "url",
		PpProp: "disambiguation",
		Format: "json",
	}
	return c.page(ctx, r, opts...)
}

func (c *Client) page(ctx context.Context, request *pageRequest, opts ...GetPageOption) (*Page, error) {
	o := defaultGetPageOptions()
	for _, opt := range opts {
		opt(o)
	}

	response, err := c.do(ctx, request)
	if err != nil {
		return nil, err
	}

	if len(response.Query.Pages) == 0 {
		return nil, fmt.Errorf("go-wikipedia: page not found")
	}

	page, err := response.Query.oneOfPage()
	if err != nil {
		return nil, err
	}

	if len(page.Missing) > 0 {
		return nil, fmt.Errorf("go-wikipedia: page not found")
	}

	if len(response.Query.Redirect) > 0 {
		if !o.Redirects {
			return nil, fmt.Errorf("go-wikipedia: page is a redirect, set Redirects option to true to follow redirects")
		}
		return c.redirect(ctx, page.Title, response.Query)
	}

	if _, ok := page.PageProps["disambiguation"]; ok {
		return c.disambiguate(ctx, &page)
	}

	var rev revision
	if len(page.Revisions) > 0 {
		rev = page.Revisions[0]
	}

	return &Page{
		PageID:     page.PageID,
		Title:      page.Title,
		URL:        page.FullURL,
		RevisionID: rev.RevID,
		ParentID:   rev.ParentID,
	}, nil
}

type disambiguationRequest struct {
	Action  Action   `url:"action" json:"action"`
	Props   []string `url:"prop" del:"|" json:"prop"`
	Titles  string   `url:"titles" json:"titles"`
	RvProp  string   `url:"rvprop" json:"rvprop"`
	RvLimit int      `url:"rvlimit" json:"rvlimit"`
	Format  string   `url:"format"`
}

func (c *Client) disambiguate(ctx context.Context, page *innerPage) (*Page, error) {
	r := &disambiguationRequest{
		Action:  ActionQuery,
		Titles:  page.Title,
		Props:   []string{"revisions"},
		RvProp:  "content",
		RvLimit: 1,
		Format:  "json",
	}

	response, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}

	var html string
	if v, ok := response.Query.Pages[strconv.Itoa(page.PageID)]; ok {
		if len(v.Revisions) > 0 {
			html = v.Revisions[0].Star
		}
	}
	if len(html) == 0 {
		return nil, fmt.Errorf("go-wikipedia: disambiguation page not found")
	}

	doc := soup.HTMLParse(html)
	var d []string
	for _, link := range doc.FindAll("li") {
		for _, l := range link.FindAll("a") {
			if ref, ok := l.Attrs()["title"]; ok {
				if len(ref) >= 1 && !lo.Contains(d, ref) {
					d = append(d, ref)
				}
			}
		}
	}

	var rev revision
	if len(page.Revisions) > 0 {
		rev = page.Revisions[0]
	}
	return &Page{
		PageID:         page.PageID,
		Title:          page.Title,
		URL:            page.FullURL,
		RevisionID:     rev.RevID,
		ParentID:       rev.ParentID,
		Disambiguation: d,
	}, nil
}

func (c *Client) redirect(ctx context.Context, title string, rq responseQuery) (*Page, error) {
	var (
		t = title
		r = rq.Redirect[0]
	)
	if len(rq.Normalize) > 0 {
		n := rq.Normalize[0]
		if n.From != title {
			return nil, fmt.Errorf("go-wikipedia: unexpected normalize response")
		}
		t = n.To
	}

	if r.From == t {
		return nil, fmt.Errorf("go-wikipedia: unexpected redirect response")
	}

	return c.GetPageByTitle(ctx, r.To, WithGetPageRedirects(true))
}

type pageContentRequest struct {
	Action Action   `url:"action" json:"action"`
	Props  []string `url:"prop" del:"|" json:"prop"`
	RvProp string   `url:"rvprop" json:"rvprop"`
	Titles string   `url:"titles" json:"titles"`
	Format string   `url:"format"`
}

// GetPageContent returns a wikipedia page content from the wikipedia API endpoint by given page id.
func (c *Client) GetPageContent(ctx context.Context, id int, opts ...GetPageOption) (*PageContent, error) {
	p, err := c.GetPage(ctx, id, opts...)
	if err != nil {
		return nil, err
	}

	return c.pageContent(ctx, p)
}

// GetPageContentByTitle returns a wikipedia page content from the wikipedia API endpoint by given page title.
func (c *Client) GetPageContentByTitle(ctx context.Context, title string, opts ...GetPageOption) (*PageContent, error) {
	p, err := c.GetPageByTitle(ctx, title, opts...)
	if err != nil {
		return nil, err
	}

	return c.pageContent(ctx, p)
}

func (c *Client) pageContent(ctx context.Context, p *Page) (*PageContent, error) {
	r := &pageContentRequest{
		Action: ActionQuery,
		Props:  []string{"extracts", "revisions"},
		RvProp: "ids",
		Titles: p.Title,
		Format: "json",
	}
	response, err := c.do(ctx, r)
	if err != nil {
		return nil, err
	}

	pv, ok := response.Query.Pages[strconv.Itoa(p.PageID)]
	if !ok {
		return nil, fmt.Errorf("go-wikipedia: page not found: %d", p.PageID)
	}

	return &PageContent{
		Page:       p,
		Content:    pv.Extract,
		RevisionID: lo.Ternary(len(pv.Revisions) > 0, pv.Revisions[0].RevID, 0),
		ParentID:   lo.Ternary(len(pv.Revisions) > 0, pv.Revisions[0].ParentID, 0),
	}, nil
}

// PageContent represents a wikipedia page content.
type PageContent struct {
	Page       *Page
	Content    string
	RevisionID int
	ParentID   int
}

// Page represents a wikipedia page info.
type Page struct {
	PageID         int              `json:"pageid"`
	Title          string           `json:"title"`
	HTML           string           `json:"html"`
	URL            string           `json:"fullurl"`
	RevisionID     int              `json:"revid"`
	ParentID       int              `json:"parentid"`
	Summary        string           `json:"summary"`
	CheckedImage   bool             `json:"checkedimage"`
	Images         []string         `json:"images"`
	Coordinate     []float64        `json:"coordinates"`
	Reference      []string         `json:"references"`
	Link           []string         `json:"links"`
	Category       []string         `json:"categories"`
	Section        []string         `json:"sections"`
	SectionOffset  map[string][]int `json:"sectionoffset"`
	Disambiguation []string         `json:"disambiguation"`
}
