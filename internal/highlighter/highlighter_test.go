package highlighter

import "testing"

func TestHighlightGoKeywords(t *testing.T) {
	lines := GoHighlighter{}.Highlight("package main\n\tfunc main() {\n\t\treturn\n\t}")

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

func TestHighlightGoTokenKinds(t *testing.T) {
	lines := GoHighlighter{}.Highlight("// comment\nvar value string = 42\nvalue = \"text\"\nvalue = nil")

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

func TestHighlighterForExtensionUsesRegistry(t *testing.T) {
	original := HighlightersByExtension[".test"]
	HighlightersByExtension[".test"] = GoHighlighter{}
	defer func() {
		if original == nil {
			delete(HighlightersByExtension, ".test")
			return
		}
		HighlightersByExtension[".test"] = original
	}()

	if _, ok := HighlighterForExtension(".test").(GoHighlighter); !ok {
		t.Fatal("extension registry did not return the assigned highlighter")
	}
	if HighlighterForExtension(".unknown") != nil {
		t.Fatal("unknown extension should not have a highlighter")
	}
}
