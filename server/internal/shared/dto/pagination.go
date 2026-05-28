package shareddto

import "github.com/Watari995/musclead/internal/pagination"

type PaginationDTO struct {
	CurrentPage  int `json:"current_page"`
	ItemsPerPage int `json:"items_per_page"`
	TotalItems   int `json:"total_items"`
	TotalPages   int `json:"total_pages"`
}

func NewPaginationDTO(paginator pagination.OffsetPaginator) PaginationDTO {
	return PaginationDTO{
		CurrentPage:  paginator.CurrentPage,
		ItemsPerPage: paginator.ItemsPerPage,
		TotalItems:   paginator.TotalItems,
		TotalPages:   paginator.TotalPages,
	}
}
