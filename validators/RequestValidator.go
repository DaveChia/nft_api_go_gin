package validators

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

type ApiError struct {
    Field string `json:"field"`
    Error   string `json:"error"`
}

func GenerateSplitValidatorErrorMessages(err error) []ApiError {
	var ve validator.ValidationErrors

	errors.As(err, &ve) 
		output := make([]ApiError, len(ve))
		for i, fe := range ve {
			fmt.Println(fe.Param())
			output[i] = ApiError{fe.Field(), getMsgByValidationErrorType(fe.Tag(), fe.Field(), fe.Param())}
		}

		return output

}

func getMsgByValidationErrorType(tag string, field string, parameter string) string {
    switch tag {
		case "required":
			return "The " + field + " field is required"
		case "max":
			return "The " + field + " field may not be greater than " + parameter + " characters."
		case "min":
			return "The " + field + " field must be at least " + parameter + " characters."
		case "alphanum":
			return "The " + field + " field may only contain letters or numbers."
		}
    return "The " + field + " field is invalid."
}