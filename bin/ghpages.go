package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
	"os"
)

func getClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

func enableOrUpdatePages(client *github.Client, owner, repo, branch, path string) error {
	ctx := context.Background()
	source := &github.PagesSource{
		Branch: github.String(branch),
		Path:   github.String(path),
	}
	opt := &github.Pages{
		Source: source,
	}

	_, resp, err := client.Repositories.EnablePages(ctx, owner, repo, opt)
	if err != nil {
		return fmt.Errorf("failed to enable/update Pages: %v (%v)", err, resp.Status)
	}

	fmt.Println("✅ GitHub Pages enabled or updated")
	return nil
}

func disablePages(client *github.Client, owner, repo string) error {
	ctx := context.Background()
	resp, err := client.Repositories.DisablePages(ctx, owner, repo)
	if err != nil {
		return fmt.Errorf("failed to disable Pages: %v (%v)", err, resp.Status)
	}

	fmt.Println("🚫 GitHub Pages disabled")
	return nil
}

func main() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: ghpages <enable|update|disable> <owner> <repo> <token> [branch] [path]")
		return
	}

	action := os.Args[1]
	owner := os.Args[2]
	repo := os.Args[3]
	token := os.Args[4]
	branch := "gh-pages"
	path := "/"

	if len(os.Args) > 5 {
		branch = os.Args[5]
	}
	if len(os.Args) > 6 {
		path = os.Args[6]
	}

	client := getClient(token)

	switch action {
	case "enable", "update":
		if err := enableOrUpdatePages(client, owner, repo, branch, path); err != nil {
			fmt.Println("Error:", err)
		}
	case "disable":
		if err := disablePages(client, owner, repo); err != nil {
			fmt.Println("Error:", err)
		}
	default:
		fmt.Println("Invalid action. Use enable, update, or disable.")
	}
}
