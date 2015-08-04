package validation

import (
	"fmt"
	"reflect"
	"strings"
)

import (
	bvalidation "github.com/astaxie/beego/validation"
)

// Recursively validates a struct against `valid:"..." tags`
// The arg `parentType` is actually only used as an optional single argument
// Anything more than one is ignored
func Validate(obj interface{}, parentType ...string) (errs []error) {
	// Actual validation happens here
	validator := bvalidation.Validation{}
	valid, err := validator.Valid(obj)
	if err != nil {
		errs = append(errs, err)
		return
	}

	// Figure out if this object is a pointer
	// If so, derefernce it
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// The actual type name of this object
	Type := v.Type().String()

	// If validation above failed
	if !valid {
		// If this function was called with at least one argument
		// Prepend that argument to the type name
		// Helps us find context in recursive calls
		if len(parentType) > 0 {
			parts := strings.Split(Type, ".")
			Type = fmt.Sprintf("%s.%s", parentType[0], parts[len(parts)-1])
		}

		// Collect all the validation errors and return them
		for _, e := range validator.Errors {
			errFmt := "%s validation failed: `%s` %s (actual value: %#v)"
			errs = append(errs,
				fmt.Errorf(errFmt, Type, e.Field, e.Message, e.Value))
		}
	}

	// Find sub-structs and validate them, too
	for i := range make([]struct{}, v.NumField()) {
		field := v.Field(i)
		if field.Kind() == reflect.Struct {
			errs = append(errs, Validate(field.Interface(), Type)...)
		} else if field.Kind() == reflect.Slice {
			if field.Len() > 0 {
				firstItem := field.Index(0)
				if firstItem.Kind() == reflect.Struct {
					for j := range make([]struct{}, field.Len()) {
						item := field.Index(j).Interface()
						errs = append(errs, Validate(item, Type)...)
					}
				}
			}
		}
	}
	return
}
