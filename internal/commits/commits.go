package commits

import (
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/AlejandroSuero/go-commitlint/internal/repo"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
)

var shaLen = 7

// Author represents the author of a commit
type Author struct {
	Name  string
	Email string
}

// Commit represents a single commit
type Commit struct {
	Hash       string
	Message    string
	Date       time.Time
	NumParents int
	Author     *Author
}

// Commits is a function that returns a slice of Commits
type Commits func() []*Commit

// ID returns the commit hash
func (c *Commit) ID() string {
	return c.Hash
}

// ShortID returns the first 7 characters of the commit hash
func (c *Commit) ShortID() string {
	return c.Hash[:shaLen]
}

// Subject returns the subject of the commit message
func (c *Commit) Subject() string {
	return strings.Split(c.Message, "\n")[0]
}

func In(repository repo.Repo) Commits {
	return func() []*Commit {
		r, err := repository()
		if err != nil {
			panic(err)
		}
		ref, err := r.Head()
		if err != nil {
			panic(err)
		}
		iter, err := r.Log(&git.LogOptions{
			From: ref.Hash(),
		})
		if err != nil {
			panic(err)
		}
		commits := make([]*Commit, 0)
		err = iter.ForEach(func(c *object.Commit) error {
			commits = append(
				commits,
				&Commit{
					Hash:       c.Hash.String(),
					Message:    c.Message,
					Date:       c.Author.When,
					NumParents: len(c.ParentHashes),
					Author: &Author{
						Name:  c.Author.Name,
						Email: c.Author.Email,
					},
				},
			)
			return nil
		})
		if err != nil {
			panic(err)
		}
		return commits
	}
}

// Body returns the body of the commit message
func (c *Commit) Body() string {
	body := ""
	bodyParts := strings.Split(c.Message, "\n\n")

	if len(bodyParts) > 1 {
		body = strings.Join(bodyParts[1:], "")
	}

	return body
}

// filtered returns a Commits function that returns only commits that
// satisfy the given filter
func filtered(filter func(*Commit) bool, in Commits) (out Commits) {
	return func() []*Commit {
		f := make([]*Commit, 0)

		for _, c := range in() {
			if filter(c) {
				f = append(f, c)
			}
		}

		return f
	}
}

// Since returns a Commits function that returns only commits that are
// after the given time
func Since(timeFormat string, commits Commits) Commits {
	return filtered(func(c *Commit) bool {
		start, err := time.Parse("2006-01-02", timeFormat)
		if err != nil {
			panic(err)
		}
		return !c.Date.Before(start)
	},
		commits,
	)
}

// NotAuthoredByName returns a Commits function that returns only commits
// that are not authored by the given name pattern
func NotAuthoredByNames(patterns []string, commits Commits) Commits {
	return filtered(func(c *Commit) bool {
		for _, pattern := range patterns {
			match, err := regexp.MatchString(pattern, c.Author.Name)
			if err != nil {
				panic(err)
			}
			if match {
				return false
			}
		}
		return true
	},
		commits,
	)
}

// NotAuthoredByEmails returns a Commits function that returns only commits
// that are not authored by the given email patterns
func NotAuthoredByEmails(patterns []string, commits Commits) Commits {
	return filtered(func(c *Commit) bool {
		for _, pattern := range patterns {
			match, err := regexp.MatchString(pattern, c.Author.Email)
			if err != nil {
				panic(err)
			}
			if match {
				return false
			}
		}
		return true
	},
		commits,
	)
}

// WithMaxParents returns a Commits function that returns only commits with
// a maximum number of parents
func WithMaxParents(n int, commits Commits) Commits {
	return filtered(func(c *Commit) bool {
		return c.NumParents <= n
	},
		commits,
	)
}

// FakeCommit returns a Commits function that returns a single fake commit
func FakeCommit(reader io.Reader) Commits {
	return func() []*Commit {
		b, err := io.ReadAll(reader)
		if err != nil {
			panic(err)
		}
		return []*Commit{
			{
				Hash:    "fakesha",
				Message: string(b),
				Date:    time.Now(),
			},
		}
	}
}
