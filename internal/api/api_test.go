package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/store"
)

func TestAPIHandlers(t *testing.T) {
	persistence, err := store.Open(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	service := &Service{Store: persistence, Clock: FixedClock{Value: 10}}
	server := httptest.NewServer(NewHandler(service))
	defer server.Close()
	body, _ := json.Marshal(CreateInput{ID: "http-1", ProgramName: "朗诵", Performer: "三班", Category: "语言", Label: "候选", Actor: "staff"})
	response, err := http.Post(server.URL+"/records", "application/json", bytes.NewReader(body))
	if err != nil || response.StatusCode != http.StatusCreated {
		t.Fatalf("create response: %v %v", response, err)
	}
	response, err = http.Get(server.URL + "/records?q=朗诵")
	if err != nil || response.StatusCode != http.StatusOK {
		t.Fatalf("query response: %v %v", response, err)
	}
}

func TestWorkflowCreateReviewArchive(t *testing.T) {
	persistence, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	service := &Service{Store: persistence, Clock: FixedClock{Value: 20}}
	if _, err := service.CreateRecord(CreateInput{ID: "wf-1", ProgramName: "舞台剧", Performer: "剧社", Category: "戏剧", Label: "初选", Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := service.BeginReview("wf-1", "a"); err != nil {
		t.Fatal(err)
	}
	if _, err := service.ReviewRecord("wf-1", ReviewInput{Score: 91, Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	archived, err := service.ArchiveRecord("wf-1", "a")
	if err != nil || archived.Status.String() != "archived" {
		t.Fatalf("archive: %#v %v", archived, err)
	}
}
