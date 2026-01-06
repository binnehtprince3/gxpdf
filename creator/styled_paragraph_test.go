package creator

import (
	"testing"
)

func TestStyledParagraph_Creation(t *testing.T) {
	sp := NewStyledParagraph()

	if sp == nil {
		t.Fatal("NewStyledParagraph() returned nil")
	}

	if len(sp.chunks) != 0 {
		t.Errorf("Expected 0 chunks, got %d", len(sp.chunks))
	}

	if sp.alignment != AlignLeft {
		t.Errorf("Expected AlignLeft, got %v", sp.alignment)
	}

	if sp.lineSpacing != 1.2 {
		t.Errorf("Expected lineSpacing 1.2, got %f", sp.lineSpacing)
	}
}

func TestStyledParagraph_Append(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("Hello World")

	if len(sp.chunks) != 1 {
		t.Fatalf("Expected 1 chunk, got %d", len(sp.chunks))
	}

	chunk := sp.chunks[0]
	if chunk.Text != "Hello World" {
		t.Errorf("Expected 'Hello World', got '%s'", chunk.Text)
	}

	// Should use default style.
	defaultStyle := DefaultTextStyle()
	if chunk.Style.Font != defaultStyle.Font {
		t.Errorf("Expected font %s, got %s", defaultStyle.Font, chunk.Style.Font)
	}
	if chunk.Style.Size != defaultStyle.Size {
		t.Errorf("Expected size %f, got %f", defaultStyle.Size, chunk.Style.Size)
	}
}

func TestStyledParagraph_AppendStyled(t *testing.T) {
	sp := NewStyledParagraph()

	customStyle := TextStyle{
		Font:  HelveticaBold,
		Size:  14,
		Color: Red,
	}

	sp.AppendStyled("Bold Red", customStyle)

	if len(sp.chunks) != 1 {
		t.Fatalf("Expected 1 chunk, got %d", len(sp.chunks))
	}

	chunk := sp.chunks[0]
	if chunk.Text != "Bold Red" {
		t.Errorf("Expected 'Bold Red', got '%s'", chunk.Text)
	}

	if chunk.Style.Font != HelveticaBold {
		t.Errorf("Expected HelveticaBold, got %s", chunk.Style.Font)
	}
	if chunk.Style.Size != 14 {
		t.Errorf("Expected size 14, got %f", chunk.Style.Size)
	}
	if chunk.Style.Color != Red {
		t.Errorf("Expected Red color, got %v", chunk.Style.Color)
	}
}

func TestStyledParagraph_MultipleFonts(t *testing.T) {
	sp := NewStyledParagraph()

	sp.Append("Normal ")
	sp.AppendStyled("Bold ", TextStyle{Font: HelveticaBold, Size: 12, Color: Black})
	sp.AppendStyled("Italic ", TextStyle{Font: HelveticaOblique, Size: 12, Color: Black})
	sp.Append("Normal again")

	if len(sp.chunks) != 4 {
		t.Fatalf("Expected 4 chunks, got %d", len(sp.chunks))
	}

	// Verify each chunk.
	tests := []struct {
		index int
		text  string
		font  FontName
	}{
		{0, "Normal ", Helvetica},
		{1, "Bold ", HelveticaBold},
		{2, "Italic ", HelveticaOblique},
		{3, "Normal again", Helvetica},
	}

	for _, tt := range tests {
		chunk := sp.chunks[tt.index]
		if chunk.Text != tt.text {
			t.Errorf("Chunk %d: expected '%s', got '%s'", tt.index, tt.text, chunk.Text)
		}
		if chunk.Style.Font != tt.font {
			t.Errorf("Chunk %d: expected font %s, got %s", tt.index, tt.font, chunk.Style.Font)
		}
	}
}

