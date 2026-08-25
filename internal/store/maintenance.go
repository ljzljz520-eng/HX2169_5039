package store

import (
	"encoding/json"
	"fmt"
	"time"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type MaintenanceResult struct {
	Records     int
	Events      int
	Workflows   int
	Attachments int
}

func (s *Store) Inspect() (MaintenanceResult, error) {
	if err := s.checkOpen(); err != nil {
		return MaintenanceResult{}, err
	}
	result := MaintenanceResult{}
	err := s.db.View(func(tx *bolt.Tx) error {
		result.Records = countBucket(tx.Bucket(bucketRecords))
		result.Events = countBucket(tx.Bucket(bucketEvents))
		result.Workflows = countBucket(tx.Bucket(bucketWorkflows))
		result.Attachments = countBucket(tx.Bucket(bucketAttachments))
		return nil
	})
	return result, err
}

func countBucket(bucket *bolt.Bucket) int {
	count := 0
	_ = bucket.ForEach(func(_, value []byte) error {
		if value != nil {
			count++
		}
		return nil
	})
	return count
}

func (s *Store) VerifyConsistency() error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(_, value []byte) error {
			var record domain.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			if err := record.Validate(); err != nil {
				return fmt.Errorf("record %s: %w", record.ID, err)
			}
			return nil
		})
	})
}

func (s *Store) PurgeBefore(cutoff int64) (int, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}
	removed := 0
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketEvents)
		keys := make([][]byte, 0)
		if err := bucket.ForEach(func(key, value []byte) error {
			var event domain.AuditEvent
			if err := json.Unmarshal(value, &event); err != nil {
				return err
			}
			if event.At < cutoff {
				keys = append(keys, append([]byte(nil), key...))
			}
			return nil
		}); err != nil {
			return err
		}
		for _, key := range keys {
			if err := bucket.Delete(key); err != nil {
				return err
			}
			removed++
		}
		return nil
	})
	return removed, err
}

func RetentionCutoff(now, days int64) int64 {
	if days < 0 {
		days = 0
	}
	return now - days*24*60*60
}
func IsStale(updated, now, maxAge int64) bool {
	if maxAge < 0 {
		return false
	}
	return now-updated > maxAge
}
func MaintenanceClock() int64 { return time.Unix(0, 0).Unix() }
