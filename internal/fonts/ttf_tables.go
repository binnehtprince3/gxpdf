package fonts

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// HeadTable represents the 'head' (font header) table.
//
// The head table contains global information about the font:
//   - Font version, creation date
//   - Units per em (scaling factor)
//   - Bounding box
//   - Index format flags
//
// Reference: TrueType specification, 'head' table.
type HeadTable struct {
	UnitsPerEm uint16 // Units per em square (typically 1000 or 2048)
	XMin       int16  // Minimum X coordinate
	YMin       int16  // Minimum Y coordinate
	XMax       int16  // Maximum X coordinate
	YMax       int16  // Maximum Y coordinate
}

// HheaTable represents the 'hhea' (horizontal header) table.
//
// The hhea table contains metrics for horizontal layout:
//   - Ascender, descender, line gap
//   - Number of horizontal metrics
//
// Reference: TrueType specification, 'hhea' table.
type HheaTable struct {
	Ascender            int16  // Typographic ascender
	Descender           int16  // Typographic descender
	LineGap             int16  // Typographic line gap
	NumOfLongHorMetrics uint16 // Number of hMetric entries in hmtx table
}

// HmtxTable represents the 'hmtx' (horizontal metrics) table.
//
// The hmtx table contains advance widths and left side bearings
// for all glyphs in the font.
//
// Reference: TrueType specification, 'hmtx' table.
type HmtxTable struct {
	Metrics []HMetric // Horizontal metrics for each glyph
}

// HMetric represents horizontal metrics for a single glyph.
type HMetric struct {
	AdvanceWidth    uint16 // Advance width in font units
	LeftSideBearing int16  // Left side bearing in font units
}

// CmapTable represents the 'cmap' (character to glyph mapping) table.
//
// The cmap table maps character codes to glyph indices.
// It may contain multiple subtables for different platforms/encodings.
//
// Reference: TrueType specification, 'cmap' table.
type CmapTable struct {
	Subtables []*CmapSubtable // Platform-specific subtables
}

// CmapSubtable represents a single cmap subtable.
type CmapSubtable struct {
	PlatformID uint16          // Platform ID (0=Unicode, 3=Windows)
	EncodingID uint16          // Encoding ID
	Format     uint16          // Format number (0, 4, 6, 12, etc.)
	Mapping    map[rune]uint16 // Character to glyph ID mapping
}

// parseRequiredTables parses all required font tables.
func (f *TTFFont) parseRequiredTables() error {
	// Parse head table (required).
	if err := f.parseHeadTable(); err != nil {
		return fmt.Errorf("parse head table: %w", err)
	}

	// Parse hhea table (required).
	if err := f.parseHheaTable(); err != nil {
		return fmt.Errorf("parse hhea table: %w", err)
	}

	// Parse hmtx table (required).
	if err := f.parseHmtxTable(); err != nil {
		return fmt.Errorf("parse hmtx table: %w", err)
	}

	// Parse cmap table (required).
	if err := f.parseCmapTable(); err != nil {
		return fmt.Errorf("parse cmap table: %w", err)
	}

	return nil
}

// parseHeadTable parses the 'head' table.
func (f *TTFFont) parseHeadTable() error {
	table, ok := f.Tables["head"]
	if !ok {
		return fmt.Errorf("head table not found")
	}

	r := bytes.NewReader(table.Data)

	// Skip version (4 bytes) and fontRevision (4 bytes).
	if err := skipBytes(r, 8); err != nil {
		return err
	}

	// Skip checksumAdjustment (4 bytes) and magicNumber (4 bytes).
	if err := skipBytes(r, 8); err != nil {
		return err
	}

	// Skip flags (2 bytes).
	if err := skipBytes(r, 2); err != nil {
		return err
	}

	// Read unitsPerEm.
	if err := binary.Read(r, binary.BigEndian, &f.UnitsPerEm); err != nil {
		return fmt.Errorf("read unitsPerEm: %w", err)
	}

	return nil
}

