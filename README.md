# go-wikipedia
go-wikipedia is a Go client library for accessing the [Wikipedia API](https://en.wikipedia.org/api/rest_v1/#/).

## Requisites

- Go 1.18 or higher

## Installation

Using Go modules to install go-wikipedia:

```bash
go get github.com/scottzhlin/go-wikipedia
```

Alternatively, you can import the package directly and then run `go mod tidy` to install the dependency:

```go
import "github.com/scottzhlin/go-wikipedia"
```

## Usage

```go
package main

import "github.com/scottzhlin/go-wikipedia"

func main() {
    // Create a new Wikipedia client
    client := wikipedia.NewClient()

    // Get page content
    page, err := client.GetPage(wikipedia.PageOptions{
        Title: "Wikipedia",
    })
    if err != nil {
        panic(err)
    }

    // Print the page summary
    fmt.Println(page.Summary())
}
```