package creator

import (
	"testing"
)

// TestNewDivision_Defaults verifies the default state of a new division.
func TestNewDivision_Defaults(t *testing.T) {
	div := NewDivision()

	if div == nil {
		t.Fatal("NewDivision() returned nil")
	}

	if div.background != nil {
		t.Error("expected nil background (transparent)")
	}

	if div.border != nil {
		t.Error("expected nil border")
	}

	if div.width != 0 {
		t.Errorf("expected width = 0 (auto), got %f", div.width)
	}

	if div.minHeight != 0 {
		t.Errorf("expected minHeight = 0, got %f", div.minHeight)
	}

	if len(div.drawables) != 0 {
		t.Errorf("expected 0 drawables, got %d", len(div.drawables))
	}
}

// TestDivision_SetBackground verifies background color setting.
func TestDivision_SetBackground(t *testing.T) {
	div := NewDivision()
	color := Red

	result := div.SetBackground(color)

	// Verify fluent API (method chaining).
	if result != div {
		t.Error("SetBackground() should return the division for chaining")
	}

	if div.background == nil {
		t.Fatal("background should not be nil after setting")
	}

	if *div.background != color {
		t.Errorf("expected background %v, got %v", color, *div.background)
	}
}

// TestDivision_SetBorder verifies border setting.
func TestDivision_SetBorder(t *testing.T) {
	div := NewDivision()
	border := Border{Width: 2.0, Color: Black}

	result := div.SetBorder(border)

	// Verify fluent API.
	if result != div {
		t.Error("SetBorder() should return the division for chaining")
	}

	if div.border == nil {
		t.Fatal("border should not be nil after setting")
	}

	if div.border.Width != 2.0 {
		t.Errorf("expected border width 2.0, got %f", div.border.Width)
	}

	if div.border.Color != Black {
		t.Errorf("expected border color Black, got %v", div.border.Color)
	}
}

// TestDivision_SetIndividualBorders verifies individual border setting.
func TestDivision_SetIndividualBorders(t *testing.T) {
	div := NewDivision()

	topBorder := Border{Width: 1.0, Color: Red}
	rightBorder := Border{Width: 2.0, Color: Green}
	bottomBorder := Border{Width: 3.0, Color: Blue}
	leftBorder := Border{Width: 4.0, Color: Yellow}

	div.SetBorderTop(topBorder)
	div.SetBorderRight(rightBorder)
	div.SetBorderBottom(bottomBorder)
	div.SetBorderLeft(leftBorder)

	// Verify all borders.
	verifyBorderTop(t, div)
	verifyBorderRight(t, div)
	verifyBorderBottom(t, div)
	verifyBorderLeft(t, div)
}

// verifyBorderTop checks the top border is set correctly.
func verifyBorderTop(t *testing.T, div *Division) {
	t.Helper()
	if div.borderTop == nil {
		t.Fatal("borderTop should not be nil")
	}
	if div.borderTop.Width != 1.0 || div.borderTop.Color != Red {
		t.Error("borderTop not set correctly")
	}
}

// verifyBorderRight checks the right border is set correctly.
func verifyBorderRight(t *testing.T, div *Division) {
	t.Helper()
	if div.borderRight == nil {
		t.Fatal("borderRight should not be nil")
	}
	if div.borderRight.Width != 2.0 || div.borderRight.Color != Green {
		t.Error("borderRight not set correctly")
	}
}

// verifyBorderBottom checks the bottom border is set correctly.
func verifyBorderBottom(t *testing.T, div *Division) {
	t.Helper()
	if div.borderBottom == nil {
		t.Fatal("borderBottom should not be nil")
	}
	if div.borderBottom.Width != 3.0 || div.borderBottom.Color != Blue {
		t.Error("borderBottom not set correctly")
	}
}

// verifyBorderLeft checks the left border is set correctly.
func verifyBorderLeft(t *testing.T, div *Division) {
	t.Helper()
	if div.borderLeft == nil {
		t.Fatal("borderLeft should not be nil")
	}
	if div.borderLeft.Width != 4.0 || div.borderLeft.Color != Yellow {
		t.Error("borderLeft not set correctly")
	}
}

