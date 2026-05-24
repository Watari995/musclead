package pagination

type OffsetPaginator struct {
	CurrentPage  int
	ItemsPerPage int
	TotalItems   int
	TotalPages   int
}
