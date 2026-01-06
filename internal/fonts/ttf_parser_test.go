package fonts

import (
	"bytes"
	"encoding/binary"
	"testing"
)

// TestParseFontDirectory tests parsing of font directory.
func TestParseFontDirectory(t *testing.T) {
	// Create minimal font directory for testing.
	var buf bytes.Buffer

	// Write sfnt version (TrueType = 0x00010000).
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x00010000))

	// Write numTables = 2.
	_ = binary.Write(&buf, binary.BigEndian, uint16(2))

	// Write searchRange, entrySelector, rangeShift.
	_ = binary.Write(&buf, binary.BigEndian, uint16(32)) // searchRange.
	_ = binary.Write(&buf, binary.BigEndian, uint16(1))  // entrySelector.
	_ = binary.Write(&buf, binary.BigEndian, uint16(0))  // rangeShift.

	// Write table entry 1: "head".
	buf.WriteString("head")
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x12345678)) // checksum.
	_ = binary.Write(&buf, binary.BigEndian, uint32(100))        // offset.
	_ = binary.Write(&buf, binary.BigEndian, uint32(54))         // length.

	// Write table entry 2: "hhea".
	buf.WriteString("hhea")
	_ = binary.Write(&buf, binary.BigEndian, uint32(0x87654321)) // checksum.
	_ = binary.Write(&buf, binary.BigEndian, uint32(200))        // offset.
	_ = binary.Write(&buf, binary.BigEndian, uint32(36))         // length.

	// Parse font directory.
	font := &TTFFont{
		Tables:      make(map[string]*TTFTable),
		GlyphWidths: make(map[uint16]uint16),
		CharToGlyph: make(map[rune]uint16),
	}

	err := font.parseFontDirectory(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("parseFontDirectory failed: %v", err)
	}

	// Verify tables were parsed.
	if len(font.Tables) != 2 {
		t.Errorf("expected 2 tables, got %d", len(font.Tables))
	}

	// Verify "head" table.
	headTable, ok := font.Tables["head"]
	if !ok {
		t.Fatal("head table not found")
	}
	if headTable.Tag != "head" {
		t.Errorf("expected tag 'head', got %q", headTable.Tag)
	}
	if headTable.Offset != 100 {
		t.Errorf("expected offset 100, got %d", headTable.Offset)
	}
	if headTable.Length != 54 {
		t.Errorf("expected length 54, got %d", headTable.Length)
	}

	// Verify "hhea" table.
	hheaTable, ok := font.Tables["hhea"]
	if !ok {
		t.Fatal("hhea table not found")
	}
	if hheaTable.Tag != "hhea" {
		t.Errorf("expected tag 'hhea', got %q", hheaTable.Tag)
	}
}

// TestParseTableEntry tests parsing of a single table entry.
func TestParseTableEntry(t *testing.T) {
	var buf bytes.Buffer

	// Write table entry: "test".
	buf.WriteString("test")
	_ = binary.Write(&buf, binary.BigEndian, uint32(0xAABBCCDD)) // checksum.
	_ = binary.Write(&buf, binary.BigEndian, uint32(1000))       // offset.
	_ = binary.Write(&buf, binary.BigEndian, uint32(500))        // length.

	// Parse entry.
	font := &TTFFont{}
	entry, err := font.parseTableEntry(bytes.NewReader(buf.Bytes()))
	if err != nil {
		t.Fatalf("parseTableEntry failed: %v", err)
	}

	// Verify fields.
	if entry.Tag != "test" {
		t.Errorf("expected tag 'test', got %q", entry.Tag)
	}
	if entry.Checksum != 0xAABBCCDD {
		t.Errorf("expected checksum 0xAABBCCDD, got 0x%08X", entry.Checksum)
	}
	if entry.Offset != 1000 {
		t.Errorf("expected offset 1000, got %d", entry.Offset)
	}
	if entry.Length != 500 {
		t.Errorf("expected length 500, got %d", entry.Length)
	}
}

// TestLoadTable tests loading table data.
func TestLoadTable(t *testing.T) {
	// Create test data.
	data := []byte{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07,
		0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F,
	}

	// Create table entry.
	table := &TTFTable{
		Tag:    "test",
		Offset: 4,
		Length: 8,
	}

	// Load table data.
	font := &TTFFont{}
	err := font.loadTable(data, table)
	if err != nil {
		t.Fatalf("loadTable failed: %v", err)
	}

	// Verify loaded data.
	expected := []byte{0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B}
	if !bytes.Equal(table.Data, expected) {
		t.Errorf("expected %v, got %v", expected, table.Data)
	}
}

// TestLoadTableOutOfBounds tests error handling for invalid offsets.
func TestLoadTableOutOfBounds(t *testing.T) {
	data := []byte{0x00, 0x01, 0x02, 0x03}

	tests := []struct {
		name   string
		offset uint32
		length uint32
	}{
		{"offset too large", 100, 10},
		{"length too large", 0, 100},
		{"offset + length overflow", 2, 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			table := &TTFTable{
				Tag:    "test",
				Offset: tt.offset,
				Length: tt.length,
			}

			font := &TTFFont{}
			err := font.loadTable(data, table)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
