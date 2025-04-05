package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"
	"os"
	"os/exec"
)

func getClient(token string) *github.Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)

	return github.NewClient(tc)
}

// buildMkdocsSite builds the mkdocs site using the `mkdocs build` command.
func buildMkdocsSite() error {
	cmd := exec.Command("python", "-m", "mkdocs", "build")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// commitToGhPages uses go-github to push changes to the gh-pages branch via GitHub API.
func commitToGhPages(client *github.Client, owner, repo, branch string, path string) error {
	ctx := context.Background()

	// Navigate to the site directory (usually "site" after mkdocs build)
	err := os.Chdir("site")
	if err != nil {
		return fmt.Errorf("failed to change directory to 'site': %w", err)
	}

	// Check if branch exists and create if it doesn't
	_, _, err = client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		enableOrUpdatePages(client, owner, repo, branch, path)
	}

	// Commit and push the changes to gh-pages branch
	commitMessage := "Update mkdocs site"

	// Push the new site contents to gh-pages
	err = pushToGitHub(branch, commitMessage)
	if err != nil {
		return fmt.Errorf("failed to push to gh-pages: %w", err)
	}

	return nil
}

// pushToGitHub commits and pushes to the gh-pages branch using GitHub API
func pushToGitHub(branch, commitMessage string) error {
	cmd := exec.Command("git", "add", ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git add failed: %w", err)
	}

	cmd = exec.Command("git", "commit", "-m", commitMessage)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git commit failed: %w", err)
	}

	cmd = exec.Command("git", "push", "-f", "origin", branch)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git push failed: %w", err)
	}

	return nil
}

func enableOrUpdatePages(client *github.Client, owner, repo, branch, path string) error {
	ctx := context.Background()

	// Check if the gh-pages branch exists
	// NOTE: The gh-pages branch must exist before GitHub Pages can be built
	_, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+branch)
	if err != nil {
		fmt.Printf("%s branch not found, creating it\n", branch)

		// Get default branch (e.g., main/master) from the repo
		getRepo, _, err := client.Repositories.Get(ctx, owner, repo)
		if err != nil {
			fmt.Errorf("ERROR: Unable to get repo info: %s", err.Error())
			return err
		}
		defaultBranch := getRepo.GetDefaultBranch()

		// Get SHA of the latest commit on the default branch
		baseRef, _, err := client.Git.GetRef(ctx, owner, repo, "refs/heads/"+defaultBranch)
		if err != nil {
			fmt.Errorf("ERROR: Unable to get default branch ref: %s", err.Error())
			return err
		}
		baseSha := baseRef.Object.GetSHA()

		// Create the gh-pages branch pointing to the same commit
		ref := &github.Reference{
			Ref: github.String("refs/heads/gh-pages"),
			Object: &github.GitObject{
				SHA: github.String(baseSha),
			},
		}
		_, _, err = client.Git.CreateRef(ctx, owner, repo, ref)
		if err != nil {
			fmt.Errorf("ERROR: Unable to create gh-pages branch: %s", err.Error())
			return err
		}

		fmt.Println("Created gh-pages branch")
	}

	source := &github.PagesSource{
		Branch: github.String(branch),
		Path:   github.String(path),
	}
	enableOpts := &github.Pages{
		Source: source,
	}

	updateOpts := &github.PagesUpdate{
		Source: source,
	}

	_, resp, err := client.Repositories.EnablePages(ctx, owner, repo, enableOpts)
	if err != nil {
		fmt.Println("ℹ️ Pages already enabled. Updating instead")
		_, err := client.Repositories.UpdatePages(ctx, owner, repo, updateOpts)
		if err != nil {
			return fmt.Errorf("failed to enable/update Pages: %v (%v)", err, resp.Status)
		}

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
		fmt.Println("Usage: mkdocs-ghpages <enable|update|disable|mkdocs-commit> <owner> <repo> <token> [branch] [path]")
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
	case "mkdocs-commit":
		err := buildMkdocsSite()
		if err != nil {
			fmt.Println("Error: mkdocs build ", err)
		}
		err = commitToGhPages(client, owner, repo, branch, path)
		if err != nil {
			fmt.Println("Error: git commit ", err)
		}
	default:
		fmt.Println("Invalid action. Use enable, update, or disable.")
	}
}
