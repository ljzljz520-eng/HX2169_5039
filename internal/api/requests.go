package api

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

type Page struct {
	Offset int
	Limit  int
}

func ParsePage(request *http.Request) Page {
	offset, _ := strconv.Atoi(request.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(request.URL.Query().Get("limit"))
	if offset < 0 {
		offset = 0
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	return Page{Offset: offset, Limit: limit}
}

func ParseQuery(request *http.Request) QueryInput {
	query := request.URL.Query()
	return QueryInput{Text: strings.TrimSpace(query.Get("q")), Label: strings.TrimSpace(query.Get("label")), Status: domain.Status(strings.TrimSpace(query.Get("status"))), IncludeArchived: query.Get("archived") == "true"}
}

func ValidateCreateInput(input CreateInput) error {
	if strings.TrimSpace(input.ID) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(input.ProgramName) == "" {
		return errors.New("program name is required")
	}
	if strings.TrimSpace(input.Performer) == "" {
		return errors.New("performer is required")
	}
	if strings.TrimSpace(input.Category) == "" {
		return errors.New("category is required")
	}
	if strings.TrimSpace(input.Label) == "" {
		return errors.New("label is required")
	}
	return nil
}

func ValidateUpdateInput(input UpdateInput) error {
	if strings.TrimSpace(input.Label) == "" {
		return errors.New("label is required")
	}
	return nil
}

func ClampPage(records int, page Page) (int, int) {
	start := page.Offset
	if start > records {
		start = records
	}
	end := start + page.Limit
	if end > records {
		end = records
	}
	return start, end
}

func PaginateRecords(records []domain.Record, page Page) []domain.Record {
	start, end := ClampPage(len(records), page)
	return records[start:end]
}
