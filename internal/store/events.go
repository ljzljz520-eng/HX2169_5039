package store

import (
	"encoding/json"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) ListEvents(recordID string) ([]domain.AuditEvent, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	var events []domain.AuditEvent
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketEvents).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}
			if recordID == "" || event.RecordID == recordID {
				events = append(events, event)
			}
			return nil
		})
	})
	sort.Slice(events, func(i, j int) bool {
		if events[i].At != events[j].At {
			return events[i].At < events[j].At
		}
		return events[i].ID < events[j].ID
	})
	return events, err
}

func (s *Store) FindEvent(id string) (domain.AuditEvent, error) {
	if err := s.checkOpen(); err != nil {
		return domain.AuditEvent{}, err
	}
	var event domain.AuditEvent
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(bucketEvents), id, &event) })
	return event, err
}

func (s *Store) EventActions(recordID string) (map[string]int, error) {
	events, err := s.ListEvents(recordID)
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, event := range events {
		counts[strings.ToLower(event.Action)]++
	}
	return counts, nil
}

func (s *Store) SaveEvents(events []domain.AuditEvent) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketEvents)
		for _, event := range events {
			if event.ID == "" || event.RecordID == "" {
				return ErrNotFound
			}
			data, err := json.Marshal(event)
			if err != nil {
				return err
			}
			if err := bucket.Put([]byte(event.ID), data); err != nil {
				return err
			}
		}
		return nil
	})
}
