package highlighter

import (
	"go/scanner"
	"go/token"
)

type GoLang struct{}

func (GoLang) Highlight(source string) []HighlightedLine {
	lines := splitLines(source)
	fileSet := token.NewFileSet()
	file := fileSet.AddFile("source.go", -1, len(source))
	var lexer scanner.Scanner
	lexer.Init(file, []byte(source), nil, scanner.ScanComments)

	for {
		position, tokenType, literal := lexer.Scan()
		if tokenType == token.EOF {
			break
		}
		if literal == "" {
			continue
		}
		kind, highlight := classifyToken(tokenType, literal)
		if !highlight {
			continue
		}
		location := fileSet.Position(position)
		line := location.Line - 1
		if line < 0 || line >= len(lines) {
			continue
		}
		lineStart := 0
		if location.Offset > 0 {
			for lineStart = location.Offset - 1; lineStart >= 0 && source[lineStart] != '\n'; lineStart-- {
			}
			lineStart++
		}
		appendTokenSpans(lines, line, VisualColumn(source[lineStart:location.Offset]), literal, kind)
	}
	return lines
}

func classifyToken(tokenType token.Token, literal string) (HighlightKind, bool) {
	switch {
	case tokenType == token.COMMENT:
		return HighlightComment, true
	case tokenType == token.STRING || tokenType == token.CHAR:
		return HighlightString, true
	case tokenType == token.INT || tokenType == token.FLOAT || tokenType == token.IMAG:
		return HighlightNumber, true
	case tokenType.IsKeyword():
		return HighlightKeyword, true
	case tokenType == token.IDENT && predeclaredIdentifiers[literal]:
		return HighlightKeyword, true
	default:
		return 0, false
	}
}

var predeclaredIdentifiers = map[string]bool{
	"bool": true, "byte": true, "complex64": true, "complex128": true,
	"error": true, "float32": true, "float64": true, "int": true,
	"int8": true, "int16": true, "int32": true, "int64": true,
	"rune": true, "string": true, "uint": true, "uint8": true,
	"uint16": true, "uint32": true, "uint64": true, "uintptr": true,
	"any": true, "comparable": true, "false": true, "iota": true,
	"nil": true, "true": true,
}

func appendTokenSpans(lines []HighlightedLine, line, column int, literal string, kind HighlightKind) {
	segmentStart := column
	for _, character := range literal {
		if character == '\n' {
			if line >= 0 && line < len(lines) && segmentStart < column {
				lines[line].Spans = append(lines[line].Spans, TokenSpan{Start: segmentStart, End: column, Kind: kind})
			}
			line++
			column = 0
			segmentStart = 0
			continue
		}
		if character == '\t' {
			column += TabWidth - column%TabWidth
		} else {
			column++
		}
	}
	if line >= 0 && line < len(lines) && segmentStart < column {
		lines[line].Spans = append(lines[line].Spans, TokenSpan{Start: segmentStart, End: column, Kind: kind})
	}
}
