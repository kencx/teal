package teal

import (
	"errors"
	"testing"

	"github.com/kencx/teal/validator"
	"golang.org/x/crypto/bcrypt"
)

var (
	testPw      = "abcd1234"
	testPwShort = "abc"
	testPwEmpty = ""
)

func TestValidateUser(t *testing.T) {
	tests := []struct {
		name string
		user *InputUser
		err  map[string]string
	}{{
		name: "success",
		user: &InputUser{
			Name:     "John Doe",
			Username: "johndoe",
			Password: testPw,
		},
		err: nil,
	}, {
		name: "no name",
		user: &InputUser{
			Name:     "",
			Username: "johndoe",
			Password: testPw,
		},
		err: map[string]string{"name": "value is missing"},
	}, {
		name: "no username",
		user: &InputUser{
			Name:     "John Doe",
			Username: "",
			Password: testPw,
		},
		err: map[string]string{"username": "value is missing"},
	}, {
		name: "no password",
		user: &InputUser{
			Name:     "John Doe",
			Username: "johndoe",
			Password: testPwEmpty,
		},
		err: map[string]string{"password": "value is missing"},
	}, {
		name: "no username",
		user: &InputUser{
			Name:     "John Doe",
			Username: "johndoe",
			Password: testPwShort,
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

func TestSetPassword(t *testing.T) {
	u := &User{
		Name:     "John Doe",
		Username: "johndoe",
	}

	err := u.SetPassword(testPw)
	if err != nil {
		t.Fatalf("unexpected err: %v", err)
	}
	err = bcrypt.CompareHashAndPassword(u.HashedPassword, []byte(testPw))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			t.Errorf("hashed password does not match")
		default:
			t.Fatalf("unexpected err: %v", err)
		}
	}
}

func TestPasswordMatches(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		u := &User{
			Name:     "John Doe",
			Username: "johndoe",
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(testPw), 12)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		u.HashedPassword = hash

		matches, err := u.PasswordMatches(testPw)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}

		if !matches {
			t.Errorf("got %v, want %v", matches, true)
		}
	})

	t.Run("fail", func(t *testing.T) {
		u := &User{
			Name:     "John Doe",
			Username: "johndoe",
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(testPw), 12)
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		u.HashedPassword = hash

		matches, err := u.PasswordMatches("wrongPassword")
		if err != nil {
			t.Fatalf("unexpected err: %v", err)
		}
		if matches {
			t.Errorf("got %v, want %v", matches, false)
		}
	})
}
