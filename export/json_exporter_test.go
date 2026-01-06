package export

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJSONExporter(t *testing.T) {
	exporter := NewJSONExporter()
	assert.NotNil(t, exporter)
	assert.NotNil(t, exporter.options)
}

func TestJSONExporter_Export(t *testing.T) {
	tbl := createTestTable(t)
	exporter := NewJSONExporter()

	var buf bytes.Buffer
	err := exporter.Export(tbl, &buf)
	require.NoError(t, err)

	// Parse JSON
	var result tableJSON
	err = json.Unmarshal(buf.Bytes(), &result)
	require.NoError(t, err)

	assert.Equal(t, 3, result.Rows)
	assert.Equal(t, 3, result.Columns)
	require.Len(t, result.Data, 3)
	require.Len(t, result.Data[0], 3)

	// Check first row
	assert.Equal(t, "Name", result.Data[0][0].Text)
	assert.Equal(t, "Age", result.Data[0][1].Text)
	assert.Equal(t, "City", result.Data[0][2].Text)

	// Check second row
	assert.Equal(t, "Alice", result.Data[1][0].Text)
	assert.Equal(t, "30", result.Data[1][1].Text)
	assert.Equal(t, "NYC", result.Data[1][2].Text)
}

func TestJSONExporter_ExportToString(t *testing.T) {
	tbl := createTestTable(t)
	exporter := NewJSONExporter()

	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	// Should be valid JSON
	var data tableJSON
	err = json.Unmarshal([]byte(result), &data)
	require.NoError(t, err)
	assert.Equal(t, 3, data.Rows)
}

func TestJSONExporter_WithPrettyPrint(t *testing.T) {
	tbl := createTestTable(t)

	// Without pretty print
	exporter1 := NewJSONExporter().WithPrettyPrint(false)
	result1, err := exporter1.ExportToString(tbl)
	require.NoError(t, err)

	// With pretty print
	exporter2 := NewJSONExporter().WithPrettyPrint(true)
	result2, err := exporter2.ExportToString(tbl)
	require.NoError(t, err)

	// Pretty print should be longer (has indentation)
	assert.Greater(t, len(result2), len(result1))
	assert.Contains(t, result2, "\n  ") // Has indentation
}

func TestJSONExporter_WithMetadata(t *testing.T) {
	tbl := createTestTable(t)

	// Without metadata
	exporter1 := NewJSONExporter().WithMetadata(false)
	result1, err := exporter1.ExportToString(tbl)
	require.NoError(t, err)

	var data1 tableJSON
	err = json.Unmarshal([]byte(result1), &data1)
	require.NoError(t, err)
	assert.Nil(t, data1.Metadata)

	// With metadata
	exporter2 := NewJSONExporter().WithMetadata(true)
	result2, err := exporter2.ExportToString(tbl)
	require.NoError(t, err)

	var data2 tableJSON
	err = json.Unmarshal([]byte(result2), &data2)
	require.NoError(t, err)
	require.NotNil(t, data2.Metadata)
	assert.Equal(t, "Lattice", data2.Metadata.Method)
	assert.Equal(t, 0, data2.Metadata.Page)
}

func TestJSONExporter_MergedCells(t *testing.T) {
	tbl := createTestTable(t)

	// Add merged cell
	mergedCell := tbl.GetCell(0, 0).WithRowSpan(2).WithColSpan(2)
	tbl.SetCell(0, 0, mergedCell)

	exporter := NewJSONExporter()
	result, err := exporter.ExportToString(tbl)
	require.NoError(t, err)

	var data tableJSON
	err = json.Unmarshal([]byte(result), &data)
	require.NoError(t, err)

	// Check merged cell has spans
	assert.Equal(t, 2, data.Data[0][0].RowSpan)
	assert.Equal(t, 2, data.Data[0][0].ColSpan)
}

func TestJSONExporter_Alignment(t *testing.T) {
	tbl := createTestTable(t)

	// Set cell alignment
	centerCell := tbl.GetCell(1, 1).WithAlignment(0) // AlignCenter
	tbl.SetCell(1, 1, centerCell)

	exporter := NewJSONExporter()
	var buf bytes.Buffer
	err := exporter.Export(tbl, &buf)
	require.NoError(t, err)

	var data tableJSON
	err = json.Unmarshal(buf.Bytes(), &data)
	require.NoError(t, err)

	// Note: AlignLeft (default) is omitted, only non-default alignments are included
	// Cell 1,1 should have alignment if it's not left
}

func TestJSONExporter_ContentType(t *testing.T) {
	exporter := NewJSONExporter()
	assert.Equal(t, "application/json", exporter.ContentType())
}

func TestJSONExporter_FileExtension(t *testing.T) {
	exporter := NewJSONExporter()
	assert.Equal(t, ".json", exporter.FileExtension())
}

func TestJSONExporter_NilTable(t *testing.T) {
	exporter := NewJSONExporter()

	var buf bytes.Buffer
	err := exporter.Export(nil, &buf)
	assert.Error(t, err)
}
