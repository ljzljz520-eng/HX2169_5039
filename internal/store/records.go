package store

import (
	"encoding/json"
	"sort"
	"strings"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type Query struct {
	Text            string
	Label           string
	Status          domain.Status
	IncludeArchived bool
}

func (s *Store) ListRecords(query Query) ([]domain.Record, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	var records []domain.Record
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			if matchesQuery(record, query) {
				records = append(records, record)
			}
			return nil
		})
	})
	sort.Slice(records, func(i, j int) bool {
		if records[i].UpdatedAt != records[j].UpdatedAt {
			return records[i].UpdatedAt < records[j].UpdatedAt
		}
		return records[i].ID < records[j].ID
	})
	return records, err
}

func matchesQuery(record domain.Record, query Query) bool {
	if !query.IncludeArchived && !record.IsVisible() {
		return false
	}
	if query.Status != "" && record.Status != query.Status {
		return false
	}
	if query.Label != "" && !strings.EqualFold(record.Label, query.Label) {
		return false
	}
	if query.Text == "" {
		return true
	}
	text := strings.ToLower(query.Text)
	return strings.Contains(strings.ToLower(record.ID), text) || strings.Contains(strings.ToLower(record.ProgramName), text) || strings.Contains(strings.ToLower(record.Performer), text)
}

func (s *Store) DeleteRecord(id string) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		if bucket.Get([]byte(id)) == nil {
			return ErrNotFound
		}
		return bucket.Delete([]byte(id))
	})
}

func (s *Store) CountRecords(query Query) (int, error) {
	records, err := s.ListRecords(query)
	return len(records), err
}

func (s *Store) ReplaceRecords(records []domain.Record) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		for _, record := range records {
			if err := record.Validate(); err != nil {
				return err
			}
			if err := putJSON(bucket, record.ID, record); err != nil {
				return err
			}
		}
		return nil
	})
}
