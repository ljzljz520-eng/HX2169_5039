package domain

import (
	"sort"
	"strings"
)

type RecordFilter struct {
	Text         string
	Categories   []string
	Labels       []string
	Statuses     []Status
	MinScore     int
	MaxScore     int
	ShowArchived bool
}

func (f RecordFilter) Matches(record Record) bool {
	if !f.ShowArchived && record.Status == StatusArchived {
		return false
	}
	if f.Text != "" {
		text := strings.ToLower(strings.TrimSpace(f.Text))
		if !strings.Contains(strings.ToLower(record.ProgramName), text) && !strings.Contains(strings.ToLower(record.Performer), text) && !strings.Contains(strings.ToLower(record.ID), text) {
			return false
		}
	}
	if len(f.Categories) > 0 && !containsString(f.Categories, record.Category) {
		return false
	}
	if len(f.Labels) > 0 && !containsStringFold(f.Labels, record.Label) {
		return false
	}
	if len(f.Statuses) > 0 && !containsStatus(f.Statuses, record.Status) {
		return false
	}
	if record.Score < f.MinScore {
		return false
	}
	if f.MaxScore > 0 && record.Score > f.MaxScore {
		return false
	}
	return true
}

func FilterRecords(records []Record, filter RecordFilter) []Record {
	result := make([]Record, 0, len(records))
	for _, record := range records {
		if filter.Matches(record) {
			result = append(result, record)
		}
	}
	return result
}

func SortByProgram(records []Record) []Record {
	result := append([]Record(nil), records...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].ProgramName != result[j].ProgramName {
			return result[i].ProgramName < result[j].ProgramName
		}
		return result[i].ID < result[j].ID
	})
	return result
}

func GroupByCategory(records []Record) map[string][]Record {
	groups := make(map[string][]Record)
	for _, record := range records {
		groups[record.Category] = append(groups[record.Category], record)
	}
	for key := range groups {
		groups[key] = SortByProgram(groups[key])
	}
	return groups
}

func UniqueLabels(records []Record) []string {
	seen := make(map[string]bool)
	for _, record := range records {
		seen[record.Label] = true
	}
	labels := make([]string, 0, len(seen))
	for label := range seen {
		labels = append(labels, label)
	}
	sort.Strings(labels)
	return labels
}

func containsString(values []string, target string) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}
func containsStringFold(values []string, target string) bool {
	for _, value := range values {
		if strings.EqualFold(value, target) {
			return true
		}
	}
	return false
}
func containsStatus(values []Status, target Status) bool {
	for _, value := range values {
		if value == target {
			return true
		}
	}
	return false
}

type TimelineEntry struct {
	At      int64
	Label   string
	Status  Status
	Version int
}

func BuildTimeline(records []Record) []TimelineEntry {
	entries := make([]TimelineEntry, 0, len(records))
	for _, record := range records {
		entries = append(entries, TimelineEntry{At: record.UpdatedAt, Label: record.Label, Status: record.Status, Version: record.Version})
	}
	sort.Slice(entries, func(i, j int) bool {
		if entries[i].At != entries[j].At {
			return entries[i].At < entries[j].At
		}
		return entries[i].Version < entries[j].Version
	})
	return entries
}
