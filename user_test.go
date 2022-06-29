package teal

import (
	"testing"

	"github.com/kencx/teal/validator"
)

var (
	testPassword      = "abcd1234"
	testPasswordFail  = "abc"
	testPasswordEmpty = ""
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name string
		user *User
		err  map[string]string
	}{{
		name: "success",
		user: &User{
			Name:     "John Doe",
			Username: "johndoe",
			HashedPassword: Password{
				Text: &testPassword,
				Hash: []byte(""),
			},
		},
		err: nil,
	}, {
		name: "no name",
		user: &User{
			Name:     "",
			Username: "johndoe",
			HashedPassword: Password{
				Text: &testPassword,
				Hash: []byte(""),
			},
		},
		err: map[string]string{"name": "value is missing"},
	}, {
		name: "no username",
		user: &User{
			Name:     "John Doe",
			Username: "",
			HashedPassword: Password{
				Text: &testPassword,
				Hash: []byte(""),
			},
		},
		err: map[string]string{"username": "value is missing"},
	}, {
		name: "no password",
		user: &User{
			Name:     "John Doe",
			Username: "johndoe",
			HashedPassword: Password{
				Text: &testPasswordEmpty,
				Hash: []byte(""),
			},
		},
		err: map[string]string{"password": "value is missing"},
	}, {
		name: "password too short",
		user: &User{
			Name:     "John Doe",
			Username: "johndoe",
			HashedPassword: Password{
				Text: &testPasswordFail,
				Hash: []byte(""),
			},
		},
		err: map[string]string{"password": "must be at least 8 chars long"},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := validator.New()
			tt.user.Validate(v)

			if !v.Valid() && tt.err == nil {
				t.Fatalf("expected no err, got %v", v.Errors)
			}

			if v.Valid() && tt.err != nil {
				t.Fatalf("expected err with %q, got nil", tt.err)
			}

			if !v.Valid() && tt.err != nil {
				if len(v.Errors) != len(tt.err) {
					t.Fatalf("got %d errs, want %d errs", len(v.Errors), len(tt.err))
				}

				for k, v := range v.Errors {
					s, ok := tt.err[k]
					if !ok {
						t.Fatalf("err field missing %q", k)
					}

					if v != s {
						t.Fatalf("got %v, want %v error", v, s)
					}
				}
			}
		})
	}
}
