package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"

	playground "github.com/go-playground/validator/v10"
)

var (
	once sync.Once
	v    *playground.Validate
)

func instance() *playground.Validate {
	once.Do(func() {
		v = playground.New()
	})
	return v
}

// Validate returns a map of field -> message when validation fails.
// If the returned map is non-empty, the input is invalid.
func Validate(input any) map[string]string {
	err := instance().Struct(input)
	if err == nil {
		return nil
	}

	ve, ok := err.(playground.ValidationErrors)
	if !ok {
		return map[string]string{"_error": err.Error()}
	}

	out := make(map[string]string, len(ve))
	typ := reflect.TypeOf(input)
	if typ != nil && typ.Kind() == reflect.Pointer {
		typ = typ.Elem()
	}

	for _, fe := range ve {
		key := jsonFieldName(typ, fe.StructField())
		out[key] = messageFor(fe)
	}

	return out
}

func jsonFieldName(typ reflect.Type, structField string) string {
	if typ != nil && typ.Kind() == reflect.Struct {
		if f, ok := typ.FieldByName(structField); ok {
			tag := f.Tag.Get("json")
			name := strings.TrimSpace(strings.Split(tag, ",")[0])
			if name != "" && name != "-" {
				return name
			}
		}
	}
	return strings.ToLower(structField)
}

func messageFor(fe playground.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	default:
		return "is invalid"
	}
}

