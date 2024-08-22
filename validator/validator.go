package validator

import (
	"fmt"
	"strings"

	pkgError "github.com/adityaeka26/go-pkg/error"

	"github.com/go-playground/validator/v10"
)

type XValidator struct {
	Validator *validator.Validate
}

var validate = validator.New()

func (v XValidator) Validate(data any) error {
	var message string
	if errs := validate.Struct(data); errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			message = message + strings.ToLower(fmt.Sprintf("%s %s;", err.Field(), err.Tag()))
		}
		return pkgError.BadRequest(message)
	}

	return nil
}
