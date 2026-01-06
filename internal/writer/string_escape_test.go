package writer

import "testing"

// TestEscapePDFString_Basic tests basic cases without special characters.
func TestEscapePDFString_Basic(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"no special characters", "Hello World", "Hello World"},
		{"alphanumeric with spaces", "The quick brown fox 123", "The quick brown fox 123"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_Backslash tests backslash escaping.
func TestEscapePDFString_Backslash(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"single backslash", "\\", "\\\\"},
		{"windows path", "C:\\path\\to\\file", "C:\\\\path\\\\to\\\\file"},
		{"multiple backslashes", "\\\\\\", "\\\\\\\\\\\\"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_Parentheses tests parentheses escaping.
func TestEscapePDFString_Parentheses(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"left parenthesis", "(", "\\("},
		{"right parenthesis", ")", "\\)"},
		{"both parentheses", "(text)", "\\(text\\)"},
		{"parentheses in text", "Price: $50 (USD)", "Price: $50 \\(USD\\)"},
		{"nested parentheses", "((nested))", "\\(\\(nested\\)\\)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_ControlChars tests control character escaping.
func TestEscapePDFString_ControlChars(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"newline", "Line1\nLine2", "Line1\\nLine2"},
		{"carriage return", "Line1\rLine2", "Line1\\rLine2"},
		{"tab", "Col1\tCol2", "Col1\\tCol2"},
		{"backspace", "Text\b", "Text\\b"},
		{"form feed", "Page1\fPage2", "Page1\\fPage2"},
		{"newline and tab", "Line1\nLine2\tColumn", "Line1\\nLine2\\tColumn"},
		{"all control characters", "\n\r\t\b\f", "\\n\\r\\t\\b\\f"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_Combined tests combined special characters.
func TestEscapePDFString_Combined(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"backslash and parentheses", "\\(escape\\)", "\\\\\\(escape\\\\\\)"},
		{"backslash and newline", "Line1\\\nLine2", "Line1\\\\\\nLine2"},
		{"all special characters", "\\(text)\n\r\t", "\\\\\\(text\\)\\n\\r\\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_RealWorld tests real-world examples.
func TestEscapePDFString_RealWorld(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"file path with spaces", "C:\\Program Files\\App", "C:\\\\Program Files\\\\App"},
		{"price with currency", "Total: $100 (including tax)", "Total: $100 \\(including tax\\)"},
		{"multi-line address", "123 Main St\nNew York, NY\n10001", "123 Main St\\nNew York, NY\\n10001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_Unicode tests Unicode and Cyrillic text (pass through).
func TestEscapePDFString_Unicode(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"cyrillic text", "Привет мир", "Привет мир"},
		{"cyrillic with special chars", "Привет (мир)", "Привет \\(мир\\)"},
		{"mixed unicode", "Hello Привет 世界", "Hello Привет 世界"},
		{"emoji", "Hello 👋 World", "Hello 👋 World"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFString_EdgeCases tests edge cases.
func TestEscapePDFString_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"only special characters", "\\()[]{}",
			"\\\\\\(\\)[]{}"},
		{"repeated special characters", "((()))",
			"\\(\\(\\(\\)\\)\\)"},
		{"backslash at end", "path\\", "path\\\\"},
		{"backslash at start", "\\path", "\\\\path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := EscapePDFString(tt.input)
			if got != tt.expected {
				t.Errorf("EscapePDFString(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

// TestEscapePDFStringIdempotence verifies that escaping is correct by checking
// that double-escaping produces the expected result.
func TestEscapePDFStringIdempotence(t *testing.T) {
	input := "Test\\(text)"

	// First escape.
	escaped1 := EscapePDFString(input)
	expected1 := "Test\\\\\\(text\\)"
	if escaped1 != expected1 {
		t.Errorf("First escape: got %q, want %q", escaped1, expected1)
	}

	// Second escape (should escape the already-escaped string).
	escaped2 := EscapePDFString(escaped1)
	expected2 := "Test\\\\\\\\\\\\\\(text\\\\\\)"
	if escaped2 != expected2 {
		t.Errorf("Second escape: got %q, want %q", escaped2, expected2)
	}
}

// BenchmarkEscapePDFString benchmarks the escaping function.
func BenchmarkEscapePDFString(b *testing.B) {
	testCases := []struct {
		name  string
		input string
	}{
		{"no_escapes", "Hello World this is a long string without special characters"},
		{"some_escapes", "Price: $100 (including tax)\nNext line here"},
		{"many_escapes", "\\\\\\(((())))\n\r\t\b\f"},
	}

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = EscapePDFString(tc.input)
			}
		})
	}
}
