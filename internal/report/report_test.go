package report

import (
	"strings"
	"testing"

	"example.com/graduation-showcase/internal/domain"
)

func TestReportSummary(t *testing.T) {
	records := []domain.Record{domain.NewRecord("r1", "甲", "a", "声乐", "推荐", 1), domain.NewRecord("r2", "乙", "b", "舞蹈", "候选", 2)}
	records[0].Score = 90
	records[1].Score = 80
	summary := BuildSummary(records)
	if summary.Total != 2 || summary.AverageScore != 85 {
		t.Fatalf("summary: %#v", summary)
	}
	text := FormatSummary(summary)
	if !strings.Contains(text, "声乐=1") || !strings.Contains(text, "推荐=1") {
		t.Fatalf("report: %s", text)
	}
}
