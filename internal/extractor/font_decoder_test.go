package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFontDecoder_WithCMap(t *testing.T) {
	t.Run("Decode with CMap", func(t *testing.T) {
		// Create CMap with Cyrillic mappings
		cmap := NewCMapTable("TestCMap")
		cmap.AddMapping(0x01, 'В') // U+0412
		cmap.AddMapping(0x02, 'ы') // U+044B
		cmap.AddMapping(0x03, 'п') // U+043F

		decoder := NewFontDecoderWithCMap(cmap)

		// Single-byte glyphs
		result := decoder.DecodeString([]byte{0x01, 0x02, 0x03})
		assert.Equal(t, "Вып", result)
	})

	t.Run("Decode 2-byte glyphs with CMap", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		// Use glyphs that won't trigger UTF-16 heuristic
		cmap.AddMapping(0x1001, 'В')
		cmap.AddMapping(0x1002, 'ы')
		cmap.AddMapping(0x1003, 'п')

		// Force 2-byte mode
		decoder := &FontDecoder{
			cmap:           cmap,
			encoding:       "",
			use2ByteGlyphs: true,
		}

		// Big-endian 2-byte glyphs (0x1001, 0x1002, 0x1003)
		result := decoder.DecodeString([]byte{0x10, 0x01, 0x10, 0x02, 0x10, 0x03})
		assert.Equal(t, "Вып", result)
	})

	t.Run("Decode with unmapped glyphs", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		cmap.AddMapping(0x01, 'A')

		decoder := NewFontDecoderWithCMap(cmap)

		// 0x01 is mapped, 0xFF is not
		result := decoder.DecodeString([]byte{0x01, 0xFF})
		// Should have 'A' and replacement character
		assert.Contains(t, result, "A")
		assert.Len(t, []rune(result), 2)
	})
}

func TestFontDecoder_WithoutCMap(t *testing.T) {
	t.Run("Decode without CMap (Latin-1 fallback)", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", false)

		// ASCII text
		result := decoder.DecodeString([]byte("Hello"))
		assert.Equal(t, "Hello", result)

		// Latin-1 extended characters
		result = decoder.DecodeString([]byte{0xE9}) // é
		assert.Equal(t, "é", result)
	})

	t.Run("Decode with WinAnsiEncoding", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "WinAnsiEncoding", false)

		// ASCII
		result := decoder.DecodeString([]byte("Test"))
		assert.Equal(t, "Test", result)

		// Windows-1252 specific characters (0x80-0x9F)
		result = decoder.DecodeString([]byte{0x80}) // Euro sign
		assert.Equal(t, "€", result)

		result = decoder.DecodeString([]byte{0x99}) // Trademark
		assert.Equal(t, "™", result)
	})
}

func TestFontDecoder_UTF16Detection(t *testing.T) {
	t.Run("Detect UTF-16BE with BOM", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", false)

		// UTF-16BE with BOM: "Привет" (Hello in Russian)
		utf16Data := []byte{
			0xFE, 0xFF, // BOM
			0x04, 0x1F, // П
			0x04, 0x40, // р
			0x04, 0x38, // и
			0x04, 0x32, // в
			0x04, 0x35, // е
			0x04, 0x42, // т
		}

		result := decoder.DecodeString(utf16Data)
		assert.Equal(t, "Привет", result)
	})

	t.Run("Detect UTF-16BE without BOM (heuristic)", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", false)

		// UTF-16BE without BOM (many null bytes suggest UTF-16)
		utf16Data := []byte{
			0x00, 0x41, // A
			0x00, 0x42, // B
			0x00, 0x43, // C
		}

		result := decoder.DecodeString(utf16Data)
		// Should detect as UTF-16 due to null byte pattern
		assert.Contains(t, result, "A")
	})

	t.Run("Do not detect UTF-16 for normal single-byte", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", false)

		// Regular single-byte data (no null bytes)
		result := decoder.DecodeString([]byte{0x41, 0x42, 0x43})
		assert.Equal(t, "ABC", result)
	})
}

