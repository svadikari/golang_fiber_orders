package middleware

import (
	"log/slog"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	StructValidator struct {
		validate *validator.Validate
	}

	GlobalErrorHandlerResp struct {
		Code   int      `json:"code"`
		Errors []string `json:"errors"`
	}
)

var validate = validator.New()

func NewStructValidator() *StructValidator {
	validate.RegisterTagNameFunc(func(fld reflect.StructField) string {

		name := strings.SplitN(fld.Tag.Get("message"), ",", 2)[0]
		if name == "" {
			return fld.Name
		}
		return name
	})
	return &StructValidator{validate: validate}
}

func (v StructValidator) Validate(data interface{}) []string {
	validationErrors := []string{}

	errs := validate.Struct(data)
	slog.Info("Validation errors: ", "errors", errs)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			validationErrors = append(validationErrors, err.Field())
		}
	}
	return validationErrors
}
