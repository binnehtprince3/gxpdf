package creator

import (
	"testing"
)

const testTextHelloWorld = "Hello World"

func TestNewParagraph(t *testing.T) {
	p := NewParagraph(testTextHelloWorld)

	if p.Text() != testTextHelloWorld {
		t.Errorf("Text() = %q, want %q", p.Text(), testTextHelloWorld)
	}

	if p.Font() != Helvetica {
		t.Errorf("Font() = %v, want Helvetica", p.Font())
	}

	if p.FontSize() != 12 {
		t.Errorf("FontSize() = %v, want 12", p.FontSize())
	}

	if p.Color() != Black {
		t.Errorf("Color() = %v, want Black", p.Color())
	}

	if p.Alignment() != AlignLeft {
		t.Errorf("Alignment() = %v, want AlignLeft", p.Alignment())
	}

	if p.LineSpacing() != 1.2 {
		t.Errorf("LineSpacing() = %v, want 1.2", p.LineSpacing())
	}
}

func TestParagraph_SetFont(t *testing.T) {
	p := NewParagraph("Test")

	result := p.SetFont(TimesBold, 18)

	// Check method chaining
	if result != p {
		t.Error("SetFont should return the paragraph for chaining")
	}

	if p.Font() != TimesBold {
		t.Errorf("Font() = %v, want TimesBold", p.Font())
	}

	if p.FontSize() != 18 {
		t.Errorf("FontSize() = %v, want 18", p.FontSize())
	}
}

func TestParagraph_SetColor(t *testing.T) {
	p := NewParagraph("Test")
	red := Color{R: 1, G: 0, B: 0}

	result := p.SetColor(red)

	if result != p {
		t.Error("SetColor should return the paragraph for chaining")
	}

	if p.Color() != red {
		t.Errorf("Color() = %v, want red", p.Color())
	}
}

func TestParagraph_SetAlignment(t *testing.T) {
	p := NewParagraph("Test")

	result := p.SetAlignment(AlignCenter)

	if result != p {
		t.Error("SetAlignment should return the paragraph for chaining")
	}

	if p.Alignment() != AlignCenter {
		t.Errorf("Alignment() = %v, want AlignCenter", p.Alignment())
	}
}

func TestParagraph_SetLineSpacing(t *testing.T) {
	p := NewParagraph("Test")

	result := p.SetLineSpacing(1.5)

	if result != p {
		t.Error("SetLineSpacing should return the paragraph for chaining")
	}

	if p.LineSpacing() != 1.5 {
		t.Errorf("LineSpacing() = %v, want 1.5", p.LineSpacing())
	}
}

func TestParagraph_SetText(t *testing.T) {
	p := NewParagraph("Original")

	result := p.SetText("Updated")

	if result != p {
		t.Error("SetText should return the paragraph for chaining")
	}

	if p.Text() != "Updated" {
		t.Errorf("Text() = %q, want %q", p.Text(), "Updated")
	}
}

func TestParagraph_MethodChaining(t *testing.T) {
	p := NewParagraph("Test").
		SetFont(HelveticaBold, 14).
		SetColor(Blue).
		SetAlignment(AlignRight).
		SetLineSpacing(1.5)

	if p.Font() != HelveticaBold {
		t.Errorf("Font() = %v, want HelveticaBold", p.Font())
	}
	if p.FontSize() != 14 {
		t.Errorf("FontSize() = %v, want 14", p.FontSize())
	}
	if p.Color() != Blue {
		t.Errorf("Color() = %v, want Blue", p.Color())
	}
	if p.Alignment() != AlignRight {
		t.Errorf("Alignment() = %v, want AlignRight", p.Alignment())
	}
	if p.LineSpacing() != 1.5 {
		t.Errorf("LineSpacing() = %v, want 1.5", p.LineSpacing())
	}
}

func TestParagraph_WrapTextLines_EmptyText(t *testing.T) {
	p := NewParagraph("")
	lines := p.WrapTextLines(500)

	if len(lines) != 0 {
		t.Errorf("WrapTextLines for empty text should return empty slice, got %v", lines)
	}
}

func TestParagraph_WrapTextLines_SingleWord(t *testing.T) {
	p := NewParagraph("Hello")
	lines := p.WrapTextLines(500)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if lines[0] != "Hello" {
		t.Errorf("Line[0] = %q, want %q", lines[0], "Hello")
	}
}

func TestParagraph_WrapTextLines_MultipleWords(t *testing.T) {
	p := NewParagraph("Hello World")
	lines := p.WrapTextLines(500)

	if len(lines) != 1 {
		t.Errorf("Expected 1 line, got %d", len(lines))
	}
	if lines[0] != "Hello World" {
		t.Errorf("Line[0] = %q, want %q", lines[0], "Hello World")
	}
}

