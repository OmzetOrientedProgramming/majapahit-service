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
		FirstURL:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastURL:     fmt.Sprintf("%s%s?limit=10&page=4", os.Getenv("BASE_URL"), path),
		NextURL:     fmt.Sprintf("%s%s?limit=10&page=2", os.Getenv("BASE_URL"), path),
		PreviousURL: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
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
		FirstURL:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastURL:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		NextURL:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		PreviousURL: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
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
		FirstURL:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastURL:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		NextURL:     fmt.Sprintf("%s%s?limit=10&page=3", os.Getenv("BASE_URL"), path),
		PreviousURL: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
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
		FirstURL:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastURL:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		NextURL:     fmt.Sprintf("%s%s?limit=10&page=5", os.Getenv("BASE_URL"), path),
		PreviousURL: fmt.Sprintf("%s%s?limit=10&page=4", os.Getenv("BASE_URL"), path),
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
		FirstURL:    fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		LastURL:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		NextURL:     fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		PreviousURL: fmt.Sprintf("%s%s?limit=10&page=1", os.Getenv("BASE_URL"), path),
		TotalPage:   1,
	}

	trueResult := GeneratePagination(0, 10, 1, path)

	assert.Equal(t, expectedResult, trueResult)
}

func TestValidateParams(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		pageString := "10"
		limitString := "1"
		var expectedErrorList []string

		page, limit, errorList := ValidateParams(pageString, limitString)

		assert.Equal(t, 10, page)
		assert.Equal(t, 1, limit)
		assert.Equal(t, expectedErrorList, errorList)
	})

	t.Run("limit and page not valid", func(t *testing.T) {
		pageString := "test"
		limitString := "test"
		expectedErrorList := []string{
			"limit should be positive integer",
			"page should be positive integer",
		}

		page, limit, errorList := ValidateParams(pageString, limitString)

		assert.Equal(t, 0, page)
		assert.Equal(t, 0, limit)
		assert.Equal(t, expectedErrorList, errorList)
	})

	t.Run("limit and page not is not inputed", func(t *testing.T) {
		pageString := ""
		limitString := ""
		var expectedErrorList []string

		page, limit, errorList := ValidateParams(pageString, limitString)

		assert.Equal(t, 1, page)
		assert.Equal(t, 10, limit)
		assert.Equal(t, expectedErrorList, errorList)
	})
}
