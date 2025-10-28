package validation

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// validate holds the settings and caches for validating request struct values.
var validate *validator.Validate

// translator is a cache of locale and translation information.
var translator ut.Translator

func init() {
	// Instantiate a validator.
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Create a translator for English so the error messages are
	// more human-readable than technical.
	translator, _ = ut.New(en.New(), en.New()).GetTranslator("en")

	// Register the English error messages for use.
	en_translations.RegisterDefaultTranslations(validate, translator)

	// Use JSON tag names for errors instead of Go struct names.
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})
}
