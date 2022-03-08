package util

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeneratePaginationTotalPageGreaterThanPage(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	path := "api/v1/test-pagination"

	expectedResult := Pagination{
		Limit:       10,
		Page:        1,
		FirstUrl:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:     fmt.Sprintf("%s%s?limit=10&page=4", os.Getenv("BASE_URL"), path),
		NextUrl:     fmt.Sprintf("%s%s?limit=10&page=2", os.Getenv("BASE_URL"), path),
		PreviousUrl: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		TotalPage:   4,
	}

	trueResult := GeneratePagination(40, 10, 1, path)

	assert.Equal(t, expectedResult, trueResult)
}

func TestGeneratePaginationTotalPageLowerThanEqualPage(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	path := "api/v1/test-pagination"

	expectedResult := Pagination{
		Limit:       10,
		Page:        1,
		FirstUrl:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		NextUrl:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		PreviousUrl: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		TotalPage:   1,
	}

	trueResult := GeneratePagination(5, 10, 1, path)

	assert.Equal(t, expectedResult, trueResult)
}

func TestGeneratePaginationInPageNumberTwo(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	path := "api/v1/test-pagination"

	expectedResult := Pagination{
		Limit:       10,
		Page:        2,
		FirstUrl:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		NextUrl:     fmt.Sprintf("%s%s?limit=10&page=3", os.Getenv("BASE_URL"), path),
		PreviousUrl: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		TotalPage:   5,
	}

	trueResult := GeneratePagination(50, 10, 2, path)

	assert.Equal(t, expectedResult, trueResult)
}

func TestGeneratePaginationInLastPage(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	path := "api/v1/test-pagination"

	expectedResult := Pagination{
		Limit:       10,
		Page:        5,
		FirstUrl:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		NextUrl:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		PreviousUrl: fmt.Sprintf("%s%s?limit=10&page=4", os.Getenv("BASE_URL"), path),
		TotalPage:   5,
	}

	trueResult := GeneratePagination(50, 10, 5, path)

	assert.Equal(t, expectedResult, trueResult)
}

func TestGeneratePaginationTotalPageZero(t *testing.T) {
	t.Setenv("BASE_URL", "localhost:8080")
	path := "api/v1/test-pagination"

	expectedResult := Pagination{
		Limit:       10,
		Page:        1,
		FirstUrl:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastUrl:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		NextUrl:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		PreviousUrl: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		TotalPage:   1,
	}

	trueResult := GeneratePagination(0, 10, 1, path)

	assert.Equal(t, expectedResult, trueResult)
}