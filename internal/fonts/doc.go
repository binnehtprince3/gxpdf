// Package fonts provides the Standard 14 Type 1 fonts that are built into all PDF readers.
//
// The 14 Standard Fonts do not require embedding - they are guaranteed to be available
// in every compliant PDF viewer. This package provides font definitions, metrics, and
// PDF object generation for these fonts.
//
// Standard 14 Fonts:
//   - Serif: Times-Roman, Times-Bold, Times-Italic, Times-BoldItalic
//   - Sans-Serif: Helvetica, Helvetica-Bold, Helvetica-Oblique, Helvetica-BoldOblique
//   - Monospace: Courier, Courier-Bold, Courier-Oblique, Courier-BoldOblique
//   - Symbol: Symbol, ZapfDingbats
//
// Usage:
//
//	font := fonts.Helvetica
//	var buf bytes.Buffer
//	font.WriteFontObject(5, &buf)
//	// Result: "5 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica /Encoding /WinAnsiEncoding >>\nendobj\n"
package fonts
