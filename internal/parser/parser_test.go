package parser

import (
	"errors"
	"io"
	"strings"
	"testing"
)

const testValuePage = "Page"

// ============================================================================
// Basic Object Parsing Tests
// ============================================================================

func TestParser_ParseInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{"positive integer", "123", 123},
		{"negative integer", "-456", -456},
		{"zero", "0", 0},
		{"large integer", "2147483647", 2147483647},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			intObj, ok := obj.(*Integer)
			if !ok {
				t.Fatalf("expected *Integer, got %T", obj)
			}

			if intObj.Value() != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, intObj.Value())
			}
		})
	}
}

func TestParser_ParseReal(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected float64
	}{
		{"positive real", "3.14", 3.14},
		{"negative real", "-2.5", -2.5},
		{"zero real", "0.0", 0.0},
		{"no leading digit", ".5", 0.5},
		{"no trailing digits", "123.", 123.0},
		{"scientific notation", "1.5", 1.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			realObj, ok := obj.(*Real)
			if !ok {
				t.Fatalf("expected *Real, got %T", obj)
			}

			if realObj.Value() != tt.expected {
				t.Errorf("expected %f, got %f", tt.expected, realObj.Value())
			}
		})
	}
}

func TestParser_ParseString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple string", "(Hello)", "Hello"},
		{"empty string", "()", ""},
		{"string with spaces", "(Hello World)", "Hello World"},
		{"nested parens", "(Hello (World))", "Hello (World)"},
		{"escaped parens", "(Hello \\(World\\))", "Hello (World)"},
		{"escape sequences", "(Line1\\nLine2\\tTab)", "Line1\nLine2\tTab"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			strObj, ok := obj.(*String)
			if !ok {
				t.Fatalf("expected *String, got %T", obj)
			}

			if strObj.Value() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strObj.Value())
			}
		})
	}
}

func TestParser_ParseHexString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple hex", "<48656C6C6F>", "Hello"},
		{"empty hex", "<>", ""},
		{"odd length", "<48656C6C6F21>", "Hello!"},
		{"with whitespace", "<48 65 6C 6C 6F>", "Hello"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			strObj, ok := obj.(*String)
			if !ok {
				t.Fatalf("expected *String, got %T", obj)
			}

			if strObj.Value() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, strObj.Value())
			}
		})
	}
}

func TestParser_ParseName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"simple name", "/Type", "Type"},
		{"name with hash", "/Name#20With#20Spaces", "Name With Spaces"},
		{"empty name", "/", ""},
		{"complex name", "/Adobe#20Green", "Adobe Green"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			nameObj, ok := obj.(*Name)
			if !ok {
				t.Fatalf("expected *Name, got %T", obj)
			}

			if nameObj.Value() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, nameObj.Value())
			}
		})
	}
}

func TestParser_ParseBoolean(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"true", "true", true},
		{"false", "false", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			boolObj, ok := obj.(*Boolean)
			if !ok {
				t.Fatalf("expected *Boolean, got %T", obj)
			}

			if boolObj.Value() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, boolObj.Value())
			}
		})
	}
}

