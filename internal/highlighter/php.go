package highlighter

import (
	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
)

type PHP struct{}

func (PHP) Highlight(source string) []HighlightedLine {
	lines := splitLines(source)
	lexer := lexers.Get("php")
	if lexer == nil {
		return lines
	}
	iterator, err := lexer.Tokenise(nil, source)
	if err != nil {
		return lines
	}

	offset := 0
	for token := iterator(); token != chroma.EOF; token = iterator() {
		if token.Value == "" {
			continue
		}
		if kind, highlight := classifyPHPToken(token.Type); highlight {
			line, column := sourcePosition(source, offset)
			appendTokenSpans(lines, line, column, token.Value, kind)
		}
		offset += len(token.Value)
	}
	return lines
}

func classifyPHPToken(tokenType chroma.TokenType) (HighlightKind, bool) {
	switch {
	case tokenType >= chroma.Comment && tokenType < chroma.CommentPreproc:
		return HighlightComment, true
	case tokenType >= chroma.LiteralString && tokenType < chroma.LiteralNumber:
		return HighlightString, true
	case tokenType >= chroma.LiteralNumber && tokenType < chroma.Operator:
		return HighlightNumber, true
	case tokenType >= chroma.Keyword && tokenType < chroma.Name:
		return HighlightKeyword, true
	case tokenType >= chroma.NameBuiltin && tokenType < chroma.NameVariable:
		return HighlightKeyword, true
	case tokenType == chroma.NameFunction || tokenType == chroma.NameOther:
		return HighlightKeyword, true
	default:
		return 0, false
	}
}

func sourcePosition(source string, offset int) (int, int) {
	line := 0
	lineStart := 0
	for index, character := range source[:offset] {
		if character == '\n' {
			line++
			lineStart = index + 1
		}
	}
	return line, VisualColumn(source[lineStart:offset])
}
