package formatter

import "testing"

func TestFormatJSON(t *testing.T) {
	formatted, err := Format("settings.json", `{"name":"Gode","enabled":true,"items":[1,2]}`)
	if err != nil {
		t.Fatal(err)
	}

	want := "{\n\t\"name\": \"Gode\",\n\t\"enabled\": true,\n\t\"items\": [\n\t\t1,\n\t\t2\n\t]\n}"
	if formatted != want {
		t.Fatalf("formatted JSON = %q, want %q", formatted, want)
	}
}

func TestFormatRejectsUnsupportedFile(t *testing.T) {
	if _, err := Format("notes.txt", "text"); err == nil {
		t.Fatal("formatting an unsupported file did not return an error")
	}
}
