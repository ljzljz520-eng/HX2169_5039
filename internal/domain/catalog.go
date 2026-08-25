package domain

import (
	"errors"
	"sort"
	"strings"
)

type Catalog struct {
	Categories []string
	Labels     []string
	Statuses   []Status
}

func NewCatalog() Catalog {
	return Catalog{Categories: []string{}, Labels: []string{}, Statuses: []Status{StatusDraft, StatusReview, StatusApproved, StatusReturned, StatusArchived}}
}

func (c *Catalog) AddCategory(category string) error {
	category = strings.TrimSpace(category)
	if category == "" {
		return errors.New("category is required")
	}
	if containsString(c.Categories, category) {
		return errors.New("category already exists")
	}
	c.Categories = append(c.Categories, category)
	sort.Strings(c.Categories)
	return nil
}

func (c *Catalog) AddLabel(label string) error {
	label = strings.TrimSpace(label)
	if label == "" {
		return errors.New("label is required")
	}
	if containsString(c.Labels, label) {
		return errors.New("label already exists")
	}
	c.Labels = append(c.Labels, label)
	sort.Strings(c.Labels)
	return nil
}

func (c Catalog) HasCategory(category string) bool { return containsString(c.Categories, category) }
func (c Catalog) HasLabel(label string) bool       { return containsStringFold(c.Labels, label) }

func ParseTags(value string) []string {
	parts := strings.Split(value, ",")
	seen := make(map[string]bool)
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		normalized := strings.TrimSpace(part)
		if normalized != "" && !seen[normalized] {
			seen[normalized] = true
			result = append(result, normalized)
		}
	}
	sort.Strings(result)
	return result
}

func EncodeTags(tags []string) string {
	copyTags := append([]string(nil), tags...)
	sort.Strings(copyTags)
	return strings.Join(copyTags, ",")
}

func StatusOrder(status Status) int {
	switch status {
	case StatusDraft:
		return 1
	case StatusReview:
		return 2
	case StatusApproved:
		return 3
	case StatusReturned:
		return 4
	case StatusArchived:
		return 5
	default:
		return 99
	}
}