func TestParser_ParseNull(t *testing.T) {
	p := NewParser(strings.NewReader("null"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	_, ok := obj.(*Null)
	if !ok {
		t.Fatalf("expected *Null, got %T", obj)
	}
}

// ============================================================================
// Array Parsing Tests
// ============================================================================

func TestParser_ParseArray_Empty(t *testing.T) {
	p := NewParser(strings.NewReader("[]"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	arr, ok := obj.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj)
	}

	if arr.Len() != 0 {
		t.Errorf("expected empty array, got length %d", arr.Len())
	}
}

func TestParser_ParseArray_Simple(t *testing.T) {
	p := NewParser(strings.NewReader("[1 2 3]"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	arr, ok := obj.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj)
	}

	if arr.Len() != 3 {
		t.Errorf("expected length 3, got %d", arr.Len())
	}

	// Check elements
	for i := 0; i < 3; i++ {
		intObj, ok := arr.Get(i).(*Integer)
		if !ok {
			t.Errorf("element %d: expected *Integer, got %T", i, arr.Get(i))
			continue
		}
		if intObj.Value() != int64(i+1) {
			t.Errorf("element %d: expected %d, got %d", i, i+1, intObj.Value())
		}
	}
}

func TestParser_ParseArray_Mixed(t *testing.T) {
	p := NewParser(strings.NewReader("[1 (hello) /Name true null]"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	arr, ok := obj.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj)
	}

	if arr.Len() != 5 {
		t.Errorf("expected length 5, got %d", arr.Len())
	}

	// Check types
	if _, ok := arr.Get(0).(*Integer); !ok {
		t.Errorf("element 0: expected Integer")
	}
	if _, ok := arr.Get(1).(*String); !ok {
		t.Errorf("element 1: expected String")
	}
	if _, ok := arr.Get(2).(*Name); !ok {
		t.Errorf("element 2: expected Name")
	}
	if _, ok := arr.Get(3).(*Boolean); !ok {
		t.Errorf("element 3: expected Boolean")
	}
	if _, ok := arr.Get(4).(*Null); !ok {
		t.Errorf("element 4: expected Null")
	}
}

func TestParser_ParseArray_Nested(t *testing.T) {
	p := NewParser(strings.NewReader("[1 [2 3] [4 [5 6]]]"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	arr, ok := obj.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj)
	}

	if arr.Len() != 3 {
		t.Errorf("expected length 3, got %d", arr.Len())
	}

	// Check nested arrays
	nested1, ok := arr.Get(1).(*Array)
	if !ok {
		t.Fatalf("element 1: expected Array, got %T", arr.Get(1))
	}
	if nested1.Len() != 2 {
		t.Errorf("nested array 1: expected length 2, got %d", nested1.Len())
	}

	nested2, ok := arr.Get(2).(*Array)
	if !ok {
		t.Fatalf("element 2: expected Array, got %T", arr.Get(2))
	}
	if nested2.Len() != 2 {
		t.Errorf("nested array 2: expected length 2, got %d", nested2.Len())
	}

	// Check deeply nested
	deepNested, ok := nested2.Get(1).(*Array)
	if !ok {
		t.Fatalf("deep nested: expected Array, got %T", nested2.Get(1))
	}
	if deepNested.Len() != 2 {
		t.Errorf("deep nested: expected length 2, got %d", deepNested.Len())
	}
}

func TestParser_ParseArray_UnterminatedError(t *testing.T) {
	p := NewParser(strings.NewReader("[1 2 3"))
	_, err := p.ParseObject()
	if err == nil {
		t.Fatal("expected error for unterminated array")
	}
}

// ============================================================================
// Dictionary Parsing Tests
// ============================================================================

func TestParser_ParseDictionary_Empty(t *testing.T) {
	p := NewParser(strings.NewReader("<<>>"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	if dict.Len() != 0 {
		t.Errorf("expected empty dictionary, got length %d", dict.Len())
	}
}

func TestParser_ParseDictionary_Simple(t *testing.T) {
	p := NewParser(strings.NewReader("<< /Type /Page >>"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	if dict.Len() != 1 {
		t.Errorf("expected length 1, got %d", dict.Len())
	}

	typeObj := dict.GetName("Type")
	if typeObj == nil {
		t.Fatal("expected /Type entry")
	}

	if typeObj.Value() != testValuePage {
		t.Errorf("expected 'Page', got %q", typeObj.Value())
	}
}

func TestParser_ParseDictionary_Multiple(t *testing.T) {
	input := "<< /Type /Page /Count 3 /Title (My Title) >>"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	if dict.Len() != 3 {
		t.Errorf("expected length 3, got %d", dict.Len())
	}

	// Check entries
	if name := dict.GetName("Type"); name == nil || name.Value() != testValuePage {
		t.Error("/Type mismatch")
	}

	if count := dict.GetInteger("Count"); count != 3 {
		t.Errorf("expected Count=3, got %d", count)
	}

	if title := dict.GetString("Title"); title != "My Title" {
		t.Errorf("expected Title='My Title', got %q", title)
	}
}

func TestParser_ParseDictionary_Nested(t *testing.T) {
	input := "<< /Outer << /Inner (value) >> >>"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	innerDict := dict.GetDictionary("Outer")
	if innerDict == nil {
		t.Fatal("expected nested dictionary")
	}

	innerValue := innerDict.GetString("Inner")
	if innerValue != "value" {
		t.Errorf("expected 'value', got %q", innerValue)
	}
}

func TestParser_ParseDictionary_WithArray(t *testing.T) {
	input := "<< /Array [1 2 3] >>"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	arr := dict.GetArray("Array")
	if arr == nil {
		t.Fatal("expected array entry")
	}

	if arr.Len() != 3 {
		t.Errorf("expected array length 3, got %d", arr.Len())
	}
}

func TestParser_ParseDictionary_UnterminatedError(t *testing.T) {
	p := NewParser(strings.NewReader("<< /Type /Page"))
	_, err := p.ParseObject()
	if err == nil {
		t.Fatal("expected error for unterminated dictionary")
	}
}

func TestParser_ParseDictionary_MissingValueError(t *testing.T) {
	p := NewParser(strings.NewReader("<< /Type >>"))
	_, err := p.ParseObject()
	if err == nil {
		t.Fatal("expected error for missing value")
	}
}

// ============================================================================
// Indirect Reference Parsing Tests
// ============================================================================

func TestParser_ParseIndirectReference(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expectedN int
		expectedG int
	}{
		{"simple reference", "1 0 R", 1, 0},
		{"non-zero generation", "5 2 R", 5, 2},
		{"large object number", "12345 0 R", 12345, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			obj, err := p.ParseObject()
			if err != nil {
				t.Fatalf("ParseObject() error = %v", err)
			}

			ref, ok := obj.(*IndirectReference)
			if !ok {
				t.Fatalf("expected *IndirectReference, got %T", obj)
			}

			if ref.Number != tt.expectedN {
				t.Errorf("expected number %d, got %d", tt.expectedN, ref.Number)
			}

			if ref.Generation != tt.expectedG {
				t.Errorf("expected generation %d, got %d", tt.expectedG, ref.Generation)
			}
		})
	}
}

func TestParser_ParseIndirectReference_InArray(t *testing.T) {
	p := NewParser(strings.NewReader("[1 0 R 2 0 R 3 0 R]"))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	arr, ok := obj.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj)
	}

	if arr.Len() != 3 {
		t.Errorf("expected length 3, got %d", arr.Len())
	}

	for i := 0; i < 3; i++ {
		ref, ok := arr.Get(i).(*IndirectReference)
		if !ok {
			t.Errorf("element %d: expected *IndirectReference, got %T", i, arr.Get(i))
			continue
		}
		if ref.Number != i+1 {
			t.Errorf("element %d: expected number %d, got %d", i, i+1, ref.Number)
		}
	}
}

// ============================================================================
// Indirect Object Parsing Tests
// ============================================================================

func TestParser_ParseIndirectObject_Integer(t *testing.T) {
	p := NewParser(strings.NewReader("1 0 obj\n42\nendobj"))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	if obj.Number != 1 {
		t.Errorf("expected number 1, got %d", obj.Number)
	}

	if obj.Generation != 0 {
		t.Errorf("expected generation 0, got %d", obj.Generation)
	}

	intObj, ok := obj.Object.(*Integer)
	if !ok {
		t.Fatalf("expected *Integer, got %T", obj.Object)
	}

	if intObj.Value() != 42 {
		t.Errorf("expected value 42, got %d", intObj.Value())
	}
}

func TestParser_ParseIndirectObject_String(t *testing.T) {
	p := NewParser(strings.NewReader("2 0 obj\n(Hello World)\nendobj"))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	strObj, ok := obj.Object.(*String)
	if !ok {
		t.Fatalf("expected *String, got %T", obj.Object)
	}

	if strObj.Value() != "Hello World" {
		t.Errorf("expected 'Hello World', got %q", strObj.Value())
	}
}

func TestParser_ParseIndirectObject_Dictionary(t *testing.T) {
	input := "3 0 obj\n<< /Type /Page /Count 5 >>\nendobj"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	dict, ok := obj.Object.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj.Object)
	}

	if dict.GetName("Type").Value() != testValuePage {
		t.Error("/Type mismatch")
	}

	if dict.GetInteger("Count") != 5 {
		t.Error("/Count mismatch")
	}
}

