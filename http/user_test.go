package http

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/kencx/teal"
	"github.com/kencx/teal/mock"
	"github.com/kencx/teal/util"
)

var (
	testUser1 = &teal.User{
		Name:     "John Doe",
		Username: "john doe",
		Role:     "user",
	}
	inputTestUser1 = &teal.InputUser{
		Name:     "John Doe",
		Username: "john doe",
		Password: "abcd1234",
		Role:     "user",
	}
	inputTestUser2 = &teal.InputUser{
		Name:     "John Doe",
		Username: "john doe",
		Password: "abc",
		Role:     "user",
	}
)

func TestGetUser(t *testing.T) {
	testServer.Users = &mock.UserStore{
		GetUserFn: func(id int64) (*teal.User, error) {
			return testUser1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/users/1",
		data:   nil,
		params: map[string]string{"id": "1"},
		fn:     testServer.GetUser,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.User
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["users"]
	assertEqual(t, got.Name, testUser1.Name)
	assertEqual(t, got.Username, testUser1.Username)
	assertEqual(t, got.Role, testUser1.Role)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestGetUserByUsername(t *testing.T) {
	testServer.Users = &mock.UserStore{
		GetUserByUsernameFn: func(username string) (*teal.User, error) {
			return testUser1, nil
		},
	}

	tc := &testCase{
		method: http.MethodGet,
		url:    "/api/users/johndoe",
		data:   nil,
		params: map[string]string{"username": "johndoe"},
		fn:     testServer.GetUserByUsername,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.User
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["users"]
	assertEqual(t, got.Name, testUser1.Name)
	assertEqual(t, got.Username, testUser1.Username)
	assertEqual(t, got.Role, testUser1.Role)
	assertEqual(t, w.Code, http.StatusOK)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestUserRegister(t *testing.T) {
	want, err := util.ToJSON(inputTestUser1)
	checkErr(t, err)

	testServer.Users = &mock.UserStore{
		CreateUserFn: func(u *teal.User) (*teal.User, error) {
			return testUser1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/users/register/",
		data:   want,
		params: nil,
		fn:     testServer.Register,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)

	var env map[string]*teal.User
	err = json.NewDecoder(w.Body).Decode(&env)
	checkErr(t, err)

	got := env["users"]
	assertEqual(t, got.Name, testUser1.Name)
	assertEqual(t, got.Username, testUser1.Username)
	assertEqual(t, got.Role, testUser1.Role)
	assertEqual(t, w.Code, http.StatusCreated)
	assertEqual(t, w.HeaderMap.Get("Content-Type"), "application/json")
}

func TestUserRegisterPasswordFail(t *testing.T) {
	want, err := util.ToJSON(inputTestUser2)
	checkErr(t, err)

	testServer.Users = &mock.UserStore{
		CreateUserFn: func(u *teal.User) (*teal.User, error) {
			return testUser1, nil
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/users/register/",
		data:   want,
		params: nil,
		fn:     testServer.Register,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)
	assertValidationError(t, w, "password", "must be at least 8 chars long")
}

func TestUserRegisterDuplicateUsername(t *testing.T) {
	want, err := util.ToJSON(inputTestUser1)
	checkErr(t, err)

	testServer.Users = &mock.UserStore{
		CreateUserFn: func(u *teal.User) (*teal.User, error) {
			return nil, teal.ErrDuplicateUsername
		},
	}

	tc := &testCase{
		method: http.MethodPost,
		url:    "/api/users/register/",
		data:   want,
		params: nil,
		fn:     testServer.Register,
	}
	w, err := testResponse(t, tc)
	checkErr(t, err)
	assertValidationError(t, w, "username", "this username already exists")
}
