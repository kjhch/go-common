package restful

type Pagination struct {
	CurrentPage *int  `json:"current-page"`
	PageSize    *int  `json:"page-size"`
	PageCount   *int  `json:"page-count"`
	Total       *int  `json:"total"`
	Data        []any `json:"data"`
}
