package store

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"sync"
	"time"

	"example.com/graduation-showcase/internal/domain"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrNotFound = errors.New("entity not found")
	ErrClosed   = errors.New("store is closed")
)

var bucketRecords = []byte("records")
var bucketEvents = []byte("events")
var bucketWorkflows = []byte("workflows")
var bucketAttachments = []byte("attachments")

type Store struct {
	db     *bolt.DB
	mu     sync.RWMutex
	closed bool
}

func Open(path string) (*Store, error) {
	if path == "" {
		return nil, errors.New("database path is required")
	}
	db, err := bolt.Open(filepath.Clean(path), 0o600, &bolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, fmt.Errorf("open bolt: %w", err)
	}
	store := &Store{db: db}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (s *Store) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, bucket := range [][]byte{bucketRecords, bucketEvents, bucketWorkflows, bucketAttachments} {
			if _, err := tx.CreateBucketIfNotExists(bucket); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	return s.db.Close()
}

func (s *Store) checkOpen() error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.closed {
		return ErrClosed
	}
	return nil
}

func putJSON(bucket *bolt.Bucket, key string, value any) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return bucket.Put([]byte(key), data)
}

func getJSON(bucket *bolt.Bucket, key string, target any) error {
	data := bucket.Get([]byte(key))
	if data == nil {
		return ErrNotFound
	}
	return json.Unmarshal(data, target)
}

func cloneJSON(source any, target any) error {
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

func (s *Store) SaveRecord(record domain.Record) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if err := record.Validate(); err != nil {
		return err
	}
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(bucketRecords), record.ID, record) })
}

func (s *Store) GetRecord(id string) (domain.Record, error) {
	if err := s.checkOpen(); err != nil {
		return domain.Record{}, err
	}
	var record domain.Record
	err := s.db.View(func(tx *bolt.Tx) error { return getJSON(tx.Bucket(bucketRecords), id, &record) })
	return record, err
}

func (s *Store) SaveEvent(event domain.AuditEvent) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if event.ID == "" || event.RecordID == "" {
		return errors.New("event id and record id are required")
	}
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(bucketEvents), event.ID, event) })
}

func (s *Store) SaveWorkflow(flow domain.Workflow) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if flow.ID == "" || flow.Name == "" {
		return errors.New("workflow id and name are required")
	}
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(bucketWorkflows), flow.ID, flow) })
}

func (s *Store) SaveAttachment(attachment domain.Attachment) error {
	if err := s.checkOpen(); err != nil {
		return err
	}
	if attachment.ID == "" || attachment.RecordID == "" || attachment.Filename == "" {
		return errors.New("attachment identity is required")
	}
	return s.db.Update(func(tx *bolt.Tx) error { return putJSON(tx.Bucket(bucketAttachments), attachment.ID, attachment) })
}

func (s *Store) Snapshot() (map[string]domain.Record, error) {
	if err := s.checkOpen(); err != nil {
		return nil, err
	}
	result := make(map[string]domain.Record)
	err := s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(bucketRecords).ForEach(func(k, v []byte) error {
			var record domain.Record
			if err := json.Unmarshal(v, &record); err != nil {
				return err
			}
			result[string(k)] = record
			return nil
		})
	})
	return result, err
}
