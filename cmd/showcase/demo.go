package main

import (
	"example.com/graduation-showcase/internal/api"
	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/report"
	"example.com/graduation-showcase/internal/store"
	"example.com/graduation-showcase/internal/workflow"
)

func runDemo(service *api.Service, persistence *store.Store) error {
	flowEngine := workflow.NewEngine(persistence, 202501010900)
	if _, err := flowEngine.Start("demo-flow", "毕业晚会评选", []string{"登记", "审核", "确认", "归档"}); err != nil {
		return err
	}
	if _, err := service.CreateRecord(api.CreateInput{ID: "demo-1", ProgramName: "星河合唱", Performer: "毕业班", Category: "声乐", Label: "待评", Actor: "demo"}); err != nil {
		return err
	}
	if _, err := service.BeginReview("demo-1", "demo"); err != nil {
		return err
	}
	if _, err := service.ReviewRecord("demo-1", api.ReviewInput{Score: 88, Actor: "demo"}); err != nil {
		return err
	}
	if _, err := service.UpdateLabel("demo-1", api.UpdateInput{Label: "重点推荐", Actor: "demo"}); err != nil {
		return err
	}
	if _, err := service.ArchiveRecord("demo-1", "demo"); err != nil {
		return err
	}
	if _, err := flowEngine.Complete("demo-flow"); err != nil {
		return err
	}
	if _, err := persistence.GetRecord("demo-1"); err != nil {
		return err
	}
	_, _ = report.BuildReport(persistence, true)
	_ = domain.StatusApproved
	return nil
}
