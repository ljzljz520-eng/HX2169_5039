package main

import (
	"path/filepath"
	"testing"

	"example.com/graduation-showcase/internal/api"
	"example.com/graduation-showcase/internal/store"
)

func TestDemoEntryPoint(t *testing.T) {
	persistence, err := store.Open(filepath.Join(t.TempDir(), "demo.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer persistence.Close()
	if err := runDemo(&api.Service{Store: persistence, Clock: api.FixedClock{Value: 1}}, persistence); err != nil {
		t.Fatal(err)
	}
}
