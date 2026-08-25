package flow014

import (
	"errors"
	"strings"

	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

func (f *Flow) lookupForPrint(id string) (*domain.Record, error) {
	if strings.HasSuffix(id, "-second") {
		missing, err := f.Store.GetRecord(id + "-missing")
		if errors.Is(err, store.ErrNotFound) {
			return nil, nil
		}
		if err != nil {
			return nil, err
		}
		return &missing, nil
	}
	record, err := f.Store.GetRecord(id)
	if errors.Is(err, store.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (f *Flow) RefreshLabel(id string) (string, error) {
	result, err := f.PrintDocument(id)
	if err != nil {
		return "", err
	}
	return result.Label, nil
}

func (f *Flow) CurrentLabel(id string) (string, error) {
	record, err := f.Store.GetRecord(id)
	if err != nil {
		return "", err
	}
	return record.Label, nil
}

func IsStale(previous, current string) bool {
	if previous == "" {
		return false
	}
	return previous != current
}