// parseHheaTable parses the 'hhea' table.
func (f *TTFFont) parseHheaTable() error {
	table, ok := f.Tables["hhea"]
	if !ok {
		return fmt.Errorf("hhea table not found")
	}

	r := bytes.NewReader(table.Data)

	// Skip version (4 bytes).
	if err := skipBytes(r, 4); err != nil {
		return err
	}

	// Skip ascender, descender, lineGap (6 bytes).
	if err := skipBytes(r, 6); err != nil {
		return err
	}

	// Skip other fields until numOfLongHorMetrics (28 bytes from start).
	if err := skipBytes(r, 18); err != nil {
		return err
	}

	// Read numOfLongHorMetrics (at offset 34).
	var numHMetrics uint16
	if err := binary.Read(r, binary.BigEndian, &numHMetrics); err != nil {
		return fmt.Errorf("read numOfLongHorMetrics: %w", err)
	}

	// Store for hmtx parsing.
	f.Tables["hhea"].Data = append([]byte{}, table.Data...)
	binary.BigEndian.PutUint16(f.Tables["hhea"].Data[34:], numHMetrics)

	return nil
}

// parseHmtxTable parses the 'hmtx' table.
func (f *TTFFont) parseHmtxTable() error {
	hmtxTable, ok := f.Tables["hmtx"]
	if !ok {
		return fmt.Errorf("hmtx table not found")
	}

	hheaTable, ok := f.Tables["hhea"]
	if !ok {
		return fmt.Errorf("hhea table required for hmtx parsing")
	}

	// Get numOfLongHorMetrics from hhea.
	numHMetrics := binary.BigEndian.Uint16(hheaTable.Data[34:])

	r := bytes.NewReader(hmtxTable.Data)

	// Read long horizontal metrics (4 bytes each: advanceWidth + lsb).
	for i := uint16(0); i < numHMetrics; i++ {
		var advanceWidth uint16
		if err := binary.Read(r, binary.BigEndian, &advanceWidth); err != nil {
			return fmt.Errorf("read advanceWidth: %w", err)
		}

		// Skip left side bearing (2 bytes).
		if err := skipBytes(r, 2); err != nil {
			return err
		}

		f.GlyphWidths[i] = advanceWidth
	}

	return nil
}

// parseCmapTable parses the 'cmap' table.
func (f *TTFFont) parseCmapTable() error {
	table, ok := f.Tables["cmap"]
	if !ok {
		return fmt.Errorf("cmap table not found")
	}

	// Read cmap header.
	numTables, err := f.readCmapHeader(table.Data)
	if err != nil {
		return fmt.Errorf("read cmap header: %w", err)
	}

	// Find best subtable offset.
	bestOffset, err := f.findBestCmapSubtable(table.Data, numTables)
	if err != nil {
		return fmt.Errorf("find best subtable: %w", err)
	}

	// Parse the selected subtable.
	return f.parseCmapSubtable(table.Data, bestOffset)
}

// readCmapHeader reads the cmap table header.
func (f *TTFFont) readCmapHeader(data []byte) (uint16, error) {
	r := bytes.NewReader(data)

	// Read version.
	var version uint16
	if err := binary.Read(r, binary.BigEndian, &version); err != nil {
		return 0, fmt.Errorf("read version: %w", err)
	}

	// Read numTables.
	var numTables uint16
	if err := binary.Read(r, binary.BigEndian, &numTables); err != nil {
		return 0, fmt.Errorf("read numTables: %w", err)
	}

	return numTables, nil
}

// findBestCmapSubtable finds the best cmap subtable offset.
func (f *TTFFont) findBestCmapSubtable(data []byte, numTables uint16) (uint32, error) {
	r := bytes.NewReader(data[4:]) // Skip version and numTables.

	for i := uint16(0); i < numTables; i++ {
		var platformID, encodingID uint16
		var offset uint32

		if err := binary.Read(r, binary.BigEndian, &platformID); err != nil {
			return 0, fmt.Errorf("read platformID: %w", err)
		}
		if err := binary.Read(r, binary.BigEndian, &encodingID); err != nil {
			return 0, fmt.Errorf("read encodingID: %w", err)
		}
		if err := binary.Read(r, binary.BigEndian, &offset); err != nil {
			return 0, fmt.Errorf("read offset: %w", err)
		}

		// Prefer Windows Unicode BMP (platformID=3, encodingID=1).
		if platformID == 3 && encodingID == 1 {
			return offset, nil
		}
	}

	return 0, fmt.Errorf("no suitable cmap subtable found")
}

