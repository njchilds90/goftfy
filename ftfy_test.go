package goftfy

import "testing"

func TestFixMojibake(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"SÃ£o Paulo", "São Paulo"},
		{"cafÃ©", "café"},
		{"rÃ©sumÃ©", "résumé"},
		{"naÃ¯ve", "naïve"},
	}
	for _, tt := range tests {
		got := Fix(tt.input)
		if got != tt.expected {
			t.Errorf("Fix(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFixHTMLEntities(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"AT&amp;T", "AT&T"},
		{"&lt;b&gt;bold&lt;/b&gt;", "<b>bold</b>"},
		{"It&#8217;s fine", "It\u2019s fine"},
	}
	for _, tt := range tests {
		got := Fix(tt.input)
		if got != tt.expected {
			t.Errorf("Fix(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}

func TestFixLineBreaks(t *testing.T) {
	got := Fix("line1\r\nline2\rline3")
	expected := "line1\nline2\nline3"
	if got != expected {
		t.Errorf("Fix line breaks: got %q, want %q", got, expected)
	}
}

func TestFixControlChars(t *testing.T) {
	got := Fix("hello\x01world\x07!")
	expected := "helloworld!"
	if got != expected {
		t.Errorf("Fix control chars: got %q, want %q", got, expected)
	}
}

func TestFixCurlyQuotes(t *testing.T) {
	opts := DefaultOptions()
	opts.FixCurlyQuotes = true
	got := FixWithOptions("\u201CHello\u201D \u2018world\u2019", opts)
	expected := `"Hello" 'world'`
	if got != expected {
		t.Errorf("Fix curly quotes: got %q, want %q", got, expected)
	}
}

func TestIsValid(t *testing.T) {
	if !IsValid("Hello, world!") {
		t.Error("expected clean text to be valid")
	}
	if IsValid("SÃ£o Paulo") {
		t.Error("expected mojibake to be invalid")
	}
}

func TestFixSlice(t *testing.T) {
	input := []string{"cafÃ©", "hello", "rÃ©sumÃ©"}
	result := FixSlice(input)
	if result[0] != "café" {
		t.Errorf("FixSlice[0]: got %q, want %q", result[0], "café")
	}
	if result[1] != "hello" {
		t.Errorf("FixSlice[1]: got %q, want %q", result[1], "hello")
	}
}

func TestFixMap(t *testing.T) {
	input := map[string]string{
		"city": "SÃ£o Paulo",
		"name": "Alice",
	}
	result := FixMap(input)
	if result["city"] != "São Paulo" {
		t.Errorf("FixMap city: got %q, want %q", result["city"], "São Paulo")
	}
}

func TestQuickFix(t *testing.T) {
	got := QuickFix("SÃ£o Paulo")
	if got != "São Paulo" {
		t.Errorf("QuickFix: got %q, want %q", got, "São Paulo")
	}
}

func TestExplain(t *testing.T) {
	explanation := Explain("SÃ£o Paulo", "São Paulo")
	if explanation == "No changes needed." {
		t.Error("expected explanation to note changes")
	}
}

func TestHasSurrogates(t *testing.T) {
	if HasSurrogates("clean text") {
		t.Error("should not detect surrogates in clean text")
	}
}

func TestHasReplacementChars(t *testing.T) {
	if !HasReplacementChars("bad\uFFFDtext") {
		t.Error("should detect replacement char")
	}
}

func TestRemoveTerminalEscapes(t *testing.T) {
	opts := DefaultOptions()
	opts.RemoveTerminalEscapes = true
	got := FixWithOptions("\x1b[31mred\x1b[0m", opts)
	if got != "red" {
		t.Errorf("terminal escape removal: got %q, want %q", got, "red")
	}
}
