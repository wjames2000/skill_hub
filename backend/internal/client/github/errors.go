package github

import "errors"

var (
	ErrNotFound       = errors.New("not found")
	ErrRateLimited    = errors.New("rate limited")
	ErrTokenExhausted = errors.New("all tokens exhausted")
	ErrRepoNotFound   = errors.New("repository not found")
)