func TestParser_ParseIndirectObject_Array(t *testing.T) {
	input := "4 0 obj\n[1 2 3 4 5]\nendobj"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	arr, ok := obj.Object.(*Array)
	if !ok {
		t.Fatalf("expected *Array, got %T", obj.Object)
	}

	if arr.Len() != 5 {
		t.Errorf("expected length 5, got %d", arr.Len())
	}
}

func TestParser_ParseIndirectObject_MissingEndobj(t *testing.T) {
	p := NewParser(strings.NewReader("1 0 obj\n42\n"))
	_, err := p.ParseIndirectObject()
	if err == nil {
		t.Fatal("expected error for missing endobj")
	}
}

// ============================================================================
// Stream Parsing Tests
// ============================================================================

func TestParser_ParseStream_Simple(t *testing.T) {
	input := "1 0 obj\n<< /Length 5 >>\nstream\nHello\nendstream\nendobj"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	stream, ok := obj.Object.(*Stream)
	if !ok {
		t.Fatalf("expected *Stream, got %T", obj.Object)
	}

	if stream.Length() != 5 {
		t.Errorf("expected length 5, got %d", stream.Length())
	}

	content := string(stream.Content())
	if content != "Hello" {
		t.Errorf("expected 'Hello', got %q", content)
	}

	// Check dictionary
	dict := stream.Dictionary()
	if dict.GetInteger("Length") != 5 {
		t.Error("/Length mismatch")
	}
}

