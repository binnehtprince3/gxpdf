package creator

import (
	"testing"
)

func TestLayoutContext_AvailableWidth(t *testing.T) {
	tests := []struct {
		name      string
		pageWidth float64
		margins   Margins
		wantWidth float64
	}{
		{
			name:      "A4 with 72pt margins",
			pageWidth: 595,
			margins:   Margins{Left: 72, Right: 72},
			wantWidth: 451,
		},
		{
			name:      "Letter with 36pt margins",
			pageWidth: 612,
			margins:   Margins{Left: 36, Right: 36},
			wantWidth: 540,
		},
		{
			name:      "Zero margins",
			pageWidth: 500,
			margins:   Margins{Left: 0, Right: 0},
			wantWidth: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := &LayoutContext{
				PageWidth: tt.pageWidth,
				Margins:   tt.margins,
			}
			if got := ctx.AvailableWidth(); got != tt.wantWidth {
				t.Errorf("AvailableWidth() = %v, want %v", got, tt.wantWidth)
			}
		})
	}
}

func TestLayoutContext_AvailableHeight(t *testing.T) {
	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Bottom: 72},
		CursorY:    100,
	}

	// Content height = 842 - 72 - 72 = 698
	// Available = 698 - 100 = 598
	want := 598.0
	if got := ctx.AvailableHeight(); got != want {
		t.Errorf("AvailableHeight() = %v, want %v", got, want)
	}
}

func TestLayoutContext_ContentEdges(t *testing.T) {
	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72, Right: 50, Bottom: 36, Left: 72},
	}

	if got := ctx.ContentLeft(); got != 72 {
		t.Errorf("ContentLeft() = %v, want 72", got)
	}

	if got := ctx.ContentRight(); got != 545 {
		t.Errorf("ContentRight() = %v, want 545", got)
	}

	if got := ctx.ContentTop(); got != 770 {
		t.Errorf("ContentTop() = %v, want 770", got)
	}

	if got := ctx.ContentBottom(); got != 36 {
		t.Errorf("ContentBottom() = %v, want 36", got)
	}
}

func TestLayoutContext_CurrentPDFY(t *testing.T) {
	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins:    Margins{Top: 72},
		CursorY:    100,
	}

	// PDF Y = ContentTop - CursorY = (842 - 72) - 100 = 670
	want := 670.0
	if got := ctx.CurrentPDFY(); got != want {
		t.Errorf("CurrentPDFY() = %v, want %v", got, want)
	}
}

func TestLayoutContext_MoveCursor(t *testing.T) {
	ctx := &LayoutContext{
		CursorX: 100,
		CursorY: 50,
	}

	ctx.MoveCursor(20, 30)

	if ctx.CursorX != 120 {
		t.Errorf("CursorX = %v, want 120", ctx.CursorX)
	}
	if ctx.CursorY != 80 {
		t.Errorf("CursorY = %v, want 80", ctx.CursorY)
	}
}

func TestLayoutContext_SetCursor(t *testing.T) {
	ctx := &LayoutContext{
		CursorX: 100,
		CursorY: 50,
	}

	ctx.SetCursor(200, 150)

	if ctx.CursorX != 200 {
		t.Errorf("CursorX = %v, want 200", ctx.CursorX)
	}
	if ctx.CursorY != 150 {
		t.Errorf("CursorY = %v, want 150", ctx.CursorY)
	}
}

func TestLayoutContext_NewLine(t *testing.T) {
	ctx := &LayoutContext{
		PageWidth: 595,
		Margins:   Margins{Left: 72},
		CursorX:   300,
		CursorY:   100,
	}

	ctx.NewLine(14.4)

	if ctx.CursorX != 72 {
		t.Errorf("CursorX = %v, want 72", ctx.CursorX)
	}
	if ctx.CursorY != 114.4 {
		t.Errorf("CursorY = %v, want 114.4", ctx.CursorY)
	}
}

func TestLayoutContext_ResetX(t *testing.T) {
	ctx := &LayoutContext{
		Margins: Margins{Left: 72},
		CursorX: 300,
	}

	ctx.ResetX()

	if ctx.CursorX != 72 {
		t.Errorf("CursorX = %v, want 72", ctx.CursorX)
	}
}

func TestLayoutContext_CanFit(t *testing.T) {
	ctx := &LayoutContext{
		PageHeight: 842,
		Margins:    Margins{Top: 72, Bottom: 72},
		CursorY:    600,
	}

	// Content height = 698, Available = 698 - 600 = 98
	if !ctx.CanFit(50) {
		t.Error("CanFit(50) should return true")
	}

	if !ctx.CanFit(98) {
		t.Error("CanFit(98) should return true")
	}

	if ctx.CanFit(100) {
		t.Error("CanFit(100) should return false")
	}
}

func TestAlignment_Constants(t *testing.T) {
	// Verify alignment constants have expected values
	if AlignLeft != 0 {
		t.Errorf("AlignLeft = %v, want 0", AlignLeft)
	}
	if AlignCenter != 1 {
		t.Errorf("AlignCenter = %v, want 1", AlignCenter)
	}
	if AlignRight != 2 {
		t.Errorf("AlignRight = %v, want 2", AlignRight)
	}
	if AlignJustify != 3 {
		t.Errorf("AlignJustify = %v, want 3", AlignJustify)
	}
}
