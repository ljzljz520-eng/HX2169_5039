package api

import (
	"errors"
	"fmt"

	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

type Clock interface{ Now() int64 }

type FixedClock struct{ Value int64 }

func (c FixedClock) Now() int64 { return c.Value }

type Service struct {
	Store *store.Store
	Clock Clock
}

type CreateInput struct {
	ID          string `json:"id"`
	ProgramName string `json:"program_name"`
	Performer   string `json:"performer"`
	Category    string `json:"category"`
	Label       string `json:"label"`
	Actor       string `json:"actor"`
}

type UpdateInput struct {
	Label string `json:"label"`
	Actor string `json:"actor"`
}

type ReviewInput struct {
	Score int    `json:"score"`
	Actor string `json:"actor"`
}

type QueryInput struct {
	Text            string
	Label           string
	Status          domain.Status
	IncludeArchived bool
}

func (s *Service) timestamp() int64 {
	if s.Clock == nil {
		return 0
	}
	return s.Clock.Now()
}

func (s *Service) CreateRecord(input CreateInput) (domain.Record, error) {
	if s.Store == nil {
		return domain.Record{}, errors.New("store is required")
	}
	record := domain.NewRecord(input.ID, input.ProgramName, input.Performer, input.Category, input.Label, s.timestamp())
	if err := record.Validate(); err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(record); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(record.ID, "created", input.Actor, "record registered"); err != nil {
		return domain.Record{}, err
	}
	return record, nil
}

func (s *Service) BeginReview(id, actor string) (domain.Record, error) {
	record, err := s.mustRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := domain.BeginReview(record, s.timestamp())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(id, "review_started", actor, "record entered review"); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) ReviewRecord(id string, input ReviewInput) (domain.Record, error) {
	record, err := s.mustRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := domain.ReviewRecord(record, input.Score, s.timestamp())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(id, "reviewed", input.Actor, fmt.Sprintf("score=%d status=%s", input.Score, updated.Status)); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) UpdateLabel(id string, input UpdateInput) (domain.Record, error) {
	record, err := s.mustRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := domain.UpdateLabel(record, domain.NormalizeLabel(input.Label), s.timestamp())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(id, "label_updated", input.Actor, "label changed to "+updated.Label); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) ArchiveRecord(id, actor string) (domain.Record, error) {
	record, err := s.mustRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	updated, err := domain.ArchiveRecord(record, s.timestamp())
	if err != nil {
		return domain.Record{}, err
	}
	if err := s.Store.SaveRecord(updated); err != nil {
		return domain.Record{}, err
	}
	if err := s.audit(id, "archived", actor, "record archived"); err != nil {
		return domain.Record{}, err
	}
	return updated, nil
}

func (s *Service) QueryRecords(input QueryInput) ([]domain.Record, error) {
	return s.Store.ListRecords(store.Query{Text: input.Text, Label: input.Label, Status: input.Status, IncludeArchived: input.IncludeArchived})
}

func (s *Service) GetRecord(id string) (domain.Record, error) { return s.mustRecord(id) }

func (s *Service) mustRecord(id string) (domain.Record, error) {
	if id == "" {
		return domain.Record{}, errors.New("record id is required")
	}
	return s.Store.GetRecord(id)
}

func (s *Service) audit(recordID, action, actor, detail string) error {
	event := domain.AuditEvent{ID: fmt.Sprintf("%s-%s-%d", recordID, action, s.timestamp()), RecordID: recordID, Action: action, Actor: actor, Detail: detail, At: s.timestamp()}
	return s.Store.SaveEvent(event)
}
