package util

import (
	"fmt"
	"math"
	"os"
	"strconv"
)

// Pagination struct for pagination data response
type Pagination struct {
	Limit       int    `json:"limit"`
	Page        int    `json:"page"`
	FirstURL    string `json:"first_url"`
	LastURL     string `json:"last_url"`
	NextURL     string `json:"next_url"`
	PreviousURL string `json:"previous_url"`
	TotalPage   int    `json:"total_page"`
}

// GeneratePagination function will generate the pagination given the parameter
func GeneratePagination(totalCount, limit, page int, path string) Pagination {
	totalPage := int(math.Ceil(float64(totalCount) / float64(limit)))
	if totalPage == 0 {
		totalPage = 1
	}

	pagination := Pagination{
		Limit:     limit,
		Page:      page,
		FirstURL:  getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1),
		LastURL:   getPaginationURL(os.Getenv("BASE_URL"), path, limit, totalPage),
		TotalPage: totalPage,
	}

	if limit < totalCount {
		if page == 1 {
			pagination.PreviousURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
			pagination.NextURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page+1)
		} else if page == totalPage {
			pagination.PreviousURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page-1)
			pagination.NextURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, totalPage)
		} else {
			pagination.PreviousURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page-1)
			pagination.NextURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, page+1)
		}
	} else {
		pagination.NextURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
		pagination.PreviousURL = getPaginationURL(os.Getenv("BASE_URL"), path, limit, 1)
	}

	return pagination
}

func getPaginationURL(baseURL, path string, limit, page int) string {
	return fmt.Sprintf("%s%s?limit=%d&page=%d", baseURL, path, limit, page)
}

// ValidateParams for validating pagination params
func ValidateParams(pageString string, limitString string) (int, int, []string) {
	var (
		limit, page int
		errorList   []string
		err         error
	)

	limit, err = strconv.Atoi(limitString)
	if err != nil {
		if limitString == "" {
			limit = DefaultLimit
		} else {
			errorList = append(errorList, "limit should be positive integer")
		}
	}

	page, err = strconv.Atoi(pageString)
	if err != nil {
		if pageString == "" {
			page = DefaultPage
		} else {
			errorList = append(errorList, "page should be positive integer")
		}
	}

	return page, limit, errorList
}
