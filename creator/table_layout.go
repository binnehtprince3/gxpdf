package creator

import (
	"github.com/coregx/gxpdf/internal/fonts"
)

// TableCell represents a cell in a table.
type TableCell struct {
	// Content is the text content of the cell.
	Content string

	// Font is the font for the cell content.
	Font FontName

	// FontSize is the font size in points.
	FontSize float64

	// Color is the text color.
	Color Color

	// Align is the horizontal alignment within the cell.
	Align Alignment

	// ColSpan is the number of columns this cell spans (future use).
	ColSpan int
}

// NewTableCell creates a new table cell with text content and default styling.
func NewTableCell(content string) TableCell {
	return TableCell{
		Content:  content,
		Font:     Helvetica,
		FontSize: 10,
		Color:    Black,
		Align:    AlignLeft,
		ColSpan:  1,
	}
}

// TableRow represents a row in a table.
type TableRow struct {
	Cells []TableCell
}

// TableLayout represents a table that can be drawn on a page.
//
// Tables support automatic column width calculation, borders,
// and header rows with bold styling.
//
// Example:
//
//	table := NewTableLayout(3)
//	table.SetBorder(0.5, Black)
//	table.AddHeaderRow("Name", "Age", "City")
//	table.AddRow("Alice", "30", "New York")
//	table.AddRow("Bob", "25", "Los Angeles")
//	page.Draw(table)
type TableLayout struct {
	columns      int
	columnWidths []float64 // nil = auto
	rows         []TableRow
	borderWidth  float64
	borderColor  *Color
	headerRows   int
	cellPadding  float64 // padding inside cells
}

// NewTableLayout creates a new table with the specified number of columns.
func NewTableLayout(columns int) *TableLayout {
	if columns < 1 {
		columns = 1
	}
	return &TableLayout{
		columns:     columns,
		rows:        make([]TableRow, 0),
		borderWidth: 0,
		borderColor: nil,
		headerRows:  0,
		cellPadding: 4.0, // default padding
	}
}

// SetColumnWidths sets explicit widths for each column.
// If not all widths are provided, remaining columns use auto width.
// Returns the table for method chaining.
func (t *TableLayout) SetColumnWidths(widths ...float64) *TableLayout {
	t.columnWidths = widths
	return t
}

// SetBorder enables table borders with the specified width and color.
// Returns the table for method chaining.
func (t *TableLayout) SetBorder(width float64, color Color) *TableLayout {
	t.borderWidth = width
	t.borderColor = &color
	return t
}

// SetCellPadding sets the padding inside cells.
// Returns the table for method chaining.
func (t *TableLayout) SetCellPadding(padding float64) *TableLayout {
	t.cellPadding = padding
	return t
}

// AddHeaderRow adds a header row with the given cell texts.
// Header rows use bold font by default.
// Returns the table for method chaining.
func (t *TableLayout) AddHeaderRow(cells ...string) *TableLayout {
	row := TableRow{
		Cells: make([]TableCell, len(cells)),
	}
	for i, content := range cells {
		row.Cells[i] = TableCell{
			Content:  content,
			Font:     HelveticaBold,
			FontSize: 10,
			Color:    Black,
			Align:    AlignLeft,
			ColSpan:  1,
		}
	}
	t.rows = append(t.rows, row)
	t.headerRows++
	return t
}

// AddRow adds a row with the given cell texts using default styling.
// Returns the table for method chaining.
func (t *TableLayout) AddRow(cells ...string) *TableLayout {
	row := TableRow{
		Cells: make([]TableCell, len(cells)),
	}
	for i, content := range cells {
		row.Cells[i] = NewTableCell(content)
	}
	t.rows = append(t.rows, row)
	return t
}

// AddRowCells adds a row with fully-configured cells.
// Returns the table for method chaining.
func (t *TableLayout) AddRowCells(cells ...TableCell) *TableLayout {
	row := TableRow{
		Cells: cells,
	}
	t.rows = append(t.rows, row)
	return t
}

// ColumnCount returns the number of columns.
func (t *TableLayout) ColumnCount() int {
	return t.columns
}

// RowCount returns the number of rows (including header rows).
func (t *TableLayout) RowCount() int {
	return len(t.rows)
}

// HeaderRowCount returns the number of header rows.
func (t *TableLayout) HeaderRowCount() int {
	return t.headerRows
}

// Height calculates the total height of the table when rendered.
func (t *TableLayout) Height(_ *LayoutContext) float64 {
	if len(t.rows) == 0 {
		return 0
	}

	rowHeight := t.calculateRowHeight()
	totalHeight := float64(len(t.rows)) * rowHeight

	// Add border widths if borders are enabled.
	if t.borderWidth > 0 {
		totalHeight += t.borderWidth // bottom border
	}

	return totalHeight
}

