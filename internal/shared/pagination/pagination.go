package pagination

import "strconv"

type Params struct {
	Page  int
	Limit int
}

type Result[T any] struct {
	Items []T
	Total int64
	Page  int
	Limit int
}

func Parse(pageStr, limitStr string, defaultLimit, maxLimit int) Params {
	page := 1
	limit := defaultLimit

	if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
		page = v
	}
	if v, err := strconv.Atoi(limitStr); err == nil && v > 0 {
		limit = v
	}
	if limit > maxLimit {
		limit = maxLimit
	}
	return Params{Page: page, Limit: limit}
}

func (p Params) Offset() int {
	return (p.Page - 1) * p.Limit
}
