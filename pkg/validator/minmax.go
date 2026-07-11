package validator

import (
	"strconv"
	"strings"
)

func (f *Form) Min(field string, params ...string) *Form {
	if len(params) == 0 {
		return f
	}
	minStr := params[0]
	min, err := strconv.Atoi(minStr)
	if err != nil {
		return f
	}

	value, ok := f.data[field]
	if !ok {
		return f
	}

	str, ok := value.(string)
	if !ok {
		return f
	}

	if len(strings.TrimSpace(str)) < min {
		f.addError(field, "The "+field+" field must be at least "+minStr+" characters.")
	}

	return f
}

func (f *Form) Max(field string, params ...string) *Form {
	if len(params) == 0 {
		return f
	}
	maxStr := params[0]
	max, err := strconv.Atoi(maxStr)
	if err != nil {
		return f
	}

	value, ok := f.data[field]
	if !ok {
		return f
	}

	str, ok := value.(string)
	if !ok {
		return f
	}

	if len(strings.TrimSpace(str)) > max {
		f.addError(field, "The "+field+" field must not exceed "+maxStr+" characters.")
	}

	return f
}