// TestDivision_SetPadding verifies padding setting.
func TestDivision_SetPadding(t *testing.T) {
	div := NewDivision()

	result := div.SetPadding(10, 20, 30, 40)

	// Verify fluent API.
	if result != div {
		t.Error("SetPadding() should return the division for chaining")
	}

	if div.padding.Top != 10 {
		t.Errorf("expected padding.Top = 10, got %f", div.padding.Top)
	}
	if div.padding.Right != 20 {
		t.Errorf("expected padding.Right = 20, got %f", div.padding.Right)
	}
	if div.padding.Bottom != 30 {
		t.Errorf("expected padding.Bottom = 30, got %f", div.padding.Bottom)
	}
	if div.padding.Left != 40 {
		t.Errorf("expected padding.Left = 40, got %f", div.padding.Left)
	}
}

// TestDivision_SetPaddingAll verifies uniform padding setting.
func TestDivision_SetPaddingAll(t *testing.T) {
	div := NewDivision()

	result := div.SetPaddingAll(15)

	// Verify fluent API.
	if result != div {
		t.Error("SetPaddingAll() should return the division for chaining")
	}

	if div.padding.Top != 15 || div.padding.Right != 15 ||
		div.padding.Bottom != 15 || div.padding.Left != 15 {
		t.Error("SetPaddingAll() should set all sides to the same value")
	}
}

// TestDivision_SetMargins verifies margin setting.
func TestDivision_SetMargins(t *testing.T) {
	div := NewDivision()
	margins := Margins{Top: 5, Right: 10, Bottom: 15, Left: 20}

	result := div.SetMargins(margins)

	// Verify fluent API.
	if result != div {
		t.Error("SetMargins() should return the division for chaining")
	}

	if div.margins != margins {
		t.Errorf("expected margins %v, got %v", margins, div.margins)
	}
}

// TestDivision_SetWidth verifies width setting.
func TestDivision_SetWidth(t *testing.T) {
	div := NewDivision()

	result := div.SetWidth(300)

	// Verify fluent API.
	if result != div {
		t.Error("SetWidth() should return the division for chaining")
	}

	if div.width != 300 {
		t.Errorf("expected width = 300, got %f", div.width)
	}
}

// TestDivision_SetMinHeight verifies minimum height setting.
func TestDivision_SetMinHeight(t *testing.T) {
	div := NewDivision()

	result := div.SetMinHeight(100)

	// Verify fluent API.
	if result != div {
		t.Error("SetMinHeight() should return the division for chaining")
	}

	if div.minHeight != 100 {
		t.Errorf("expected minHeight = 100, got %f", div.minHeight)
	}
}

// TestDivision_Add verifies adding drawables.
func TestDivision_Add(t *testing.T) {
	div := NewDivision()
	para := NewParagraph("test")

	result := div.Add(para)

	// Verify fluent API.
	if result != div {
		t.Error("Add() should return the division for chaining")
	}

	if len(div.drawables) != 1 {
		t.Fatalf("expected 1 drawable, got %d", len(div.drawables))
	}

	if div.drawables[0] != para {
		t.Error("drawable not stored correctly")
	}
}

// TestDivision_Clear verifies clearing drawables.
func TestDivision_Clear(t *testing.T) {
	div := NewDivision()
	div.Add(NewParagraph("test1"))
	div.Add(NewParagraph("test2"))

	if len(div.drawables) != 2 {
		t.Fatalf("expected 2 drawables before clear, got %d", len(div.drawables))
	}

	result := div.Clear()

	// Verify fluent API.
	if result != div {
		t.Error("Clear() should return the division for chaining")
	}

	if len(div.drawables) != 0 {
		t.Errorf("expected 0 drawables after clear, got %d", len(div.drawables))
	}
}

// TestDivision_Height_Empty verifies height calculation with no content.
func TestDivision_Height_Empty(t *testing.T) {
	div := NewDivision()
	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	height := div.Height(ctx)

	// Empty division should have 0 height.
	if height != 0 {
		t.Errorf("expected height = 0 for empty division, got %f", height)
	}
}