func TestStyledParagraph_MultipleColors(t *testing.T) {
	sp := NewStyledParagraph()

	sp.AppendStyled("Black ", TextStyle{Font: Helvetica, Size: 12, Color: Black})
	sp.AppendStyled("Red ", TextStyle{Font: Helvetica, Size: 12, Color: Red})
	sp.AppendStyled("Blue", TextStyle{Font: Helvetica, Size: 12, Color: Blue})

	if len(sp.chunks) != 3 {
		t.Fatalf("Expected 3 chunks, got %d", len(sp.chunks))
	}

	// Verify colors.
	expectedColors := []Color{Black, Red, Blue}
	for i, expected := range expectedColors {
		chunk := sp.chunks[i]
		if chunk.Style.Color != expected {
			t.Errorf("Chunk %d: expected color %v, got %v", i, expected, chunk.Style.Color)
		}
	}
}

func TestStyledParagraph_TextWrapping_SingleChunk(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("This is a long text that should wrap across multiple lines when the available width is limited.")

	// Mock layout context with limited width.
	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	// Available width = 612 - 72 - 72 = 468 points.
	// With 12pt Helvetica, this should wrap into multiple lines.
	lines := sp.wrapText(ctx.AvailableWidth())

	if len(lines) == 0 {
		t.Fatal("Expected at least 1 line, got 0")
	}

	// Should have multiple lines due to limited width.
	if len(lines) == 1 {
		t.Logf("Warning: Expected multiple lines, got 1. Total width: %f", lines[0].totalWidth)
	}

	// Verify all lines fit within available width.
	for i, line := range lines {
		if line.totalWidth > ctx.AvailableWidth() {
			t.Errorf("Line %d exceeds available width: %f > %f", i, line.totalWidth, ctx.AvailableWidth())
		}
	}
}

func TestStyledParagraph_TextWrapping_AcrossChunks(t *testing.T) {
	sp := NewStyledParagraph()

	// Create text that will wrap, with different styles.
	sp.Append("This is normal text and ")
	sp.AppendStyled("this is bold text ", TextStyle{Font: HelveticaBold, Size: 12, Color: Black})
	sp.Append("and this is normal again and it should wrap properly across multiple styles.")

	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	lines := sp.wrapText(ctx.AvailableWidth())

	if len(lines) == 0 {
		t.Fatal("Expected at least 1 line, got 0")
	}

	// Verify all lines fit within available width.
	for i, line := range lines {
		if line.totalWidth > ctx.AvailableWidth() {
			t.Errorf("Line %d exceeds available width: %f > %f", i, line.totalWidth, ctx.AvailableWidth())
		}

		// Verify line has words.
		if len(line.words) == 0 {
			t.Errorf("Line %d has no words", i)
		}
	}

	// Should have multiple lines.
	if len(lines) < 2 {
		t.Logf("Warning: Expected at least 2 lines, got %d", len(lines))
	}
}

func TestStyledParagraph_Height(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("Single line text")

	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	height := sp.Height(ctx)

	if height <= 0 {
		t.Errorf("Expected height > 0, got %f", height)
	}

	// Height should be approximately fontSize * lineSpacing for single line.
	// Default: 12pt * 1.2 = 14.4pt (approximate, depends on font metrics).
	expectedMin := 12.0
	expectedMax := 20.0

	if height < expectedMin || height > expectedMax {
		t.Errorf("Height %f outside expected range [%f, %f]", height, expectedMin, expectedMax)
	}
}

func TestStyledParagraph_Height_Empty(t *testing.T) {
	sp := NewStyledParagraph()

	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	height := sp.Height(ctx)

	if height != 0 {
		t.Errorf("Expected height 0 for empty paragraph, got %f", height)
	}
}

func TestStyledParagraph_Alignment(t *testing.T) {
	tests := []struct {
		name      string
		alignment Alignment
	}{
		{"Left", AlignLeft},
		{"Center", AlignCenter},
		{"Right", AlignRight},
		{"Justify", AlignJustify},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sp := NewStyledParagraph()
			sp.Append("Test text")
			sp.SetAlignment(tt.alignment)

			if sp.alignment != tt.alignment {
				t.Errorf("Expected alignment %v, got %v", tt.alignment, sp.alignment)
			}
		})
	}
}

func TestStyledParagraph_LineSpacing(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("Test")
	sp.SetLineSpacing(1.5)

	if sp.lineSpacing != 1.5 {
		t.Errorf("Expected lineSpacing 1.5, got %f", sp.lineSpacing)
	}
}

