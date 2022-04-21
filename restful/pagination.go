package restful

type Pagination struct {
	CurrentPage *int  `json:"current_page"`
	PageSize    *int  `json:"page_size"`
	PageCount   *int  `json:"page_count"`
	Total       *int  `json:"total"`
	Data        []any `json:"data"`
}
