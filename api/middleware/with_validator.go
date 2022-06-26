package middleware

// import (
// 	"context"
// 	"encoding/json"
// 	"net/http"
// 	"reflect"

// 	"github.com/ashalfarhan/realworld/api/response"
// 	"github.com/go-playground/validator/v10"
// )

// func (m *ConduitMiddleware) WithValidator(next http.HandlerFunc, schema interface{}) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		// d is generic pointer of struct with its validation
// 		// s is the new instance of d, so that validation is going to created every time
// 		// and decoding request body to the struct is only for the current request context
// 		// @see: https://stackoverflow.com/a/10211940/14855798
// 		//
// 		// if d is not the pointer of the struct, then use:
// 		// s := reflect.New(reflect.TypeOf(d)).Interface()
// 		s := reflect.New(reflect.TypeOf(schema).Elem()).Interface() // new pointer of d

// 		if err := json.NewDecoder(r.Body).Decode(s); err != nil {
// 			response.ClientError(w, err)
// 			return
// 		}

// 		v := validator.New()
// 		if err := v.Struct(s); err != nil {
// 			m.logger.Infof("Validation Error %#v, %v", s, err)
// 			response.EntityError(w, err)
// 			return
// 		}

// 		ctx := context.WithValue(r.Context(), "", s)
// 		next(w, r.WithContext(ctx))
// 	}
// }
