package report

import (
	"fmt"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

type Summary struct {
	Total        int            `json:"total"`
	Visible      int            `json:"visible"`
	Approved     int            `json:"approved"`
	Archived     int            `json:"archived"`
	AverageScore int            `json:"average_score"`
	ByCategory   map[string]int `json:"by_category"`
	ByLabel      map[string]int `json:"by_label"`
}

func BuildSummary(records []domain.Record) Summary {
	summary := Summary{ByCategory: make(map[string]int), ByLabel: make(map[string]int), Total: len(records)}
	totalScore := 0
	for _, record := range records {
		if record.IsVisible() {
			summary.Visible++
		}
		if record.Status == domain.StatusApproved {
			summary.Approved++
		}
		if record.Status == domain.StatusArchived {
			summary.Archived++
		}
		summary.ByCategory[record.Category]++
		summary.ByLabel[record.Label]++
		totalScore += record.Score
	}
	if len(records) > 0 {
		summary.AverageScore = totalScore / len(records)
	}
	return summary
}

func BuildReport(persistence *store.Store, includeArchived bool) (string, error) {
	records, err := persistence.ListRecords(store.Query{IncludeArchived: includeArchived})
	if err != nil {
		return "", err
	}
	summary := BuildSummary(records)
	return FormatSummary(summary), nil
}

func FormatSummary(summary Summary) string {
	categoryKeys := sortedKeys(summary.ByCategory)
	labelKeys := sortedKeys(summary.ByLabel)
	var builder strings.Builder
	fmt.Fprintf(&builder, "节目评选汇总\n总数:%d 可见:%d 已确认:%d 已归档:%d 平均分:%d\n", summary.Total, summary.Visible, summary.Approved, summary.Archived, summary.AverageScore)
	builder.WriteString("分类:")
	for _, key := range categoryKeys {
		fmt.Fprintf(&builder, "%s=%d ", key, summary.ByCategory[key])
	}
	builder.WriteString("\n标签:")
	for _, key := range labelKeys {
		fmt.Fprintf(&builder, "%s=%d ", key, summary.ByLabel[key])
	}
	return strings.TrimSpace(builder.String())
}

func sortedKeys(values map[string]int) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
