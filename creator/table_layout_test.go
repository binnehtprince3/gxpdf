package creator

import (
	"testing"
)

func TestNewTableLayout(t *testing.T) {
	table := NewTableLayout(3)

	if table.ColumnCount() != 3 {
		t.Errorf("ColumnCount() = %v, want 3", table.ColumnCount())
	}

	if table.RowCount() != 0 {
		t.Errorf("RowCount() = %v, want 0", table.RowCount())
	}

	if table.HeaderRowCount() != 0 {
		t.Errorf("HeaderRowCount() = %v, want 0", table.HeaderRowCount())
	}
}

func TestNewTableLayout_MinColumns(t *testing.T) {
	table := NewTableLayout(0)

	// Should default to 1 column minimum
	if table.ColumnCount() != 1 {
		t.Errorf("ColumnCount() = %v, want 1", table.ColumnCount())
	}
}

func TestTableLayout_SetColumnWidths(t *testing.T) {
	table := NewTableLayout(3)

	result := table.SetColumnWidths(100, 150, 200)

	// Check method chaining
	if result != table {
		t.Error("SetColumnWidths should return the table for chaining")
	}
}

func TestTableLayout_SetBorder(t *testing.T) {
	table := NewTableLayout(2)

	result := table.SetBorder(1.0, Black)

	if result != table {
		t.Error("SetBorder should return the table for chaining")
	}

	if table.borderWidth != 1.0 {
		t.Errorf("borderWidth = %v, want 1.0", table.borderWidth)
	}

	if table.borderColor == nil {
		t.Error("borderColor should not be nil")
	}
}

func TestTableLayout_SetCellPadding(t *testing.T) {
	table := NewTableLayout(2)

	result := table.SetCellPadding(8.0)

	if result != table {
		t.Error("SetCellPadding should return the table for chaining")
	}

	if table.cellPadding != 8.0 {
		t.Errorf("cellPadding = %v, want 8.0", table.cellPadding)
	}
}

func TestTableLayout_AddHeaderRow(t *testing.T) {
	table := NewTableLayout(3)

	result := table.AddHeaderRow("Name", "Age", "City")

	if result != table {
		t.Error("AddHeaderRow should return the table for chaining")
	}

	if table.RowCount() != 1 {
		t.Errorf("RowCount() = %v, want 1", table.RowCount())
	}

	if table.HeaderRowCount() != 1 {
		t.Errorf("HeaderRowCount() = %v, want 1", table.HeaderRowCount())
	}

	// Check header row uses bold font
	if len(table.rows) > 0 && len(table.rows[0].Cells) > 0 {
		if table.rows[0].Cells[0].Font != HelveticaBold {
			t.Errorf("Header font = %v, want HelveticaBold", table.rows[0].Cells[0].Font)
		}
	}
}

func TestTableLayout_AddRow(t *testing.T) {
	table := NewTableLayout(3)

	result := table.AddRow("Alice", "30", "New York")

	if result != table {
		t.Error("AddRow should return the table for chaining")
	}

	if table.RowCount() != 1 {
		t.Errorf("RowCount() = %v, want 1", table.RowCount())
	}

	if table.HeaderRowCount() != 0 {
		t.Errorf("HeaderRowCount() = %v, want 0", table.HeaderRowCount())
	}

	// Check regular row uses normal font
	if len(table.rows) > 0 && len(table.rows[0].Cells) > 0 {
		if table.rows[0].Cells[0].Font != Helvetica {
			t.Errorf("Row font = %v, want Helvetica", table.rows[0].Cells[0].Font)
		}
	}
}

