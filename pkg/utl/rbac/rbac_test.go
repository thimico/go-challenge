package rbac_test

import (
	"testing"

	"my-chat-jobsity-challenge"
	"my-chat-jobsity-challenge/pkg/utl/mock"
	"my-chat-jobsity-challenge/pkg/utl/rbac"

	"github.com/labstack/echo"

	"github.com/stretchr/testify/assert"
)

func TestUser(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{
		"id", "company_id", "location_id", "username", "email", "role"},
		9, 15, 52, "ribice", "ribice@gmail.com", jobsity.SuperAdminRole)
	wantUser := jobsity.AuthUser{
		ID:         9,
		Username:   "ribice",
		CompanyID:  15,
		LocationID: 52,
		Email:      "ribice@gmail.com",
		Role:       jobsity.SuperAdminRole,
	}
	rbacSvc := rbac.Service{}
	assert.Equal(t, wantUser, rbacSvc.User(ctx))
}

func TestEnforceRole(t *testing.T) {
	type args struct {
		ctx  echo.Context
		role jobsity.AccessRole
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, jobsity.CompanyAdminRole), role: jobsity.SuperAdminRole},
			wantErr: true,
		},
		{
			name:    "Authorized",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"role"}, jobsity.SuperAdminRole), role: jobsity.CompanyAdminRole},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.Service{}
			res := rbacSvc.EnforceRole(tt.args.ctx, tt.args.role)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceUser(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same user, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 15, jobsity.LocationAdminRole), id: 122},
			wantErr: true,
		},
		{
			name:    "Not same user, but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 22, jobsity.SuperAdminRole), id: 44},
			wantErr: false,
		},
		{
			name:    "Same user",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"id", "role"}, 8, jobsity.AdminRole), id: 8},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.Service{}
			res := rbacSvc.EnforceUser(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceCompany(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same company, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 7, jobsity.UserRole), id: 9},
			wantErr: true,
		},
		{
			name:    "Same company, not company admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 22, jobsity.UserRole), id: 22},
			wantErr: true,
		},
		{
			name:    "Same company, company admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 5, jobsity.CompanyAdminRole), id: 5},
			wantErr: false,
		},
		{
			name:    "Not same company but admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "role"}, 8, jobsity.AdminRole), id: 9},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.Service{}
			res := rbacSvc.EnforceCompany(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestEnforceLocation(t *testing.T) {
	type args struct {
		ctx echo.Context
		id  int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Not same location, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 7, jobsity.UserRole), id: 9},
			wantErr: true,
		},
		{
			name:    "Same location, not company admin or admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 22, jobsity.UserRole), id: 22},
			wantErr: true,
		},
		{
			name:    "Same location, company admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 5, jobsity.CompanyAdminRole), id: 5},
			wantErr: false,
		},
		{
			name:    "Location admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"location_id", "role"}, 5, jobsity.LocationAdminRole), id: 5},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.Service{}
			res := rbacSvc.EnforceLocation(tt.args.ctx, tt.args.id)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestAccountCreate(t *testing.T) {
	type args struct {
		ctx        echo.Context
		roleID     jobsity.AccessRole
		companyID  int
		locationID int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.UserRole), roleID: 500, companyID: 7, locationID: 8},
			wantErr: true,
		},
		{
			name:    "Same location, not company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.UserRole), roleID: 500, companyID: 2, locationID: 8},
			wantErr: true,
		},
		{
			name:    "Different location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.CompanyAdminRole), roleID: 400, companyID: 2, locationID: 4},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, not an admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.CompanyAdminRole), roleID: 500, companyID: 2, locationID: 3},
			wantErr: false,
		},
		{
			name:    "Same location, company, creating user role, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.CompanyAdminRole), roleID: 500, companyID: 2, locationID: 3},
			wantErr: false,
		},
		{
			name:    "Different everything, admin",
			args:    args{ctx: mock.EchoCtxWithKeys([]string{"company_id", "location_id", "role"}, 2, 3, jobsity.AdminRole), roleID: 200, companyID: 7, locationID: 4},
			wantErr: false,
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			rbacSvc := rbac.Service{}
			res := rbacSvc.AccountCreate(tt.args.ctx, tt.args.roleID, tt.args.companyID, tt.args.locationID)
			assert.Equal(t, tt.wantErr, res == echo.ErrForbidden)
		})
	}
}

func TestIsLowerRole(t *testing.T) {
	ctx := mock.EchoCtxWithKeys([]string{"role"}, jobsity.CompanyAdminRole)
	rbacSvc := rbac.Service{}
	if rbacSvc.IsLowerRole(ctx, jobsity.LocationAdminRole) != nil {
		t.Error("The requested user is higher role than the user requesting it")
	}
	if rbacSvc.IsLowerRole(ctx, jobsity.AdminRole) == nil {
		t.Error("The requested user is lower role than the user requesting it")
	}
}
