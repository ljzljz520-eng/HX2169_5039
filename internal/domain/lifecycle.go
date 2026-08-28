package domain

import (
	"errors"
	"fmt"
)

type Transition struct {
	From   Status
	To     Status
	Reason string
}

func AllowedTransitions() []Transition {
	return []Transition{{From: StatusDraft, To: StatusReview, Reason: "资料登记完成"}, {From: StatusReturned, To: StatusReview, Reason: "退回后再次审核"}, {From: StatusReview, To: StatusApproved, Reason: "评分达到标准"}, {From: StatusReview, To: StatusReturned, Reason: "评分未达标准"}, {From: StatusApproved, To: StatusArchived, Reason: "确认后归档"}}
}

func CanTransition(from, to Status) bool {
	for _, transition := range AllowedTransitions() {
		if transition.From == from && transition.To == to {
			return true
		}
	}
	return false
}

func TransitionReason(from, to Status) (string, error) {
	for _, transition := range AllowedTransitions() {
		if transition.From == from && transition.To == to {
			return transition.Reason, nil
		}
	}
	return "", fmt.Errorf("%w: %s -> %s", ErrInvalidTransition, from, to)
}

func ApplyTransition(record Record, to Status, at int64) (Record, error) {
	if !CanTransition(record.Status, to) {
		return record, errors.New("transition is not allowed")
	}
	record.Status = to
	record.UpdatedAt = at
	record.Version++
	return record, nil
}

func LifecycleLabel(record Record) string {
	switch record.Status {
	case StatusDraft:
		return "待登记"
	case StatusReview:
		return "审核中"
	case StatusApproved:
		return "已确认"
	case StatusReturned:
		return "需补充"
	case StatusArchived:
		return "已归档"
	default:
		return "未知"
	}
}

func StatusSummary(status Status) string {
	if !ValidStatus(status) {
		return "未知状态"
	}
	return string(status)
}

func RequiresReview(record Record) bool {
	return record.Status == StatusDraft || record.Status == StatusReturned
}
