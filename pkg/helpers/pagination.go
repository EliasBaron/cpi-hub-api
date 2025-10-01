package helpers

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetPaginationValues(c *gin.Context) (int, int) {
	page := c.Query("page")
	pageSize := c.Query("page_size")

	pageInt := 1
	pageSizeInt := 25

	if page != "" {
		if p, err := parseInt(page); err == nil && p > 0 {
			pageInt = p
		}
	}

	if pageSize != "" {
		if ps, err := parseInt(pageSize); err == nil && ps > 0 {
			pageSizeInt = ps
		}
	}

	return pageInt, pageSizeInt
}

func GetSortValues(c *gin.Context) (string, string) {
	orderBy := c.Query("order_by")
	sortDirection := c.Query("sort_direction")

	if orderBy == "" {
		orderBy = "created_at"
	}

	if sortDirection != "asc" && sortDirection != "desc" {
		sortDirection = "desc"
	}

	return orderBy, sortDirection
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}
