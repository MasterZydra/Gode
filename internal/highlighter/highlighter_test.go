package highlighter

import "testing"

func TestHighlightGoKeywords(t *testing.T) {
	lines := GoLang{}.Highlight("package main\n\tfunc main() {\n\t\treturn\n\t}")

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
	lines := GoLang{}.Highlight("// comment\nvar value string = 42\nvalue = \"text\"\nvalue = nil")

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

func TestHighlightJSONLiteralsAndKeywords(t *testing.T) {
	lines := JSON{}.Highlight("{\n\t\"message\": \"hello \\\"world\\\"\",\n\t\"count\": -12.5e+2,\n\t\"enabled\": true,\n\t\"nothing\": null\n}")

	if len(lines) != 6 {
		t.Fatalf("got %d lines, want 6", len(lines))
	}
	if len(lines[1].Spans) != 2 || lines[1].Spans[0].Kind != HighlightString || lines[1].Spans[1].Kind != HighlightString {
		t.Fatalf("string spans = %#v", lines[1].Spans)
	}
	if lines[2].Spans[0] != (TokenSpan{Start: 4, End: 11, Kind: HighlightString}) || lines[2].Spans[1] != (TokenSpan{Start: 13, End: 21, Kind: HighlightNumber}) {
		t.Fatalf("number spans = %#v", lines[2].Spans)
	}
	if len(lines[3].Spans) != 2 || lines[3].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("boolean line spans = %#v", lines[3].Spans)
	}
	if len(lines[4].Spans) != 2 || lines[4].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("null line spans = %#v", lines[4].Spans)
	}
}

func TestHighlighterForExtensionRecognizesJSON(t *testing.T) {
	if _, ok := HighlighterForExtension(".json").(JSON); !ok {
		t.Fatal("JSON extension did not use the JSON highlighter")
	}
}

func TestHighlightPHPTokenKinds(t *testing.T) {
	lines := PHP{}.Highlight("<?php\n// comment\nfunction greet(string $name) {\n\treturn \"Hello\" . 42;\n}")

	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5", len(lines))
	}
	if len(lines[1].Spans) == 0 || lines[1].Spans[0].Kind != HighlightComment {
		t.Fatalf("comment was not highlighted: %#v", lines[1].Spans)
	}
	if len(lines[2].Spans) < 2 || lines[2].Spans[0].Kind != HighlightKeyword || lines[2].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("PHP keywords were not highlighted: %#v", lines[2].Spans)
	}
	if len(lines[3].Spans) < 2 || lines[3].Spans[0].Kind != HighlightKeyword || lines[3].Spans[1].Kind != HighlightString {
		t.Fatalf("PHP return/string spans were not highlighted: %#v", lines[3].Spans)
	}
	if lines[3].Spans[len(lines[3].Spans)-1].Kind != HighlightNumber {
		t.Fatalf("PHP number was not highlighted: %#v", lines[3].Spans)
	}
}

func TestHighlighterForExtensionRecognizesPHP(t *testing.T) {
	if _, ok := HighlighterForExtension(".php").(PHP); !ok {
		t.Fatal("PHP extension did not use the PHP highlighter")
	}
}

func TestHighlightDockerfileKeywordsAndComments(t *testing.T) {
	lines := Dockerfile{}.Highlight("FROM alpine:latest\n# comment\nRUN echo hello # inline comment\nENV VALUE=\"#not a comment\"")

	if lines[0].Spans[0] != (TokenSpan{Start: 0, End: 4, Kind: HighlightKeyword}) {
		t.Fatalf("FROM span = %#v", lines[0].Spans)
	}
	if lines[1].Spans[0].Kind != HighlightComment {
		t.Fatalf("comment span = %#v", lines[1].Spans)
	}
	if len(lines[2].Spans) != 2 || lines[2].Spans[0].Kind != HighlightKeyword || lines[2].Spans[1].Kind != HighlightComment {
		t.Fatalf("RUN/comment spans = %#v", lines[2].Spans)
	}
	if len(lines[3].Spans) != 2 || lines[3].Spans[0].Kind != HighlightKeyword || lines[3].Spans[1].Kind != HighlightString {
		t.Fatalf("quoted value was not highlighted as a string: %#v", lines[3].Spans)
	}
}

func TestHighlightDockerfileStringLiterals(t *testing.T) {
	lines := Dockerfile{}.Highlight(`ENV GREETING="hello world" PATH='a\'b'`)

	if len(lines[0].Spans) != 3 {
		t.Fatalf("got spans = %#v, want keyword and two strings", lines[0].Spans)
	}
	for _, span := range lines[0].Spans[1:] {
		if span.Kind != HighlightString {
			t.Fatalf("span = %#v, want string", span)
		}
	}
}

func TestHighlighterForFileRecognizesDockerfile(t *testing.T) {
	if _, ok := HighlighterForFile("/tmp/Dockerfile").(Dockerfile); !ok {
		t.Fatal("Dockerfile did not use the Dockerfile highlighter")
	}
	if _, ok := HighlighterForFile("/tmp/dockerfile").(Dockerfile); !ok {
		t.Fatal("lowercase dockerfile did not use the Dockerfile highlighter")
	}
}

func TestHighlighterForExtensionUsesRegistry(t *testing.T) {
	original := HighlightersByExtension[".test"]
	HighlightersByExtension[".test"] = GoLang{}
	defer func() {
		if original == nil {
			delete(HighlightersByExtension, ".test")
			return
		}
		HighlightersByExtension[".test"] = original
	}()

	if _, ok := HighlighterForExtension(".test").(GoLang); !ok {
		t.Fatal("extension registry did not return the assigned highlighter")
	}
	if HighlighterForExtension(".unknown") != nil {
		t.Fatal("unknown extension should not have a highlighter")
	}
}

func TestDiffHighlighterAddsLineTokensAndPreservesSyntax(t *testing.T) {
	lines := Diff{Base: GoLang{}}.Highlight("@@ -1 +1 @@\n-package old\n+package new\n+\treturn")

	if lines[1].Spans[0] != (TokenSpan{Start: 0, End: 12, Kind: HighlightDiffRemoved}) {
		t.Fatalf("removed line span = %#v", lines[1].Spans)
	}
	if lines[2].Spans[0] != (TokenSpan{Start: 0, End: 12, Kind: HighlightDiffAdded}) {
		t.Fatalf("added line span = %#v", lines[2].Spans)
	}
	if len(lines[2].Spans) < 2 || lines[2].Spans[1].Start != 1 || lines[2].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("added line syntax spans = %#v", lines[2].Spans)
	}
	if len(lines[3].Spans) < 2 || lines[3].Spans[1].Start != 4 || lines[3].Spans[1].Kind != HighlightKeyword {
		t.Fatalf("context line syntax spans = %#v", lines[3].Spans)
	}
}