func TestStyledParagraph_MethodChaining(t *testing.T) {
	// Test that all setter methods support method chaining.
	sp := NewStyledParagraph().
		Append("Text 1 ").
		AppendStyled("Text 2", TextStyle{Font: HelveticaBold, Size: 12, Color: Black}).
		SetAlignment(AlignCenter).
		SetLineSpacing(1.5)

	if len(sp.chunks) != 2 {
		t.Errorf("Expected 2 chunks, got %d", len(sp.chunks))
	}

	if sp.alignment != AlignCenter {
		t.Errorf("Expected AlignCenter, got %v", sp.alignment)
	}

	if sp.lineSpacing != 1.5 {
		t.Errorf("Expected lineSpacing 1.5, got %f", sp.lineSpacing)
	}
}

func TestStyledParagraph_ImplementsDrawable(_ *testing.T) {
	var _ Drawable = (*StyledParagraph)(nil)
}

func TestStyledParagraph_SplitChunksIntoWords(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("Hello World")
	sp.AppendStyled("Bold Text", TextStyle{Font: HelveticaBold, Size: 12, Color: Black})

	words := sp.splitChunksIntoWords()

	// Should have 4 words: "Hello", " World", " Bold", " Text".
	// First word has no leading space, others do.
	expectedCount := 4
	if len(words) != expectedCount {
		t.Fatalf("Expected %d words, got %d", expectedCount, len(words))
	}

	// Verify first word has no leading space.
	if words[0].text != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", words[0].text)
	}

	// Verify subsequent words have leading space.
	expectedWords := []string{"Hello", " World", " Bold", " Text"}
	for i, expected := range expectedWords {
		if words[i].text != expected {
			t.Errorf("Word %d: expected '%s', got '%s'", i, expected, words[i].text)
		}
	}

	// Verify styles.
	if words[0].style.Font != Helvetica {
		t.Errorf("Word 0: expected Helvetica, got %s", words[0].style.Font)
	}
	if words[2].style.Font != HelveticaBold {
		t.Errorf("Word 2: expected HelveticaBold, got %s", words[2].style.Font)
	}
}

func TestStyledParagraph_BuildLines(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("Word1 Word2 Word3 Word4")

	words := sp.splitChunksIntoWords()

	// Test with very limited width (forces one word per line).
	lines := sp.buildLines(words, 50.0)

	if len(lines) < 2 {
		t.Errorf("Expected at least 2 lines with limited width, got %d", len(lines))
	}

	// Test with unlimited width (should fit on one line).
	lines = sp.buildLines(words, 1000.0)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line with unlimited width, got %d", len(lines))
	}
}

func TestDefaultTextStyle(t *testing.T) {
	style := DefaultTextStyle()

	if style.Font != Helvetica {
		t.Errorf("Expected Helvetica, got %s", style.Font)
	}

	if style.Size != 12 {
		t.Errorf("Expected size 12, got %f", style.Size)
	}

	if style.Color != Black {
		t.Errorf("Expected Black, got %v", style.Color)
	}
}

func TestStyledParagraph_EmptyText(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("")

	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	lines := sp.wrapText(ctx.AvailableWidth())

	if len(lines) != 0 {
		t.Errorf("Expected 0 lines for empty text, got %d", len(lines))
	}

	height := sp.Height(ctx)
	if height != 0 {
		t.Errorf("Expected height 0 for empty text, got %f", height)
	}
}

func TestStyledParagraph_WhitespaceOnly(t *testing.T) {
	sp := NewStyledParagraph()
	sp.Append("   ")

	ctx := &LayoutContext{
		PageWidth:  612,
		PageHeight: 792,
		Margins: Margins{
			Top:    72,
			Bottom: 72,
			Left:   72,
			Right:  72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	// strings.Fields() removes whitespace, so should have no words.
	words := sp.splitChunksIntoWords()

	if len(words) != 0 {
		t.Errorf("Expected 0 words for whitespace-only text, got %d", len(words))
	}

	lines := sp.wrapText(ctx.AvailableWidth())

	if len(lines) != 0 {
		t.Errorf("Expected 0 lines for whitespace-only text, got %d", len(lines))
	}
}
