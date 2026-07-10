package response

import "math"

type Paginated struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PerPage    int         `json:"per_page"`
	TotalPages int         `json:"total_pages"`
}

func Paginate(page, perPage int, total int64) Paginated {
	if perPage < 1 || perPage > 100 {
		perPage = 20
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	return Paginated{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
