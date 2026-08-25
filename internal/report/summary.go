package report

import (
	"sort"

	"example.com/graduation-showcase/internal/domain"
)

type CategoryResult struct {
	Category string
	Count    int
	Average  int
}

func ByCategory(records []domain.Record) []CategoryResult {
	type aggregate struct {
		count int
		score int
	}
	groups := make(map[string]aggregate)
	for _, record := range records {
		value := groups[record.Category]
		value.count++
		value.score += record.Score
		groups[record.Category] = value
	}
	result := make([]CategoryResult, 0, len(groups))
	for category, value := range groups {
		average := 0
		if value.count > 0 {
			average = value.score / value.count
		}
		result = append(result, CategoryResult{Category: category, Count: value.count, Average: average})
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Average != result[j].Average {
			return result[i].Average > result[j].Average
		}
		return result[i].Category < result[j].Category
	})
	return result
}

func TopRecords(records []domain.Record, limit int) []domain.Record {
	if limit < 1 {
		return []domain.Record{}
	}
	ranked := domain.RankRecords(records)
	if len(ranked) < limit {
		return ranked
	}
	return ranked[:limit]
}

func StatusCounts(records []domain.Record) map[domain.Status]int {
	counts := make(map[domain.Status]int)
	for _, record := range records {
		counts[record.Status]++
	}
	return counts
}
