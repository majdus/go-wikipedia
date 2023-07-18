package wikipedia

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scottzhlin/go-wikipedia/internal/testhelper"
)

func checkQuery(u url.Values, k, v string) bool {
	return u.Get(k) == v
}

func TestClient_GetPage(t *testing.T) {
	ts := testhelper.NewTestHTTPServer()
	ts.RegisterHandler("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		if !checkQuery(r.Form, "action", "query") {
			http.Error(w, "invalid action", http.StatusBadRequest)
		}
		if !checkQuery(r.Form, "prop", "info|pageprops") {
			http.Error(w, "invalid prop", http.StatusBadRequest)
		}
		if !checkQuery(r.Form, "inprop", "url") {
			http.Error(w, "invalid inprop", http.StatusBadRequest)
		}
		if !checkQuery(r.Form, "ppprop", "disambiguation") {
			http.Error(w, "invalid ppprop", http.StatusBadRequest)
		}
		if !checkQuery(r.Form, "format", "json") {
			http.Error(w, "invalid format", http.StatusBadRequest)
		}
		if !checkQuery(r.Form, "pageids", "534366") {
			http.Error(w, "invalid pageids", http.StatusBadRequest)
		}
		fmt.Fprintf(w, `
{
    "batchcomplete": "",
    "query": {
        "pages": {
            "534366": {
                "pageid": 534366,
                "ns": 0,
                "title": "Barack Obama",
                "contentmodel": "wikitext",
                "pagelanguage": "en",
                "pagelanguagehtmlcode": "en",
                "pagelanguagedir": "ltr",
                "touched": "2023-07-18T16:26:29Z",
                "lastrevid": 1165884406,
                "length": 346245,
                "fullurl": "https://en.wikipedia.org/wiki/Barack_Obama",
                "editurl": "https://en.wikipedia.org/w/index.php?title=Barack_Obama&action=edit",
                "canonicalurl": "https://en.wikipedia.org/wiki/Barack_Obama"
            }
        }
    }
}
`)
	})

	ts.Start()
	defer ts.Stop()

	c, err := NewClient()
	require.NoError(t, err)

	c.url = ts.URL() + "/w/api.php"
	got, err := c.GetPage(context.TODO(), 534366)
	require.NoError(t, err)
	require.Equal(
		t,
		&Page{
			PageID: 534366,
			Title:  "Barack Obama",
			URL:    "https://en.wikipedia.org/wiki/Barack_Obama",
		},
		got,
	)
}

func TestClient_GetPageByTitle(t *testing.T) {
	ts := testhelper.NewTestHTTPServer()
	ts.RegisterHandler("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		if !checkQuery(r.Form, "action", "query") {
			http.Error(w, "invalid action", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "prop", "info|pageprops") {
			http.Error(w, "invalid prop", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "inprop", "url") {
			http.Error(w, "invalid inprop", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "ppprop", "disambiguation") {
			http.Error(w, "invalid ppprop", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "format", "json") {
			http.Error(w, "invalid format", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "titles", "Barack Obama") {
			http.Error(w, "invalid pageids", http.StatusBadRequest)
			return
		}
		fmt.Fprintf(w, `
{
    "batchcomplete": "",
    "query": {
        "pages": {
            "534366": {
                "pageid": 534366,
                "ns": 0,
                "title": "Barack Obama",
                "contentmodel": "wikitext",
                "pagelanguage": "en",
                "pagelanguagehtmlcode": "en",
                "pagelanguagedir": "ltr",
                "touched": "2023-07-18T16:26:29Z",
                "lastrevid": 1165884406,
                "length": 346245,
                "fullurl": "https://en.wikipedia.org/wiki/Barack_Obama",
                "editurl": "https://en.wikipedia.org/w/index.php?title=Barack_Obama&action=edit",
                "canonicalurl": "https://en.wikipedia.org/wiki/Barack_Obama"
            }
        }
    }
}
`)
	})

	ts.Start()
	defer ts.Stop()

	c, err := NewClient()
	require.NoError(t, err)

	c.url = ts.URL() + "/w/api.php"
	got, err := c.GetPageByTitle(context.TODO(), "Barack Obama")
	require.NoError(t, err)
	require.Equal(
		t,
		&Page{
			PageID: 534366,
			Title:  "Barack Obama",
			URL:    "https://en.wikipedia.org/wiki/Barack_Obama",
		},
		got,
	)
}

func TestClient_GetPageContent(t *testing.T) {
	ts := testhelper.NewTestHTTPServer()
	ts.RegisterHandler("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}
		if !checkQuery(r.Form, "action", "query") {
			http.Error(w, "invalid action", http.StatusBadRequest)
			return
		}
		if !checkQuery(r.Form, "format", "json") {
			http.Error(w, "invalid format", http.StatusBadRequest)
			return
		}
		if checkQuery(r.Form, "pageids", "534366") {
			if !checkQuery(r.Form, "inprop", "url") {
				http.Error(w, "invalid inprop", http.StatusBadRequest)
				return
			}
			if !checkQuery(r.Form, "ppprop", "disambiguation") {
				http.Error(w, "invalid ppprop", http.StatusBadRequest)
				return
			}
			if !checkQuery(r.Form, "prop", "info|pageprops") {
				http.Error(w, "invalid prop", http.StatusBadRequest)
				return
			}

			fmt.Fprintf(w, `
{
    "batchcomplete": "",
    "query": {
        "pages": {
            "534366": {
                "pageid": 534366,
                "ns": 0,
                "title": "Barack Obama",
                "contentmodel": "wikitext",
                "pagelanguage": "en",
                "pagelanguagehtmlcode": "en",
                "pagelanguagedir": "ltr",
                "touched": "2023-07-18T16:26:29Z",
                "lastrevid": 1165884406,
                "length": 346245,
                "fullurl": "https://en.wikipedia.org/wiki/Barack_Obama",
                "editurl": "https://en.wikipedia.org/w/index.php?title=Barack_Obama&action=edit",
                "canonicalurl": "https://en.wikipedia.org/wiki/Barack_Obama"
            }
        }
    }
}
`)
			return
		}
		if checkQuery(r.Form, "rvprop", "ids") {
			if !checkQuery(r.Form, "prop", "extracts|revisions") {
				http.Error(w, "invalid prop", http.StatusBadRequest)
				return
			}
			fmt.Fprintf(w, `
{
    "batchcomplete": "",
    "warnings": {
        "extracts": {
            "*": "HTML may be malformed and/or unbalanced and may omit inline images. Use at your own risk. "
        }
    },
    "query": {
        "pages": {
            "534366": {
                "pageid": 534366,
                "ns": 0,
                "title": "Barack Obama",
                "extract": "xxx",
                "revisions": [
                    {
                        "revid": 1165884406,
                        "parentid": 1165765677
                    }
                ]
            }
        }
    }
}
`)
			return
		}
		http.Error(w, "invalid request", http.StatusBadRequest)
	})

	ts.Start()
	defer ts.Stop()

	c, err := NewClient()
	require.NoError(t, err)

	c.url = ts.URL() + "/w/api.php"
	got, err := c.GetPageContent(context.TODO(), 534366)
	require.NoError(t, err)
	require.Equal(
		t,
		&PageContent{
			Page: &Page{
				PageID: 534366,
				Title:  "Barack Obama",
				URL:    "https://en.wikipedia.org/wiki/Barack_Obama",
			},
			Content:    "xxx",
			RevisionID: 1165884406,
			ParentID:   1165765677,
		},
		got,
	)
}
