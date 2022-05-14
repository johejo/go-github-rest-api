# go-github-rest-api

Auto Generated GitHub's v3 REST API Client for Go

## Description

This package was generated from [GitHub's REST API OpenAPI Description](https://github.com/github/rest-api-description) using [deepmap/oapi-codegen](https://github.com/deepmap/oapi-codegen).

This project aims to ensure that keep up with the latest API updates and keep the Go Client available.

## Example

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/johejo/go-github-rest-api/github"
)

func main() {
	c, err := github.NewClientWithResponses("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resp, err := c.ReposgetWithResponse(ctx, "golang", "go")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(resp.JSON200.FullName)
}
```

## Other

[google/go-github](https://github.com/google/go-github) is well known Go client library for accessing the GitHub API v3.
