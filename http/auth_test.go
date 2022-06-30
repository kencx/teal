package http

import (
	"testing"

	"github.com/kencx/teal"
	"github.com/kencx/teal/mock"
)

func TestAuthenticate(t *testing.T) {
	testUser1.SetPassword(inputTestUser1.Password)

	testServer.Users = &mock.UserStore{
		GetUserByUsernameFn: func(username string) (*teal.User, error) {
			return testUser1, nil
		},
	}

	t.Run("success", func(t *testing.T) {
		authenticated, user, err := testServer.authenticate(testUser1.Username, inputTestUser1.Password)
		checkErr(t, err)

		assertEqual(t, authenticated, true)
		assertEqual(t, user.Name, testUser1.Name)
		assertEqual(t, user.Username, testUser1.Username)
		assertEqual(t, user.Role, testUser1.Role)
	})

	t.Run("no username password", func(t *testing.T) {
		authenticated, user, err := testServer.authenticate("", "")
		checkErr(t, err)
		assertEqual(t, authenticated, false)
		assertEqual(t, user, nil)
	})

	t.Run("password does not match", func(t *testing.T) {
		authenticated, _, err := testServer.authenticate(testUser1.Username, "wrongPassword")
		checkErr(t, err)
		assertEqual(t, authenticated, false)
	})

	t.Run("user does not exist", func(t *testing.T) {
		testServer.Users = &mock.UserStore{
			GetUserByUsernameFn: func(username string) (*teal.User, error) {
				return nil, teal.ErrDoesNotExist
			},
		}

		authenticated, user, err := testServer.authenticate("wrongUser", inputTestUser1.Password)
		if err == nil {
			t.Errorf("expected err: user does not exist")
		}
		if err != teal.ErrDoesNotExist {
			t.Errorf("unexpected err: %v", err)
		}
		assertEqual(t, authenticated, false)
		assertEqual(t, user, nil)
	})
}
