package github_test

import (
	"context"
	"fmt"
	"log"

	"github.com/johejo/go-github-rest-api/github"
)

func Example() {
	c, err := github.NewClientWithResponses("https://api.github.com")
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	resp, err := c.ReposgetWithResponse(ctx, "golang", "go")
	if err != nil {
		log.Fatal(err)
	}
	// Output: golang/go
	fmt.Println(resp.JSON200.FullName)
}
