package utils

import (
	"math"
	"net/http"
	"simple-backend/internal/interactor/page"
	"strconv"
)

type Response struct {
	HttpStatusCode int         `json:"-"`
	Code           string      `json:"code"`
	Message        string      `json:"message"`
	Items          interface{} `json:"items,omitempty"`
}

func (r *Response) AppendPagination(pageInfo *page.Pagination) *Response {
	r.Items = map[string]interface{}{
		"data":        r.Items,
		"currentPage": pageInfo.CurrentPage,
		"perPage":     pageInfo.PerPage,
		"totalCount":  pageInfo.TotalCount,
		"totalPage":   int64(math.Abs(float64(pageInfo.TotalCount / pageInfo.PerPage))),
	}

	return r // return self for easy to put into c.JSON
}

func MakeCommonResponse(body interface{}, httpStatusCode *int, code *string, message *string) *Response {
	res := &Response{
		HttpStatusCode: http.StatusOK,
		Code:           strconv.Itoa(http.StatusOK),
		Items:          body,
	}

	if httpStatusCode != nil {
		res.HttpStatusCode = *httpStatusCode
		res.Code = strconv.Itoa(*httpStatusCode)
	}
	if code != nil {
		res.Code = *code
	}
	if message != nil {
		res.Message = *message
	}

	return res
}
