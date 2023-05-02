package query_test

import (
	"testing"

	"github.com/labstack/echo"

	"my-chat-jobsity-challenge"

	"github.com/stretchr/testify/assert"

	"my-chat-jobsity-challenge/pkg/utl/query"
)

func TestList(t *testing.T) {
	type args struct {
		user jobsity.AuthUser
	}
	cases := []struct {
		name     string
		args     args
		wantData *jobsity.ListQuery
		wantErr  error
	}{
		{
			name: "Super admin user",
			args: args{user: jobsity.AuthUser{
				Role: jobsity.SuperAdminRole,
			}},
		},
		{
			name: "Company admin user",
			args: args{user: jobsity.AuthUser{
				Role:      jobsity.CompanyAdminRole,
				CompanyID: 1,
			}},
			wantData: &jobsity.ListQuery{
				Query: "company_id = ?",
				ID:    1},
		},
		{
			name: "Location admin user",
			args: args{user: jobsity.AuthUser{
				Role:       jobsity.LocationAdminRole,
				CompanyID:  1,
				LocationID: 2,
			}},
			wantData: &jobsity.ListQuery{
				Query: "location_id = ?",
				ID:    2},
		},
		{
			name: "Normal user",
			args: args{user: jobsity.AuthUser{
				Role: jobsity.UserRole,
			}},
			wantErr: echo.ErrForbidden,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			q, err := query.List(tt.args.user)
			assert.Equal(t, tt.wantData, q)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
