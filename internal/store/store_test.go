package store

import (
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/domain"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	path := filepath.Join(t.TempDir(), "records.db")
	first, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	record := domain.NewRecord("persist-1", "合唱", "一班", "声乐", "推荐", 100)
	if err := first.SaveRecord(record); err != nil {
		t.Fatal(err)
	}
	if err := first.SaveEvent(domain.AuditEvent{ID: "event-1", RecordID: record.ID, Action: "created", At: 100}); err != nil {
		t.Fatal(err)
	}
	if err := first.SaveWorkflow(domain.Workflow{ID: "workflow-1", Name: "导入", Status: "running", Steps: []string{"a", "b", "c", "d"}}); err != nil {
		t.Fatal(err)
	}
	if err := first.SaveAttachment(domain.Attachment{ID: "attachment-1", RecordID: record.ID, Filename: "节目.txt", Checksum: "abc", Size: 3}); err != nil {
		t.Fatal(err)
	}
	if err := first.Close(); err != nil {
		t.Fatal(err)
	}
	second, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer second.Close()
	loaded, err := second.GetRecord(record.ID)
	if err != nil || loaded.Label != "推荐" {
		t.Fatalf("loaded: %#v %v", loaded, err)
	}
	events, err := second.ListEvents(record.ID)
	if err != nil || len(events) != 1 {
		t.Fatalf("events: %#v %v", events, err)
	}
	flow, err := second.GetWorkflow("workflow-1")
	if err != nil || flow.Name != "导入" {
		t.Fatalf("workflow: %#v %v", flow, err)
	}
	attachment, err := second.GetAttachment("attachment-1")
	if err != nil || attachment.Filename != "节目.txt" {
		t.Fatalf("attachment: %#v %v", attachment, err)
	}
}

func TestStoreQueries(t *testing.T) {
	persistence, err := Open(filepath.Join(t.TempDir(), "query.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	for _, record := range []domain.Record{domain.NewRecord("a", "星河", "甲", "声乐", "推荐", 1), domain.NewRecord("b", "青春", "乙", "舞蹈", "待定", 2)} {
		if err := persistence.SaveRecord(record); err != nil {
			t.Fatal(err)
		}
	}
	items, err := persistence.ListRecords(Query{Label: "推荐"})
	if err != nil || len(items) != 1 {
		t.Fatalf("items: %#v %v", items, err)
	}
}
