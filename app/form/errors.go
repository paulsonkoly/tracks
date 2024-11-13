package form

type errors struct {
	FieldErrors map[string][]string `form:"-"`
	Errors      []string            `form:"-"`
}

func (f errors) valid() bool {
	return len(f.FieldErrors) == 0 && len(f.Errors) == 0
}

// AddFieldError adds an error to the form object that is specific to a
// particular field. A field error will be rendered next to the field that has
// the problem.
func (f *errors) AddFieldError(field string, message string) {
	if f.FieldErrors == nil {
		f.FieldErrors = make(map[string]([]string))
	}

	existing, ok := f.FieldErrors[field]
	if ok {
		f.FieldErrors[field] = append(existing, message)
	} else {
		f.FieldErrors[field] = []string{message}
	}
}

// AddError adds and error to the form object that cannot be tied to any single
// field value, but indicates a more generic problem with the form values. Can
// also be used when a combination of multiple fields cause the error. A non
// field error will be rendered above the form.
func (f *errors) AddError(message string) {
	f.Errors = append(f.Errors, message)
}
