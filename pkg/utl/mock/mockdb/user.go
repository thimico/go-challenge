package mockdb

import (
	"github.com/go-pg/pg/v9/orm"

	"my-chat-jobsity-challenge"
)

// User database mock
type User struct {
	CreateFn         func(orm.DB, jobsity.User) (jobsity.User, error)
	ViewFn           func(orm.DB, int) (jobsity.User, error)
	FindByUsernameFn func(orm.DB, string) (jobsity.User, error)
	FindByTokenFn    func(orm.DB, string) (jobsity.User, error)
	ListFn           func(orm.DB, *jobsity.ListQuery, jobsity.Pagination) ([]jobsity.User, error)
	DeleteFn         func(orm.DB, jobsity.User) error
	UpdateFn         func(orm.DB, jobsity.User) error
}

// Create mock
func (u *User) Create(db orm.DB, usr jobsity.User) (jobsity.User, error) {
	return u.CreateFn(db, usr)
}

// View mock
func (u *User) View(db orm.DB, id int) (jobsity.User, error) {
	return u.ViewFn(db, id)
}

// FindByUsername mock
func (u *User) FindByUsername(db orm.DB, uname string) (jobsity.User, error) {
	return u.FindByUsernameFn(db, uname)
}

// FindByToken mock
func (u *User) FindByToken(db orm.DB, token string) (jobsity.User, error) {
	return u.FindByTokenFn(db, token)
}

// List mock
func (u *User) List(db orm.DB, lq *jobsity.ListQuery, p jobsity.Pagination) ([]jobsity.User, error) {
	return u.ListFn(db, lq, p)
}

// Delete mock
func (u *User) Delete(db orm.DB, usr jobsity.User) error {
	return u.DeleteFn(db, usr)
}

// Update mock
func (u *User) Update(db orm.DB, usr jobsity.User) error {
	return u.UpdateFn(db, usr)
}
