package user_test

import (
	"testing"

	"github.com/rezaAmiri123/test-project/internal/users/domain/user"
	"github.com/stretchr/testify/require"
)

func TestUser_Validate(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		Name      string
		GetUserFn func() *user.User
		HasError  bool
	}{
		{
			Name: "valid_user",
			GetUserFn: func() *user.User {
				return newUser(t)
			},
			HasError: false,
		},
		{
			Name: "invalid_empty_email",
			GetUserFn: func() *user.User {
				u := newUser(t)
				u.Email = ""
				return u
			},
			HasError: true,
		},
		{
			Name: "invalid_email",
			GetUserFn: func() *user.User {
				u := newUser(t)
				u.Email = "123456789"
				return u
			},
			HasError: true,
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()
			u := tc.GetUserFn()
			err := u.Validate()
			if tc.HasError {
				require.True(t, err != nil)
			} else {
				require.True(t, err == nil)
			}
		})
	}
}

func newUser(t testing.TB) *user.User {
	t.Helper()
	return user.NewUser("username", "password", "email@example.com", "", "")
}
