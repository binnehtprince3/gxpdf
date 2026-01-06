package export

import (
	"bytes"
	"strings"
	"testing"

	"github.com/coregx/gxpdf/internal/models/table"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createTestTable(t *testing.T) *table.Table {
	t.Helper()

	tbl, err := table.NewTable(3, 3)
	require.NoError(t, err)

	// Row 0
	tbl.SetCell(0, 0, table.NewCell("Name", 0, 0))
	tbl.SetCell(0, 1, table.NewCell("Age", 0, 1))
	tbl.SetCell(0, 2, table.NewCell("City", 0, 2))

	// Row 1
	tbl.SetCell(1, 0, table.NewCell("Alice", 1, 0))
	tbl.SetCell(1, 1, table.NewCell("30", 1, 1))
	tbl.SetCell(1, 2, table.NewCell("NYC", 1, 2))

	// Row 2
	tbl.SetCell(2, 0, table.NewCell("Bob", 2, 0))
	tbl.SetCell(2, 1, table.NewCell("25", 2, 1))
	tbl.SetCell(2, 2, table.NewCell("LA", 2, 2))

	tbl.Method = "Lattice"
	tbl.PageNum = 0
	tbl.Bounds = table.NewRectangle(100, 200, 400, 300)

	return tbl
}

func TestNewCSVExporter(t *testing.T) {
	exporter := NewCSVExporter()
	assert.NotNil(t, exporter)
	assert.NotNil(t, exporter.options)
	assert.Equal(t, ",", exporter.options.Delimiter)
}

func TestCSVExporter_Export(t *testing.T) {
	tbl := createTestTable(t)
	exporter := NewCSVExporter()

	var buf bytes.Buffer
	err := exporter.Export(tbl, &buf)
	require.NoError(t, err)

	result := buf.String()
	lines := strings.Split(strings.TrimSpace(result), "\n")

	require.Len(t, lines, 3)
	assert.Equal(t, "Name,Age,City", lines[0])
	assert.Equal(t, "Alice,30,NYC", lines[1])
	assert.Equal(t, "Bob,25,LA", lines[2])
}

func TestCSVExporter_ExportToString(t *testing.T) {
	tbl := createTestTable(t)
	exporter := NewCSVExporter()

	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(result), "\n")
	require.Len(t, lines, 3)
	assert.Contains(t, lines[0], "Name")
	assert.Contains(t, lines[1], "Alice")
	assert.Contains(t, lines[2], "Bob")
}

func TestCSVExporter_WithDelimiter(t *testing.T) {
	tbl := createTestTable(t)
	exporter := NewCSVExporter().WithDelimiter(";")

	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	// Should use semicolon delimiter
	assert.Contains(t, result, "Name;Age;City")
	assert.Contains(t, result, "Alice;30;NYC")
}

func TestCSVExporter_EmptyTable(t *testing.T) {
	tbl, err := table.NewTable(2, 2)
	require.NoError(t, err)

	exporter := NewCSVExporter()
	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	lines := strings.Split(strings.TrimSpace(result), "\n")
	require.Len(t, lines, 2)
	assert.Equal(t, ",", lines[0])
	assert.Equal(t, ",", lines[1])
}

func TestCSVExporter_WithQuotes(t *testing.T) {
	tbl, err := table.NewTable(1, 2)
	require.NoError(t, err)

	// Cell with comma (should be quoted)
	tbl.SetCell(0, 0, table.NewCell("Last, First", 0, 0))
	tbl.SetCell(0, 1, table.NewCell("Age", 0, 1))

	exporter := NewCSVExporter()
	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	// encoding/csv automatically quotes fields with commas
	assert.Contains(t, result, "\"Last, First\"")
}

func TestCSVExporter_WithNewlines(t *testing.T) {
	tbl, err := table.NewTable(1, 2)
	require.NoError(t, err)

	// Cell with newline (should be quoted)
	tbl.SetCell(0, 0, table.NewCell("Line1\nLine2", 0, 0))
	tbl.SetCell(0, 1, table.NewCell("Value", 0, 1))

	exporter := NewCSVExporter()
	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	// Should contain the newline in quoted field
	assert.Contains(t, result, "\"Line1\nLine2\"")
}

func TestCSVExporter_NilTable(t *testing.T) {
	exporter := NewCSVExporter()

	var buf bytes.Buffer
	err := exporter.Export(nil, &buf)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")
}

func TestCSVExporter_ContentType(t *testing.T) {
	exporter := NewCSVExporter()
	assert.Equal(t, "text/csv", exporter.ContentType())
}

func TestCSVExporter_FileExtension(t *testing.T) {
	// CSV
	exporter := NewCSVExporter()
	assert.Equal(t, ".csv", exporter.FileExtension())

	// TSV (tab-separated)
	tsvExporter := NewCSVExporter().WithDelimiter("\t")
	assert.Equal(t, ".tsv", tsvExporter.FileExtension())
}
