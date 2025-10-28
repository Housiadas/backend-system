// Package validation contains the support for validating models.
package validation

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

// Check validates the provided model against it's declared tags.
func Check(val any) error {
	if err := validate.Struct(val); err != nil {
		var vErrors validator.ValidationErrors
		ok := errors.As(err, &vErrors)
		if !ok {
			return err
		}

		var fields FieldErrors
		for _, verror := range vErrors {
			fields.Add(
				verror.Field(),
				errors.New(verror.Translate(translator)),
			)
		}

		return &fields
	}

	return nil
}