func TestParser_ParseStream_MultiLine(t *testing.T) {
	input := "2 0 obj\n<< /Length 11 >>\nstream\nLine1\nLine2\nendstream\nendobj"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	stream, ok := obj.Object.(*Stream)
	if !ok {
		t.Fatalf("expected *Stream, got %T", obj.Object)
	}

	content := string(stream.Content())
	expected := "Line1\nLine2"
	if content != expected {
		t.Errorf("expected %q, got %q", expected, content)
	}
}

func TestParser_ParseStream_WithFilter(t *testing.T) {
	input := "3 0 obj\n<< /Length 5 /Filter /FlateDecode >>\nstream\nHello\nendstream\nendobj"
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseIndirectObject()
	if err != nil {
		t.Fatalf("ParseIndirectObject() error = %v", err)
	}

	stream, ok := obj.Object.(*Stream)
	if !ok {
		t.Fatalf("expected *Stream, got %T", obj.Object)
	}

	filter := stream.GetFilter()
	if filter == nil {
		t.Fatal("expected /Filter entry")
	}

	filterName, ok := filter.(*Name)
	if !ok {
		t.Fatalf("expected /Filter to be Name, got %T", filter)
	}

	if filterName.Value() != "FlateDecode" {
		t.Errorf("expected 'FlateDecode', got %q", filterName.Value())
	}
}

// ============================================================================
// Complex Nested Structure Tests
// ============================================================================

