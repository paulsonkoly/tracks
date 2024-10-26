package form

type errors struct {
	FieldErrors map[string][]string `form:"-"`
	Errors      []string            `form:"-"`
}

func (f errors) Valid() bool {
	return len(f.FieldErrors) == 0 && len(f.Errors) == 0
}

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

func (f *errors) AddError(message string) {
	f.Errors = append(f.Errors, message)
}
