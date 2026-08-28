package api

import (
	"encoding/json"
	"net/http"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

type RecordResponse struct {
	ID          string `json:"id"`
	ProgramName string `json:"program_name"`
	Performer   string `json:"performer"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Status      string `json:"status"`
	Score       int    `json:"score"`
	Version     int    `json:"version"`
}

func NewRecordResponse(record domain.Record) RecordResponse {
	return RecordResponse{ID: record.ID, ProgramName: record.ProgramName, Performer: record.Performer, Category: record.Category, Label: record.Label, Status: string(record.Status), Score: record.Score, Version: record.Version}
}

func NewRecordResponses(records []domain.Record) []RecordResponse {
	result := make([]RecordResponse, 0, len(records))
	for _, record := range records {
		result = append(result, NewRecordResponse(record))
	}
	return result
}

func WriteRecords(w http.ResponseWriter, records []domain.Record, page Page) {
	responses := NewRecordResponses(records)
	writeJSON(w, http.StatusOK, PaginateResponses(responses, page))
}

func PaginateResponses(records []RecordResponse, page Page) []RecordResponse {
	start := page.Offset
	if start > len(records) {
		start = len(records)
	}
	end := start + page.Limit
	if end > len(records) {
		end = len(records)
	}
	return records[start:end]
}

func SortResponses(records []RecordResponse) []RecordResponse {
	result := append([]RecordResponse(nil), records...)
	sort.Slice(result, func(i, j int) bool {
		if result[i].ProgramName != result[j].ProgramName {
			return result[i].ProgramName < result[j].ProgramName
		}
		return result[i].ID < result[j].ID
	})
	return result
}

func EncodeResponse(value any) ([]byte, error) { return json.Marshal(value) }

func NormalizeHeader(value string) string { return strings.TrimSpace(strings.ToLower(value)) }

func AcceptsJSON(request *http.Request) bool {
	header := NormalizeHeader(request.Header.Get("Accept"))
	return header == "" || strings.Contains(header, "application/json")
}
