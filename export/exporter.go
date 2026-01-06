// Package export provides table export functionality for various formats.
//
// Supported formats:
//   - CSV (Comma-Separated Values)
//   - JSON (JavaScript Object Notation)
//   - Excel (.xlsx)
//
// Example:
//
//	tables := doc.ExtractTables()
//	export.ToCSV(tables[0], "output.csv")
//	export.ToJSON(tables[0], "output.json")
package export

import (
	"io"

	"github.com/coregx/gxpdf/internal/models/table"
)

// TableExporter is the interface for exporting tables to different formats.
//
// This interface enables:
//   - Multiple export formats (CSV, JSON, Excel, etc.)
//   - Custom exporter implementations
//   - Easy testing with mocks
//   - Dependency injection
//
// Example usage:
//
//	exporter := export.NewCSVExporter()
//	err := exporter.Export(table, writer)
type TableExporter interface {
	// Export writes the table to the writer in the format implemented by the exporter.
	//
	// Parameters:
	//   - tbl: The table to export
	//   - w: The writer to write to (file, buffer, network, etc.)
	//
	// Returns an error if export fails.
	Export(tbl *table.Table, w io.Writer) error

	// ExportToString exports the table to a string.
	//
	// This is a convenience method for formats that produce text output.
	// Returns the exported string, or error.
	ExportToString(tbl *table.Table) (string, error)

	// ContentType returns the MIME content type of the exported format.
	//
	// Examples:
	//   - CSV: "text/csv"
	//   - JSON: "application/json"
	//   - Excel: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	ContentType() string

	// FileExtension returns the recommended file extension for the format.
	//
	// Examples:
	//   - CSV: ".csv"
	//   - JSON: ".json"
	//   - Excel: ".xlsx"
	FileExtension() string
}

// ExportOptions contains options for table export.
type ExportOptions struct {
	// IncludeEmpty indicates whether to include empty cells in export.
	// Default: true
	IncludeEmpty bool

	// PreserveSpans indicates whether to preserve merged cells (row/col spans).
	// Not all formats support this (e.g., CSV doesn't).
	// Default: false
	PreserveSpans bool

	// Delimiter is the field delimiter for CSV export (e.g., ",", ";", "\t").
	// Default: ","
	Delimiter string

	// IncludeMetadata indicates whether to include table metadata in export.
	// Applicable to JSON and other metadata-aware formats.
	// Default: false
	IncludeMetadata bool

	// PrettyPrint indicates whether to format output for readability.
	// Applicable to JSON export.
	// Default: false
	PrettyPrint bool
}

// DefaultExportOptions returns default export options.
func DefaultExportOptions() *ExportOptions {
	return &ExportOptions{
		IncludeEmpty:    true,
		PreserveSpans:   false,
		Delimiter:       ",",
		IncludeMetadata: false,
		PrettyPrint:     false,
	}
}
