package highlighter

import (
	"path/filepath"
	"strings"
)

const TabWidth = 4

type Highlighter interface {
	Highlight(source string) []HighlightedLine
}

type TokenSpan struct {
	Start int
	End   int
	Kind  HighlightKind
}

type HighlightKind int

const (
	HighlightKeyword HighlightKind = iota
	HighlightComment
	HighlightString
	HighlightNumber
	HighlightDiffAdded
	HighlightDiffRemoved
)

type HighlightedLine struct {
	Text  string
	Spans []TokenSpan
}

var HighlightersByExtension = map[string]Highlighter{
	".go":   GoLang{},
	".json": JSON{},
	".php":  PHP{},
}

var highlightersByFileName = map[string]Highlighter{
	"dockerfile": Dockerfile{},
}

func HighlighterForExtension(extension string) Highlighter {
	return HighlightersByExtension[extension]
}

func HighlighterForFile(path string) Highlighter {
	if highlighter := HighlighterForExtension(filepath.Ext(path)); highlighter != nil {
		return highlighter
	}
	return highlightersByFileName[strings.ToLower(filepath.Base(path))]
}

func PlainLines(source string) []HighlightedLine {
	return splitLines(source)
}

func splitLines(source string) []HighlightedLine {
	parts := make([]HighlightedLine, 0)
	start := 0
	for i, r := range source {
		if r == '\n' {
			parts = append(parts, HighlightedLine{Text: expandTabs(source[start:i])})
			start = i + 1
		}
	}
	parts = append(parts, HighlightedLine{Text: expandTabs(source[start:])})
	return parts
}

func expandTabs(line string) string {
	result := make([]rune, 0, len([]rune(line)))
	column := 0
	for _, character := range line {
		if character == '\t' {
			spaces := TabWidth - column%TabWidth
			for range spaces {
				result = append(result, ' ')
			}
			column += spaces
			continue
		}
		result = append(result, character)
		column++
	}
	return string(result)
}

func VisualColumn(text string) int {
	column := 0
	for _, character := range text {
		if character == '\t' {
			column += TabWidth - column%TabWidth
		} else {
			column++
		}
	}
	return column
}
