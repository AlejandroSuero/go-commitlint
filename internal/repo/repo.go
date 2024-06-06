package repo

import (
	git "github.com/go-git/go-git/v5"
)

type Repo func() (*git.Repository, error)

func FileSystem(directory string) Repo {
	return func() (*git.Repository, error) {
		return git.PlainOpen(directory)
	}
}