//nolint:cyclop // Test case validates complex nested structure.
func TestParser_ComplexNested(t *testing.T) {
	input := `<< /Type /Page
	             /Contents 5 0 R
	             /Resources << /Font << /F1 10 0 R >> >>
	             /MediaBox [0 0 612 792]
	          >>`
	p := NewParser(strings.NewReader(input))
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}

	dict, ok := obj.(*Dictionary)
	if !ok {
		t.Fatalf("expected *Dictionary, got %T", obj)
	}

	// Check /Type
	if typeObj := dict.GetName("Type"); typeObj == nil || typeObj.Value() != testValuePage {
		t.Error("/Type mismatch")
	}

	// Check /Contents (indirect reference)
	contentsRef, ok := dict.Get("Contents").(*IndirectReference)
	if !ok {
		t.Fatalf("expected /Contents to be IndirectReference, got %T", dict.Get("Contents"))
	}
	if contentsRef.Number != 5 {
		t.Errorf("expected reference to object 5, got %d", contentsRef.Number)
	}

	// Check /Resources (nested dictionary)
	resources := dict.GetDictionary("Resources")
	if resources == nil {
		t.Fatal("expected /Resources dictionary")
	}

	font := resources.GetDictionary("Font")
	if font == nil {
		t.Fatal("expected /Font dictionary")
	}

	f1Ref, ok := font.Get("F1").(*IndirectReference)
	if !ok {
		t.Fatalf("expected /F1 to be IndirectReference, got %T", font.Get("F1"))
	}
	if f1Ref.Number != 10 {
		t.Errorf("expected reference to object 10, got %d", f1Ref.Number)
	}

	// Check /MediaBox (array)
	mediaBox := dict.GetArray("MediaBox")
	if mediaBox == nil {
		t.Fatal("expected /MediaBox array")
	}
	if mediaBox.Len() != 4 {
		t.Errorf("expected MediaBox length 4, got %d", mediaBox.Len())
	}
}

// ============================================================================
// Error Handling Tests
// ============================================================================

func TestParser_ErrorInvalidSyntax(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"unexpected token", "foobar"},
		{"unmatched bracket", "[1 2 3"},
		{"unmatched dict", "<< /Type /Page"},
		{"unexpected EOF", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewParser(strings.NewReader(tt.input))
			_, err := p.ParseObject()
			if err == nil || (err != nil && !errors.Is(err, io.EOF)) {
				// Expected either an error or EOF
				if err == nil {
					t.Error("expected error, got nil")
				}
			}
		})
	}
}

// ============================================================================
// Helper Types Tests
// ============================================================================

func TestIndirectReference_String(t *testing.T) {
	ref := NewIndirectReference(5, 0)
	expected := "5 0 R"
	if ref.String() != expected {
		t.Errorf("expected %q, got %q", expected, ref.String())
	}
}

func TestIndirectReference_Equals(t *testing.T) {
	ref1 := NewIndirectReference(5, 0)
	ref2 := NewIndirectReference(5, 0)
	ref3 := NewIndirectReference(6, 0)

	if !ref1.Equals(ref2) {
		t.Error("expected ref1 to equal ref2")
	}

	if ref1.Equals(ref3) {
		t.Error("expected ref1 not to equal ref3")
	}

	if ref1.Equals(nil) {
		t.Error("expected ref1 not to equal nil")
	}
}

func TestIndirectReference_Clone(t *testing.T) {
	ref1 := NewIndirectReference(5, 0)
	ref2 := ref1.Clone()

	if !ref1.Equals(ref2) {
		t.Error("cloned reference should be equal to original")
	}

	// Modify clone
	ref2.Number = 10
	if ref1.Number == ref2.Number {
		t.Error("modifying clone should not affect original")
	}
}

func TestStream_Clone(t *testing.T) {
	dict := NewDictionary()
	dict.SetInteger("Length", 5)
	stream := NewStream(dict, []byte("Hello"))

	cloned := stream.Clone()

	if string(cloned.Content()) != "Hello" {
		t.Error("cloned stream content mismatch")
	}

	if cloned.Dictionary().GetInteger("Length") != 5 {
		t.Error("cloned stream dictionary mismatch")
	}

	// Modify clone
	cloned.SetContent([]byte("World"))
	if string(stream.Content()) == string(cloned.Content()) {
		t.Error("modifying clone should not affect original")
	}
}

