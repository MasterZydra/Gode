package formatter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"path/filepath"
)

func Format(path, source string) (string, error) {
	switch filepath.Ext(path) {
	case ".go":
		formatted, err := format.Source([]byte(source))
		if err != nil {
			return "", err
		}
		return string(formatted), nil
	case ".json":
		var formatted bytes.Buffer
		if err := json.Indent(&formatted, []byte(source), "", "\t"); err != nil {
			return "", err
		}
		return formatted.String(), nil
	default:
		return "", fmt.Errorf("formatting is not supported for %s files", filepath.Ext(path))
	}
}
