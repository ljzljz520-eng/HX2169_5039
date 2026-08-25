package workflow

import (
	"errors"
	"fmt"
	"strings"
)

type Plan struct {
	Name    string
	Steps   []string
	Current int
}

func NewPlan(name string, steps []string) (Plan, error) {
	if strings.TrimSpace(name) == "" {
		return Plan{}, errors.New("plan name is required")
	}
	if len(steps) < 4 {
		return Plan{}, errors.New("plan requires four steps")
	}
	return Plan{Name: name, Steps: append([]string(nil), steps...), Current: 0}, nil
}

func (p Plan) CurrentStep() string {
	if p.Current < 0 || p.Current >= len(p.Steps) {
		return ""
	}
	return p.Steps[p.Current]
}

func (p *Plan) Advance() error {
	if p.Current >= len(p.Steps) {
		return errors.New("plan already complete")
	}
	p.Current++
	return nil
}
func (p Plan) Complete() bool { return p.Current >= len(p.Steps) }
func (p Plan) Remaining() int {
	remaining := len(p.Steps) - p.Current
	if remaining < 0 {
		return 0
	}
	return remaining
}
func (p Plan) String() string { return fmt.Sprintf("%s:%d/%d", p.Name, p.Current, len(p.Steps)) }

func ValidatePlan(p Plan) error {
	if strings.TrimSpace(p.Name) == "" {
		return errors.New("plan name is required")
	}
	if len(p.Steps) < 4 {
		return errors.New("plan requires four steps")
	}
	if p.Current < 0 || p.Current > len(p.Steps) {
		return errors.New("plan cursor out of range")
	}
	return nil
}

func PlanProgress(p Plan) float64 {
	if len(p.Steps) == 0 {
		return 0
	}
	return float64(p.Current) / float64(len(p.Steps))
}
