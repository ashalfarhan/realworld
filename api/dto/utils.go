package dto

import (
	"github.com/ashalfarhan/realworld/conduit"
	"github.com/go-playground/validator/v10"
)

func ValidateDto(s interface{}, v *validator.Validate) interface{} {
	if err := v.Struct(s); err != nil {
		e, ok := err.(validator.ValidationErrors)
		if !ok {
			return conduit.ErrInternal
		}

		errs := map[string][]string{}
		for _, field := range e {
			errs[field.Field()] = append(errs[field.Field()], field.Error())
		}

		return errs
	}

	return nil
}
