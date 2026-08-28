package importer

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"example.com/graduation-showcase/internal/domain"
	"example.com/graduation-showcase/internal/store"
)

type Row struct {
	ID          string
	ProgramName string
	Performer   string
	Category    string
	Label       string
	Filename    string
	Content     string
}

type Result struct {
	Records     []domain.Record
	Attachments []domain.Attachment
	Rejected    []string
}

func ImportRows(rows []Row, persistence *store.Store, at int64) (Result, error) {
	if persistence == nil {
		return Result{}, errors.New("store is required")
	}
	result := Result{Records: make([]domain.Record, 0, len(rows)), Attachments: make([]domain.Attachment, 0, len(rows)), Rejected: []string{}}
	for index, row := range rows {
		record, attachment, err := BuildRow(row, at)
		if err != nil {
			result.Rejected = append(result.Rejected, strconv.Itoa(index)+":"+err.Error())
			continue
		}
		if err := persistence.SaveRecord(record); err != nil {
			return Result{}, err
		}
		if err := persistence.SaveAttachment(attachment); err != nil {
			return Result{}, err
		}
		result.Records = append(result.Records, record)
		result.Attachments = append(result.Attachments, attachment)
	}
	return result, nil
}

func BuildRow(row Row, at int64) (domain.Record, domain.Attachment, error) {
	if err := ValidateRow(row); err != nil {
		return domain.Record{}, domain.Attachment{}, err
	}
	record := domain.NewRecord(row.ID, row.ProgramName, row.Performer, row.Category, row.Label, at)
	checksum := sha256.Sum256([]byte(row.Content))
	attachment := domain.Attachment{ID: row.ID + "-attachment", RecordID: row.ID, Filename: row.Filename, Checksum: hex.EncodeToString(checksum[:]), Size: int64(len(row.Content))}
	return record, attachment, nil
}

func ValidateRow(row Row) error {
	if strings.TrimSpace(row.ID) == "" {
		return errors.New("id is required")
	}
	if strings.TrimSpace(row.ProgramName) == "" {
		return errors.New("program name is required")
	}
	if strings.TrimSpace(row.Performer) == "" {
		return errors.New("performer is required")
	}
	if strings.TrimSpace(row.Category) == "" {
		return errors.New("category is required")
	}
	if strings.TrimSpace(row.Label) == "" {
		return errors.New("label is required")
	}
	if strings.TrimSpace(row.Filename) == "" {
		return errors.New("filename is required")
	}
	if len(row.Content) == 0 {
		return errors.New("attachment content is required")
	}
	return nil
}

func ParseLines(lines []string) ([]Row, error) {
	rows := make([]Row, 0, len(lines))
	for lineNumber, line := range lines {
		parts := strings.Split(line, "|")
		if len(parts) != 7 {
			return nil, fmt.Errorf("line %d requires seven fields", lineNumber+1)
		}
		rows = append(rows, Row{ID: parts[0], ProgramName: parts[1], Performer: parts[2], Category: parts[3], Label: parts[4], Filename: parts[5], Content: parts[6]})
	}
	return rows, nil
}

func NormalizeRows(rows []Row) []Row {
	result := make([]Row, 0, len(rows))
	for _, row := range rows {
		row.ID = strings.TrimSpace(row.ID)
		row.ProgramName = strings.TrimSpace(row.ProgramName)
		row.Performer = strings.TrimSpace(row.Performer)
		row.Category = strings.TrimSpace(row.Category)
		row.Label = strings.TrimSpace(row.Label)
		row.Filename = strings.TrimSpace(row.Filename)
		result = append(result, row)
	}
	return result
}
