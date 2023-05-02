package user

import (
	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"

	"my-chat-jobsity-challenge"
	"my-chat-jobsity-challenge/pkg/api/user/platform/pgsql"
)

// Service represents user application interface
type Service interface {
	Create(echo.Context, jobsity.User) (jobsity.User, error)
	List(echo.Context, jobsity.Pagination) ([]jobsity.User, error)
	View(echo.Context, int) (jobsity.User, error)
	Delete(echo.Context, int) error
	Update(echo.Context, Update) (jobsity.User, error)
}

// New creates new user application service
func New(db *pg.DB, udb UDB, rbac RBAC, sec Securer) *User {
	return &User{db: db, udb: udb, rbac: rbac, sec: sec}
}

// Initialize initalizes User application service with defaults
func Initialize(db *pg.DB, rbac RBAC, sec Securer) *User {
	return New(db, pgsql.User{}, rbac, sec)
}

// User represents user application service
type User struct {
	db   *pg.DB
	udb  UDB
	rbac RBAC
	sec  Securer
}

// Securer represents security interface
type Securer interface {
	Hash(string) string
}

// UDB represents user repository interface
type UDB interface {
	Create(orm.DB, jobsity.User) (jobsity.User, error)
	View(orm.DB, int) (jobsity.User, error)
	List(orm.DB, *jobsity.ListQuery, jobsity.Pagination) ([]jobsity.User, error)
	Update(orm.DB, jobsity.User) error
	Delete(orm.DB, jobsity.User) error
}

// RBAC represents role-based-access-control interface
type RBAC interface {
	User(echo.Context) jobsity.AuthUser
	EnforceUser(echo.Context, int) error
	AccountCreate(echo.Context, jobsity.AccessRole, int, int) error
	IsLowerRole(echo.Context, jobsity.AccessRole) error
}
