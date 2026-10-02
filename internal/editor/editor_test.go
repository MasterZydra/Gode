package editor

import (
	"os"
	"path/filepath"
	"testing"

	"fyne.io/fyne/v2/app"
)

func TestMain(main *testing.M) {
	app.NewWithID("com.gode.editor.editor-tests")
	os.Exit(main.Run())
}

func TestFormatGoFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "main.go")
	if err := os.WriteFile(path, []byte("package main\nfunc main(){println(\"hello\")}\n"), 0600); err != nil {
		t.Fatal(err)
	}

	editor := New()
	if err := editor.Load(path); err != nil {
		t.Fatal(err)
	}
	if err := editor.Format(); err != nil {
		t.Fatal(err)
	}

	if got, want := editor.Widget.Text(), "package main\n\nfunc main() { println(\"hello\") }\n"; got != want {
		t.Fatalf("formatted text = %q, want %q", got, want)
	}
	if !editor.Dirty {
		t.Fatal("formatted editor is not dirty")
	}
}

func TestLoadAndSavePreservesLineEndings(t *testing.T) {
	tests := []struct {
		name       string
		lineEnding string
	}{
		{name: "LF", lineEnding: "\n"},
		{name: "CRLF", lineEnding: "\r\n"},
		{name: "CR", lineEnding: "\r"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "text.txt")
			original := "first" + test.lineEnding + "second" + test.lineEnding
			if err := os.WriteFile(path, []byte(original), 0600); err != nil {
				t.Fatal(err)
			}

			editor := New()
			if err := editor.Load(path); err != nil {
				t.Fatal(err)
			}
			if got, want := editor.Widget.Text(), "first\nsecond\n"; got != want {
				t.Fatalf("editor text = %q, want %q", got, want)
			}

			editor.Widget.SetText("first\nsecond\nadded\n")
			if err := editor.Save(); err != nil {
				t.Fatal(err)
			}

			if got, want := readFile(t, path), "first"+test.lineEnding+"second"+test.lineEnding+"added"+test.lineEnding; got != want {
				t.Fatalf("saved text = %q, want %q", got, want)
			}
		})
	}
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(contents)
}
