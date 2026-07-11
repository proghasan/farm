package validator

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

var rules = map[string]func(*Form, string, ...string) *Form{
	"required": (*Form).Required,
	"min":      (*Form).Min,
	"max":      (*Form).Max,
	"unique":   (*Form).Unique,
}

type Form struct {
	ctx    fiber.Ctx
	db     *gorm.DB
	rules  map[string]string
	data   map[string]any
	errors map[string][]string
}

func New(c fiber.Ctx, db *gorm.DB) *Form {
	return &Form{
		ctx:   c,
		db:    db,
		rules: make(map[string]string),
	}
}

func (f *Form) Rules(rules map[string]string) *Form {
	f.rules = rules
	return f
}

func (f *Form) Validate(dest any) error {
	if err := f.ctx.Bind().Body(dest); err != nil {
		return err
	}

	if err := f.ctx.Bind().Body(&f.data); err != nil {
		return err
	}

	for field, ruleStr := range f.rules {
		for _, rule := range strings.Split(ruleStr, "|") {
			parts := strings.SplitN(rule, ":", 2)
			handler, ok := rules[parts[0]]
			if !ok {
				continue
			}

			if len(parts) > 1 {
				_ = handler(f, field, strings.Split(parts[1], ",")...)
			} else {
				_ = handler(f, field)
			}
		}
	}

	if len(f.errors) > 0 {
		return &ValidationError{Errors: f.errors}
	}

	return nil
}