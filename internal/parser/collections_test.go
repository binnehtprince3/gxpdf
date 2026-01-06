package parser

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ============================================================================
// Array Tests
// ============================================================================

func TestArray_NewArray(t *testing.T) {
	arr := NewArray()
	assert.NotNil(t, arr)
	assert.Equal(t, 0, arr.Len())
}

func TestArray_Append(t *testing.T) {
	arr := NewArray()

	arr.Append(NewInteger(42))
	arr.Append(NewString("hello"))
	arr.Append(NewBoolean(true))

	assert.Equal(t, 3, arr.Len())
	assert.Equal(t, int64(42), arr.Get(0).(*Integer).Value())
	assert.Equal(t, "hello", arr.Get(1).(*String).Value())
	assert.True(t, arr.Get(2).(*Boolean).Value())
}

func TestArray_AppendAll(t *testing.T) {
	arr := NewArray()

	arr.AppendAll(
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	)

	assert.Equal(t, 3, arr.Len())
}

func TestArray_Get(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	// Valid indices
	assert.Equal(t, int64(1), arr.Get(0).(*Integer).Value())
	assert.Equal(t, int64(2), arr.Get(1).(*Integer).Value())
	assert.Equal(t, int64(3), arr.Get(2).(*Integer).Value())

	// Out of bounds
	assert.Nil(t, arr.Get(-1))
	assert.Nil(t, arr.Get(3))
}

func TestArray_Set(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	// Valid set
	err := arr.Set(1, NewInteger(42))
	require.NoError(t, err)
	assert.Equal(t, int64(42), arr.Get(1).(*Integer).Value())

	// Out of bounds
	err = arr.Set(3, NewInteger(99))
	assert.Error(t, err)

	err = arr.Set(-1, NewInteger(99))
	assert.Error(t, err)
}

func TestArray_Insert(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(3),
	})

	// Insert in middle
	err := arr.Insert(1, NewInteger(2))
	require.NoError(t, err)
	assert.Equal(t, 3, arr.Len())
	assert.Equal(t, int64(1), arr.Get(0).(*Integer).Value())
	assert.Equal(t, int64(2), arr.Get(1).(*Integer).Value())
	assert.Equal(t, int64(3), arr.Get(2).(*Integer).Value())

	// Insert at beginning
	err = arr.Insert(0, NewInteger(0))
	require.NoError(t, err)
	assert.Equal(t, 4, arr.Len())
	assert.Equal(t, int64(0), arr.Get(0).(*Integer).Value())

	// Insert at end
	err = arr.Insert(4, NewInteger(4))
	require.NoError(t, err)
	assert.Equal(t, 5, arr.Len())
	assert.Equal(t, int64(4), arr.Get(4).(*Integer).Value())

	// Out of bounds
	err = arr.Insert(10, NewInteger(99))
	assert.Error(t, err)
}

func TestArray_Remove(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	// Remove middle element
	err := arr.Remove(1)
	require.NoError(t, err)
	assert.Equal(t, 2, arr.Len())
	assert.Equal(t, int64(1), arr.Get(0).(*Integer).Value())
	assert.Equal(t, int64(3), arr.Get(1).(*Integer).Value())

	// Out of bounds
	err = arr.Remove(5)
	assert.Error(t, err)
}

func TestArray_Clear(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	arr.Clear()
	assert.Equal(t, 0, arr.Len())
}

func TestArray_Elements(t *testing.T) {
	original := []PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	}
	arr := NewArrayFromSlice(original)

	elements := arr.Elements()

	// Should be equal
	assert.Equal(t, len(original), len(elements))
	for i := range original {
		assert.Equal(t, original[i].String(), elements[i].String())
	}

	// Modifying returned slice should not affect array
	elements[0] = NewInteger(99)
	assert.Equal(t, int64(1), arr.Get(0).(*Integer).Value())
}

