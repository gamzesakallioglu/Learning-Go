package pagination

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// Decleare default and max page size
var defaultPageSize = 10
var maxPageSize = 10
var pageParam = "page"
var pageSizeParam = "pageSize"

type Pagination struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	Sorting    string      `json:"sorting"`
	TotalRows  int64       `json:"total_rows"`
	TotalPages int         `json:"total_pages"`
	Rows       interface{} `json:"rows"`
}

func GetPaginationParametersFromRequest(g *gin.Context) (pageIndex, pageSize int) {
	pageIndex = paramsToIntOrDefault(g.Query(pageParam), 1)
	pageSize = paramsToIntOrDefault(g.Query(pageSizeParam), defaultPageSize)
	return pageIndex, pageSize
}

func paramsToIntOrDefault(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return defaultValue
}

func NewPagination(page int, pageSize int, total int64, sorting string) *Pagination {
	if page <= 0 {
		page = 0
	}

	if pageSize <= 0 {
		pageSize = defaultPageSize
	}

	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}
	pageCount := -1
	if total >= 0 {
		pageCount = (int(total) + pageSize - 1) / pageSize
	}

	return &Pagination{
		Page:       page,
		PageSize:   pageSize,
		TotalRows:  total,
		TotalPages: pageCount,
		Sorting:    sorting,
	}
}

func (p *Pagination) GetPage() int {
	if p.Page == 0 {
		p.Page = 1
	}
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	if p.PageSize == 0 {
		p.PageSize = 5
	}
	return p.PageSize
}

func (p *Pagination) GetOffset() int { // starting point
	return (p.GetPage() - 1) * p.GetPageSize()
}

func (p *Pagination) GetSorting() string {
	if p.Sorting == "" {
		p.Sorting = "ID desc"
	}
	return p.Sorting
}
