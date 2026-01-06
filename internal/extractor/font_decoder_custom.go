package extractor

// glyphNameToUnicode maps Adobe Glyph List (AGL) glyph names to Unicode code points.
//
// This table is used to decode custom font encodings that specify glyph names
// via the /Encoding/Differences array.
//
// Reference: Adobe Glyph List Specification v2.0
// https://github.com/adobe-type-tools/agl-aglfn
//
// This is a subset containing the most common glyphs:
// - ASCII digits and punctuation
// - Basic Latin letters
// - Common symbols
//
// Full AGL has ~4300 entries; we include ~200 most common ones for Phase 1.
var glyphNameToUnicode = map[string]rune{
	// Digits
	"zero":  '0', // U+0030
	"one":   '1', // U+0031
	"two":   '2', // U+0032
	"three": '3', // U+0033
	"four":  '4', // U+0034
	"five":  '5', // U+0035
	"six":   '6', // U+0036
	"seven": '7', // U+0037
	"eight": '8', // U+0038
	"nine":  '9', // U+0039

	// Basic punctuation
	"period":     '.',  // U+002E
	"comma":      ',',  // U+002C
	"colon":      ':',  // U+003A
	"semicolon":  ';',  // U+003B
	"slash":      '/',  // U+002F
	"backslash":  '\\', // U+005C
	"hyphen":     '-',  // U+002D
	"minus":      '-',  // U+002D (same as hyphen in most contexts)
	"plus":       '+',  // U+002B
	"equal":      '=',  // U+003D
	"underscore": '_',  // U+005F
	"space":      ' ',  // U+0020
	"exclam":     '!',  // U+0021
	"question":   '?',  // U+003F
	"numbersign": '#',  // U+0023
	"percent":    '%',  // U+0025
	"ampersand":  '&',  // U+0026
	"asterisk":   '*',  // U+002A
	"at":         '@',  // U+0040

	// Brackets and braces
	"parenleft":    '(', // U+0028
	"parenright":   ')', // U+0029
	"bracketleft":  '[', // U+005B
	"bracketright": ']', // U+005D
	"braceleft":    '{', // U+007B
	"braceright":   '}', // U+007D
	"less":         '<', // U+003C
	"greater":      '>', // U+003E

	// Quotes
	"quotesingle":    '\'',     // U+0027
	"quotedbl":       '"',      // U+0022
	"quotedblleft":   '\u201C', // U+201C "
	"quotedblright":  '\u201D', // U+201D "
	"quoteleft":      '\u2018', // U+2018 '
	"quoteright":     '\u2019', // U+2019 '
	"guillemotleft":  '\u00AB', // U+00AB «
	"guillemotright": '\u00BB', // U+00BB »
	"guilsinglleft":  '\u2039', // U+2039 ‹
	"guilsinglright": '\u203A', // U+203A ›

	// Lowercase Latin letters
	"a": 'a', "b": 'b', "c": 'c', "d": 'd', "e": 'e',
	"f": 'f', "g": 'g', "h": 'h', "i": 'i', "j": 'j',
	"k": 'k', "l": 'l', "m": 'm', "n": 'n', "o": 'o',
	"p": 'p', "q": 'q', "r": 'r', "s": 's', "t": 't',
	"u": 'u', "v": 'v', "w": 'w', "x": 'x', "y": 'y',
	"z": 'z',

	// Uppercase Latin letters
	"A": 'A', "B": 'B', "C": 'C', "D": 'D', "E": 'E',
	"F": 'F', "G": 'G', "H": 'H', "I": 'I', "J": 'J',
	"K": 'K', "L": 'L', "M": 'M', "N": 'N', "O": 'O',
	"P": 'P', "Q": 'Q', "R": 'R', "S": 'S', "T": 'T',
	"U": 'U', "V": 'V', "W": 'W', "X": 'X', "Y": 'Y',
	"Z": 'Z',

	// Common accented characters (Latin-1)
	"aacute":      'á', // U+00E1
	"eacute":      'é', // U+00E9
	"iacute":      'í', // U+00ED
	"oacute":      'ó', // U+00F3
	"uacute":      'ú', // U+00FA
	"agrave":      'à', // U+00E0
	"egrave":      'è', // U+00E8
	"igrave":      'ì', // U+00EC
	"ograve":      'ò', // U+00F2
	"ugrave":      'ù', // U+00F9
	"acircumflex": 'â', // U+00E2
	"ecircumflex": 'ê', // U+00EA
	"icircumflex": 'î', // U+00EE
	"ocircumflex": 'ô', // U+00F4
	"ucircumflex": 'û', // U+00FB
	"adieresis":   'ä', // U+00E4
	"edieresis":   'ë', // U+00EB
	"idieresis":   'ï', // U+00EF
	"odieresis":   'ö', // U+00F6
	"udieresis":   'ü', // U+00FC
	"ntilde":      'ñ', // U+00F1
	"ccedilla":    'ç', // U+00E7
	"aring":       'å', // U+00E5
	"ae":          'æ', // U+00E6
	"oslash":      'ø', // U+00F8

	// Currency and symbols
	"dollar":     '$', // U+0024
	"cent":       '¢', // U+00A2
	"sterling":   '£', // U+00A3
	"yen":        '¥', // U+00A5
	"Euro":       '€', // U+20AC
	"currency":   '¤', // U+00A4
	"degree":     '°', // U+00B0
	"mu":         'µ', // U+00B5
	"section":    '§', // U+00A7
	"paragraph":  '¶', // U+00B6
	"copyright":  '©', // U+00A9
	"registered": '®', // U+00AE
	"trademark":  '™', // U+2122
	"bullet":     '•', // U+2022
	"dagger":     '†', // U+2020
	"daggerdbl":  '‡', // U+2021
	"ellipsis":   '…', // U+2026

	// Math symbols
	"multiply":      '×', // U+00D7
	"divide":        '÷', // U+00F7
	"plusminus":     '±', // U+00B1
	"onehalf":       '½', // U+00BD
	"onequarter":    '¼', // U+00BC
	"threequarters": '¾', // U+00BE

	// Dashes
	"endash": '–', // U+2013
	"emdash": '—', // U+2014

	// Special spaces
	"nbspace": '\u00A0', // Non-breaking space
	"emspace": '\u2003', // Em space
	"enspace": '\u2002', // En space

	// Common ligatures
	"fi":  'ﬁ', // U+FB01
	"fl":  'ﬂ', // U+FB02
	"ff":  'ﬀ', // U+FB00
	"ffi": 'ﬃ', // U+FB03
	"ffl": 'ﬄ', // U+FB04

	// Arrows (subset)
	"arrowleft":  '←', // U+2190
	"arrowup":    '↑', // U+2191
	"arrowright": '→', // U+2192
	"arrowdown":  '↓', // U+2193

	// Card suits (for completeness)
	"club":    '♣', // U+2663
	"diamond": '♦', // U+2666
	"heart":   '♥', // U+2665
	"spade":   '♠', // U+2660
}

