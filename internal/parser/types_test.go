package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Null Tests
// ============================================================================

func TestNull(t *testing.T) {
	n := NewNull()

	assert.Equal(t, "null", n.String())

	var buf bytes.Buffer
	written, err := n.WriteTo(&buf)
	require.NoError(t, err)
	assert.Equal(t, int64(4), written)
	assert.Equal(t, "null", buf.String())
}

// ============================================================================
// Boolean Tests
// ============================================================================

func TestBoolean(t *testing.T) {
	tests := []struct {
		name  string
		value bool
		want  string
	}{
		{"true value", true, "true"},
		{"false value", false, "false"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBoolean(tt.value)

			assert.Equal(t, tt.value, b.Value())
			assert.Equal(t, tt.want, b.String())

			var buf bytes.Buffer
			written, err := b.WriteTo(&buf)
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.want)), written)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

// ============================================================================
// Integer Tests
// ============================================================================

func TestInteger(t *testing.T) {
	tests := []struct {
		name  string
		value int64
		want  string
	}{
		{"zero", 0, "0"},
		{"positive", 42, "42"},
		{"negative", -42, "-42"},
		{"large positive", 2147483647, "2147483647"},
		{"large negative", -2147483648, "-2147483648"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewInteger(tt.value)

			assert.Equal(t, tt.value, i.Value())
			assert.Equal(t, int(tt.value), i.Int())
			assert.Equal(t, tt.want, i.String())

			var buf bytes.Buffer
			written, err := i.WriteTo(&buf)
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.want)), written)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

// ============================================================================
// Real Tests
// ============================================================================

func TestReal(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  string
	}{
		{"zero", 0.0, "0"},
		{"integer-like", 42.0, "42"},
		{"positive decimal", 3.14, "3.14"},
		{"negative decimal", -3.14, "-3.14"},
		{"trailing zeros removed", 3.14000, "3.14"},
		{"small number", 0.001, "0.001"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewReal(tt.value)

			assert.Equal(t, tt.value, r.Value())
			assert.Equal(t, tt.want, r.String())

			var buf bytes.Buffer
			written, err := r.WriteTo(&buf)
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.want)), written)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

// ============================================================================
// String Tests
// ============================================================================

func TestString_Literal(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"simple text", "Hello", "(Hello)"},
		{"with spaces", "Hello World", "(Hello World)"},
		{"empty", "", "()"},
		{"with parentheses", "Hello (World)", "(Hello \\(World\\))"},
		{"with backslash", "Hello\\World", "(Hello\\\\World)"},
		{"with newline", "Hello\nWorld", "(Hello\\nWorld)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewString(tt.value)

			assert.Equal(t, tt.value, s.Value())
			assert.False(t, s.IsHex())

			var buf bytes.Buffer
			written, err := s.WriteTo(&buf)
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.want)), written)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestString_Hexadecimal(t *testing.T) {
	tests := []struct {
		name  string
		value string // Input text
		want  string // Expected hex output
	}{
		{"simple text", "Hello", "<48656C6C6F>"},
		{"empty", "", "<>"},
		{"binary data", "\x00\x01\x02", "<000102>"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewHexString(tt.value)

			assert.Equal(t, tt.value, s.Value())
			assert.True(t, s.IsHex())

			var buf bytes.Buffer
			_, err := s.WriteTo(&buf)
			require.NoError(t, err)
			// Note: WriteTo converts to uppercase hex
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestString_Bytes(t *testing.T) {
	data := []byte{0x48, 0x65, 0x6C, 0x6C, 0x6F} // "Hello"
	s := NewStringBytes(data)

	assert.Equal(t, data, s.Bytes())
	assert.Equal(t, "Hello", s.Value())
}

// ============================================================================
// Name Tests
// ============================================================================

func TestName(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  string
	}{
		{"simple name", "Type", "/Type"},
		{"with leading slash", "/Type", "/Type"},
		{"CamelCase", "MediaBox", "/MediaBox"},
		{"with number", "Font1", "/Font1"},
		{"special chars", "A#B", "/A#23B"}, // # becomes #23
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := NewName(tt.value)

			// Value should not have leading /
			expectedValue := tt.value
			if expectedValue[0] == '/' {
				expectedValue = expectedValue[1:]
			}
			assert.Equal(t, expectedValue, n.Value())

			var buf bytes.Buffer
			written, err := n.WriteTo(&buf)
			require.NoError(t, err)
			assert.Equal(t, int64(len(tt.want)), written)
			assert.Equal(t, tt.want, buf.String())
		})
	}
}

func TestName_Equals(t *testing.T) {
	n1 := NewName("Type")
	n2 := NewName("Type")
	n3 := NewName("Font")

	assert.True(t, n1.Equals(n2))
	assert.False(t, n1.Equals(n3))
	assert.False(t, n1.Equals(nil))
}

// ============================================================================
// Type Tests
// ============================================================================

func TestTypeOf(t *testing.T) {
	tests := []struct {
		name string
		obj  PdfObject
		want Type
	}{
		{"null", NewNull(), TypeNull},
		{"boolean", NewBoolean(true), TypeBoolean},
		{"integer", NewInteger(42), TypeInteger},
		{"real", NewReal(3.14), TypeReal},
		{"string", NewString("hello"), TypeString},
		{"name", NewName("Type"), TypeName},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := TypeOf(tt.obj)
			assert.Equal(t, tt.want, got)
			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

// ============================================================================
// Clone Tests
// ============================================================================

func TestClone(t *testing.T) {
	tests := []struct {
		name string
		obj  PdfObject
	}{
		{"null", NewNull()},
		{"boolean", NewBoolean(true)},
		{"integer", NewInteger(42)},
		{"real", NewReal(3.14)},
		{"string", NewString("hello")},
		{"name", NewName("Type")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cloned := Clone(tt.obj)
			require.NotNil(t, cloned)

			// Values should be equal
			assert.Equal(t, tt.obj.String(), cloned.String())

			// For primitive types, they should be different instances
			// (except Null which is a singleton-like)
			if _, ok := tt.obj.(*Null); !ok {
				assert.NotSame(t, tt.obj, cloned)
			}
		})
	}
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkInteger_WriteTo(b *testing.B) {
	i := NewInteger(42)
	var buf bytes.Buffer

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_, _ = i.WriteTo(&buf)
	}
}

func BenchmarkReal_WriteTo(b *testing.B) {
	r := NewReal(3.14159)
	var buf bytes.Buffer

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_, _ = r.WriteTo(&buf)
	}
}

func BenchmarkString_WriteTo(b *testing.B) {
	s := NewString("Hello, World!")
	var buf bytes.Buffer

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		buf.Reset()
		_, _ = s.WriteTo(&buf)
	}
}

func BenchmarkName_WriteTo(b *testing.B) {
	n := NewName("MediaBox")
	var buf bytes.Buffer

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = n.WriteTo(&buf)
	}
}