func TestArray_String(t *testing.T) {
	tests := []struct {
		name  string
		array *Array
		want  string
	}{
		{
			name:  "empty array",
			array: NewArray(),
			want:  "[]",
		},
		{
			name: "single element",
			array: NewArrayFromSlice([]PdfObject{
				NewInteger(42),
			}),
			want: "[42]",
		},
		{
			name: "multiple elements",
			array: NewArrayFromSlice([]PdfObject{
				NewInteger(1),
				NewString("hello"),
				NewBoolean(true),
			}),
			want: "[1 (hello) true]",
		},
		{
			name: "nested array",
			array: NewArrayFromSlice([]PdfObject{
				NewInteger(1),
				NewArrayFromSlice([]PdfObject{
					NewInteger(2),
					NewInteger(3),
				}),
			}),
			want: "[1 [2 3]]",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.array.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestArray_WriteTo(t *testing.T) {
	arr := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewString("hello"),
		NewBoolean(true),
	})

	var buf bytes.Buffer
	written, err := arr.WriteTo(&buf)
	require.NoError(t, err)
	assert.Greater(t, written, int64(0))
	assert.Equal(t, "[1 (hello) true]", buf.String())
}

func TestArray_Clone(t *testing.T) {
	original := NewArrayFromSlice([]PdfObject{
		NewInteger(1),
		NewInteger(2),
		NewInteger(3),
	})

	cloned := original.Clone()

	// Should have same values
	assert.Equal(t, original.Len(), cloned.Len())
	for i := 0; i < original.Len(); i++ {
		assert.Equal(t, original.Get(i).String(), cloned.Get(i).String())
	}

	// Modifying clone should not affect original
	err := cloned.Set(0, NewInteger(99))
	require.NoError(t, err)
	assert.Equal(t, int64(1), original.Get(0).(*Integer).Value())
	assert.Equal(t, int64(99), cloned.Get(0).(*Integer).Value())
}

// ============================================================================
// Dictionary Tests
// ============================================================================

func TestDictionary_NewDictionary(t *testing.T) {
	dict := NewDictionary()
	assert.NotNil(t, dict)
	assert.Equal(t, 0, dict.Len())
}

func TestDictionary_Set_Get(t *testing.T) {
	dict := NewDictionary()

	dict.Set("Type", NewName("Page"))
	dict.Set("Count", NewInteger(42))
	dict.Set("Title", NewString("Hello"))

	assert.Equal(t, 3, dict.Len())
	assert.Equal(t, "Page", dict.GetName("Type").Value())
	assert.Equal(t, int64(42), dict.GetInteger("Count"))
	assert.Equal(t, "Hello", dict.GetString("Title"))
}

func TestDictionary_ConvenienceMethods(t *testing.T) {
	dict := NewDictionary()

	dict.SetName("Type", "Page")
	dict.SetInteger("Count", 42)
	dict.SetReal("Width", 595.0)
	dict.SetBoolean("Visible", true)
	dict.SetString("Title", "Test")

	assert.Equal(t, "Page", dict.GetName("Type").Value())
	assert.Equal(t, int64(42), dict.GetInteger("Count"))
	assert.Equal(t, 595.0, dict.GetReal("Width"))
	assert.True(t, dict.GetBoolean("Visible"))
	assert.Equal(t, "Test", dict.GetString("Title"))
}

func TestDictionary_Has(t *testing.T) {
	dict := NewDictionary()
	dict.Set("Type", NewName("Page"))

	assert.True(t, dict.Has("Type"))
	assert.False(t, dict.Has("NotExists"))
}

func TestDictionary_Remove(t *testing.T) {
	dict := NewDictionary()
	dict.Set("Type", NewName("Page"))
	dict.Set("Count", NewInteger(42))

	assert.Equal(t, 2, dict.Len())

	dict.Remove("Type")

	assert.Equal(t, 1, dict.Len())
	assert.False(t, dict.Has("Type"))
	assert.True(t, dict.Has("Count"))
}

func TestDictionary_Keys(t *testing.T) {
	dict := NewDictionary()

	// Insert in specific order
	dict.Set("Type", NewName("Page"))
	dict.Set("MediaBox", NewArray())
	dict.Set("Count", NewInteger(1))

	keys := dict.Keys()

	// Should maintain insertion order
	assert.Equal(t, []string{"Type", "MediaBox", "Count"}, keys)

	// KeysSorted should return alphabetical order
	sorted := dict.KeysSorted()
	assert.Equal(t, []string{"Count", "MediaBox", "Type"}, sorted)
}

func TestDictionary_Clear(t *testing.T) {
	dict := NewDictionary()
	dict.Set("Type", NewName("Page"))
	dict.Set("Count", NewInteger(1))

	dict.Clear()

	assert.Equal(t, 0, dict.Len())
	assert.False(t, dict.Has("Type"))
}

func TestDictionary_NestedStructures(t *testing.T) {
	dict := NewDictionary()

	// Nested dictionary
	resources := NewDictionary()
	resources.Set("Font", NewName("F1"))
	dict.Set("Resources", resources)

	// Nested array
	mediaBox := NewArrayFromSlice([]PdfObject{
		NewInteger(0),
		NewInteger(0),
		NewInteger(595),
		NewInteger(842),
	})
	dict.Set("MediaBox", mediaBox)

	// Retrieve nested structures
	assert.NotNil(t, dict.GetDictionary("Resources"))
	assert.Equal(t, "F1", dict.GetDictionary("Resources").GetName("Font").Value())

	assert.NotNil(t, dict.GetArray("MediaBox"))
	assert.Equal(t, 4, dict.GetArray("MediaBox").Len())
}

func TestDictionary_String(t *testing.T) {
	tests := []struct {
		name string
		dict *Dictionary
		want string
	}{
		{
			name: "empty dictionary",
			dict: NewDictionary(),
			want: "<<>>",
		},
		{
			name: "single entry",
			dict: func() *Dictionary {
				d := NewDictionary()
				d.Set("Type", NewName("Page"))
				return d
			}(),
			want: "<</Type /Page>>",
		},
		{
			name: "multiple entries",
			dict: func() *Dictionary {
				d := NewDictionary()
				d.Set("Type", NewName("Page"))
				d.Set("Count", NewInteger(1))
				return d
			}(),
			want: "<</Type /Page /Count 1>>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.dict.String()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDictionary_WriteTo(t *testing.T) {
	dict := NewDictionary()
	dict.Set("Type", NewName("Page"))
	dict.Set("Count", NewInteger(1))

	var buf bytes.Buffer
	written, err := dict.WriteTo(&buf)
	require.NoError(t, err)
	assert.Greater(t, written, int64(0))
	assert.Equal(t, "<</Type /Page /Count 1>>", buf.String())
}

func TestDictionary_Clone(t *testing.T) {
	original := NewDictionary()
	original.Set("Type", NewName("Page"))
	original.Set("Count", NewInteger(1))

	cloned := original.Clone()

	// Should have same values
	assert.Equal(t, original.Len(), cloned.Len())
	assert.Equal(t, original.GetName("Type").Value(), cloned.GetName("Type").Value())

	// Modifying clone should not affect original
	cloned.Set("Count", NewInteger(99))
	assert.Equal(t, int64(1), original.GetInteger("Count"))
	assert.Equal(t, int64(99), cloned.GetInteger("Count"))
}

func TestDictionary_Merge(t *testing.T) {
	dict1 := NewDictionary()
	dict1.Set("Type", NewName("Page"))
	dict1.Set("Count", NewInteger(1))

	dict2 := NewDictionary()
	dict2.Set("Count", NewInteger(2)) // Overwrites
	dict2.Set("Title", NewString("Test"))

	dict1.Merge(dict2)

	assert.Equal(t, 3, dict1.Len())
	assert.Equal(t, "Page", dict1.GetName("Type").Value())
	assert.Equal(t, int64(2), dict1.GetInteger("Count")) // Overwritten
	assert.Equal(t, "Test", dict1.GetString("Title"))

	// Merge with nil should not crash
	dict1.Merge(nil)
	assert.Equal(t, 3, dict1.Len())
}

// ============================================================================
// Integration Tests
// ============================================================================

func TestComplexStructure(t *testing.T) {
	// Build a structure similar to PDF page object
	page := NewDictionary()
	page.SetName("Type", "Page")

	// MediaBox array
	mediaBox := NewArrayFromSlice([]PdfObject{
		NewInteger(0),
		NewInteger(0),
		NewReal(595.276),
		NewReal(841.890),
	})
	page.Set("MediaBox", mediaBox)

	// Resources dictionary
	resources := NewDictionary()
	fonts := NewDictionary()
	fonts.SetName("F1", "Font1")
	fonts.SetName("F2", "Font2")
	resources.Set("Font", fonts)
	page.Set("Resources", resources)

	// Verify structure
	assert.Equal(t, "Page", page.GetName("Type").Value())
	assert.Equal(t, 4, page.GetArray("MediaBox").Len())
	assert.NotNil(t, page.GetDictionary("Resources"))
	assert.NotNil(t, page.GetDictionary("Resources").GetDictionary("Font"))

	// Write to buffer
	var buf bytes.Buffer
	_, err := page.WriteTo(&buf)
	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "/Type /Page")
	assert.Contains(t, output, "/MediaBox")
	assert.Contains(t, output, "/Resources")
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkArray_Append(b *testing.B) {
	arr := NewArray()
	obj := NewInteger(42)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		arr.Append(obj)
	}
}

func BenchmarkArray_Get(b *testing.B) {
	arr := NewArray()
	for i := 0; i < 100; i++ {
		arr.Append(NewInteger(int64(i)))
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = arr.Get(50)
	}
}

func BenchmarkDictionary_Set(b *testing.B) {
	dict := NewDictionary()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dict.Set("Key", NewInteger(42))
	}
}

func BenchmarkDictionary_Get(b *testing.B) {
	dict := NewDictionary()
	dict.Set("Key", NewInteger(42))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dict.Get("Key")
	}
}

func BenchmarkArray_WriteTo(b *testing.B) {
	arr := NewArray()
	for i := 0; i < 10; i++ {
		arr.Append(NewInteger(int64(i)))
	}

	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = arr.WriteTo(&buf)
	}
}

func BenchmarkDictionary_WriteTo(b *testing.B) {
	dict := NewDictionary()
	for i := 0; i < 10; i++ {
		dict.Set("Key"+string(rune(i)), NewInteger(int64(i)))
	}

	var buf bytes.Buffer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		_, _ = dict.WriteTo(&buf)
	}
}
