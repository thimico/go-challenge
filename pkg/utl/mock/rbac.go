package mock

import (
	"github.com/labstack/echo"

	"my-chat-jobsity-challenge"
)

// RBAC Mock
type RBAC struct {
	UserFn            func(echo.Context) jobsity.AuthUser
	EnforceRoleFn     func(echo.Context, jobsity.AccessRole) error
	EnforceUserFn     func(echo.Context, int) error
	EnforceCompanyFn  func(echo.Context, int) error
	EnforceLocationFn func(echo.Context, int) error
	AccountCreateFn   func(echo.Context, jobsity.AccessRole, int, int) error
	IsLowerRoleFn     func(echo.Context, jobsity.AccessRole) error
}

// User mock
func (a RBAC) User(c echo.Context) jobsity.AuthUser {
	return a.UserFn(c)
}

// EnforceRole mock
func (a RBAC) EnforceRole(c echo.Context, role jobsity.AccessRole) error {
	return a.EnforceRoleFn(c, role)
}

// EnforceUser mock
func (a RBAC) EnforceUser(c echo.Context, id int) error {
	return a.EnforceUserFn(c, id)
}

// EnforceCompany mock
func (a RBAC) EnforceCompany(c echo.Context, id int) error {
	return a.EnforceCompanyFn(c, id)
}

// EnforceLocation mock
func (a RBAC) EnforceLocation(c echo.Context, id int) error {
	return a.EnforceLocationFn(c, id)
}

// AccountCreate mock
func (a RBAC) AccountCreate(c echo.Context, roleID jobsity.AccessRole, companyID, locationID int) error {
	return a.AccountCreateFn(c, roleID, companyID, locationID)
}

// IsLowerRole mock
func (a RBAC) IsLowerRole(c echo.Context, role jobsity.AccessRole) error {
	return a.IsLowerRoleFn(c, role)
}
