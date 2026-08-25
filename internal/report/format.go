package report

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

func RenderCSV(records []domain.Record) (string, error) {
	ordered := domain.RankRecords(records)
	var builder strings.Builder
	writer := csv.NewWriter(&builder)
	if err := writer.Write([]string{"id", "program", "performer", "category", "label", "status", "score"}); err != nil {
		return "", err
	}
	for _, record := range ordered {
		if err := writer.Write([]string{record.ID, record.ProgramName, record.Performer, record.Category, record.Label, string(record.Status), fmt.Sprintf("%d", record.Score)}); err != nil {
			return "", err
		}
	}
	writer.Flush()
	return strings.TrimSpace(builder.String()), writer.Error()
}

func RenderMarkdown(records []domain.Record) string {
	ordered := domain.SortByProgram(records)
	var builder strings.Builder
	builder.WriteString("| 节目 | 表演者 | 分类 | 标签 | 分数 |\n|---|---|---|---|---|\n")
	for _, record := range ordered {
		fmt.Fprintf(&builder, "| %s | %s | %s | %s | %d |\n", record.ProgramName, record.Performer, record.Category, record.Label, record.Score)
	}
	return strings.TrimSpace(builder.String())
}

func ValidateSummary(summary Summary) error {
	if summary.Total < 0 || summary.Visible < 0 || summary.Approved < 0 || summary.Archived < 0 {
		return fmt.Errorf("negative summary value")
	}
	if summary.Visible+summary.Archived > summary.Total {
		return fmt.Errorf("summary totals do not balance")
	}
	return nil
}

func MergeSummaries(left, right Summary) Summary {
	merged := Summary{Total: left.Total + right.Total, Visible: left.Visible + right.Visible, Approved: left.Approved + right.Approved, Archived: left.Archived + right.Archived, ByCategory: make(map[string]int), ByLabel: make(map[string]int)}
	for key, value := range left.ByCategory {
		merged.ByCategory[key] += value
	}
	for key, value := range right.ByCategory {
		merged.ByCategory[key] += value
	}
	for key, value := range left.ByLabel {
		merged.ByLabel[key] += value
	}
	for key, value := range right.ByLabel {
		merged.ByLabel[key] += value
	}
	if merged.Total > 0 {
		merged.AverageScore = (left.AverageScore*left.Total + right.AverageScore*right.Total) / merged.Total
	}
	return merged
}

func SortedCategories(records []domain.Record) []string {
	set := make(map[string]bool)
	for _, record := range records {
		set[record.Category] = true
	}
	result := make([]string, 0, len(set))
	for key := range set {
		result = append(result, key)
	}
	sort.Strings(result)
	return result
}
