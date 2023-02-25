package utils

import (
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

// TODO: implement reason struct, that allowed to constraints reason field
// type TReason struct{}

// simple error description could use message field
// complex error description could use reasons field
type CustomError struct {
	HttpStatusCode int            `json:"-"`
	Code           string         `json:"code"`
	Message        string         `json:"message"`
	Reasons        []*interface{} `json:"reasons,omitempty"`
}

func (e *CustomError) AppendReason(reason *interface{}) *CustomError {
	e.Reasons = append(e.Reasons, reason)
	return e // return self for easy way to chain append reason
}

func NewCustomError(
	httpStatusCode int,
	message string,
	code *string,
) *CustomError {
	customError := &CustomError{
		HttpStatusCode: httpStatusCode,
		Code:           strconv.Itoa(httpStatusCode),
		Message:        message,
	}

	if code != nil {
		customError.Code = *code
	}

	return customError
}

func CheckErrAndConvert(err error, httpStatusCode int, code *string, message *string) *CustomError {
	if err != nil {
		if message != nil {
			return NewCustomError(httpStatusCode, *message, code)
		}

		return NewCustomError(httpStatusCode, err.Error(), code)
	}

	return nil
}

// Database error, please use in repository, do not use in service or deliver
// If database changing, only to add new database error check func
func CheckGormError(err error) *CustomError {
	var httpStatusCode int
	var code, message string

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

	return NewCustomError(httpStatusCode, message, &code)
}
