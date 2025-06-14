// Package validation contains the support for validating models.
package validation

import (
	"errors"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

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