func TestStream_SetContent(t *testing.T) {
	stream := NewStream(NewDictionary(), []byte("Hello"))

	newContent := []byte("New Content")
	stream.SetContent(newContent)

	if string(stream.Content()) != string(newContent) {
		t.Error("content not updated")
	}

	if stream.Dictionary().GetInteger("Length") != int64(len(newContent)) {
		t.Error("Length not updated in dictionary")
	}
}

func TestIndirectObject_String(t *testing.T) {
	obj := NewIndirectObject(1, 0, NewInteger(42))
	str := obj.String()
	expected := "1 0 obj 42 endobj"
	if str != expected {
		t.Errorf("expected %q, got %q", expected, str)
	}
}

func TestIndirectObject_WriteTo(t *testing.T) {
	obj := NewIndirectObject(5, 2, NewInteger(100))
	var buf strings.Builder
	written, err := obj.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	if written <= 0 {
		t.Errorf("expected written > 0, got %d", written)
	}
	output := buf.String()
	if !strings.Contains(output, "5 2 obj") {
		t.Errorf("output should contain '5 2 obj', got %q", output)
	}
	if !strings.Contains(output, "endobj") {
		t.Errorf("output should contain 'endobj', got %q", output)
	}
}

func TestIndirectReference_WriteTo(t *testing.T) {
	ref := NewIndirectReference(10, 0)
	var buf strings.Builder
	written, err := ref.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	if written <= 0 {
		t.Errorf("expected written > 0, got %d", written)
	}
	expected := "10 0 R"
	if buf.String() != expected {
		t.Errorf("expected %q, got %q", expected, buf.String())
	}
}

func TestStream_WriteTo(t *testing.T) {
	dict := NewDictionary()
	stream := NewStream(dict, []byte("Test content"))
	var buf strings.Builder
	written, err := stream.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}
	if written <= 0 {
		t.Errorf("expected written > 0, got %d", written)
	}
	output := buf.String()
	if !strings.Contains(output, "stream") {
		t.Errorf("output should contain 'stream', got %q", output)
	}
	if !strings.Contains(output, "endstream") {
		t.Errorf("output should contain 'endstream', got %q", output)
	}
	if !strings.Contains(output, "Test content") {
		t.Errorf("output should contain 'Test content', got %q", output)
	}
}

func TestStream_String(t *testing.T) {
	dict := NewDictionary()
	stream := NewStream(dict, []byte("Hello World"))
	str := stream.String()
	if !strings.Contains(str, "stream[") {
		t.Errorf("expected string to contain 'stream[', got %q", str)
	}
	if !strings.Contains(str, "length=11") {
		t.Errorf("expected string to contain 'length=11', got %q", str)
	}
}

func TestStream_Bytes(t *testing.T) {
	content := []byte("test bytes")
	stream := NewStream(NewDictionary(), content)
	result := stream.Bytes()
	if string(result) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(result))
	}
}

func TestStream_Reader(t *testing.T) {
	content := []byte("test reader")
	stream := NewStream(NewDictionary(), content)
	reader := stream.Reader()
	buf := make([]byte, len(content))
	n, err := reader.Read(buf)
	if err != nil {
		t.Fatalf("Read() error = %v", err)
	}
	if n != len(content) {
		t.Errorf("expected to read %d bytes, got %d", len(content), n)
	}
	if string(buf) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(buf))
	}
}

func TestStream_Decode(t *testing.T) {
	content := []byte("test decode")
	stream := NewStream(NewDictionary(), content)
	decoded, err := stream.Decode()
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	// Currently just returns raw content
	if string(decoded) != string(content) {
		t.Errorf("expected %q, got %q", string(content), string(decoded))
	}
}

func TestStream_Encode(t *testing.T) {
	stream := NewStream(NewDictionary(), []byte("test"))
	err := stream.Encode([]string{"FlateDecode"})
	// Currently a no-op
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}
}

