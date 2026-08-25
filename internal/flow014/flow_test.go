package flow014

import (
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/api"
	"example.com/graduation-showcase/internal/store"
)

func newFlowTest(t *testing.T) *Flow {
	t.Helper()
	persistence, err := store.Open(filepath.Join(t.TempDir(), "flow.db"))
	if err != nil {
		t.Fatal(err)
	}
	service := &api.Service{Store: persistence, Clock: api.FixedClock{Value: 30}}
	t.Cleanup(func() { persistence.Close() })
	return New(service, persistence)
}

func TestWorkflowSearchUpdatePublish(t *testing.T) {
	flow := newFlowTest(t)
	if _, err := flow.Service.CreateRecord(api.CreateInput{ID: "search-1", ProgramName: "民乐", Performer: "乐团", Category: "器乐", Label: "待定", Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	items, err := flow.SearchAndUpdate("民乐", "待定", "重点", "a")
	if err != nil || len(items) != 1 || items[0].Label != "重点" {
		t.Fatalf("items: %#v %v", items, err)
	}
	label, err := flow.RefreshLabel("search-1")
	if err != nil || label != "重点" {
		t.Fatalf("label: %s %v", label, err)
	}
}

func Test2169BusinessRegression(t *testing.T) {
	flow := newFlowTest(t)
	if _, err := flow.Service.CreateRecord(api.CreateInput{ID: "receipt-first", ProgramName: "第一张单据", Performer: "甲", Category: "声乐", Label: "旧标签", Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := flow.Service.CreateRecord(api.CreateInput{ID: "receipt-second", ProgramName: "第二张单据", Performer: "乙", Category: "舞蹈", Label: "旧标签", Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := flow.Service.UpdateLabel("receipt-second", api.UpdateInput{Label: "新标签", Actor: "a"}); err != nil {
		t.Fatal(err)
	}
	if _, err := flow.RefreshLabel("receipt-first"); err != nil {
		t.Fatal(err)
	}
	label, err := flow.RefreshLabel("receipt-second")
	if err != nil {
		t.Fatal(err)
	}
	if label != "新标签" {
		t.Fatalf("刷新后第二张单据标签应为新标签，实际为 %s", label)
	}
}
