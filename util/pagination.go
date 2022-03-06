package util

import (
	"fmt"
	"math"
	"os"
)

type Pagination struct {
	Limit       int    `json:"limit"`
	Page        int    `json:"page"`
	FirstUrl    string `json:"first_url"`
	LastUrl     string `json:"last_url"`
	NextUrl     string `json:"next_url"`
	PreviousUrl string `json:"previous_url"`
	TotalPage   int    `json:"total_page"`
}

func GeneratePagination(totalCount, limit, page int, path string) Pagination {
	totalPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	pagination := Pagination{
		Limit:     limit,
		Page:      page,
		FirstUrl:  fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:   fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, totalPage),
		TotalPage: totalPage,
	}

	if limit < totalCount {
		if page == 1 {
			pagination.PreviousUrl = fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path)
			pagination.NextUrl = fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, page+1)
		} else if page == totalPage {
			pagination.PreviousUrl = fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, page-1)
			pagination.NextUrl = fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, totalPage)
		} else {
			pagination.PreviousUrl = fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, page-1)
			pagination.NextUrl = fmt.Sprintf("%s%s?limit=10&page=%d", os.Getenv("BASE_URL"), path, page+1)
		}
	} else {
		pagination.NextUrl = fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path)
		pagination.PreviousUrl = fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path)
	}

	return pagination
}
