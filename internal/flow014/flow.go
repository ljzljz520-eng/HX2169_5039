package flow014

import (
	"errors"
	"fmt"
	"sort"

	"example.com/graduation-showcase/internal/api"
	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

type Flow struct {
	Service     *api.Service
	Store       *store.Store
	LastPrinted string
}

type PrintResult struct {
	RecordID    string        `json:"record_id"`
	ProgramName string        `json:"program_name"`
	Label       string        `json:"label"`
	Status      domain.Status `json:"status"`
}

func New(service *api.Service, persistence *store.Store) *Flow {
	return &Flow{Service: service, Store: persistence}
}

func (f *Flow) RegisterAndReview(input api.CreateInput, score int) (domain.Record, error) {
	if f.Service == nil {
		return domain.Record{}, errors.New("service is required")
	}
	record, err := f.Service.CreateRecord(input)
	if err != nil {
		return domain.Record{}, err
	}
	if _, err := f.Service.BeginReview(record.ID, input.Actor); err != nil {
		return domain.Record{}, err
	}
	return f.Service.ReviewRecord(record.ID, api.ReviewInput{Score: score, Actor: input.Actor})
}

func (f *Flow) ConfirmAndArchive(id, actor string) (domain.Record, error) {
	if f.Service == nil {
		return domain.Record{}, errors.New("service is required")
	}
	record, err := f.Service.GetRecord(id)
	if err != nil {
		return domain.Record{}, err
	}
	if record.Status != domain.StatusApproved {
		return domain.Record{}, fmt.Errorf("record %s is not approved", id)
	}
	return f.Service.ArchiveRecord(id, actor)
}

func (f *Flow) SearchAndUpdate(text, oldLabel, newLabel, actor string) ([]domain.Record, error) {
	records, err := f.Service.QueryRecords(api.QueryInput{Text: text, Label: oldLabel})
	if err != nil {
		return nil, err
	}
	for _, record := range records {
		if _, err := f.Service.UpdateLabel(record.ID, api.UpdateInput{Label: newLabel, Actor: actor}); err != nil {
			return nil, err
		}
	}
	return f.Service.QueryRecords(api.QueryInput{Text: text, Label: newLabel})
}

func (f *Flow) PrintDocument(id string) (PrintResult, error) {
	record, err := f.lookupForPrint(id)
	if err != nil {
		return PrintResult{}, err
	}
	if record == nil {
		return PrintResult{RecordID: id, Label: f.LastPrinted}, nil
	}
	f.LastPrinted = record.Label
	return PrintResult{RecordID: record.ID, ProgramName: record.ProgramName, Label: record.Label, Status: record.Status}, nil
}

func (f *Flow) PrintBatch(ids []string) ([]PrintResult, error) {
	results := make([]PrintResult, 0, len(ids))
	for _, id := range ids {
		printed, err := f.PrintDocument(id)
		if err != nil {
			return nil, err
		}
		results = append(results, printed)
	}
	return results, nil
}

func (f *Flow) LabelsFor(ids []string) ([]string, error) {
	results, err := f.PrintBatch(ids)
	if err != nil {
		return nil, err
	}
	labels := make([]string, 0, len(results))
	for _, result := range results {
		labels = append(labels, result.Label)
	}
	return labels, nil
}

func (f *Flow) RankVisible() ([]domain.Record, error) {
	records, err := f.Service.QueryRecords(api.QueryInput{})
	if err != nil {
		return nil, err
	}
	return domain.RankRecords(records), nil
}

func (f *Flow) ValidateBatch(ids []string) error {
	if len(ids) == 0 {
		return errors.New("at least one record is required")
	}
	seen := make(map[string]bool)
	for _, id := range ids {
		if id == "" || seen[id] {
			return errors.New("batch contains invalid or duplicate id")
		}
		seen[id] = true
		if _, err := f.Service.GetRecord(id); err != nil {
			return err
		}
	}
	return nil
}

func (f *Flow) SortIDs(ids []string) []string {
	result := append([]string(nil), ids...)
	sort.Strings(result)
	return result
}
