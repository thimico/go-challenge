package query

import (
	"github.com/labstack/echo"

	"my-chat-jobsity-challenge"
)

// List prepares data for list queries
func List(u jobsity.AuthUser) (*jobsity.ListQuery, error) {
	switch true {
	case u.Role <= jobsity.AdminRole: // user is SuperAdmin or Admin
		return nil, nil
	case u.Role == jobsity.CompanyAdminRole:
		return &jobsity.ListQuery{Query: "company_id = ?", ID: u.CompanyID}, nil
	case u.Role == jobsity.LocationAdminRole:
		return &jobsity.ListQuery{Query: "location_id = ?", ID: u.LocationID}, nil
	default:
		return nil, echo.ErrForbidden
	}
}
