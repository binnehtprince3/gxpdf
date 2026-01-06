package fonts

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// TTFFont represents a parsed TrueType/OpenType font.
//
// TrueType fonts (.ttf) and OpenType fonts with TrueType outlines (.otf)
// share the same basic structure and can be parsed using the same logic.
//
// The font file contains:
//   - Font directory with table entries
//   - Multiple tables (head, hhea, hmtx, cmap, glyf, loca, etc.)
//   - Glyph outlines
//
// Reference: TrueType specification, Microsoft Typography.
type TTFFont struct {
	// FilePath is the path to the font file.
	FilePath string

	// PostScriptName is the font's PostScript name (from name table).
	PostScriptName string

	// Tables contains all parsed font tables.
	Tables map[string]*TTFTable

	// UnitsPerEm is the number of font units per em square (from head table).
	UnitsPerEm uint16

	// GlyphWidths maps glyph IDs to their advance widths.
	GlyphWidths map[uint16]uint16

	// CharToGlyph maps Unicode code points to glyph IDs.
	CharToGlyph map[rune]uint16

	// FontData is the raw font file data (for embedding).
	FontData []byte
}

// TTFTable represents a single table in the font file.
type TTFTable struct {
	Tag      string // 4-character tag (e.g., "head", "hhea")
	Checksum uint32 // Table checksum
	Offset   uint32 // Offset from beginning of file
	Length   uint32 // Length of table in bytes
	Data     []byte // Raw table data
}

// LoadTTF loads and parses a TrueType/OpenType font file.
//
// This function:
//  1. Reads the entire font file
//  2. Parses the font directory
//  3. Loads all required tables
//  4. Extracts glyph metrics
//  5. Builds character-to-glyph mapping
//
// Returns an error if the file is not a valid TTF/OTF font.
func LoadTTF(path string) (*TTFFont, error) {
	//nolint:gosec // Font file path is provided by user, not arbitrary.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read font file: %w", err)
	}

	font := &TTFFont{
		FilePath:    path,
		Tables:      make(map[string]*TTFTable),
		GlyphWidths: make(map[uint16]uint16),
		CharToGlyph: make(map[rune]uint16),
		FontData:    data,
	}

	if err := font.parse(data); err != nil {
		return nil, fmt.Errorf("parse TTF: %w", err)
	}

	return font, nil
}

// parse parses the font file structure.
func (f *TTFFont) parse(data []byte) error {
	r := bytes.NewReader(data)

	// Parse font directory.
	if err := f.parseFontDirectory(r); err != nil {
		return fmt.Errorf("parse font directory: %w", err)
	}

	// Load all table data.
	if err := f.loadTables(data); err != nil {
		return fmt.Errorf("load tables: %w", err)
	}

	// Parse required tables.
	if err := f.parseRequiredTables(); err != nil {
		return fmt.Errorf("parse required tables: %w", err)
	}

	return nil
}

// parseFontDirectory parses the font directory (table of contents).
//
// Font directory format:
//   - sfntVersion (4 bytes): 0x00010000 for TrueType, "OTTO" for CFF
//   - numTables (2 bytes): Number of tables
//   - searchRange (2 bytes): (maximum power of 2 <= numTables) * 16
//   - entrySelector (2 bytes): log2(maximum power of 2 <= numTables)
//   - rangeShift (2 bytes): numTables * 16 - searchRange
//
// Followed by table directory entries (16 bytes each).
func (f *TTFFont) parseFontDirectory(r io.Reader) error {
	// Read sfnt version (4 bytes).
	var version uint32
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return fmt.Errorf("read sfnt version: %w", err)
	}

	// Check version (0x00010000 for TrueType).
	if version != 0x00010000 {
		return fmt.Errorf("unsupported font format: 0x%08X", version)
	}

	// Read number of tables.
	var numTables uint16
	if err := binary.Read(r, binary.BigEndian, &numTables); err != nil {
		return fmt.Errorf("read num tables: %w", err)
	}

	// Skip searchRange, entrySelector, rangeShift (6 bytes total).
	if _, err := io.CopyN(io.Discard, r, 6); err != nil {
		return fmt.Errorf("skip font directory fields: %w", err)
	}

	// Read table directory entries.
	for i := uint16(0); i < numTables; i++ {
		entry, err := f.parseTableEntry(r)
		if err != nil {
			return fmt.Errorf("parse table entry %d: %w", i, err)
		}
		f.Tables[entry.Tag] = entry
	}

	return nil
}

// parseTableEntry parses a single table directory entry.
//
// Table directory entry format (16 bytes):
//   - tag (4 bytes): Table identifier (ASCII)
//   - checksum (4 bytes): Table checksum
//   - offset (4 bytes): Offset from beginning of file
//   - length (4 bytes): Length of table in bytes
func (f *TTFFont) parseTableEntry(r io.Reader) (*TTFTable, error) {
	var entry TTFTable

	// Read tag (4 ASCII characters).
	tagBytes := make([]byte, 4)
	if _, err := io.ReadFull(r, tagBytes); err != nil {
		return nil, fmt.Errorf("read tag: %w", err)
	}
	entry.Tag = string(tagBytes)

	// Read checksum, offset, length.
	if err := binary.Read(r, binary.BigEndian, &entry.Checksum); err != nil {
		return nil, fmt.Errorf("read checksum: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &entry.Offset); err != nil {
		return nil, fmt.Errorf("read offset: %w", err)
	}
	if err := binary.Read(r, binary.BigEndian, &entry.Length); err != nil {
		return nil, fmt.Errorf("read length: %w", err)
	}

	return &entry, nil
}

// loadTables loads the data for all tables.
func (f *TTFFont) loadTables(data []byte) error {
	for _, table := range f.Tables {
		if err := f.loadTable(data, table); err != nil {
			return fmt.Errorf("load table %s: %w", table.Tag, err)
		}
	}
	return nil
}

// loadTable loads data for a single table.
func (f *TTFFont) loadTable(data []byte, table *TTFTable) error {
	//nolint:gosec // len(data) from file size, typically < 2GB.
	if table.Offset+table.Length > uint32(len(data)) {
		return fmt.Errorf("table offset/length out of bounds")
	}
	table.Data = data[table.Offset : table.Offset+table.Length]
	return nil
}
