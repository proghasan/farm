package validator

import "strings"

func (f *Form) Required(field string, _ ...string) *Form {
	value, ok := f.data[field]

	if !ok || value == nil {
		f.addError(field, "The "+field+" field is required.")
		return f
	}

	if str, ok := value.(string); ok && strings.TrimSpace(str) == "" {
		f.addError(field, "The "+field+" field is required.")
	}

	return f
}