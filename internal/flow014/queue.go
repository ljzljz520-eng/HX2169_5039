package flow014

import (
	"errors"
	"sort"

	"example.com/graduation-showcase/internal/domain"
)

type PrintQueue struct {
	IDs     []string
	Current int
}

func NewPrintQueue(ids []string) PrintQueue {
	ordered := append([]string(nil), ids...)
	sort.Strings(ordered)
	return PrintQueue{IDs: ordered}
}
func (q PrintQueue) Empty() bool { return len(q.IDs) == 0 }
func (q PrintQueue) Done() bool  { return q.Current >= len(q.IDs) }
func (q PrintQueue) CurrentID() string {
	if q.Done() {
		return ""
	}
	return q.IDs[q.Current]
}

func (q *PrintQueue) Next() error {
	if q.Done() {
		return errors.New("print queue is complete")
	}
	q.Current++
	return nil
}
func (q PrintQueue) Remaining() int {
	remaining := len(q.IDs) - q.Current
	if remaining < 0 {
		return 0
	}
	return remaining
}

func QueueRecords(records []domain.Record) PrintQueue {
	ids := make([]string, 0, len(records))
	for _, record := range records {
		if record.IsVisible() {
			ids = append(ids, record.ID)
		}
	}
	return NewPrintQueue(ids)
}

func (f *Flow) PrintQueue(queue PrintQueue) ([]PrintResult, error) {
	results := make([]PrintResult, 0, len(queue.IDs))
	for !queue.Done() {
		result, err := f.PrintDocument(queue.CurrentID())
		if err != nil {
			return nil, err
		}
		results = append(results, result)
		if err := queue.Next(); err != nil {
			return nil, err
		}
	}
	return results, nil
}

func ValidatePrintResult(result PrintResult) error {
	if result.RecordID == "" {
		return errors.New("record id is required")
	}
	if result.Label == "" {
		return errors.New("label is required")
	}
	return nil
}
