package dto

import (
	"github.com/go-playground/validator/v10"
)

func ValidateDto(s interface{}, v *validator.Validate) interface{} {
	if err := v.Struct(s); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			return err.Error()
		}
		errs := map[string][]string{}
		for _, f := range e {
			errs[f.Field()] = append(errs[f.Field()], f.Tag())
		}

		return errs
	}

	return nil
}