func TestFontDecoder_WinAnsiTable(t *testing.T) {
	tests := []struct {
		name     string
		input    byte
		expected rune
	}{
		{"ASCII space", 0x20, ' '},
		{"ASCII A", 0x41, 'A'},
		{"Euro sign", 0x80, '€'},
		{"Single low-9 quote", 0x82, '‚'},
		{"Ellipsis", 0x85, '…'},
		{"Dagger", 0x86, '†'},
		{"Per mille", 0x89, '‰'},
		{"Trade mark", 0x99, '™'},
		{"Latin-1 é", 0xE9, 'é'},
		{"Latin-1 ñ", 0xF1, 'ñ'},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := decodeWinAnsi(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFontDecoder_EmptyInput(t *testing.T) {
	decoder := NewFontDecoderWithCMap(nil)
	result := decoder.DecodeString([]byte{})
	assert.Equal(t, "", result)
}

func TestFontDecoder_MixedContent(t *testing.T) {
	t.Run("CMap with ASCII fallback", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		cmap.AddMapping(0x01, 'Ф') // Cyrillic Ф

		decoder := NewFontDecoderWithCMap(cmap)

		// Mixed: mapped glyph (0x01) and unmapped ASCII-range bytes
		result := decoder.DecodeString([]byte{0x01, 0x20, 0x41}) // Ф, space, A
		assert.Contains(t, result, "Ф")
		// Unmapped bytes should still produce some output (fallback)
		assert.Len(t, []rune(result), 3)
	})
}

func TestFontDecoder_AutoDetectGlyphSize(t *testing.T) {
	t.Run("Auto-detect 1-byte glyphs", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		cmap.AddMapping(0x01, 'A')
		cmap.AddMapping(0xFF, 'Z')

		decoder := NewFontDecoderWithCMap(cmap)
		assert.False(t, decoder.use2ByteGlyphs, "Should use 1-byte glyphs")
	})

	t.Run("Auto-detect 2-byte glyphs", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		cmap.AddMapping(0x0100, 'A') // Requires 2 bytes

		decoder := NewFontDecoderWithCMap(cmap)
		assert.True(t, decoder.use2ByteGlyphs, "Should use 2-byte glyphs")
	})
}

func TestFontDecoder_HasCMap(t *testing.T) {
	t.Run("With CMap", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		decoder := NewFontDecoderWithCMap(cmap)
		assert.True(t, decoder.HasCMap())
	})

	t.Run("Without CMap", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", false)
		assert.False(t, decoder.HasCMap())
	})
}

func TestFontDecoder_Encoding(t *testing.T) {
	decoder := NewFontDecoder(nil, "WinAnsiEncoding", false)
	assert.Equal(t, "WinAnsiEncoding", decoder.Encoding())
}

func TestFontDecoder_String(t *testing.T) {
	t.Run("CMap with encoding", func(t *testing.T) {
		cmap := NewCMapTable("TestCMap")
		decoder := NewFontDecoder(cmap, "WinAnsiEncoding", false)
		str := decoder.String()
		assert.Contains(t, str, "CMap:TestCMap")
		assert.Contains(t, str, "Encoding:WinAnsiEncoding")
		assert.Contains(t, str, "1-byte-glyphs")
	})

	t.Run("2-byte glyphs", func(t *testing.T) {
		decoder := NewFontDecoder(nil, "", true)
		str := decoder.String()
		assert.Contains(t, str, "2-byte-glyphs")
	})
}

func TestFontDecoder_RealWorldCyrillic(t *testing.T) {
	t.Run("Russian text from real PDF", func(t *testing.T) {
		// Simulate a real PDF CMap for "Выписка" (Statement)
		cmap := NewCMapTable("Adobe-Identity-UCS")
		cmap.AddMapping(0x01, 'В') // U+0412
		cmap.AddMapping(0x02, 'ы') // U+044B
		cmap.AddMapping(0x03, 'п') // U+043F
		cmap.AddMapping(0x04, 'и') // U+0438
		cmap.AddMapping(0x05, 'с') // U+0441
		cmap.AddMapping(0x06, 'к') // U+043A
		cmap.AddMapping(0x07, 'а') // U+0430

		decoder := NewFontDecoderWithCMap(cmap)

		// Glyph sequence for "Выписка"
		glyphs := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07}
		result := decoder.DecodeString(glyphs)
		assert.Equal(t, "Выписка", result)
	})
}

func BenchmarkFontDecoder_DecodeCyrillic(b *testing.B) {
	cmap := NewCMapTable("TestCMap")
	for i := uint16(0); i < 256; i++ {
		cmap.AddMapping(i, rune(0x0400+i))
	}

	decoder := NewFontDecoderWithCMap(cmap)
	testData := make([]byte, 100)
	for i := range testData {
		testData[i] = byte(i)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = decoder.DecodeString(testData)
	}
}

func BenchmarkFontDecoder_DecodeUTF16(b *testing.B) {
	decoder := NewFontDecoder(nil, "", false)

	// UTF-16BE data
	utf16Data := []byte{0xFE, 0xFF} // BOM
	for i := 0; i < 50; i++ {
		utf16Data = append(utf16Data, 0x04, byte(0x10+i)) // Cyrillic characters
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = decoder.DecodeString(utf16Data)
	}
}
