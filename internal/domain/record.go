package domain

import (
	"errors"
	"fmt"
	"strings"
)

type Status string

const (
	StatusDraft    Status = "draft"
	StatusReview   Status = "review"
	StatusApproved Status = "approved"
	StatusReturned Status = "returned"
	StatusArchived Status = "archived"
)

func (s Status) String() string { return string(s) }

type Record struct {
	ID          string `json:"id"`
	ProgramName string `json:"program_name"`
	Performer   string `json:"performer"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Status      Status `json:"status"`
	Score       int    `json:"score"`
	CreatedAt   int64  `json:"created_at"`
	UpdatedAt   int64  `json:"updated_at"`
	Version     int    `json:"version"`
}

type AuditEvent struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Action   string `json:"action"`
	Actor    string `json:"actor"`
	Detail   string `json:"detail"`
	At       int64  `json:"at"`
}

type Attachment struct {
	ID       string `json:"id"`
	RecordID string `json:"record_id"`
	Filename string `json:"filename"`
	Checksum string `json:"checksum"`
	Size     int64  `json:"size"`
}

type Workflow struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Status      string   `json:"status"`
	Steps       []string `json:"steps"`
	StartedAt   int64    `json:"started_at"`
	CompletedAt int64    `json:"completed_at"`
}

func NewRecord(id, program, performer, category, label string, at int64) Record {
	return Record{ID: id, ProgramName: strings.TrimSpace(program), Performer: strings.TrimSpace(performer), Category: strings.TrimSpace(category), Label: strings.TrimSpace(label), Status: StatusDraft, CreatedAt: at, UpdatedAt: at, Version: 1}
}

func (r Record) Validate() error {
	if r.ID == "" {
		return errors.New("record id is required")
	}
	if r.ProgramName == "" || r.Performer == "" {
		return errors.New("program name and performer are required")
	}
	if r.Category == "" {
		return errors.New("category is required")
	}
	if r.Label == "" {
		return errors.New("label is required")
	}
	if r.Score < 0 || r.Score > 100 {
		return fmt.Errorf("score out of range: %d", r.Score)
	}
	if !ValidStatus(r.Status) {
		return fmt.Errorf("invalid status: %s", r.Status)
	}
	return nil
}

func ValidStatus(status Status) bool {
	switch status {
	case StatusDraft, StatusReview, StatusApproved, StatusReturned, StatusArchived:
		return true
	default:
		return false
	}
}

func (r Record) IsVisible() bool {
	return r.Status != StatusArchived
}

func (r Record) IsReadyForArchive() bool {
	return r.Status == StatusApproved && r.Score >= 60
}

func (r Record) Summary() string {
	return fmt.Sprintf("%s:%s:%s:%d", r.ID, r.ProgramName, r.Label, r.Score)
}
