package restful

type Pagination struct {
	CurrentPage *int  `json:"currentPage"`
	PageSize    *int  `json:"pageSize"`
	PageCount   *int  `json:"pageCount"`
	Total       *int  `json:"total"`
	HasMore     *bool `json:"hasMore"`
	Data        []any `json:"data"`
}
