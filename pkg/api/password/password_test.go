package password_test

import (
	"testing"

	"my-chat-jobsity-challenge"
	"my-chat-jobsity-challenge/pkg/api/password"

	"github.com/go-pg/pg/v9/orm"
	"github.com/labstack/echo"

	"my-chat-jobsity-challenge/pkg/utl/mock"
	"my-chat-jobsity-challenge/pkg/utl/mock/mockdb"

	"github.com/stretchr/testify/assert"
)

func TestChange(t *testing.T) {
	type args struct {
		oldpass string
		newpass string
		id      int
	}
	cases := []struct {
		name    string
		args    args
		wantErr bool
		udb     *mockdb.User
		rbac    *mock.RBAC
		sec     *mock.Secure
	}{
		{
			name: "Fail on EnforceUser",
			args: args{id: 1},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return jobsity.ErrGeneric
				}},
			wantErr: true,
		},
		{
			name:    "Fail on ViewUser",
			args:    args{id: 1},
			wantErr: true,
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (jobsity.User, error) {
					if id != 1 {
						return jobsity.User{}, nil
					}
					return jobsity.User{}, jobsity.ErrGeneric
				},
			},
		},
		{
			name: "Fail on PasswordMatch",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (jobsity.User, error) {
					return jobsity.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return false
				},
			},
		},
		{
			name: "Fail on InsecurePassword",
			args: args{id: 1, oldpass: "hunter123"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			wantErr: true,
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (jobsity.User, error) {
					return jobsity.User{
						Password: "HashedPassword",
					}, nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return false
				},
			},
		},
		{
			name: "Success",
			args: args{id: 1, oldpass: "hunter123", newpass: "password"},
			rbac: &mock.RBAC{
				EnforceUserFn: func(c echo.Context, id int) error {
					return nil
				}},
			udb: &mockdb.User{
				ViewFn: func(db orm.DB, id int) (jobsity.User, error) {
					return jobsity.User{
						Password: "$2a$10$udRBroNGBeOYwSWCVzf6Lulg98uAoRCIi4t75VZg84xgw6EJbFNsG",
					}, nil
				},
				UpdateFn: func(orm.DB, jobsity.User) error {
					return nil
				},
			},
			sec: &mock.Secure{
				HashMatchesPasswordFn: func(string, string) bool {
					return true
				},
				PasswordFn: func(string, ...string) bool {
					return true
				},
				HashFn: func(string) string {
					return "hash3d"
				},
			},
		},
	}
	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			s := password.New(nil, tt.udb, tt.rbac, tt.sec)
			err := s.Change(nil, tt.args.id, tt.args.oldpass, tt.args.newpass)
			assert.Equal(t, tt.wantErr, err != nil)
			// Check whether password was changed
		})
	}
}