// TestDivision_Height_WithContent verifies height calculation with content.
func TestDivision_Height_WithContent(t *testing.T) {
	div := NewDivision()
	para := NewParagraph("test")
	para.SetFont(Helvetica, 12)
	div.Add(para)

	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	height := div.Height(ctx)

	// Should be > 0 with content.
	if height <= 0 {
		t.Errorf("expected height > 0 with content, got %f", height)
	}
}

// TestDivision_Height_WithPadding verifies height includes padding.
func TestDivision_Height_WithPadding(t *testing.T) {
	div1 := NewDivision()
	para1 := NewParagraph("test")
	para1.SetFont(Helvetica, 12)
	div1.Add(para1)

	div2 := NewDivision()
	para2 := NewParagraph("test")
	para2.SetFont(Helvetica, 12)
	div2.Add(para2)
	div2.SetPaddingAll(20)

	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	height1 := div1.Height(ctx)
	height2 := div2.Height(ctx)

	// div2 should be taller by 40 (20 top + 20 bottom).
	expectedDiff := 40.0
	actualDiff := height2 - height1

	if actualDiff != expectedDiff {
		t.Errorf("expected height difference of %f, got %f", expectedDiff, actualDiff)
	}
}

// TestDivision_Height_WithBorder verifies height includes borders.
func TestDivision_Height_WithBorder(t *testing.T) {
	div1 := NewDivision()
	para1 := NewParagraph("test")
	para1.SetFont(Helvetica, 12)
	div1.Add(para1)

	div2 := NewDivision()
	para2 := NewParagraph("test")
	para2.SetFont(Helvetica, 12)
	div2.Add(para2)
	div2.SetBorder(Border{Width: 5, Color: Black})

	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	height1 := div1.Height(ctx)
	height2 := div2.Height(ctx)

	// div2 should be taller by 10 (5 top + 5 bottom).
	expectedDiff := 10.0
	actualDiff := height2 - height1

	if actualDiff != expectedDiff {
		t.Errorf("expected height difference of %f, got %f", expectedDiff, actualDiff)
	}
}

// TestDivision_Height_MinHeight verifies minimum height enforcement.
func TestDivision_Height_MinHeight(t *testing.T) {
	div := NewDivision()
	div.SetMinHeight(100)

	// No content, but minHeight = 100.
	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	height := div.Height(ctx)

	if height != 100 {
		t.Errorf("expected height = 100 (minHeight), got %f", height)
	}
}

// TestDivision_Draw verifies basic drawing without errors.
func TestDivision_Draw(t *testing.T) {
	creator := New()
	page, err := creator.NewPage()
	if err != nil {
		t.Fatalf("NewPage() failed: %v", err)
	}

	div := NewDivision()
	para := NewParagraph("test")
	para.SetFont(Helvetica, 12)
	div.Add(para)

	ctx := page.GetLayoutContext()
	err = div.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}
}

// TestDivision_Draw_WithBackground verifies background drawing.
func TestDivision_Draw_WithBackground(t *testing.T) {
	creator := New()
	page, err := creator.NewPage()
	if err != nil {
		t.Fatalf("NewPage() failed: %v", err)
	}

	div := NewDivision()
	div.SetBackground(LightGray)
	para := NewParagraph("test")
	para.SetFont(Helvetica, 12)
	div.Add(para)

	ctx := page.GetLayoutContext()
	err = div.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() with background returned error: %v", err)
	}

	// Verify graphics operation was added for background.
	graphicsOps := page.GraphicsOperations()
	if len(graphicsOps) == 0 {
		t.Error("expected graphics operations for background, got none")
	}
}

// TestDivision_Draw_WithBorder verifies border drawing.
func TestDivision_Draw_WithBorder(t *testing.T) {
	creator := New()
	page, err := creator.NewPage()
	if err != nil {
		t.Fatalf("NewPage() failed: %v", err)
	}

	div := NewDivision()
	div.SetBorder(Border{Width: 1, Color: Black})
	para := NewParagraph("test")
	para.SetFont(Helvetica, 12)
	div.Add(para)

	ctx := page.GetLayoutContext()
	err = div.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() with border returned error: %v", err)
	}

	// Verify graphics operations were added for borders (4 lines).
	graphicsOps := page.GraphicsOperations()
	if len(graphicsOps) < 4 {
		t.Errorf("expected at least 4 graphics operations for borders, got %d", len(graphicsOps))
	}
}

