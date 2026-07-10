package request

type PaginatedRequest struct {
	Page    int    `query:"page"`
	PerPage int    `query:"per_page"`
	Search  string `query:"search"`
}
