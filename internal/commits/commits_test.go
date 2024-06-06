package commits_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/AlejandroSuero/go-commitlint/internal/commits"
	"github.com/AlejandroSuero/go-commitlint/internal/repo"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCommitID(t *testing.T) {
	const ID = "1234567890"
	commitID := (&commits.Commit{Hash: ID}).ID()
	if ID != commitID {
		t.Errorf("Expected \"%s\", got \"%s\"", ID, commitID)
	}
}

func TestCommitShortID(t *testing.T) {
	const ID = "1234567890"
	commitID := (&commits.Commit{Hash: ID}).ShortID()
	if ID[:7] != commitID {
		t.Errorf("Expected \"%s\", got \"%s\"", ID[:7], commitID)
	}
}

func TestCommitSubject(t *testing.T) {
	const subject = "This is a commit subject"
	commitSubject := (&commits.Commit{Message: subject + "\n\nThis is a commit body"}).Subject()
	if subject != commitSubject {
		t.Errorf("Expected \"%s\", got \"%s\"", subject, commitSubject)
	}
}

func TestCommitBody(t *testing.T) {
	const body = "This is a commit body"
	commitBody := (&commits.Commit{Message: "This is a commit subject\n\n" + body}).Body()
	if body != commitBody {
		t.Errorf("Expected \"%s\", got \"%s\"", body, commitBody)
	}
}

func TestIn(t *testing.T) {
	msgs := []string{
		"subject1\n\nbody1",
		"subject2\n\nbody2",
		"subject3\n\nbody3",
	}
	repo := tmpRepo(t, msgs...)
	cmmts := commits.In(repo)()
	if len(cmmts) != len(msgs) {
		t.Errorf("Expected %d commits, got %d", len(msgs), len(cmmts))
	}
	for i, msg := range msgs {
		commit := cmmts[len(cmmts)-i-1]
		expected := commit.Subject() + "\n\n" + commit.Body()
		if msg != expected {
			t.Errorf("Expected \"%s\", got \"%s\"", expected, msg)
		}
	}
}

func TestSince(t *testing.T) {
	before, err := time.Parse("2006-01-02", "2020-01-01")
	if err != nil {
		t.Fatal(err)
	}
	since, err := time.Parse("2006-01-02", "2020-01-02")
	if err != nil {
		t.Fatal(err)
	}
	after, err := time.Parse("2006-01-02", "2020-01-03")
	if err != nil {
		t.Fatal(err)
	}
	cmmts := commits.Since(
		"2020-01-02",
		func() []*commits.Commit {
			return []*commits.Commit{
				{Date: before},
				{Date: since},
				{Date: after},
			}
		},
	)()
	if len(cmmts) != 2 {
		t.Errorf("Expected %d commits, got %d", 2, len(cmmts))
	}
	assert.Contains(t, cmmts, &commits.Commit{Date: since})
	assert.Contains(t, cmmts, &commits.Commit{Date: after})
}

func TestFakeCommit(t *testing.T) {
	const message = "test subject\n\ntest body"
	cmmts := commits.FakeCommit(strings.NewReader(message))()
	if len(cmmts) != 1 {
		t.Errorf("Expected %d commits, got %d", 1, len(cmmts))
	}
	if cmmts[0].Subject() != "test subject" {
		t.Errorf("Expected \"%s\", got \"%s\"", "test subject", cmmts[0].Subject())
	}
	if cmmts[0].Body() != "test body" {
		t.Errorf("Expected \"%s\", got \"%s\"", "test body", cmmts[0].Body())
	}
}

func TestWithMaxParents(t *testing.T) {
	const maxCommits = 2
	cmmts := commits.WithMaxParents(maxCommits, func() []*commits.Commit {
		return []*commits.Commit{
			{NumParents: 1},
			{NumParents: maxCommits},
			{NumParents: 3},
		}
	})()
	if len(cmmts) != 2 {
		t.Errorf("Expected %d commits, got %d", 2, len(cmmts))
	}
	if cmmts[1].NumParents != maxCommits {
		t.Errorf("Expected %d parents, got %d", maxCommits, cmmts[1].NumParents)
	}
}

func TestNotAuthored(t *testing.T) {
	filtered := &commits.Commit{Author: randomAuthor()}
	expected := []*commits.Commit{
		{Author: randomAuthor()},
		{Author: randomAuthor()},
		{Author: randomAuthor()},
		{Author: randomAuthor()},
	}

	actual := commits.NotAuthoredByNames(
		[]string{filtered.Author.Name},
		func() []*commits.Commit { return append(expected, filtered) },
	)()

	assert.Equal(t, expected, actual)

	actual = commits.NotAuthoredByEmails(
		[]string{filtered.Author.Email},
		func() []*commits.Commit { return append(expected, filtered) },
	)()

	assert.Equal(t, expected, actual)
}

func randomAuthor() *commits.Author {
	return &commits.Author{
		Name:  uuid.New().String(),
		Email: uuid.New().String() + "@example.com",
	}
}

func tmpRepo(t *testing.T, msgs ...string) repo.Repo {
	directory, err := os.MkdirTemp("", strings.ReplaceAll(uuid.New().String(), "-", ""))
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(directory)

	return func() (*git.Repository, error) {
		repo, err := git.PlainInit(directory, false)
		if err != nil {
			t.Fatal(err)
		}
		wt, err := repo.Worktree()
		if err != nil {
			t.Fatal(err)
		}
		for i, msg := range msgs {
			file := fmt.Sprintf("msg%d.txt", i)
			err = os.WriteFile(filepath.Join(directory, file), []byte(msg), 0600)
			if err != nil {
				t.Fatal(err)
			}
			_, err = wt.Add(file)
			if err != nil {
				t.Fatal(err)
			}
			_, err = wt.Commit(msg, &git.CommitOptions{
				Author: &object.Signature{
					Name:  "John Doe",
					Email: "john.doe@example.com",
					When:  time.Now(),
				},
			})
			if err != nil {
				t.Fatal(err)
			}
		}
		return repo, nil
	}
}