func TestTableLayout_AddRowCells(t *testing.T) {
	table := NewTableLayout(2)

	cell1 := TableCell{
		Content:  "Custom",
		Font:     TimesBold,
		FontSize: 14,
		Color:    Red,
		Align:    AlignCenter,
	}
	cell2 := TableCell{
		Content:  "Cell",
		Font:     TimesRoman,
		FontSize: 12,
		Color:    Blue,
		Align:    AlignRight,
	}

	result := table.AddRowCells(cell1, cell2)

	if result != table {
		t.Error("AddRowCells should return the table for chaining")
	}

	if table.RowCount() != 1 {
		t.Errorf("RowCount() = %v, want 1", table.RowCount())
	}

	// Verify cell properties preserved
	if len(table.rows) > 0 && len(table.rows[0].Cells) > 0 {
		if table.rows[0].Cells[0].Font != TimesBold {
			t.Errorf("Cell font = %v, want TimesBold", table.rows[0].Cells[0].Font)
		}
		if table.rows[0].Cells[0].Color != Red {
			t.Errorf("Cell color = %v, want Red", table.rows[0].Cells[0].Color)
		}
	}
}

func TestTableLayout_MethodChaining(t *testing.T) {
	table := NewTableLayout(3).
		SetBorder(0.5, Black).
		SetCellPadding(5).
		SetColumnWidths(100, 100, 100).
		AddHeaderRow("A", "B", "C").
		AddRow("1", "2", "3").
		AddRow("4", "5", "6")

	if table.ColumnCount() != 3 {
		t.Errorf("ColumnCount() = %v, want 3", table.ColumnCount())
	}

	if table.RowCount() != 3 {
		t.Errorf("RowCount() = %v, want 3", table.RowCount())
	}

	if table.HeaderRowCount() != 1 {
		t.Errorf("HeaderRowCount() = %v, want 1", table.HeaderRowCount())
	}
}

func TestNewTableCell(t *testing.T) {
	cell := NewTableCell("Test")

	if cell.Content != "Test" {
		t.Errorf("Content = %q, want %q", cell.Content, "Test")
	}

	if cell.Font != Helvetica {
		t.Errorf("Font = %v, want Helvetica", cell.Font)
	}

	if cell.FontSize != 10 {
		t.Errorf("FontSize = %v, want 10", cell.FontSize)
	}

	if cell.Color != Black {
		t.Errorf("Color = %v, want Black", cell.Color)
	}

	if cell.Align != AlignLeft {
		t.Errorf("Align = %v, want AlignLeft", cell.Align)
	}

	if cell.ColSpan != 1 {
		t.Errorf("ColSpan = %v, want 1", cell.ColSpan)
	}
}

func TestTableLayout_Height_EmptyTable(t *testing.T) {
	table := NewTableLayout(3)
	ctx := &LayoutContext{
		PageWidth: 595,
		Margins:   Margins{Left: 72, Right: 72},
	}

	height := table.Height(ctx)

	if height != 0 {
		t.Errorf("Height() = %v, want 0 for empty table", height)
	}
}

func TestTableLayout_Height_WithRows(t *testing.T) {
	table := NewTableLayout(3).
		AddHeaderRow("A", "B", "C").
		AddRow("1", "2", "3")

	ctx := &LayoutContext{
		PageWidth: 595,
		Margins:   Margins{Left: 72, Right: 72},
	}

	height := table.Height(ctx)

	// 2 rows, default fontSize 10 + padding 4*2 = 18 per row
	// Total = 36
	expectedHeight := 36.0

	if height != expectedHeight {
		t.Errorf("Height() = %v, want %v", height, expectedHeight)
	}
}

func TestTableLayout_Height_WithBorder(t *testing.T) {
	table := NewTableLayout(3).
		SetBorder(2.0, Black).
		AddRow("1", "2", "3")

	ctx := &LayoutContext{
		PageWidth: 595,
		Margins:   Margins{Left: 72, Right: 72},
	}

	height := table.Height(ctx)

	// 1 row = 18, plus border width = 2.0
	expectedHeight := 20.0

	if height != expectedHeight {
		t.Errorf("Height() = %v, want %v", height, expectedHeight)
	}
}

func TestTableLayout_Draw_EmptyTable(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	table := NewTableLayout(3)
	ctx := page.GetLayoutContext()

	err = table.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	// No operations should be added
	if len(page.TextOperations()) != 0 {
		t.Errorf("Expected no text operations, got %d", len(page.TextOperations()))
	}
}

