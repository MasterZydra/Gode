package commands

import "testing"

func TestSearchFiltersAndSortsCommands(t *testing.T) {
	registry := NewRegistry()
	registry.Register(Command{Name: "Format Document"})
	registry.Register(Command{Name: "Open Folder"})
	registry.Register(Command{Name: "Format Selection"})

	results := registry.Search("FORMAT")
	if len(results) != 2 {
		t.Fatalf("got %d results, want 2", len(results))
	}
	if results[0].Name != "Format Document" || results[1].Name != "Format Selection" {
		t.Fatalf("got results %q and %q, want sorted format commands", results[0].Name, results[1].Name)
	}
}
