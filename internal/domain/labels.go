package domain

import (
	"sort"
	"strings"
)

type LabelPolicy struct {
	Allowed   []string
	Default   string
	MaxLength int
}

func DefaultLabelPolicy() LabelPolicy {
	return LabelPolicy{Allowed: []string{"待评", "候选", "重点推荐", "备选", "已发布"}, Default: "待评", MaxLength: 40}
}

func (p LabelPolicy) Normalize(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return p.Default
	}
	if p.MaxLength > 0 && len(value) > p.MaxLength {
		return value[:p.MaxLength]
	}
	return value
}

func (p LabelPolicy) Allows(value string) bool {
	normalized := p.Normalize(value)
	if normalized == p.Default && strings.TrimSpace(value) == "" {
		return true
	}
	for _, allowed := range p.Allowed {
		if normalized == allowed {
			return true
		}
	}
	return false
}

func LabelOrder(label string) int {
	switch label {
	case "重点推荐":
		return 1
	case "已发布":
		return 2
	case "候选":
		return 3
	case "备选":
		return 4
	case "待评":
		return 5
	default:
		return 99
	}
}

func SortLabels(labels []string) []string {
	result := append([]string(nil), labels...)
	sort.SliceStable(result, func(i, j int) bool {
		left, right := LabelOrder(result[i]), LabelOrder(result[j])
		if left != right {
			return left < right
		}
		return result[i] < result[j]
	})
	return result
}

func LabelTransitions(previous, current string) bool {
	if previous == current {
		return false
	}
	if LabelOrder(current) < LabelOrder(previous) {
		return true
	}
	return strings.TrimSpace(current) != ""
}

func ApplyLabelPolicy(record Record, policy LabelPolicy, value string, at int64) (Record, bool) {
	normalized := policy.Normalize(value)
	if !policy.Allows(normalized) {
		return record, false
	}
	updated, err := UpdateLabel(record, normalized, at)
	if err != nil {
		return record, false
	}
	return updated, true
}
