package pgsql

import (
	"github.com/go-pg/pg/v9/orm"

	"my-chat-jobsity-challenge"
)

// User represents the client for user table
type User struct{}

// View returns single user by ID
func (u User) View(db orm.DB, id int) (jobsity.User, error) {
	user := jobsity.User{Base: jobsity.Base{ID: id}}
	err := db.Select(&user)
	return user, err
}

// Update updates user's info
func (u User) Update(db orm.DB, user jobsity.User) error {
	return db.Update(&user)
}
