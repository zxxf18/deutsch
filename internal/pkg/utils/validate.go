package utils

import (
	"reflect"
	"regexp"
	"unicode/utf8"

	"github.com/go-playground/validator/v10"
)

const (
	duration = "duration"

	nonzero = "nonzero"
	nonnil  = "nonnil"
)

var regexps = map[string]string{
	duration: "^[1-9][0-9]*(s|m|h)$",
}
var validate *validator.Validate

func init() {
	validate = validator.New()
	RegisterValidate(validate)
}

func GetValidator() *validator.Validate {
	return validate
}

func RegisterValidation(key string, fn validator.Func) {
	GetValidator().RegisterValidation(key, fn)
}

func RegisterValidate(v *validator.Validate) {
	if v != nil {
		for key, val := range regexps {
			key0, val0 := key, val
			v.RegisterValidation(key0, func(fl validator.FieldLevel) bool {
				match, _ := regexp.MatchString(val0, fl.Field().String())
				return match
			})
		}

		v.RegisterValidation(nonzero, func(fl validator.FieldLevel) bool {
			return nonzeroValid(fl.Field().Interface())
		})

		v.RegisterValidation(nonnil, func(fl validator.FieldLevel) bool {
			return nonnilValid(fl.Field().Interface())
		})
	}
}

// nonzeroValid tests whether a variable value non-zero as defined by the golang spec.
func nonzeroValid(v interface{}) bool {
	st := reflect.ValueOf(v)
	valid := true
	switch st.Kind() {
	case reflect.String:
		valid = utf8.RuneCountInString(st.String()) != 0
	case reflect.Ptr, reflect.Interface:
		valid = !st.IsNil()
	case reflect.Slice, reflect.Map, reflect.Array:
		valid = st.Len() != 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid = st.Int() != 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid = st.Uint() != 0
	case reflect.Float32, reflect.Float64:
		valid = st.Float() != 0
	case reflect.Bool:
		valid = st.Bool()
	case reflect.Invalid:
		valid = false // always invalid
	case reflect.Struct:
		valid = true // always valid since only nil pointers are empty
	default:
		valid = false
	}
	return valid
}

// nonnilValid validates that the given pointer is not nil
func nonnilValid(v interface{}) bool {
	st := reflect.ValueOf(v)
	switch st.Kind() {
	case reflect.Ptr, reflect.Interface:
		if st.IsNil() {
			return false
		}
	case reflect.Invalid:
		return false
	}
	return true
}
