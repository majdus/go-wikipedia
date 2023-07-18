package wikipedia

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/scottzhlin/go-wikipedia/internal/testhelper"
)

func TestClient_Search(t *testing.T) {
	ts := testhelper.NewTestHTTPServer()
	ts.RegisterHandler("/w/api.php", func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			panic(err)
		}

		if !checkQuery(r.Form, "format", "json") {
			http.Error(w, "format must be json", http.StatusBadRequest)
			return
		}

		fmt.Fprint(w, `
{
    "warnings": {
        "main": {
            "*": "Unrecognized parameter: limit."
        }
    },
    "batchcomplete": "",
    "continue": {
        "sroffset": 10,
        "continue": "-||"
    },
    "query": {
        "searchinfo": {
            "totalhits": 22107
        },
        "search": [
            {
                "ns": 0,
                "title": "Barack Obama",
                "pageid": 534366
            },
            {
                "ns": 0,
                "title": "Barack Obama Sr.",
                "pageid": 16136849
            },
            {
                "ns": 0,
                "title": "Family of Barack Obama",
                "pageid": 17775180
            },
            {
                "ns": 0,
                "title": "Presidency of Barack Obama",
                "pageid": 20082093
            },
            {
                "ns": 0,
                "title": "Barack Obama Presidential Center",
                "pageid": 41828619
            },
            {
                "ns": 0,
                "title": "Barack Obama religion conspiracy theories",
                "pageid": 26472604
            },
            {
                "ns": 0,
                "title": "Barack Obama Plaza",
                "pageid": 57773094
            },
            {
                "ns": 0,
                "title": "Barack Obama citizenship conspiracy theories",
                "pageid": 20617631
            },
            {
                "ns": 0,
                "title": "Early life and career of Barack Obama",
                "pageid": 16394033
            },
            {
                "ns": 0,
                "title": "Cabinet of Barack Obama",
                "pageid": 21341288
            }
        ]
    }
}`)
	})

	ts.Start()
	defer ts.Stop()

	c, err := NewClient()
	require.NoError(t, err)

	c.url = ts.URL() + "/w/api.php"
	got, err := c.Search(context.TODO(), "Barack Obama", func(searchOptions *SearchOptions) { searchOptions.Limit = 10 })
	require.NoError(t, err)
	require.Equal(
		t,
		[]string{
			"Barack Obama",
			"Barack Obama Sr.",
			"Family of Barack Obama",
			"Presidency of Barack Obama",
			"Barack Obama Presidential Center",
			"Barack Obama religion conspiracy theories",
			"Barack Obama Plaza",
			"Barack Obama citizenship conspiracy theories",
			"Early life and career of Barack Obama",
			"Cabinet of Barack Obama",
		},
		got,
	)
}
