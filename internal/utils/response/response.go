package utils

import (
	"math"
	"net/http"
	"simple-backend/internal/interactor/page"

	"github.com/gin-gonic/gin"
)

type Response struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
	Items      any    `json:"items,omitempty"`
}

type PaginationItems struct {
	Data any `json:"data"`
	page.Pagination
}

func MakeCommonResponse(body any, code ...any) Response {
	if code != nil {
		return Response{
			StatusCode: code[0].(int),
			Message:    "success",
			Items:      body,
		}
	}

	return Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Items:      body,
	}
}

func MakePaginationResponse(body any, pageInfo page.Pagination) Response {
	items := &PaginationItems{Data: body}
	items.CurrentPage = pageInfo.CurrentPage
	items.PerPage = pageInfo.PerPage
	items.TotalCount = pageInfo.TotalCount
	items.TotalPage = int64(math.Abs(float64(pageInfo.TotalCount / pageInfo.PerPage)))

	return Response{
		StatusCode: http.StatusOK,
		Message:    "success",
		Items:      items,
	}
}

func MakeErrorResponse(c *gin.Context, code int, err error) {
	if err != nil {
		c.JSON(code, Response{StatusCode: code, Message: err.Error()})
		return
	}

	c.JSON(code, Response{StatusCode: code, Message: "something wrong"})
}
