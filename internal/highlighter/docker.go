package highlighter

// This file should be called "dockerfile.go" but than VS Code uses the syntax highlighting for dockerfiles...

import "strings"

type Dockerfile struct{}

var dockerfileInstructions = map[string]bool{
	"ADD":         true,
	"ARG":         true,
	"CMD":         true,
	"COPY":        true,
	"ENTRYPOINT":  true,
	"ENV":         true,
	"EXPOSE":      true,
	"FROM":        true,
	"HEALTHCHECK": true,
	"LABEL":       true,
	"MAINTAINER":  true,
	"ONBUILD":     true,
	"RUN":         true,
	"SHELL":       true,
	"STOPSIGNAL":  true,
	"USER":        true,
	"VOLUME":      true,
	"WORKDIR":     true,
}

func (Dockerfile) Highlight(source string) []HighlightedLine {
	lines := splitLines(source)
	for lineIndex, line := range lines {
		text := line.Text
		start := 0
		for start < len(text) && (text[start] == ' ' || text[start] == '\t') {
			start++
		}
		end := start
		for end < len(text) && text[end] != ' ' && text[end] != '\t' {
			end++
		}
		if dockerfileInstructions[strings.ToUpper(text[start:end])] {
			lines[lineIndex].Spans = append(lines[lineIndex].Spans, TokenSpan{Start: start, End: end, Kind: HighlightKeyword})
		}

		lines[lineIndex].Spans = append(lines[lineIndex].Spans, dockerfileLiteralSpans(text)...)
	}
	return lines
}

func dockerfileLiteralSpans(line string) []TokenSpan {
	spans := make([]TokenSpan, 0)
	for index := 0; index < len(line); index++ {
		switch line[index] {
		case '\'', '"':
			quote := line[index]
			start := index
			index++
			for index < len(line) {
				if line[index] == '\\' {
					index += 2
					continue
				}
				if line[index] == quote {
					index++
					break
				}
				index++
			}
			spans = append(spans, TokenSpan{Start: start, End: index, Kind: HighlightString})
		case '#':
			if index == 0 || line[index-1] == ' ' || line[index-1] == '\t' {
				spans = append(spans, TokenSpan{Start: index, End: len(line), Kind: HighlightComment})
				return spans
			}
		}
	}
	return spans
}
