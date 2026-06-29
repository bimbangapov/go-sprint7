package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCafeNegative(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		request string
		status  int
		message string
	}{
		{"/cafe", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=omsk", http.StatusBadRequest, "unknown city"},
		{"/cafe?city=tula&count=na", http.StatusBadRequest, "incorrect count"},
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v.request, nil)
		handler.ServeHTTP(response, req)

		assert.Equal(t, v.status, response.Code)
		assert.Equal(t, v.message, strings.TrimSpace(response.Body.String()))
	}
}

func TestCafeSearch(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		search string
		want   int
	}{
		{"фасоль", 0},
		{"кофе", 2},
		{"вилка", 1},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&search=%s", v.search)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		count := 0
		body := strings.TrimSpace(response.Body.String())
		arrayResponse := strings.Split(body, ",")

		if body == "" {
			count = 0
		} else {
			count = len(arrayResponse)
		}

		for _, cafe := range arrayResponse {
			var expected bool

			if v.want != 0 {
				expected = true
			}

			assert.Equal(t, expected, strings.Contains(strings.ToLower(cafe), strings.ToLower(v.search)), "Кафе: %s \nне содержит поисковой запрос: %s", cafe, v.search)
		}
		assert.Equal(t, v.want, count)

	}
}

func TestCafeCount(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []struct {
		count int
		want  int
	}{
		{0, 0},
		{1, 1},
		{2, 2},
		{100, 5},
	}

	for _, v := range requests {
		response := httptest.NewRecorder()
		url := fmt.Sprintf("/cafe?city=moscow&count=%d", v.count)
		req := httptest.NewRequest("GET", url, nil)
		handler.ServeHTTP(response, req)

		require.Equal(t, http.StatusOK, response.Code)

		count := 0
		body := strings.TrimSpace(response.Body.String())
		arrayResponse := strings.Split(body, ",")

		if body == "" {
			count = 0
		} else {
			count = len(arrayResponse)
		}

		minCount := min(count, v.count)

		assert.Equal(t, v.want, minCount)
	}
}

func TestCafeWhenOk(t *testing.T) {
	handler := http.HandlerFunc(mainHandle)

	requests := []string{
		"/cafe?count=2&city=moscow",
		"/cafe?city=tula",
		"/cafe?city=moscow&search=ложка",
	}
	for _, v := range requests {
		response := httptest.NewRecorder()
		req := httptest.NewRequest("GET", v, nil)

		handler.ServeHTTP(response, req)

		assert.Equal(t, http.StatusOK, response.Code)
	}
}
