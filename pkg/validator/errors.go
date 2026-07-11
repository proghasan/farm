package validator


func (f *Form) addError(field, message string) {
	if f.errors == nil {
		f.errors = make(map[string][]string)
	}

	f.errors[field] = append(f.errors[field], message)
}




type ValidationError struct {
	Errors map[string][]string `json:"errors"`
}

func (e *ValidationError) Error() string {
	return "validation failed"
}