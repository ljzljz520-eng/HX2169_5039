package workflow

import (
	"errors"
	"fmt"

	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

type Engine struct {
	Store *store.Store
	Now   int64
}

func NewEngine(persistence *store.Store, now int64) *Engine {
	return &Engine{Store: persistence, Now: now}
}

func (e *Engine) Start(id, name string, steps []string) (domain.Workflow, error) {
	if e.Store == nil {
		return domain.Workflow{}, errors.New("store is required")
	}
	if id == "" || name == "" {
		return domain.Workflow{}, errors.New("workflow identity is required")
	}
	if len(steps) < 4 {
		return domain.Workflow{}, errors.New("workflow requires four steps")
	}
	flow := domain.Workflow{ID: id, Name: name, Status: "running", Steps: append([]string(nil), steps...), StartedAt: e.Now}
	if err := e.Store.SaveWorkflow(flow); err != nil {
		return domain.Workflow{}, err
	}
	return flow, nil
}

func (e *Engine) Complete(id string) (domain.Workflow, error) {
	flow, err := e.Store.GetWorkflow(id)
	if err != nil {
		return domain.Workflow{}, err
	}
	if flow.Status != "running" {
		return domain.Workflow{}, fmt.Errorf("workflow %s is not running", id)
	}
	flow.Status = "completed"
	flow.CompletedAt = e.Now
	if err := e.Store.SaveWorkflow(flow); err != nil {
		return domain.Workflow{}, err
	}
	return flow, nil
}

func (e *Engine) Cancel(id string) (domain.Workflow, error) {
	flow, err := e.Store.GetWorkflow(id)
	if err != nil {
		return domain.Workflow{}, err
	}
	if flow.Status != "running" {
		return domain.Workflow{}, fmt.Errorf("workflow %s is not running", id)
	}
	flow.Status = "cancelled"
	flow.CompletedAt = e.Now
	if err := e.Store.SaveWorkflow(flow); err != nil {
		return domain.Workflow{}, err
	}
	return flow, nil
}

func Progress(flow domain.Workflow) float64 {
	if len(flow.Steps) == 0 {
		return 0
	}
	if flow.Status == "completed" {
		return 1
	}
	return 0.5
}

func IsTerminal(status string) bool { return status == "completed" || status == "cancelled" }
