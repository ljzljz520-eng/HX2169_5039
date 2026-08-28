package workflow

import (
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/store"
)

func TestEngineLifecycle(t *testing.T) {
	persistence, err := store.Open(filepath.Join(t.TempDir(), "workflow.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	engine := NewEngine(persistence, 10)
	flow, err := engine.Start("w1", "登记审核归档", []string{"登记", "审核", "确认", "归档"})
	if err != nil || Progress(flow) != 0.5 {
		t.Fatalf("start: %#v %v", flow, err)
	}
	flow, err = engine.Complete("w1")
	if err != nil || !IsTerminal(flow.Status) || Progress(flow) != 1 {
		t.Fatalf("complete: %#v %v", flow, err)
	}
}
