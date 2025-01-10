[![Test Status](https://github.com/majdus/go-wikipedia/workflows/pre-commit/badge.svg)](https://github.com/majdus/go-wikipedia/actions?query=workflow%3Apre-commit)
[![Latest release](https://img.shields.io/github/release/majdus/go-wikipedia.svg)](https://github.com/majdus/go-wikipedia/releases)
[![GoDoc](https://godoc.org/github.com/majdus/go-wikipedia?status.svg)](https://godoc.org/github.com/majdus/go-wikipedia)

# go-wikipedia
go-wikipedia is a Go client library for accessing the [Wikipedia API](https://en.wikipedia.org/api/rest_v1/#/).

## Requisites

- Go 1.18 or higher

## Installation

Using Go modules to install go-wikipedia:

```bash
go get github.com/majdus/go-wikipedia
```

Alternatively, you can import the package directly and then run `go mod tidy` to install the dependency:

```go
import "github.com/majdus/go-wikipedia/wikipedia"
```

## Usage

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/majdus/go-wikipedia/wikipedia"
)

func main() {
	// create a new Wikipedia client
	client, err := wikipedia.NewClient()
	if err != nil {
		panic(err)
	}

	titles, err := client.Search(context.TODO(), "golang")
	if err != nil {
		panic(err)
	}

	fmt.Println(titles)
}
```