// Draw renders the table on the page at the current cursor position.
func (t *TableLayout) Draw(ctx *LayoutContext, page *Page) error {
	if len(t.rows) == 0 {
		return nil
	}

	colWidths := t.calculateColumnWidths(ctx.AvailableWidth())
	rowHeight := t.calculateRowHeight()
	startX := ctx.ContentLeft()
	startY := ctx.CurrentPDFY()

	// Draw rows.
	for rowIdx, row := range t.rows {
		y := startY - float64(rowIdx)*rowHeight

		if err := t.drawRow(page, row, startX, y, colWidths, rowHeight); err != nil {
			return err
		}
	}

	// Draw borders if enabled.
	if t.borderWidth > 0 && t.borderColor != nil {
		if err := t.drawBorders(page, startX, startY, colWidths, rowHeight); err != nil {
			return err
		}
	}

	// Update cursor position.
	ctx.CursorY += t.Height(ctx)

	return nil
}

// calculateRowHeight returns the height of one row.
func (t *TableLayout) calculateRowHeight() float64 {
	// Find the maximum font size across all cells.
	maxSize := 10.0
	for _, row := range t.rows {
		for _, cell := range row.Cells {
			if cell.FontSize > maxSize {
				maxSize = cell.FontSize
			}
		}
	}
	return maxSize + t.cellPadding*2
}

// calculateColumnWidths calculates widths for each column.
func (t *TableLayout) calculateColumnWidths(availableWidth float64) []float64 {
	widths := make([]float64, t.columns)

	// Use explicit widths if provided.
	explicitTotal := 0.0
	autoCount := 0

	for i := 0; i < t.columns; i++ {
		if i < len(t.columnWidths) && t.columnWidths[i] > 0 {
			widths[i] = t.columnWidths[i]
			explicitTotal += t.columnWidths[i]
		} else {
			autoCount++
		}
	}

	// Distribute remaining width to auto columns.
	if autoCount > 0 {
		remainingWidth := availableWidth - explicitTotal
		if remainingWidth < 0 {
			remainingWidth = 0
		}
		autoWidth := remainingWidth / float64(autoCount)

		for i := 0; i < t.columns; i++ {
			if widths[i] == 0 {
				widths[i] = autoWidth
			}
		}
	}

	return widths
}

// drawRow draws a single row at the specified position.
func (t *TableLayout) drawRow(
	page *Page,
	row TableRow,
	startX, y float64,
	colWidths []float64,
	_ float64, // rowHeight reserved for future multi-line cell support
) error {
	x := startX

	for colIdx := 0; colIdx < t.columns && colIdx < len(row.Cells); colIdx++ {
		cell := row.Cells[colIdx]
		colWidth := colWidths[colIdx]

		// Calculate text position within cell.
		textX := t.calculateCellTextX(x, colWidth, cell)
		textY := y - t.cellPadding - cell.FontSize // baseline

		if err := page.AddTextColor(cell.Content, textX, textY, cell.Font, cell.FontSize, cell.Color); err != nil {
			return err
		}

		x += colWidth
	}

	return nil
}

// calculateCellTextX calculates the X position for text within a cell.
func (t *TableLayout) calculateCellTextX(cellX, cellWidth float64, cell TableCell) float64 {
	textWidth := fonts.MeasureString(string(cell.Font), cell.Content, cell.FontSize)
	contentWidth := cellWidth - t.cellPadding*2

	switch cell.Align {
	case AlignCenter:
		return cellX + t.cellPadding + (contentWidth-textWidth)/2
	case AlignRight:
		return cellX + cellWidth - t.cellPadding - textWidth
	default:
		return cellX + t.cellPadding
	}
}

// drawBorders draws the table borders.
func (t *TableLayout) drawBorders(
	page *Page,
	startX, startY float64,
	colWidths []float64,
	rowHeight float64,
) error {
	totalWidth := 0.0
	for _, w := range colWidths {
		totalWidth += w
	}
	totalHeight := float64(len(t.rows)) * rowHeight

	opts := &LineOptions{
		Color: *t.borderColor,
		Width: t.borderWidth,
	}

	// Draw horizontal lines.
	for i := 0; i <= len(t.rows); i++ {
		y := startY - float64(i)*rowHeight
		if err := page.DrawLine(startX, y, startX+totalWidth, y, opts); err != nil {
			return err
		}
	}

	// Draw vertical lines.
	x := startX
	for i := 0; i <= t.columns; i++ {
		if err := page.DrawLine(x, startY, x, startY-totalHeight, opts); err != nil {
			return err
		}
		if i < t.columns {
			x += colWidths[i]
		}
	}

	return nil
}
