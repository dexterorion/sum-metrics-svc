package helpers

import "errors"

var (
	ErrNotFoundOnRepository = errors.New("mongo: no documents in result")
	ErrUserAlreadyOnboarded = errors.New("user already onboarded")
)
