package importer

import (
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/store"
)

func TestWorkflowImportReport(t *testing.T) {
	persistence, err := store.Open(filepath.Join(t.TempDir(), "import.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	rows := []Row{{ID: "import-1", ProgramName: "街舞", Performer: "社团", Category: "舞蹈", Label: "候选", Filename: "街舞.txt", Content: "deterministic"}}
	result, err := ImportRows(rows, persistence, 44)
	if err != nil || len(result.Records) != 1 || len(result.Attachments) != 1 {
		t.Fatalf("result: %#v %v", result, err)
	}
	if !ValidBatch(rows) {
		t.Fatal("rows should be valid")
	}
}

func TestImporterRejectsInvalid(t *testing.T) {
	if _, _, err := BuildRow(Row{ID: "bad"}, 1); err == nil {
		t.Fatal("expected validation error")
	}
	parsed, err := ParseLines([]string{"a|b|c|d|e|f|g"})
	if err != nil || len(parsed) != 1 {
		t.Fatalf("parsed: %#v %v", parsed, err)
	}
}
