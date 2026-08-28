package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

func EncodeRows(rows []Row) []string {
	ordered := append([]Row(nil), rows...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	lines := make([]string, 0, len(ordered))
	for _, row := range ordered {
		lines = append(lines, strings.Join([]string{row.ID, row.ProgramName, row.Performer, row.Category, row.Label, row.Filename, row.Content}, "|"))
	}
	return lines
}

func DigestRows(rows []Row) string {
	digest := sha256.Sum256([]byte(strings.Join(EncodeRows(rows), "\n")))
	return hex.EncodeToString(digest[:])
}

func RowsFromRecords(records []domain.Record) []Row {
	result := make([]Row, 0, len(records))
	for _, record := range records {
		result = append(result, Row{ID: record.ID, ProgramName: record.ProgramName, Performer: record.Performer, Category: record.Category, Label: record.Label})
	}
	return result
}

func MergeRows(left, right []Row) []Row {
	result := append([]Row(nil), left...)
	index := make(map[string]int)
	for i, row := range result {
		index[row.ID] = i
	}
	for _, row := range right {
		if position, ok := index[row.ID]; ok {
			result[position] = row
		} else {
			index[row.ID] = len(result)
			result = append(result, row)
		}
	}
	return result
}
