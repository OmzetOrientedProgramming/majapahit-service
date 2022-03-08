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
		FirstUrl:  getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1),
		LastUrl:   getPaginationURL(os.Getenv("BASE_URL"), path, limit, totalPage),
		TotalPage: totalPage,
	}

	if limit < totalCount {
		if page == 1 {
			pagination.PreviousUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
			pagination.NextUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page+1)
		} else if page == totalPage {
			pagination.PreviousUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page-1)
			pagination.NextUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, totalPage)
		} else {
			pagination.PreviousUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page-1)
			pagination.NextUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page+1)
		}
	} else {
		pagination.NextUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
		pagination.PreviousUrl = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
	}

	return pagination
}

func getPaginationURL(baseURL, path string, limit, page int) string {
	return fmt.Sprintf("%s%s?limit=%d&page=%d", baseURL, path, limit, page)
}