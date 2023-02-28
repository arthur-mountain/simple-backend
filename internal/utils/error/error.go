package utils

import (
	"errors"
	"net/http"
	"strings"

	errorConstant "simple-backend/internal/constants/error"

	"gorm.io/gorm"
)

// TODO: implement reason struct, that allowed to constraints reason field
// type TReason struct{}

// Message for simple error description, Reasons for complex error description
type CustomError struct {
	HttpStatusCode int                     `json:"-"`
	Code           errorConstant.ErrorCode `json:"code"`
	Message        string                  `json:"message"`
	Reasons        []*interface{}          `json:"reasons,omitempty"`
}

// implement error interface
func (c *CustomError) Error() string {
	return c.Message
}

// could set or update those error field, return self for easy way to chain
func (c *CustomError) SetCode(code errorConstant.ErrorCode) *CustomError {
	c.Code = code
	// Use Mapping table to set message, otherwise could use SetMessage method
	if message, ok := errorConstant.ErrorCodeMapping[code]; ok {
		c.Message = message
	}
	return c
}
func (c *CustomError) SetMessage(message string) *CustomError {
	c.Message = message
	return c
}
func (c *CustomError) AppendReason(reason *interface{}) *CustomError {
	c.Reasons = append(c.Reasons, reason)
	return c
}

func NewCustomError(err error, httpStatusCode int) *CustomError {
	return &CustomError{
		HttpStatusCode: httpStatusCode,
		Code:           errorConstant.ErrorCode(httpStatusCode),
		Message:        err.Error(),
	}
}

func CheckErrAndConvert(err error, httpStatusCode int) *CustomError {
	if err != nil {
		return NewCustomError(err, httpStatusCode)
	}

	return nil
}

// Database error, please use in repository, do not use in service or deliver
// If database changing, only to add new database error check func
func CheckRepoError(err error) *CustomError {
	var httpStatusCode int

	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		httpStatusCode = http.StatusNotFound
	case
		errors.Is(err, gorm.ErrInvalidData),
		errors.Is(err, gorm.ErrInvalidField):
		// unsupported data
		httpStatusCode = http.StatusUnprocessableEntity
	case errors.Is(err, gorm.ErrEmptySlice):
		// empty slice founded
		httpStatusCode = http.StatusNoContent
		err = errors.New("empty data was founded")
	case strings.HasPrefix(err.Error(), "Error "):
		// like "Error 1062 duplicate entry"
		httpStatusCode = http.StatusBadRequest
	default:
		panic(err) // will be recovered by log middleware
	}

	return NewCustomError(err, httpStatusCode)
}
