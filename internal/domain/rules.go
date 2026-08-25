package domain

import (
	"errors"
	"fmt"
)

var (
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrLowScore          = errors.New("score must be at least 60")
	ErrArchived          = errors.New("archived record cannot change")
)

func BeginReview(r Record, at int64) (Record, error) {
	if r.Status != StatusDraft && r.Status != StatusReturned {
		return r, fmt.Errorf("%w: %s to review", ErrInvalidTransition, r.Status)
	}
	r.Status = StatusReview
	r.UpdatedAt = at
	r.Version++
	return r, nil
}

func ReviewRecord(r Record, score int, at int64) (Record, error) {
	if r.Status != StatusReview {
		return r, fmt.Errorf("%w: review required", ErrInvalidTransition)
	}
	if score < 0 || score > 100 {
		return r, fmt.Errorf("score out of range: %d", score)
	}
	r.Score = score
	r.UpdatedAt = at
	r.Version++
	if score >= 60 {
		r.Status = StatusApproved
	} else {
		r.Status = StatusReturned
	}
	return r, nil
}

func UpdateLabel(r Record, label string, at int64) (Record, error) {
	if r.Status == StatusArchived {
		return r, ErrArchived
	}
	if label == "" {
		return r, errors.New("label is required")
	}
	r.Label = label
	r.UpdatedAt = at
	r.Version++
	return r, nil
}

func ArchiveRecord(r Record, at int64) (Record, error) {
	if !r.IsReadyForArchive() {
		return r, fmt.Errorf("%w: record is not approved", ErrInvalidTransition)
	}
	r.Status = StatusArchived
	r.UpdatedAt = at
	r.Version++
	return r, nil
}

func Publishable(r Record) bool {
	if r.Status == StatusApproved {
		return true
	}
	if r.Status == StatusReview && r.Score >= 60 {
		return true
	}
	return false
}

func NormalizeLabel(label string) string {
	if len(label) > 40 {
		return label[:40]
	}
	return label
}
