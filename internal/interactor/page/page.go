package page

type Pagination struct {
	CurrentPage int64 `form:"current_page" json:"current_page" validate:"gt=0"`
	PerPage     int64 `form:"per_page" json:"per_page" validate:"gt=0"`
	TotalCount  int64 `json:"total_count"`
	TotalPage   int64 `json:"total_page"`
}
