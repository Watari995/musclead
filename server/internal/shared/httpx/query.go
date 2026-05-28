package httpx

import (
	"net/http"
	"strconv"
)

const (
	defaultLimit  = 20
	defaultOffset = 0
	maxLimit      = 100
)

func ParseIntOr(s string, def int) int {
	if s == "" {
		return def
	}
	// ANCII to int
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}

func ParseOffsetPagination(r *http.Request) (limit, offset int) {
	q := r.URL.Query()
	limit = ParseIntOr(q.Get("limit"), defaultLimit)
	if limit > maxLimit {
		limit = maxLimit
	}
	if limit <= 0 {
		limit = defaultLimit
	}
	offset = ParseIntOr(q.Get("offset"), defaultOffset)
	if offset < 0 {
		offset = defaultOffset
	}
	return limit, offset
}
