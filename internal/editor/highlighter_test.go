package editor

import "testing"

func TestHighlightGoKeywords(t *testing.T) {
	lines := HighlightGo("package main\n\tfunc main() {\n\t\treturn\n\t}")

	if len(lines) != 4 {
		t.Fatalf("got %d lines, want 4", len(lines))
	}
	if len(lines[0].Spans) != 1 || lines[0].Spans[0].Start != 0 {
		t.Fatalf("package keyword was not highlighted: %#v", lines[0].Spans)
	}
	if len(lines[1].Spans) != 1 || lines[1].Spans[0].Start != 4 {
		t.Fatalf("tab-indented func keyword was not highlighted: %#v", lines[1].Spans)
	}
	if lines[1].Text[:4] != "    " {
		t.Fatalf("tab indentation was not expanded: %q", lines[1].Text)
	}
}
