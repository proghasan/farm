package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type FiberValidator struct {
	validate *validator.Validate
}

func (v *FiberValidator) Validate(out any) error {
	trimStrings(out)
	return v.validate.Struct(out)
}

func NewFiberValidator() *FiberValidator {
	return &FiberValidator{validate: validate}
}

func trimStrings(v any) {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}
	for i := range val.NumField() {
		field := val.Field(i)
		switch field.Kind() {
		case reflect.String:
			if field.CanSet() {
				field.SetString(strings.TrimSpace(field.String()))
			}
		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.String && field.Elem().CanSet() {
				field.Elem().SetString(strings.TrimSpace(field.Elem().String()))
			}
		}
	}
}


