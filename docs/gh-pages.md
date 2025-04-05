# GitHub Pages Management Tool

This Go-based command-line tool allows you to enable, update, or disable GitHub Pages for a repository. By leveraging the GitHub API, the tool simplifies the management of GitHub Pages configurations, providing the following functionalities:

- Enable GitHub Pages: Set up GitHub Pages for the repository with a specified branch and path.
- Update GitHub Pages: Modify the existing GitHub Pages setup, changing the branch or path.
- Disable GitHub Pages: Disable GitHub Pages for the repository.

The tool interacts with GitHub via the go-github library, making it easy to integrate into your workflows, automate repository configurations, and manage Pages without manually editing settings on GitHub.


### 📌 Dependencies

- `golang` 1.24

### 🏗️ Build

```bash
go build .
```

## 🧑‍💻 Usage

> 📌 NOTE

```bash
Usage:
    mkdocs-ghpages <enable|update|disable|mkdocs-commit> <owner> <repo> <token> [branch] [path]
```


