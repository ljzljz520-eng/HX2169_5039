package importer

import (
	"sort"
	"strings"
)

type Issue struct {
	Row     int
	Field   string
	Message string
}

func ValidateRows(rows []Row) []Issue {
	issues := make([]Issue, 0)
	seen := make(map[string]int)
	for index, row := range rows {
		if previous, ok := seen[row.ID]; ok {
			issues = append(issues, Issue{Row: index, Field: "id", Message: "duplicate of row " + string(rune(previous+1))})
		}
		seen[row.ID] = index
		if strings.TrimSpace(row.ProgramName) == "" {
			issues = append(issues, Issue{Row: index, Field: "program_name", Message: "required"})
		}
		if strings.TrimSpace(row.Performer) == "" {
			issues = append(issues, Issue{Row: index, Field: "performer", Message: "required"})
		}
		if strings.TrimSpace(row.Content) == "" {
			issues = append(issues, Issue{Row: index, Field: "content", Message: "required"})
		}
	}
	sort.Slice(issues, func(i, j int) bool {
		if issues[i].Row != issues[j].Row {
			return issues[i].Row < issues[j].Row
		}
		return issues[i].Field < issues[j].Field
	})
	return issues
}

func ValidBatch(rows []Row) bool { return len(rows) > 0 && len(ValidateRows(rows)) == 0 }

func CountValid(rows []Row) int {
	count := 0
	for _, row := range rows {
		if ValidateRow(row) == nil {
			count++
		}
	}
	return count
}
