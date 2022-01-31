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
			// fmt.Println("field.ActualTag(): ", field.ActualTag())
			// fmt.Println("field.Error(): ", field.Error())
			// fmt.Println("field.Field(): ", field.Field())
			// fmt.Println("field.Kind(): ", field.Kind())
			// fmt.Println("field.Param(): ", field.Param())
			// fmt.Println("field.Value(): ", field.Value())
			// fmt.Println("field.Kind().String(): ", field.Kind().String())
			errs[field.Field()] = append(errs[field.Field()], field.Error())
		}

		return errs
	}

	return nil
}
