// Package export provides table export functionality.
package export

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coregx/gxpdf/internal/models/table"
)

// JSONExporter exports tables to JSON format.
//
// JSON export provides a rich, structured representation of tables.
//
// Features:
//   - Includes cell metadata (position, spans, alignment)
//   - Includes table metadata (page, bounds, method)
//   - Supports merged cells
//   - Pretty-printing option
//
// Output format:
//
//	{
//	  "rows": 3,
//	  "columns": 4,
//	  "data": [
//	    [{"text": "A1", "row": 0, "col": 0}, ...],
//	    ...
//	  ],
//	  "metadata": {
//	    "page": 0,
//	    "method": "Lattice",
//	    "bounds": {...}
//	  }
//	}
//
// Example usage:
//
//	exporter := export.NewJSONExporter().WithPrettyPrint(true)
//	err := exporter.Export(table, file)
type JSONExporter struct {
	options *ExportOptions
}

// NewJSONExporter creates a new JSON exporter with default options.
func NewJSONExporter() *JSONExporter {
	return &JSONExporter{
		options: DefaultExportOptions(),
	}
}

// NewJSONExporterWithOptions creates a new JSON exporter with custom options.
func NewJSONExporterWithOptions(options *ExportOptions) *JSONExporter {
	if options == nil {
		options = DefaultExportOptions()
	}
	return &JSONExporter{
		options: options,
	}
}

// WithPrettyPrint returns a new JSONExporter with pretty printing enabled/disabled.
func (e *JSONExporter) WithPrettyPrint(pretty bool) *JSONExporter {
	opts := *e.options
	opts.PrettyPrint = pretty
	return &JSONExporter{options: &opts}
}

// WithMetadata returns a new JSONExporter with metadata inclusion enabled/disabled.
func (e *JSONExporter) WithMetadata(include bool) *JSONExporter {
	opts := *e.options
	opts.IncludeMetadata = include
	return &JSONExporter{options: &opts}
}

// tableJSON is the JSON structure for table export.
type tableJSON struct {
	Rows     int           `json:"rows"`
	Columns  int           `json:"columns"`
	Data     [][]cellJSON  `json:"data"`
	Metadata *metadataJSON `json:"metadata,omitempty"`
}

// cellJSON is the JSON structure for a cell.
type cellJSON struct {
	Text      string `json:"text"`
	Row       int    `json:"row"`
	Column    int    `json:"column"`
	RowSpan   int    `json:"rowSpan,omitempty"`
	ColSpan   int    `json:"colSpan,omitempty"`
	Alignment string `json:"alignment,omitempty"`
}

// metadataJSON is the JSON structure for table metadata.
type metadataJSON struct {
	Page   int        `json:"page"`
	Method string     `json:"method"`
	Bounds boundsJSON `json:"bounds"`
}

// boundsJSON is the JSON structure for bounding rectangle.
type boundsJSON struct {
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
}

// Export writes the table to the writer in JSON format.
func (e *JSONExporter) Export(tbl *table.Table, w io.Writer) error {
	if tbl == nil {
		return fmt.Errorf("table is nil")
	}

	if err := tbl.Validate(); err != nil {
		return fmt.Errorf("invalid table: %w", err)
	}

	// Build JSON structure
	jsonData := e.buildJSON(tbl)

	// Encode to JSON
	encoder := json.NewEncoder(w)
	if e.options.PrettyPrint {
		encoder.SetIndent("", "  ")
	}

	if err := encoder.Encode(jsonData); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}

// buildJSON builds the JSON structure from the table.
func (e *JSONExporter) buildJSON(tbl *table.Table) *tableJSON {
	jsonData := &tableJSON{
		Rows:    tbl.RowCount,
		Columns: tbl.ColCount,
		Data:    make([][]cellJSON, tbl.RowCount),
	}

	// Convert cells
	for r := 0; r < tbl.RowCount; r++ {
		jsonData.Data[r] = make([]cellJSON, tbl.ColCount)
		for c := 0; c < tbl.ColCount; c++ {
			cell := tbl.GetCell(r, c)
			jsonCell := cellJSON{
				Text:   cell.Text,
				Row:    cell.Row,
				Column: cell.Column,
			}

			// Include spans if cell is merged
			if cell.IsMerged() {
				jsonCell.RowSpan = cell.RowSpan
				jsonCell.ColSpan = cell.ColSpan
			}

			// Include alignment if not default
			if cell.TextAlign != table.AlignLeft {
				jsonCell.Alignment = cell.TextAlign.String()
			}

			jsonData.Data[r][c] = jsonCell
		}
	}

	// Include metadata if requested
	if e.options.IncludeMetadata {
		jsonData.Metadata = &metadataJSON{
			Page:   tbl.PageNum,
			Method: tbl.Method,
			Bounds: boundsJSON{
				X:      tbl.Bounds.X,
				Y:      tbl.Bounds.Y,
				Width:  tbl.Bounds.Width,
				Height: tbl.Bounds.Height,
			},
		}
	}

	return jsonData
}

// ExportToString exports the table to a JSON string.
func (e *JSONExporter) ExportToString(tbl *table.Table) (string, error) {
	var buf bytes.Buffer
	if err := e.Export(tbl, &buf); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// ContentType returns the MIME content type for JSON.
func (e *JSONExporter) ContentType() string {
	return "application/json"
}

// FileExtension returns the file extension for JSON.
func (e *JSONExporter) FileExtension() string {
	return ".json"
}
