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
