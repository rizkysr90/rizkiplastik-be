package util

type PaginationData struct {
	PageNumber    int `json:"page_number"`
	PageSize      int `json:"page_size"`
	TotalPages    int `json:"total_pages"`
	TotalElements int `json:"total_elements"`
}

func NewPaginationData(
	pageNumber string,
	pageSize string,
) (*PaginationData, error) {
	if pageNumber == "" {
		pageNumber = "1"
	}
	if pageSize == "" {
		pageSize = "10"
	}
	pageNumberInt, err := ConvertStringToInt(pageNumber)
	if err != nil {
		return nil, err
	}
	pageSizeInt, err := ConvertStringToInt(pageSize)
	if err != nil {
		return nil, err
	}
	return &PaginationData{
		PageNumber: pageNumberInt,
		PageSize:   pageSizeInt,
	}, nil
}

func (p *PaginationData) GetOffset() int {
	return (p.PageNumber - 1) * p.PageSize
}
func (p *PaginationData) SetTotalPagesAndTotalElement(totalElement int) {
	p.TotalElements = totalElement
	p.TotalPages = totalElement / p.PageSize
	if totalElement%p.PageSize > 0 {
		p.TotalPages++
	}
}
