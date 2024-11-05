package form_test

import (
	"context"
	"testing"

	"github.com/paulsonkoly/tracks/app/form"
)

type userTestDatum struct {
	testName        string
	username        string
	password        string
	passwordConfirm string
	valid           bool
}

var userTestData = [...]userTestDatum{
	{"simple valid example", "username", "password", "password", true},
	{"username too short", "op", "password", "password", false},
	{"password too short", "username", "12345", "12345", false},
	{"password mismatch", "username", "password", "<PASSWORD>", false},
}

type UserTestUnique struct{}

func (u UserTestUnique) UserUnique(_ context.Context, _ string) (bool, error) { return true, nil }

func TestUserValidate(t *testing.T) {
	for _, d := range userTestData {
		t.Run(d.testName, func(t *testing.T) {
			form := form.User{
				Username:        d.username,
				Password:        d.password,
				PasswordConfirm: d.passwordConfirm,
			}

			result, err := form.Validate(context.Background(), UserTestUnique{})
			if err != nil {
				t.Errorf("User{Usename: %q, Password: %q, PasswordConfirm: %q}.Validate() returned error: %v", d.username, d.password, d.passwordConfirm, err)
			}

			if result != d.valid {
				t.Errorf("User{Usename: %q, Password: %q, PasswordConfirm: %q}.Validate() = %v, want %v", d.username, d.password, d.passwordConfirm, result, d.valid)
				t.Errorf("Errors: %v, %v", form.Errors, form.FieldErrors)
			}
		})
	}
}
