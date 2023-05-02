package mock

import (
	"my-chat-jobsity-challenge"
)

// JWT mock
type JWT struct {
	GenerateTokenFn func(jobsity.User) (string, error)
}

// GenerateToken mock
func (j JWT) GenerateToken(u jobsity.User) (string, error) {
	return j.GenerateTokenFn(u)
}
