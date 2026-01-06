// Package export provides table export functionality.
package export

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	"github.com/coregx/gxpdf/internal/models/table"
)

// CSVExporter exports tables to CSV format.
//
// CSV (Comma-Separated Values) is a simple text format for tabular data.
//
// Features:
//   - Configurable delimiter (comma, semicolon, tab, etc.)
//   - Proper quoting and escaping
//   - Handles multi-line cells
//   - Standard RFC 4180 compliant
//
// Limitations:
//   - Does not support merged cells (spans)
//   - No cell formatting (alignment, colors, etc.)
//
// Example usage:
//
//	exporter := export.NewCSVExporter()
//	err := exporter.Export(table, file)
type CSVExporter struct {
	options *ExportOptions
}

// NewCSVExporter creates a new CSV exporter with default options.
func NewCSVExporter() *CSVExporter {
	return &CSVExporter{
		options: DefaultExportOptions(),
	}
}

// NewCSVExporterWithOptions creates a new CSV exporter with custom options.
func NewCSVExporterWithOptions(options *ExportOptions) *CSVExporter {
	if options == nil {
		options = DefaultExportOptions()
	}
	return &CSVExporter{
		options: options,
	}
}

// WithDelimiter returns a new CSVExporter with a custom delimiter.
//
// Common delimiters:
//   - "," - Comma (default)
//   - ";" - Semicolon (European standard)
//   - "\t" - Tab (TSV format)
func (e *CSVExporter) WithDelimiter(delimiter string) *CSVExporter {
	opts := *e.options
	opts.Delimiter = delimiter
	return &CSVExporter{options: &opts}
}

// Export writes the table to the writer in CSV format.
func (e *CSVExporter) Export(tbl *table.Table, w io.Writer) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}

	if err := tbl.Validate(); err != nil {
		return fmt.Errorf("invalid table: %w", err)
	}

	// Create CSV writer
	csvWriter := csv.NewWriter(w)

	// Set delimiter (default is comma)
	if len(e.options.Delimiter) > 0 {
		csvWriter.Comma = rune(e.options.Delimiter[0])
	}

	// Write rows
	for r := 0; r < tbl.RowCount; r++ {
		row := make([]string, tbl.ColCount)
		for c := 0; c < tbl.ColCount; c++ {
			cell := tbl.GetCell(r, c)
			if cell != nil {
				row[c] = cell.Text
			} else {
				row[c] = ""
			}
		}

		if err := csvWriter.Write(row); err != nil {
			return fmt.Errorf("failed to write row %d: %w", r, err)
		}
	}

	// Flush and check for errors
	csvWriter.Flush()
	if err := csvWriter.Error(); err != nil {
		return fmt.Errorf("CSV writer error: %w", err)
	}

	return nil
}

// ExportToString exports the table to a CSV string.
func (e *CSVExporter) ExportToString(tbl *table.Table) (string, error) {
	var buf bytes.Buffer
	if err := e.Export(tbl, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ContentType returns the MIME content type for CSV.
func (e *CSVExporter) ContentType() string {
	return "text/csv"
}

// FileExtension returns the file extension for CSV.
func (e *CSVExporter) FileExtension() string {
	if e.options.Delimiter == "\t" {
		return ".tsv"
	}
	return ".csv"
}
