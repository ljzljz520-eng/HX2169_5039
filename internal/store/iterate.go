package store

import (
	"encoding/json"
	"sort"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type RecordVisitor func(domain.Record) error

func (s *Store) VisitRecords(visitor RecordVisitor) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if visitor == nil {
		return nil
	}
	return s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			return visitor(record)
		})
	})
}

func (s *Store) IDs() ([]string, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	ids := []string{}
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(key, _ []byte) error { ids = append(ids, string(key)); return nil })
	})
	sort.Strings(ids)
	return ids, err
}

func (s *Store) Labels() ([]string, error) {
	records, err := s.ListRecords(Query{IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	return domain.UniqueLabels(records), nil
}

func (s *Store) CategoryCounts() (map[string]int, error) {
	records, err := s.ListRecords(Query{IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	counts := make(map[string]int)
	for _, record := range records {
		counts[record.Category]++
	}
	return counts, nil
}

func (s *Store) FindByVersion(version int) ([]domain.Record, error) {
	records, err := s.ListRecords(Query{IncludeArchived: true})
	if err != nil {
		return nil, err
	}
	result := make([]domain.Record, 0)
	for _, record := range records {
		if record.Version == version {
			result = append(result, record)
		}
	}
	return result, nil
}
