package report

import (
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
)

type AuditSummary struct {
	Actions map[string]int
	Actors  map[string]int
	FirstAt int64
	LastAt  int64
}

func BuildAuditSummary(events []domain.AuditEvent) AuditSummary {
	summary := AuditSummary{Actions: make(map[string]int), Actors: make(map[string]int)}
	ordered := append([]domain.AuditEvent(nil), events...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].At < ordered[j].At })
	if len(ordered) > 0 {
		summary.FirstAt = ordered[0].At
		summary.LastAt = ordered[len(ordered)-1].At
	}
	for _, event := range ordered {
		summary.Actions[strings.ToLower(event.Action)]++
		summary.Actors[event.Actor]++
	}
	return summary
}

func AuditActionNames(summary AuditSummary) []string {
	result := make([]string, 0, len(summary.Actions))
	for action := range summary.Actions {
		result = append(result, action)
	}
	sort.Strings(result)
	return result
}

func AuditActorNames(summary AuditSummary) []string {
	result := make([]string, 0, len(summary.Actors))
	for actor := range summary.Actors {
		result = append(result, actor)
	}
	sort.Strings(result)
	return result
}

func AuditTotal(summary AuditSummary) int {
	total := 0
	for _, count := range summary.Actions {
		total += count
	}
	return total
}

func AuditHasAction(summary AuditSummary, action string) bool {
	_, ok := summary.Actions[strings.ToLower(action)]
	return ok
}
