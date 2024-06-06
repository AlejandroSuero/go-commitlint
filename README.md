<a name="readme-top"></a>

<div align="center">

[![Stargazers][stars-shield]][stars-url]
[![Issues][issues-shield]][issues-url]

# Commit linter in Go

A commit linter based on [conventional commits](https://www.conventionalcommits.org/en/v1.0.0/)
rules written in [Go](https://go.dev/).

</div>

## What is commitlint?

Commitlint is a tool to help you maintain a
[conventional commit](https://www.conventionalcommits.org/en/v1.0.0/) message
style in your project.

General pattern:

```bash
type(scope): subject # scope is optional
```

Examples:

```bash
feat(commits): add filtering by scope
fix: fix typo in README
docs: add documentation for new features
feat!: add new breaking change feature
```

Common types according to the [conventional commits specification](https://www.conventionalcommits.org/en/v1.0.0/):

- `feat`: A new feature
- `fix`: A bug fix
- `docs`: Documentation only changes
- `style`: Changes that do not affect the meaning of the code (white-space,
  formatting, missing semi-colons, etc)
- `refactor`: A code change that neither fixes a bug nor adds a feature
- `perf`: A code change that improves performance
- `test`: Adding missing tests or correcting existing tests
- `build`: Changes that affect the build system or external dependencies
  (example scopes: gulp, broccoli, npm)
- `ci`: Changes to our CI configuration files and scripts
  (example scopes: Travis, Circle, BrowserStack, SauceLabs)
- `chore`: Other changes that don't modify src or test files

[stars-shield]: https://img.shields.io/github/stars/AlejandroSuero/go-commitlint.svg?style=for-the-badge
[stars-url]: https://github.com/AlejandroSuero/go-commitlint/stargazers
[issues-shield]: https://img.shields.io/github/issues/AlejandroSuero/go-commitlint.svg?style=for-the-badge
[issues-url]: https://github.com/AlejandroSuero/go-commitlint/issues
