package form_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/paulsonkoly/tracks/app/form"
	"github.com/stretchr/testify/assert"
)

type userTestDatum struct {
	form                form.User
	expectedResult      bool
	expectedErrors      []string
	expectedFieldErrors map[string][]string
}

var userTestData = [...]userTestDatum{
	{form.User{Username: "username", Password: "password", PasswordConfirm: "password"}, true, nil, nil},
	{form.User{Username: "op", Password: "password", PasswordConfirm: "password"}, false, nil, map[string][]string{"Username": {"Username too short. Must be at least 3 characters long."}}},
	{form.User{Username: "username", Password: "12345", PasswordConfirm: "12345"}, false, nil, map[string][]string{"Password": {"Password too short. Must be at least 6 characters long."}}},
	{form.User{Username: "username", Password: "password", PasswordConfirm: "<PASSWORD>"}, false, nil, map[string][]string{"PasswordConfirm": {"Passwords do not match."}}},
}

var userEditTestData = [...]userTestDatum{
	{form.User{Username: "username", Password: "password", PasswordConfirm: "password"}, true, nil, nil},
	{form.User{Username: "", Password: "", PasswordConfirm: ""}, true, nil, nil},
	{form.User{Username: "", Password: "password", PasswordConfirm: "password"}, true, nil, nil},
	{form.User{Username: "username", Password: "", PasswordConfirm: ""}, true, nil, nil},
	{form.User{Username: "username", Password: "", PasswordConfirm: "<PASSWORD>"}, false, nil, map[string][]string{"Password": {"Password too short. Must be at least 6 characters long."}, "PasswordConfirm": {"Passwords do not match."}}},
	{form.User{Username: "username", Password: "password", PasswordConfirm: ""}, false, nil, map[string][]string{"PasswordConfirm": {"Passwords do not match."}}},
}

type UserTestUnique struct{}

func (u UserTestUnique) UserUnique(_ context.Context, _ string) (bool, error) { return true, nil }
func (u UserTestUnique) UserUniqueExceptID(_ context.Context, _ int, _ string) (bool, error) {
	return true, nil
}

func TestUserValidate(t *testing.T) {
	testUserForm(t, "save", userTestData[:])
}

func TestUserValidateEdit(t *testing.T) {
	testUserForm(t, "edit", userEditTestData[:])
}

func testUserForm(t *testing.T, op string, testData []userTestDatum) {
	for _, d := range testData {
		f := d.form
		testName := fmt.Sprintf("%s User{Usename: %q, Password: %q, PasswordConfirm: %q}.Validate()", op, f.Username, f.Password, f.PasswordConfirm)

		t.Run(testName, func(t *testing.T) {
			var (
				result bool
				err    error
			)

			switch op {
			case "save":
				result, err = f.Validate(context.Background(), UserTestUnique{})
			case "edit":
				result, err = f.ValidateEdit(context.Background(), UserTestUnique{})
			}

			assert.NoError(t, err)
			assert.Equal(t, d.expectedResult, result)
			assert.Equal(t, d.expectedErrors, f.Errors)
			assert.Equal(t, d.expectedFieldErrors, f.FieldErrors)
		})
	}
}
