package transport

import (
	"my-chat-jobsity-challenge"
)

// User model response
// swagger:response userResp
type swaggUserResponse struct {
	// in:body
	Body struct {
		*jobsity.User
	}
}

// Users model response
// swagger:response userListResp
type swaggUserListResponse struct {
	// in:body
	Body struct {
		Users []jobsity.User `json:"users"`
		Page  int            `json:"page"`
	}
}
