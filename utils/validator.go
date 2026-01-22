package utils

import (
	"fmt"

	"github.com/asaskevich/govalidator"
)

type ValidationError struct {
	Field   string
	Message string
}

type ValidatorErrorBag struct {
	Errors map[string][]string
}

func (v *ValidatorErrorBag) Error() string {
	return "The given input was invalid"
}

func Validate(s interface{}) error {
	_, err := govalidator.ValidateStruct(s)
	if err == nil {
		return nil
	}

	bag := &ValidatorErrorBag{
		Errors: make(map[string][]string),
	}

	// Helper function to process errors
	var processErrors func(err error)
	processErrors = func(err error) {
		switch e := err.(type) {
		case govalidator.Error:
			bag.Errors[e.Name] = append(bag.Errors[e.Name], e.Err.Error())
		case govalidator.Errors:
			for _, item := range e.Errors() {
				processErrors(item)
			}
		default:
			// Fallback for unexpected error types
			bag.Errors["_error"] = append(bag.Errors["_error"], fmt.Sprintf("%v", e))
		}
	}

	processErrors(err)

	if len(bag.Errors) > 0 {
		return bag
	}

	return nil
}
