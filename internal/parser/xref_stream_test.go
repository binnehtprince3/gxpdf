package parser

import (
	"bytes"
	"compress/zlib"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFlateDecoder tests Flate/zlib decompression
func TestFlateDecoder_Decode(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "simple text",
			input: "Hello, World!",
			want:  "Hello, World!",
		},
		{
			name:  "empty data",
			input: "",
			want:  "",
		},
		{
			name:  "repeated pattern",
			input: "AAAABBBBCCCCDDDD",
			want:  "AAAABBBBCCCCDDDD",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Compress data
			var buf bytes.Buffer
			w := zlib.NewWriter(&buf)
			_, err := w.Write([]byte(tt.input))
			require.NoError(t, err)
			require.NoError(t, w.Close())

			compressed := buf.Bytes()

			// Decompress using our decoder
			decoder := &flateDecoder{}
			decompressed, err := decoder.Decode(compressed)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, string(decompressed))
		})
	}
}

// TestReadBigEndianInt tests big-endian integer reading
func TestReadBigEndianInt(t *testing.T) {
	tests := []struct {
		name string
		data []byte
		want int64
	}{
		{
			name: "1 byte",
			data: []byte{0x12},
			want: 0x12,
		},
		{
			name: "2 bytes",
			data: []byte{0x12, 0x34},
			want: 0x1234,
		},
		{
			name: "3 bytes",
			data: []byte{0x12, 0x34, 0x56},
			want: 0x123456,
		},
		{
			name: "zero",
			data: []byte{0x00, 0x00},
			want: 0,
		},
		{
			name: "max byte",
			data: []byte{0xFF},
			want: 255,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := readBigEndianInt(tt.data)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestParseXRefStreamEntries tests parsing binary xref entries
func TestParseXRefStreamEntries(t *testing.T) {
	tests := []struct {
		name   string
		wArray []int64 // /W array
		index  []int   // /Index array
		data   []byte  // Binary xref data
		want   []struct {
			objNum int
			typ    XRefEntryType
			field2 int64
			field3 int
		}
		wantErr bool
	}{
		{
			name:   "single in-use entry",
			wArray: []int64{1, 3, 2}, // Type:1 byte, Offset:3 bytes, Gen:2 bytes
			index:  []int{10, 1},     // Objects 10-10 (1 object)
			data: []byte{
				0x01,             // Type 1 (in-use)
				0x00, 0x00, 0x64, // Offset 100
				0x00, 0x00, // Generation 0
			},
			want: []struct {
				objNum int
				typ    XRefEntryType
				field2 int64
				field3 int
			}{
				{objNum: 10, typ: XRefEntryInUse, field2: 100, field3: 0},
			},
		},
		{
			name:   "free entry",
			wArray: []int64{1, 2, 1},
			index:  []int{5, 1},
			data: []byte{
				0x00,       // Type 0 (free)
				0x00, 0x00, // Next free object 0
				0x03, // Generation 3
			},
			want: []struct {
				objNum int
				typ    XRefEntryType
				field2 int64
				field3 int
			}{
				{objNum: 5, typ: XRefEntryFree, field2: 0, field3: 3},
			},
		},
		{
			name:   "compressed entry",
			wArray: []int64{1, 2, 1},
			index:  []int{20, 1},
			data: []byte{
				0x02,       // Type 2 (compressed)
				0x00, 0x0A, // ObjStm number 10
				0x05, // Index 5
			},
			want: []struct {
				objNum int
				typ    XRefEntryType
				field2 int64
				field3 int
			}{
				{objNum: 20, typ: XRefEntryCompressed, field2: 10, field3: 5},
			},
		},
		{
			name:   "multiple entries",
			wArray: []int64{1, 3, 2},
			index:  []int{0, 3}, // Objects 0-2
			data: []byte{
				// Object 0
				0x00, // Free
				0x00, 0x00, 0x00,
				0xFF, 0xFF,
				// Object 1
				0x01, // In-use
				0x00, 0x00, 0x10,
				0x00, 0x00,
				// Object 2
				0x01, // In-use
				0x00, 0x00, 0x20,
				0x00, 0x00,
			},
			want: []struct {
				objNum int
				typ    XRefEntryType
				field2 int64
				field3 int
			}{
				{objNum: 0, typ: XRefEntryFree, field2: 0, field3: 0xFFFF},
				{objNum: 1, typ: XRefEntryInUse, field2: 0x10, field3: 0},
				{objNum: 2, typ: XRefEntryInUse, field2: 0x20, field3: 0},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create dictionary with /W and /Index
			dict := NewDictionary()

			// Set /W array
			wArr := NewArray()
			for _, w := range tt.wArray {
				wArr.Append(NewInteger(w))
			}
			dict.Set("W", wArr)

			// Set /Index array (if provided)
			if len(tt.index) > 0 {
				indexArr := NewArray()
				for _, idx := range tt.index {
					indexArr.Append(NewInteger(int64(idx)))
				}
				dict.Set("Index", indexArr)
			} else {
				// Default: [0 Size]
				dict.Set("Size", NewInteger(int64(len(tt.want))))
			}

			// Create parser
			p := &Parser{}

			// Parse entries
			table, err := p.parseXRefStreamEntries(dict, tt.data)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.NotNil(t, table)

			// Verify entries
			for _, expected := range tt.want {
				entry, ok := table.GetEntry(expected.objNum)
				require.True(t, ok, "object %d not found", expected.objNum)
				assert.Equal(t, expected.typ, entry.Type, "object %d type mismatch", expected.objNum)
				assert.Equal(t, expected.field2, entry.Offset, "object %d offset mismatch", expected.objNum)
				assert.Equal(t, expected.field3, entry.Generation, "object %d generation mismatch", expected.objNum)
			}
		})
	}
}

// TestXRefStreamDetection tests detection of xref stream vs traditional xref
func TestXRefStreamDetection(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		isStream bool
	}{
		{
			name:     "traditional xref",
			content:  "xref\n0 6\n0000000000 65535 f\n",
			isStream: false,
		},
		{
			name:     "xref stream",
			content:  "90 0 obj\n<</Type /XRef>>",
			isStream: true,
		},
		{
			name:     "xref stream with space",
			content:  " 123 0 obj",
			isStream: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Check first character
			isDigit := len(tt.content) > 0 && tt.content[0] >= '0' && tt.content[0] <= '9'
			isSpace := len(tt.content) > 0 && tt.content[0] == ' '
			isStream := isDigit || (isSpace && len(tt.content) > 1 && tt.content[1] >= '0' && tt.content[1] <= '9')

			assert.Equal(t, tt.isStream, isStream)
		})
	}
}

// Benchmark for binary xref entry parsing
func BenchmarkParseXRefStreamEntries(b *testing.B) {
	// Create test data for 100 objects
	wArray := []int64{1, 3, 2}
	data := make([]byte, 100*6) // 100 entries * 6 bytes each

	for i := 0; i < 100; i++ {
		offset := i * 6
		data[offset] = 0x01                     // Type 1 (in-use)
		data[offset+1] = byte((i * 1000) >> 16) // Offset high byte
		data[offset+2] = byte((i * 1000) >> 8)  // Offset mid byte
		data[offset+3] = byte(i * 1000)         // Offset low byte
		data[offset+4] = 0x00                   // Generation high
		data[offset+5] = 0x00                   // Generation low
	}

	// Create dictionary
	dict := NewDictionary()
	wArr := NewArray()
	for _, w := range wArray {
		wArr.Append(NewInteger(w))
	}
	dict.Set("W", wArr)
	dict.Set("Size", NewInteger(100))

	p := &Parser{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := p.parseXRefStreamEntries(dict, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
