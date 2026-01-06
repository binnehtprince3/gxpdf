package fonts

import (
	"bytes"
	"fmt"
	"sort"
)

// GenerateToUnicodeCMap generates a ToUnicode CMap for text extraction.
//
// A ToUnicode CMap allows PDF viewers to extract correct Unicode text
// from documents using embedded fonts.
//
// The CMap maps character codes (as used in the PDF content stream)
// to Unicode code points.
//
// Reference: PDF 1.7 specification, Section 9.10 (ToUnicode CMaps).
func GenerateToUnicodeCMap(subset *FontSubset) ([]byte, error) {
	var buf bytes.Buffer

	// Write CMap header.
	if err := writeCMapHeader(&buf); err != nil {
		return nil, fmt.Errorf("write header: %w", err)
	}

	// Write character mappings.
	if err := writeCharMappings(&buf, subset); err != nil {
		return nil, fmt.Errorf("write mappings: %w", err)
	}

	// Write CMap footer.
	if err := writeCMapFooter(&buf); err != nil {
		return nil, fmt.Errorf("write footer: %w", err)
	}

	return buf.Bytes(), nil
}

// writeCMapHeader writes the CMap header.
func writeCMapHeader(buf *bytes.Buffer) error {
	header := `/CIDInit /ProcSet findresource begin
12 dict begin
begincmap
/CIDSystemInfo
<< /Registry (Adobe)
/Ordering (UCS)
/Supplement 0
>> def
/CMapName /Adobe-Identity-UCS def
/CMapType 2 def
1 begincodespacerange
<00> <FF>
endcodespacerange
`
	_, err := buf.WriteString(header)
	return err
}

// writeCharMappings writes character to Unicode mappings.
func writeCharMappings(buf *bytes.Buffer, subset *FontSubset) error {
	// Collect all used characters.
	chars := make([]rune, 0, len(subset.UsedChars))
	for ch := range subset.UsedChars {
		chars = append(chars, ch)
	}

	// Sort by character code.
	sort.Slice(chars, func(i, j int) bool {
		return chars[i] < chars[j]
	})

	// Write mappings in batches of 100 (PDF spec limit).
	const maxBatchSize = 100
	for i := 0; i < len(chars); i += maxBatchSize {
		end := i + maxBatchSize
		if end > len(chars) {
			end = len(chars)
		}

		if err := writeMappingBatch(buf, chars[i:end]); err != nil {
			return fmt.Errorf("write batch: %w", err)
		}
	}

	return nil
}

// writeMappingBatch writes a batch of character mappings.
func writeMappingBatch(buf *bytes.Buffer, chars []rune) error {
	// Write batch header.
	if _, err := fmt.Fprintf(buf, "%d beginbfchar\n", len(chars)); err != nil {
		return err
	}

	// Write each mapping.
	for _, ch := range chars {
		// Character code (1 byte for ASCII/Latin, 2+ bytes for Unicode).
		var charCode string
		if ch <= 0xFF {
			charCode = fmt.Sprintf("<%02X>", ch)
		} else {
			charCode = fmt.Sprintf("<%04X>", ch)
		}

		// Unicode code point (always 4 hex digits).
		unicode := fmt.Sprintf("<%04X>", ch)

		// Write mapping line.
		if _, err := fmt.Fprintf(buf, "%s %s\n", charCode, unicode); err != nil {
			return err
		}
	}

	// Write batch footer.
	if _, err := buf.WriteString("endbfchar\n"); err != nil {
		return err
	}

	return nil
}

// writeCMapFooter writes the CMap footer.
func writeCMapFooter(buf *bytes.Buffer) error {
	footer := `endcmap
CMapName currentdict /CMap defineresource pop
end
end
`
	_, err := buf.WriteString(footer)
	return err
}
