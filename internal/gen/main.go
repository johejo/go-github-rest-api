package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	matrix := []struct {
		desc   string
		target string
		pkg    string
	}{
		{desc: "descriptions", target: "api.github.com", pkg: "github"},
	}

	ctx := context.Background()
	sha := getSha(ctx)
	if err := os.WriteFile("rest-api-description_sha.txt", []byte(sha), 0o644); err != nil {
		log.Fatal(err)
	}
	for _, m := range matrix {
		gen(ctx, sha, m.desc, m.target, m.pkg)
	}
}

func gen(ctx context.Context, sha string, desc string, target string, pkg string) {
	u := fmt.Sprintf("https://raw.githubusercontent.com/github/rest-api-description/%s/%s/%s/%s.yaml", sha, desc, target, target)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	f, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(f.Name())

	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	tmpl := filepath.Join(wd, "internal", "gen", "templates")

	types := filepath.Join(wd, "github", "types.go")
	if err := os.RemoveAll(types); err != nil {
		log.Fatal(err)
	}

	if err := execRunWithStderr(ctx,
		//"oapi-codegen",
		"go", "run", "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest",
		"-package", "github",
		"-generate", "types",
		"-templates", tmpl,
		"-o", types,
		f.Name(),
	); err != nil {
		log.Fatal(err)
	}

	client := filepath.Join(wd, "github", "client.go")
	if err := os.RemoveAll(client); err != nil {
		log.Fatal(err)
	}
	if err := execRunWithStderr(ctx,
		//"oapi-codegen",
		"go", "run", "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest",
		"-package", "github",
		"-generate", "client",
		"-templates", tmpl,
		"-o", client,
		f.Name(),
	); err != nil {
		log.Fatal(err)
	}

	if err := exec.CommandContext(ctx, "go", "mod", "tidy").Run(); err != nil {
		log.Fatal(err)
	}
}

func execRunWithStderr(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func getSha(ctx context.Context) string {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/repos/github/rest-api-description/commits", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	var commits []*commit
	if err := json.NewDecoder(resp.Body).Decode(&commits); err != nil {
		log.Fatal(err)
	}

	fmt.Println(commits[0].Sha)
	return commits[0].Sha
}

type commit struct {
	Sha string `json:"sha"`
}
