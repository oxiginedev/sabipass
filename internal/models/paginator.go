package models

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Paginator struct {
	PerPage int64
	Page    int64
}

func (p Paginator) Offset() int64 {
	if p.Page <= 0 {
		return 0
	}

	return (p.Page - 1) * p.PerPage
}

func PaginatorFromContext(c *gin.Context) Paginator {
	perPage := c.DefaultQuery("per_page", "10")
	perPageInt, err := strconv.ParseInt(perPage, 10, 64)
	if err != nil || perPageInt < 1 {
		perPageInt = 10
	}

	page := c.DefaultQuery("page", "1")
	pageInt, err := strconv.ParseInt(page, 10, 64)
	if err != nil || pageInt < 1 {
		pageInt = 1
	}

	return Paginator{
		PerPage: perPageInt,
		Page:    pageInt,
	}
}

type PaginationMeta struct {
	TotalCount int64 `json:"total_count"`
	Page       int64 `json:"page"`
	PerPage    int64 `json:"per_page"`
}

type PaginatedResponse struct {
	Response
	Meta PaginationMeta `json:"meta"`
}

func NewPaginatedResponse(msg string, data any, totalCount int64, paginator Paginator) PaginatedResponse {
	return PaginatedResponse{
		Response: newResponse(true, msg, data, nil),
		Meta: PaginationMeta{
			TotalCount: totalCount,
			Page:       paginator.Page,
			PerPage:    paginator.PerPage,
		},
	}
}
