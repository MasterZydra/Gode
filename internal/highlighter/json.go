package highlighter

type JSON struct{}

func (JSON) Highlight(source string) []HighlightedLine {
	lines := splitLines(source)
	for lineIndex, line := range lines {
		lines[lineIndex].Spans = jsonLiteralSpans(line.Text)
	}
	return lines
}

func jsonLiteralSpans(line string) []TokenSpan {
	spans := make([]TokenSpan, 0)
	for index := 0; index < len(line); {
		switch {
		case line[index] == '"':
			start := index
			index++
			for index < len(line) {
				if line[index] == '\\' {
					index += 2
					continue
				}
				if line[index] == '"' {
					index++
					break
				}
				index++
			}
			spans = append(spans, jsonSpan(line, start, index, HighlightString))
		case isJSONNumberStart(line, index):
			start := index
			index = jsonNumberEnd(line, index)
			spans = append(spans, jsonSpan(line, start, index, HighlightNumber))
		case isJSONKeywordStart(line, index):
			start := index
			for index < len(line) && ((line[index] >= 'a' && line[index] <= 'z') || (line[index] >= 'A' && line[index] <= 'Z')) {
				index++
			}
			spans = append(spans, jsonSpan(line, start, index, HighlightKeyword))
		default:
			index++
		}
	}
	return spans
}

func isJSONNumberStart(line string, index int) bool {
	return line[index] >= '0' && line[index] <= '9' ||
		line[index] == '-' && index+1 < len(line) && line[index+1] >= '0' && line[index+1] <= '9'
}

func jsonNumberEnd(line string, index int) int {
	if line[index] == '-' {
		index++
	}
	if line[index] == '0' {
		index++
	} else {
		for index < len(line) && line[index] >= '0' && line[index] <= '9' {
			index++
		}
	}
	if index < len(line) && line[index] == '.' {
		index++
		for index < len(line) && line[index] >= '0' && line[index] <= '9' {
			index++
		}
	}
	if index < len(line) && (line[index] == 'e' || line[index] == 'E') {
		index++
		if index < len(line) && (line[index] == '+' || line[index] == '-') {
			index++
		}
		for index < len(line) && line[index] >= '0' && line[index] <= '9' {
			index++
		}
	}
	return index
}

func isJSONKeywordStart(line string, index int) bool {
	for _, keyword := range []string{"true", "false", "null"} {
		if len(line)-index >= len(keyword) && line[index:index+len(keyword)] == keyword {
			return true
		}
	}
	return false
}

func jsonSpan(line string, start, end int, kind HighlightKind) TokenSpan {
	return TokenSpan{Start: VisualColumn(line[:start]), End: VisualColumn(line[:end]), Kind: kind}
}
