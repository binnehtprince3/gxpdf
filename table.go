package gxpdf

import (
	"io"

	"github.com/coregx/gxpdf/export"
	internaltable "github.com/coregx/gxpdf/internal/models/table"
)

// Table represents an extracted table from a PDF document.
//
// Table provides methods to access table data and export to various formats.
//
// Example:
//
//	tables := doc.ExtractTables()
//	for _, t := range tables {
//	    rows := t.Rows()
//	    for _, row := range rows {
//	        fmt.Println(row)
//	    }
//	}
type Table struct {
	internal *internaltable.Table
}

// Rows returns the table data as a 2D string slice.
//
// This is the simplest way to access table data.
func (t *Table) Rows() [][]string {
	return t.internal.ToStringGrid()
}

// RowCount returns the number of rows in the table.
func (t *Table) RowCount() int {
	return t.internal.RowCount
}

// ColumnCount returns the number of columns in the table.
func (t *Table) ColumnCount() int {
	return t.internal.ColCount
}

// PageNumber returns the page number where the table was found (0-based).
func (t *Table) PageNumber() int {
	return t.internal.PageNum
}

// Method returns the extraction method used ("Lattice", "Stream", or "Hybrid").
func (t *Table) Method() string {
	return t.internal.Method
}

// IsEmpty returns true if all cells in the table are empty.
func (t *Table) IsEmpty() bool {
	return t.internal.IsEmpty()
}

// Cell returns the text content of a cell at the given row and column.
//
// Returns empty string if the position is out of bounds.
func (t *Table) Cell(row, col int) string {
	cell := t.internal.GetCell(row, col)
	if cell == nil {
		return ""
	}
	return cell.Text
}

// String returns a string representation of the table.
func (t *Table) String() string {
	return t.internal.String()
}

// ExportCSV exports the table to CSV format.
func (t *Table) ExportCSV(w io.Writer) error {
	return export.NewCSVExporter().Export(t.internal, w)
}

// ExportJSON exports the table to JSON format.
func (t *Table) ExportJSON(w io.Writer) error {
	return export.NewJSONExporter().WithPrettyPrint(true).Export(t.internal, w)
}

// ExportExcel exports the table to Excel format.
func (t *Table) ExportExcel(w io.Writer) error {
	return export.NewExcelExporter().Export(t.internal, w)
}

// ToCSV returns the table as a CSV string.
func (t *Table) ToCSV() (string, error) {
	return export.NewCSVExporter().ExportToString(t.internal)
}

// ToJSON returns the table as a JSON string.
func (t *Table) ToJSON() (string, error) {
	return export.NewJSONExporter().WithPrettyPrint(true).ExportToString(t.internal)
}

// Internal returns the internal table representation.
//
// This is for advanced users who need access to cell bounds,
// alignment, and other detailed information.
func (t *Table) Internal() *internaltable.Table {
	return t.internal
}