func TestTableLayout_Draw_SimpleTable(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	table := NewTableLayout(2).
		AddHeaderRow("Name", "Value").
		AddRow("A", "1")

	ctx := page.GetLayoutContext()
	err = table.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	// Should have 4 text operations (2 header + 2 data cells)
	ops := page.TextOperations()
	if len(ops) != 4 {
		t.Errorf("Expected 4 text operations, got %d", len(ops))
	}
}

func TestTableLayout_Draw_WithBorders(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	table := NewTableLayout(2).
		SetBorder(1.0, Black).
		AddRow("A", "B")

	ctx := page.GetLayoutContext()
	err = table.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	// Should have graphics operations for borders
	gops := page.GraphicsOperations()
	if len(gops) == 0 {
		t.Error("Expected graphics operations for borders")
	}
}

func TestTableLayout_Draw_CursorAdvances(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	table := NewTableLayout(2).
		AddRow("A", "B").
		AddRow("C", "D")

	ctx := page.GetLayoutContext()
	initialY := ctx.CursorY

	err = table.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	// Cursor should have advanced by the table height
	if ctx.CursorY <= initialY {
		t.Error("Cursor Y should have advanced after drawing table")
	}
}

func TestTableLayout_CalculateColumnWidths_Auto(t *testing.T) {
	table := NewTableLayout(4)
	widths := table.calculateColumnWidths(400)

	if len(widths) != 4 {
		t.Errorf("Expected 4 widths, got %d", len(widths))
	}

	// All should be equal (100 each)
	for i, w := range widths {
		if w != 100 {
			t.Errorf("Width[%d] = %v, want 100", i, w)
		}
	}
}

func TestTableLayout_CalculateColumnWidths_Explicit(t *testing.T) {
	table := NewTableLayout(3).SetColumnWidths(100, 150, 200)
	widths := table.calculateColumnWidths(500)

	if len(widths) != 3 {
		t.Errorf("Expected 3 widths, got %d", len(widths))
	}

	if widths[0] != 100 {
		t.Errorf("Width[0] = %v, want 100", widths[0])
	}
	if widths[1] != 150 {
		t.Errorf("Width[1] = %v, want 150", widths[1])
	}
	if widths[2] != 200 {
		t.Errorf("Width[2] = %v, want 200", widths[2])
	}
}

func TestTableLayout_CalculateColumnWidths_Partial(t *testing.T) {
	// Only first 2 columns have explicit widths
	table := NewTableLayout(4).SetColumnWidths(100, 100)
	widths := table.calculateColumnWidths(400)

	if widths[0] != 100 {
		t.Errorf("Width[0] = %v, want 100", widths[0])
	}
	if widths[1] != 100 {
		t.Errorf("Width[1] = %v, want 100", widths[1])
	}

	// Remaining 200 spread over 2 columns = 100 each
	if widths[2] != 100 {
		t.Errorf("Width[2] = %v, want 100", widths[2])
	}
	if widths[3] != 100 {
		t.Errorf("Width[3] = %v, want 100", widths[3])
	}
}

func TestTableLayout_ImplementsDrawable(_ *testing.T) {
	var _ Drawable = (*TableLayout)(nil)
}

func TestTableLayout_CellAlignment(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	leftCell := TableCell{Content: "Left", Align: AlignLeft, Font: Helvetica, FontSize: 10}
	centerCell := TableCell{Content: "Center", Align: AlignCenter, Font: Helvetica, FontSize: 10}
	rightCell := TableCell{Content: "Right", Align: AlignRight, Font: Helvetica, FontSize: 10}

	table := NewTableLayout(3).
		SetColumnWidths(150, 150, 150).
		AddRowCells(leftCell, centerCell, rightCell)

	ctx := page.GetLayoutContext()
	err = table.Draw(ctx, page)

	if err != nil {
		t.Errorf("Draw() returned error: %v", err)
	}

	ops := page.TextOperations()
	if len(ops) != 3 {
		t.Fatalf("Expected 3 text operations, got %d", len(ops))
	}

	// Left cell should be leftmost
	// Center cell should be more centered
	// Right cell should be rightmost within its column

	// Basic sanity check: X positions should be increasing
	if ops[0].X >= ops[1].X || ops[1].X >= ops[2].X {
		t.Error("Text X positions should increase for different columns")
	}
}
