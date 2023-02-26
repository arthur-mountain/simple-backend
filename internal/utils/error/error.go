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

// could set or update those error field, return self for easy way to chain
func (c *CustomError) SetCustomCode(code string) *CustomError {
	c.Code = code
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
		Code:           strconv.Itoa(httpStatusCode),
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
		errors.Is(err, gorm.ErrInvalidField): // unsupported data
		httpStatusCode = http.StatusUnprocessableEntity
	case errors.Is(err, gorm.ErrEmptySlice): // empty slice founded
		httpStatusCode = http.StatusOK
		err = errors.New("empty data was founded")
	default: // TODO: these is server error should logger
		httpStatusCode = http.StatusExpectationFailed
		err = errors.New("please contact the administrator")
	}

	return NewCustomError(err, httpStatusCode)
}