func TestParagraph_WrapTextLines_ForcedWrap(t *testing.T) {
	// Use Helvetica 12pt: "Hello World" is about 60 points
	// Set width to 40 to force wrap
	p := NewParagraph("Hello World").SetFont(Helvetica, 12)
	lines := p.WrapTextLines(40)

	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d: %v", len(lines), lines)
	}
	if len(lines) >= 2 {
		if lines[0] != "Hello" {
			t.Errorf("Line[0] = %q, want %q", lines[0], "Hello")
		}
		if lines[1] != "World" {
			t.Errorf("Line[1] = %q, want %q", lines[1], "World")
		}
	}
}

func TestParagraph_WrapTextLines_LongParagraph(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	p := NewParagraph(text).SetFont(Helvetica, 12)
	lines := p.WrapTextLines(200) // Force multiple lines

	if len(lines) < 2 {
		t.Errorf("Expected at least 2 lines, got %d: %v", len(lines), lines)
	}

	// Verify all words are preserved
	reconstructed := ""
	for i, line := range lines {
		if i > 0 {
			reconstructed += " "
		}
		reconstructed += line
	}

	if reconstructed != text {
		t.Errorf("Reconstructed text = %q, want %q", reconstructed, text)
	}
}

func TestParagraph_Height(t *testing.T) {
	p := NewParagraph("Hello World").SetFont(Helvetica, 12).SetLineSpacing(1.5)

	ctx := &LayoutContext{
		PageWidth: 595,
		Margins:   Margins{Left: 72, Right: 72},
	}

	// Single line, height = fontSize * lineSpacing = 12 * 1.5 = 18
	height := p.Height(ctx)
	expectedHeight := 18.0

	if height != expectedHeight {
		t.Errorf("Height() = %v, want %v", height, expectedHeight)
	}
}

func TestParagraph_Height_MultipleLines(t *testing.T) {
	text := "The quick brown fox jumps over the lazy dog"
	p := NewParagraph(text).SetFont(Helvetica, 12).SetLineSpacing(1.0)

	ctx := &LayoutContext{
		PageWidth: 200,
		Margins:   Margins{Left: 0, Right: 0},
	}

	lines := p.WrapTextLines(ctx.AvailableWidth())
	expectedHeight := float64(len(lines)) * 12.0 // fontSize * 1.0

	height := p.Height(ctx)

	if height != expectedHeight {
		t.Errorf("Height() = %v, want %v (for %d lines)", height, expectedHeight, len(lines))
	}
}

func TestParagraph_Draw(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()
	p := NewParagraph("Hello World").SetFont(Helvetica, 12)

	initialCursorY := ctx.CursorY
	err = p.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	// Cursor should have advanced
	if ctx.CursorY <= initialCursorY {
		t.Error("Cursor Y should have advanced after drawing")
	}

	// Check that text operation was added
	ops := page.TextOperations()
	if len(ops) != 1 {
		t.Errorf("Expected 1 text operation, got %d", len(ops))
	}

	if len(ops) > 0 && ops[0].Text != "Hello World" {
		t.Errorf("Text operation text = %q, want %q", ops[0].Text, "Hello World")
	}
}

func TestParagraph_Draw_Alignment_Left(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()
	p := NewParagraph("Test").SetFont(Helvetica, 12).SetAlignment(AlignLeft)

	_ = p.Draw(ctx, page)

	ops := page.TextOperations()
	if len(ops) == 0 {
		t.Fatal("No text operations")
	}

	// Left alignment should be at ContentLeft
	if ops[0].X != ctx.ContentLeft() {
		t.Errorf("X = %v, want %v (ContentLeft)", ops[0].X, ctx.ContentLeft())
	}
}

func TestParagraph_Draw_Alignment_Center(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()
	p := NewParagraph("Test").SetFont(Helvetica, 12).SetAlignment(AlignCenter)

	_ = p.Draw(ctx, page)

	ops := page.TextOperations()
	if len(ops) == 0 {
		t.Fatal("No text operations")
	}

	// Center alignment should be between left and right
	if ops[0].X <= ctx.ContentLeft() {
		t.Errorf("X = %v, should be > %v for center alignment", ops[0].X, ctx.ContentLeft())
	}
}

func TestParagraph_Draw_Alignment_Right(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()
	p := NewParagraph("Test").SetFont(Helvetica, 12).SetAlignment(AlignRight)

	_ = p.Draw(ctx, page)

	ops := page.TextOperations()
	if len(ops) == 0 {
		t.Fatal("No text operations")
	}

	// Right alignment should be near ContentRight
	if ops[0].X <= ctx.ContentLeft() {
		t.Errorf("X = %v, should be > %v for right alignment", ops[0].X, ctx.ContentLeft())
	}
}

func TestParagraph_ImplementsDrawable(_ *testing.T) {
	var _ Drawable = (*Paragraph)(nil)
}