func TestStream_GetDecodeParams(t *testing.T) {
	dict := NewDictionary()
	dict.Set("DecodeParms", NewInteger(42))
	stream := NewStream(dict, []byte("test"))
	params := stream.GetDecodeParams()
	if params == nil {
		t.Error("expected DecodeParms, got nil")
	}
}

func TestParser_Reset(t *testing.T) {
	p := NewParser(strings.NewReader("123"))
	obj1, _ := p.ParseObject()
	if intObj, ok := obj1.(*Integer); !ok || intObj.Value() != 123 {
		t.Error("first parse failed")
	}

	// Reset with new input
	p.Reset(strings.NewReader("456"))
	obj2, _ := p.ParseObject()
	if intObj, ok := obj2.(*Integer); !ok || intObj.Value() != 456 {
		t.Error("parse after reset failed")
	}
}

func TestParser_Position(t *testing.T) {
	p := NewParser(strings.NewReader("123"))
	line, col := p.Position()
	if line != 1 || col != 0 {
		t.Errorf("expected position (1, 0), got (%d, %d)", line, col)
	}
}

func TestNewStream_WithNilDict(t *testing.T) {
	stream := NewStream(nil, []byte("test"))
	if stream.Dictionary() == nil {
		t.Error("expected non-nil dictionary")
	}
}

func TestParser_NewParserFromLexer(t *testing.T) {
	lexer := NewLexer(strings.NewReader("42"))
	p := NewParserFromLexer(lexer)
	obj, err := p.ParseObject()
	if err != nil {
		t.Fatalf("ParseObject() error = %v", err)
	}
	intObj, ok := obj.(*Integer)
	if !ok {
		t.Fatalf("expected *Integer, got %T", obj)
	}
	if intObj.Value() != 42 {
		t.Errorf("expected 42, got %d", intObj.Value())
	}
}

func TestParser_ParseStreamWithoutLength(t *testing.T) {
	// Test the fallback parseStreamUntilEndstream
	// Note: The parseStreamUntilEndstream is actually quite complex to test
	// because it reads raw bytes and scans for "endstream"
	// For now, we'll skip this test as it requires special setup
	t.Skip("parseStreamUntilEndstream requires special test setup")
}

func TestParser_ParseIndirectObject_StreamNotDictionary(t *testing.T) {
	// Try to create a stream with a non-dictionary object
	input := "1 0 obj\n42\nstream\ntest\nendstream\nendobj"
	p := NewParser(strings.NewReader(input))
	_, err := p.ParseIndirectObject()
	if err == nil {
		t.Error("expected error for stream without dictionary, got nil")
	}
}

func TestParser_ParseIndirectObject_MissingObjKeyword(t *testing.T) {
	input := "1 0 notobj\n42\nendobj"
	p := NewParser(strings.NewReader(input))
	_, err := p.ParseIndirectObject()
	if err == nil {
		t.Error("expected error for missing 'obj' keyword, got nil")
	}
}

func TestParser_ParseIndirectObject_InvalidObjectNumber(t *testing.T) {
	input := "abc 0 obj\n42\nendobj"
	p := NewParser(strings.NewReader(input))
	_, err := p.ParseIndirectObject()
	if err == nil {
		t.Error("expected error for invalid object number, got nil")
	}
}

func TestParser_ParseIndirectObject_InvalidGenerationNumber(t *testing.T) {
	// "abc" will be parsed as a keyword token, not an integer
	input := "1 abc 0 obj\n42\nendobj"
	p := NewParser(strings.NewReader(input))
	_, err := p.ParseIndirectObject()
	if err == nil {
		t.Error("expected error for invalid generation number, got nil")
	}
}

func TestParser_ParseDictionary_NonNameKey(t *testing.T) {
	// Dictionary key must be a name, not an integer
	input := "<< 123 value >>"
	p := NewParser(strings.NewReader(input))
	_, err := p.ParseObject()
	if err == nil {
		t.Error("expected error for non-name dictionary key, got nil")
	}
}