// buildCustomEncoding creates a glyph ID → Unicode mapping from glyph names.
//
// This converts the Differences array (which maps glyph IDs to glyph names)
// into a direct glyph ID → Unicode mapping using the Adobe Glyph List.
//
// Parameters:
//   - differences: Map of glyph ID → glyph name (from /Encoding/Differences)
//
// Returns:
//   - Map of glyph ID → Unicode rune
//
// Example:
//
//	differences := map[uint16]string{
//	    1: "zero", 2: "one", 3: "two", // Custom digit mapping
//	}
//	encoding := buildCustomEncoding(differences)
//	// encoding[1] = '0', encoding[2] = '1', encoding[3] = '2'
func buildCustomEncoding(differences map[uint16]string) map[uint16]rune {
	encoding := make(map[uint16]rune, len(differences))

	for glyphID, glyphName := range differences {
		// Look up glyph name in Adobe Glyph List
		if unicode, ok := glyphNameToUnicode[glyphName]; ok {
			encoding[glyphID] = unicode
		} else {
			// Unknown glyph name - try to use it as-is if it's a single character
			if len(glyphName) == 1 {
				encoding[glyphID] = rune(glyphName[0])
			}
			// Otherwise, skip this glyph (will use fallback encoding)
		}
	}

	return encoding
}

// NewFontDecoderWithCustomEncoding creates a FontDecoder with custom glyph mappings.
//
// This is used when a font has a /Encoding dictionary with /Differences array
// but no ToUnicode CMap.
//
// Parameters:
//   - differences: Map of glyph ID → glyph name (from /Encoding/Differences)
//   - baseEncoding: Base encoding name (e.g., "WinAnsiEncoding")
//   - use2ByteGlyphs: true for 2-byte glyphs, false for 1-byte glyphs
//
// Returns:
//   - FontDecoder configured with custom glyph mappings
func NewFontDecoderWithCustomEncoding(differences map[uint16]string, baseEncoding string, use2ByteGlyphs bool) *FontDecoder {
	customEncoding := buildCustomEncoding(differences)

	return &FontDecoder{
		cmap:           nil,
		encoding:       baseEncoding,
		use2ByteGlyphs: use2ByteGlyphs,
		customEncoding: customEncoding,
	}
}
