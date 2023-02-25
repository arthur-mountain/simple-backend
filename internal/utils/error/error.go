package utils

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type CustomError struct {
	HttpStatusCode int            `json:"-"`
	Code           string         `json:"code"`
	Message        string         `json:"message"`
	Reasons        []*interface{} `json:"reasons,omitempty"` // error reason slice or nil
}

func (e *CustomError) AppendReason(reason *interface{}) *CustomError {
	e.Reasons = append(e.Reasons, reason)
	return e // return self for easy way to chain append reason
}

func NewErrorResponse(
	httpStatusCode int,
	message string,
	code *string,
) *CustomError {
	customError := &CustomError{
		HttpStatusCode: httpStatusCode,
		Code:           *code,
		Message:        message,
	}

	if customError.Code == "" {
		customError.Code = strconv.Itoa(customError.HttpStatusCode)
	}

	return customError
}

func CheckErrAndConvert(err error, httpStatusCode int, code *string, message *string) *CustomError {
	if err != nil {
		if message != nil {
			return NewErrorResponse(httpStatusCode, *message, code)
		}

		return NewErrorResponse(httpStatusCode, err.Error(), code)
	}

	return nil
}

// Database error, please use in repository, do not use in service or deliver
// If database changing, only to add new database error check func
func CheckGormError(err error) *CustomError {
	var httpStatusCode int
	var code string
	var message string

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		httpStatusCode = http.StatusNotFound
	case errors.Is(err, gorm.ErrInvalidData): // unsupported data
	case errors.Is(err, gorm.ErrInvalidField): // invalid field
	case errors.Is(err, gorm.ErrEmptySlice): // empty slice found
	default: // TODO: these is server error should logger
		httpStatusCode = http.StatusExpectationFailed
		message = "please contact the administrator"
	}

	return NewErrorResponse(httpStatusCode, message, &code)
}
