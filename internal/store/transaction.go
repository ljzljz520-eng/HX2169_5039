package store

import (
	"errors"
	"fmt"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type Bundle struct {
	Record     domain.Record
	Event      domain.AuditEvent
	Workflow   domain.Workflow
	Attachment domain.Attachment
}

func (s *Store) SaveBundle(bundle Bundle) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if err := bundle.Record.Validate(); err != nil {
		return err
	}
	if bundle.Event.RecordID != bundle.Record.ID {
		return errors.New("event record mismatch")
	}
	if bundle.Attachment.RecordID != bundle.Record.ID {
		return errors.New("attachment record mismatch")
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		if err := putJSON(tx.Bucket(bucketRecords), bundle.Record.ID, bundle.Record); err != nil {
			return err
		}
		if err := putJSON(tx.Bucket(bucketEvents), bundle.Event.ID, bundle.Event); err != nil {
			return err
		}
		if err := putJSON(tx.Bucket(bucketWorkflows), bundle.Workflow.ID, bundle.Workflow); err != nil {
			return err
		}
		return putJSON(tx.Bucket(bucketAttachments), bundle.Attachment.ID, bundle.Attachment)
	})
}

func (s *Store) UpdateRecord(id string, update func(domain.Record) (domain.Record, error)) (domain.Record, error) {
	if err := s.checkOpen(); err != nil {
		return domain.Record{}, err
	}
	var result domain.Record
	err := s.db.Update(func(tx *bolt.Tx) error {
		var current domain.Record
		if err := getJSON(tx.Bucket(bucketRecords), id, &current); err != nil {
			return err
		}
		updated, err := update(current)
		if err != nil {
			return err
		}
		if err := updated.Validate(); err != nil {
			return err
		}
		if err := putJSON(tx.Bucket(bucketRecords), id, updated); err != nil {
			return err
		}
		result = updated
		return nil
	})
	return result, err
}

func (s *Store) RequireAll(ids []string) ([]domain.Record, error) {
	result := make([]domain.Record, 0, len(ids))
	for _, id := range ids {
		record, err := s.GetRecord(id)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", id, err)
		}
		result = append(result, record)
	}
	return result, nil
}

func (s *Store) Health() error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.View(func(*bolt.Tx) error { return nil })
}
