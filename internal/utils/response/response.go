package utils

import (
	"math"
	"net/http"
	"simple-backend/internal/interactor/page"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode string      `json:"status_code"`
	Message    string      `json:"message"`
	Items      interface{} `json:"items,omitempty"`
}

type PaginationItems struct {
	Data interface{} `json:"data"`
	page.Pagination
}

func MakeCommonResponse(body interface{}, code ...any) Response {
	if len(code) > 0 {
		return Response{
			StatusCode: strconv.Itoa(code[0].(int)),
			Message:    "success",
			Items:      body,
		}
	}

	return Response{
		StatusCode: strconv.Itoa(http.StatusOK),
		Message:    "success",
		Items:      body,
	}
}

func MakePaginationResponse(body interface{}, pageInfo page.Pagination) Response {
	items := &PaginationItems{Data: body}
	items.CurrentPage = pageInfo.CurrentPage
	items.PerPage = pageInfo.PerPage
	items.TotalCount = pageInfo.TotalCount
	items.TotalPage = int64(math.Abs(float64(pageInfo.TotalCount / pageInfo.PerPage)))

	return Response{
		StatusCode: strconv.Itoa(http.StatusOK),
		Message:    "success",
		Items:      items,
	}
}

func MakeErrorResponse(c *gin.Context, code int, err error) {
	if err != nil {
		c.JSON(code, Response{StatusCode: strconv.Itoa(code), Message: err.Error()})
		return
	}

	c.JSON(code, Response{StatusCode: strconv.Itoa(code), Message: "something wrong"})
}
