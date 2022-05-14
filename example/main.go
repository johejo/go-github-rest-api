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
