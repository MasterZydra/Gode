package editor

import "testing"

func TestHighlightGoTokenKinds(t *testing.T) {
	lines := HighlightGo("// comment\nvar value string = 42\nvalue = \"text\"\nvalue = nil")

	if len(lines[0].Spans) == 0 || lines[0].Spans[0].Kind != HighlightComment {
		t.Fatalf("comment was not highlighted green: %#v", lines[0].Spans)
	}
	if len(lines[1].Spans) < 3 {
		t.Fatalf("expected keyword, type, and number spans: %#v", lines[1].Spans)
	}
	if lines[1].Spans[0].Kind != HighlightKeyword || lines[1].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("keyword/type spans have wrong kinds: %#v", lines[1].Spans)
	}
	if lines[2].Spans == nil || lines[2].Spans[0].Kind != HighlightString {
		t.Fatalf("string was not highlighted red: %#v", lines[2].Spans)
	}
	if lines[3].Spans == nil || lines[3].Spans[0].Kind != HighlightKeyword {
		t.Fatalf("nil was not highlighted blue: %#v", lines[3].Spans)
	}
}