// TestDivision_MethodChaining verifies fluent API method chaining.
func TestDivision_MethodChaining(t *testing.T) {
	div := NewDivision().
		SetBackground(White).
		SetBorder(Border{Width: 1, Color: Gray}).
		SetPaddingAll(10).
		SetMargins(Margins{Top: 5, Right: 5, Bottom: 5, Left: 5}).
		SetWidth(200).
		SetMinHeight(50).
		Add(NewParagraph("test"))

	// Verify all settings were applied.
	if div.background == nil {
		t.Error("background not set through chaining")
	}
	if div.border == nil {
		t.Error("border not set through chaining")
	}
	if div.padding.Top != 10 {
		t.Error("padding not set through chaining")
	}
	if div.margins.Top != 5 {
		t.Error("margins not set through chaining")
	}
	if div.width != 200 {
		t.Error("width not set through chaining")
	}
	if div.minHeight != 50 {
		t.Error("minHeight not set through chaining")
	}
	if len(div.drawables) != 1 {
		t.Error("drawables not added through chaining")
	}
}

// TestDivision_ImplementsDrawable verifies Division implements Drawable.
func TestDivision_ImplementsDrawable(_ *testing.T) {
	var _ Drawable = (*Division)(nil)
}

// TestDivision_ContentWidth verifies content width calculation.
func TestDivision_ContentWidth(t *testing.T) {
	div := NewDivision()
	div.SetPadding(5, 10, 5, 15)
	div.SetBorder(Border{Width: 2, Color: Black})

	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 72, Bottom: 72, Left: 72},
	}

	// Available width = 595 - 72 - 72 = 451.
	// Content width = 451 - 10 - 15 - 2 - 2 = 422.
	expectedWidth := 451.0 - 10 - 15 - 2 - 2
	contentWidth := div.ContentWidth(ctx)

	if contentWidth != expectedWidth {
		t.Errorf("expected content width %f, got %f", expectedWidth, contentWidth)
	}
}

// TestDivision_GetEffectiveBorder verifies border override logic.
func TestDivision_GetEffectiveBorder(t *testing.T) {
	div := NewDivision()
	defaultBorder := Border{Width: 1, Color: Black}
	topBorder := Border{Width: 3, Color: Red}

	div.SetBorder(defaultBorder)
	div.SetBorderTop(topBorder)

	// Top should use individual border.
	if got := div.getEffectiveBorderTop(); got == nil || got.Width != 3 {
		t.Error("getEffectiveBorderTop() should return individual border")
	}

	// Right should use default border.
	if got := div.getEffectiveBorderRight(); got == nil || got.Width != 1 {
		t.Error("getEffectiveBorderRight() should return default border")
	}
}

// TestDivision_Drawables verifies getting drawables.
func TestDivision_Drawables(t *testing.T) {
	div := NewDivision()
	para1 := NewParagraph("test1")
	para2 := NewParagraph("test2")

	div.Add(para1).Add(para2)

	drawables := div.Drawables()

	if len(drawables) != 2 {
		t.Fatalf("expected 2 drawables, got %d", len(drawables))
	}

	if drawables[0] != para1 || drawables[1] != para2 {
		t.Error("Drawables() did not return correct order")
	}
}

// TestDivision_Background verifies getting background.
func TestDivision_Background(t *testing.T) {
	div := NewDivision()

	// Initially nil.
	if div.Background() != nil {
		t.Error("Background() should return nil initially")
	}

	// After setting.
	div.SetBackground(Red)
	if bg := div.Background(); bg == nil || *bg != Red {
		t.Error("Background() should return set color")
	}
}

// TestDivision_Padding verifies getting padding.
func TestDivision_Padding(t *testing.T) {
	div := NewDivision()
	div.SetPaddingAll(20)

	padding := div.Padding()

	if padding.Top != 20 || padding.Right != 20 ||
		padding.Bottom != 20 || padding.Left != 20 {
		t.Error("Padding() did not return correct values")
	}
}
