package store

import (
	"encoding/json"
	"sort"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *Store) GetWorkflow(id string) (domain.Workflow, error) {
	if err := s.checkOpen(); err != nil {
		return domain.Workflow{}, err
	}
	var flow domain.Workflow
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(bucketWorkflows), id, &flow) })
	return flow, err
}

func (s *Store) ListWorkflows(status string) ([]domain.Workflow, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	var result []domain.Workflow
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketWorkflows).ForEach(func(_, value []byte) error {
			var flow domain.Workflow
			if err := json.Unmarshal(value, &flow); err != nil {
				return err
			}
			if status == "" || flow.Status == status {
				result = append(result, flow)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}

func (s *Store) GetAttachment(id string) (domain.Attachment, error) {
	if err := s.checkOpen(); err != nil {
		return domain.Attachment{}, err
	}
	var attachment domain.Attachment
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(bucketAttachments), id, &attachment) })
	return attachment, err
}

func (s *Store) ListAttachments(recordID string) ([]domain.Attachment, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	var result []domain.Attachment
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketAttachments).ForEach(func(_, value []byte) error {
			var attachment domain.Attachment
			if err := json.Unmarshal(value, &attachment); err != nil {
				return err
			}
			if recordID == "" || attachment.RecordID == recordID {
				result = append(result, attachment)
			}
			return nil
		})
	})
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, err
}