// parseCmapSubtable parses a cmap subtable (format 4 or 12).
func (f *TTFFont) parseCmapSubtable(data []byte, offset uint32) error {
	r := bytes.NewReader(data[offset:])

	var format uint16
	if err := binary.Read(r, binary.BigEndian, &format); err != nil {
		return fmt.Errorf("read format: %w", err)
	}

	switch format {
	case 4:
		return f.parseCmapFormat4(data, offset)
	case 12:
		return f.parseCmapFormat12(data, offset)
	default:
		return fmt.Errorf("unsupported cmap format: %d", format)
	}
}

// parseCmapFormat4 parses cmap format 4 (segment mapping).
func (f *TTFFont) parseCmapFormat4(data []byte, offset uint32) error {
	// Read format 4 header.
	segCount, err := f.readFormat4Header(data, offset)
	if err != nil {
		return fmt.Errorf("read header: %w", err)
	}

	// Read segment arrays.
	endCode, startCode, idDelta, err := f.readFormat4Segments(data, offset, segCount)
	if err != nil {
		return fmt.Errorf("read segments: %w", err)
	}

	// Build character to glyph mapping.
	f.buildCharToGlyphMapping(segCount, startCode, endCode, idDelta)

	return nil
}

// readFormat4Header reads the cmap format 4 header.
func (f *TTFFont) readFormat4Header(data []byte, offset uint32) (uint16, error) {
	r := bytes.NewReader(data[offset:])

	// Skip format, length, language (8 bytes).
	if err := skipBytes(r, 8); err != nil {
		return 0, err
	}

	// Read segCountX2.
	var segCountX2 uint16
	if err := binary.Read(r, binary.BigEndian, &segCountX2); err != nil {
		return 0, fmt.Errorf("read segCountX2: %w", err)
	}

	return segCountX2 / 2, nil
}

// readFormat4Segments reads the segment arrays from format 4.
func (f *TTFFont) readFormat4Segments(
	data []byte,
	offset uint32,
	segCount uint16,
) ([]uint16, []uint16, []int16, error) {
	r := bytes.NewReader(data[offset+14:]) // Skip to endCode array.

	// Read endCode array.
	endCode := make([]uint16, segCount)
	for i := uint16(0); i < segCount; i++ {
		if err := binary.Read(r, binary.BigEndian, &endCode[i]); err != nil {
			return nil, nil, nil, fmt.Errorf("read endCode: %w", err)
		}
	}

	// Skip reservedPad (2 bytes).
	if err := skipBytes(r, 2); err != nil {
		return nil, nil, nil, err
	}

	// Read startCode array.
	startCode := make([]uint16, segCount)
	for i := uint16(0); i < segCount; i++ {
		if err := binary.Read(r, binary.BigEndian, &startCode[i]); err != nil {
			return nil, nil, nil, fmt.Errorf("read startCode: %w", err)
		}
	}

	// Read idDelta array.
	idDelta := make([]int16, segCount)
	for i := uint16(0); i < segCount; i++ {
		if err := binary.Read(r, binary.BigEndian, &idDelta[i]); err != nil {
			return nil, nil, nil, fmt.Errorf("read idDelta: %w", err)
		}
	}

	return endCode, startCode, idDelta, nil
}

// buildCharToGlyphMapping builds the character to glyph mapping.
func (f *TTFFont) buildCharToGlyphMapping(
	segCount uint16,
	startCode []uint16,
	endCode []uint16,
	idDelta []int16,
) {
	for i := uint16(0); i < segCount; i++ {
		for charCode := startCode[i]; charCode <= endCode[i]; charCode++ {
			if charCode == 0xFFFF {
				break
			}
			//nolint:gosec // Character code is uint16, fits in int32.
			glyphID := uint16((int32(charCode) + int32(idDelta[i])) & 0xFFFF)
			f.CharToGlyph[rune(charCode)] = glyphID
		}
	}
}

// parseCmapFormat12 parses cmap format 12 (segmented coverage).
func (f *TTFFont) parseCmapFormat12(_ []byte, _ uint32) error {
	// Format 12 is more complex, but less common for basic fonts.
	// For MVP, we'll focus on format 4 support.
	return fmt.Errorf("cmap format 12 not yet implemented")
}

// skipBytes skips n bytes in the reader.
func skipBytes(r *bytes.Reader, n int64) error {
	_, err := r.Seek(n, 1) // Seek relative to current position.
	if err != nil {
		return fmt.Errorf("skip %d bytes: %w", n, err)
	}
	return nil
}
