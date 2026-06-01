package pagination

type OffsetPaginator struct {
	CurrentPage  int
	ItemsPerPage int
	TotalItems   int
	TotalPages   int
}

func NewOffsetPaginator(totalItems, offset, limit int) OffsetPaginator {
	if limit <= 0 {
		return OffsetPaginator{}
	}
	return OffsetPaginator{
		CurrentPage:  offset/limit + 1,
		ItemsPerPage: limit,
		TotalItems:   totalItems,
		TotalPages:   (totalItems + limit - 1) / limit,
	}
}
