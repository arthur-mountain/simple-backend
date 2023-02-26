package utils

import (
	"math"
	"net/http"
	"simple-backend/internal/interactor/page"

	"github.com/gin-gonic/gin"
)

type Response struct {
	HttpStatusCode int         `json:"-"`
	Code           int         `json:"code"`
	Message        string      `json:"message"`
	Items          interface{} `json:"items,omitempty"`
}

// return self for easy to chain set properties and put into c.JSON
func (r *Response) AppendPagination(pageInfo *page.Pagination) *Response {
	r.Items = map[string]interface{}{
		"data":        r.Items,
		"currentPage": pageInfo.CurrentPage,
		"perPage":     pageInfo.PerPage,
		"totalCount":  pageInfo.TotalCount,
		"totalPage":   int64(math.Abs(float64(pageInfo.TotalCount / pageInfo.PerPage))),
	}
	return r
}
func (r *Response) SetHttpCode(httpCode int) *Response {
	r.HttpStatusCode = httpCode
	return r
}
func (r *Response) SetCode(code int) *Response {
	r.Code = code
	return r
}
func (r *Response) SetMessage(message string) *Response {
	r.Message = message
	return r
}
func (r *Response) Done(c *gin.Context) {
	if r.Code == 0 {
		r.Code = r.HttpStatusCode
	}
	c.JSON(r.HttpStatusCode, r)
}

func New(body interface{}) *Response {
	return &Response{
		HttpStatusCode: http.StatusOK,
		Message:        "success",
		Items:          body,
	}
}
