package highlighter

import "strings"

type Diff struct {
	Base Highlighter
}

func NewDiffHighlighter(base Highlighter) Highlighter {
	return Diff{Base: base}
}

func (d Diff) Highlight(source string) []HighlightedLine {
	rawLines := strings.Split(source, "\n")
	contentLines := make([]string, len(rawLines))
	for index, line := range rawLines {
		if isDiffContentLine(line) {
			contentLines[index] = line[1:]
		} else {
			contentLines[index] = line
		}
	}

	lines := splitLines(source)
	if d.Base != nil {
		baseLines := d.Base.Highlight(strings.Join(contentLines, "\n"))
		for index := range lines {
			if index >= len(baseLines) {
				break
			}
			for _, span := range baseLines[index].Spans {
				if isDiffContentLine(rawLines[index]) {
					span.Start = diffColumn(rawLines[index], span.Start)
					span.End = diffColumn(rawLines[index], span.End)
				}
				lines[index].Spans = append(lines[index].Spans, span)
			}
		}
	}

	for index, line := range rawLines {
		kind, ok := diffLineKind(line)
		if !ok {
			continue
		}
		lines[index].Spans = append([]TokenSpan{{Start: 0, End: len([]rune(lines[index].Text)), Kind: kind}}, lines[index].Spans...)
	}
	return lines
}

func diffColumn(line string, contentColumn int) int {
	rawRunes := []rune(line)
	contentRunes := rawRunes[1:]
	rawColumn := VisualColumn(string(rawRunes[:1]))
	contentVisualColumn := 0
	for _, character := range contentRunes {
		if contentVisualColumn >= contentColumn {
			return rawColumn
		}
		contentVisualColumn = nextVisualColumn(contentVisualColumn, character)
		rawColumn = nextVisualColumn(rawColumn, character)
	}
	return rawColumn
}

func nextVisualColumn(column int, character rune) int {
	if character == '\t' {
		return column + TabWidth - column%TabWidth
	}
	return column + 1
}

func isDiffContentLine(line string) bool {
	if len(line) == 0 {
		return false
	}
	if strings.HasPrefix(line, "+++") || strings.HasPrefix(line, "---") {
		return false
	}
	return line[0] == ' ' || line[0] == '+' || line[0] == '-'
}

func diffLineKind(line string) (HighlightKind, bool) {
	if !isDiffContentLine(line) {
		return 0, false
	}
	switch line[0] {
	case '+':
		return HighlightDiffAdded, true
	case '-':
		return HighlightDiffRemoved, true
	default:
		return 0, false
	}
}
