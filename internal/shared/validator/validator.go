package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
	apperrors "github.com/lion/go-framework/internal/shared/errors"
)

var validate = validator.New()

func Struct(s any) error {
	if err := validate.Struct(s); err != nil {
		var details []map[string]string
		if ve, ok := err.(validator.ValidationErrors); ok {
			for _, fe := range ve {
				details = append(details, map[string]string{
					"field":   strings.ToLower(fe.Field()),
					"tag":     fe.Tag(),
					"message": fmt.Sprintf("field '%s' failed on '%s'", strings.ToLower(fe.Field()), fe.Tag()),
				})
			}
		}
		return apperrors.Validation("Invalid request", details)
	}
	return nil
}
