package flow014

import (
	"errors"
	"fmt"
	"sort"

	"example.com/graduation-showcase/internal/domain"
)

type Checklist struct {
	RecordID string
	Steps    []string
	Done     map[string]bool
}

func NewChecklist(recordID string) Checklist {
	return Checklist{RecordID: recordID, Steps: []string{"资料齐全", "节目已审核", "标签已确认", "归档可追溯"}, Done: make(map[string]bool)}
}

func (c *Checklist) Complete(step string) error {
	if c.RecordID == "" {
		return errors.New("record id is required")
	}
	if !containsStep(c.Steps, step) {
		return fmt.Errorf("unknown checklist step: %s", step)
	}
	c.Done[step] = true
	return nil
}

func (c Checklist) CompleteCount() int {
	count := 0
	for _, step := range c.Steps {
		if c.Done[step] {
			count++
		}
	}
	return count
}
func (c Checklist) Ready() bool { return len(c.Steps) > 0 && c.CompleteCount() == len(c.Steps) }

func BuildChecklist(record domain.Record) Checklist {
	checklist := NewChecklist(record.ID)
	if record.ProgramName != "" && record.Performer != "" {
		_ = checklist.Complete("资料齐全")
	}
	if record.Status == domain.StatusApproved || record.Status == domain.StatusArchived {
		_ = checklist.Complete("节目已审核")
	}
	if record.Label != "" {
		_ = checklist.Complete("标签已确认")
	}
	if record.Status == domain.StatusArchived {
		_ = checklist.Complete("归档可追溯")
	}
	return checklist
}

func MissingSteps(checklist Checklist) []string {
	result := make([]string, 0)
	for _, step := range checklist.Steps {
		if !checklist.Done[step] {
			result = append(result, step)
		}
	}
	return result
}

func SortedMissingSteps(checklist Checklist) []string {
	result := MissingSteps(checklist)
	sort.Strings(result)
	return result
}
func containsStep(steps []string, target string) bool {
	for _, step := range steps {
		if step == target {
			return true
		}
	}
	return false
}
