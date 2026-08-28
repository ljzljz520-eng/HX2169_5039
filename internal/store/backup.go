package store

import (
	"encoding/json"
	"sort"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

type Backup struct {
	Records     []domain.Record     `json:"records"`
	Events      []domain.AuditEvent `json:"events"`
	Workflows   []domain.Workflow   `json:"workflows"`
	Attachments []domain.Attachment `json:"attachments"`
}

func (s *Store) Export() (Backup, error) {
	if err := s.checkOpen(); err != nil {
		return Backup{}, err
	}
	backup := Backup{}
	err := s.db.View(func(tx *bolt.Tx) error {
		if err := eachJSON(tx.Bucket(bucketRecords), func(value []byte) error {
			var item domain.Record
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			backup.Records = append(backup.Records, item)
			return nil
		}); err != nil {
			return err
		}
		if err := eachJSON(tx.Bucket(bucketEvents), func(value []byte) error {
			var item domain.AuditEvent
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			backup.Events = append(backup.Events, item)
			return nil
		}); err != nil {
			return err
		}
		if err := eachJSON(tx.Bucket(bucketWorkflows), func(value []byte) error {
			var item domain.Workflow
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			backup.Workflows = append(backup.Workflows, item)
			return nil
		}); err != nil {
			return err
		}
		return eachJSON(tx.Bucket(bucketAttachments), func(value []byte) error {
			var item domain.Attachment
			if err := json.Unmarshal(value, &item); err != nil {
				return err
			}
			backup.Attachments = append(backup.Attachments, item)
			return nil
		})
	})
	sort.Slice(backup.Records, func(i, j int) bool { return backup.Records[i].ID < backup.Records[j].ID })
	sort.Slice(backup.Events, func(i, j int) bool { return backup.Events[i].ID < backup.Events[j].ID })
	sort.Slice(backup.Workflows, func(i, j int) bool { return backup.Workflows[i].ID < backup.Workflows[j].ID })
	sort.Slice(backup.Attachments, func(i, j int) bool { return backup.Attachments[i].ID < backup.Attachments[j].ID })
	return backup, err
}

func eachJSON(bucket *bolt.Bucket, callback func([]byte) error) error {
	return bucket.ForEach(func(_, value []byte) error { return callback(value) })
}

func EncodeBackup(backup Backup) ([]byte, error) { return json.MarshalIndent(backup, "", "  ") }

func DecodeBackup(data []byte) (Backup, error) {
	var backup Backup
	err := json.Unmarshal(data, &backup)
	return backup, err
}

func (s *Store) ImportBackup(backup Backup) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, record := range backup.Records {
			if err := putJSON(tx.Bucket(bucketRecords), record.ID, record); err != nil {
				return err
			}
		}
		for _, event := range backup.Events {
			if err := putJSON(tx.Bucket(bucketEvents), event.ID, event); err != nil {
				return err
			}
		}
		for _, flow := range backup.Workflows {
			if err := putJSON(tx.Bucket(bucketWorkflows), flow.ID, flow); err != nil {
				return err
			}
		}
		for _, attachment := range backup.Attachments {
			if err := putJSON(tx.Bucket(bucketAttachments), attachment.ID, attachment); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) DeleteArchived() (int, error) {
	if err := s.checkOpen(); err != nil {
		return 0, err
	}
	removed := 0
	err := s.db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucketRecords)
		return bucket.ForEach(func(key, value []byte) error {
			var record domain.Record
			if err := json.Unmarshal(value, &record); err != nil {
				return err
			}
			if record.Status == domain.StatusArchived {
				if err := bucket.Delete(key); err != nil {
					return err
				}
				removed++
			}
			return nil
		})
	})
	return removed, err
}
