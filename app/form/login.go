package form

// Login represents the login form.
type Login struct {
	Username string `form:"username"`
	Password string `form:"password"`
	errors   `form:"-"`
}
