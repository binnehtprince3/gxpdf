package writer

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewIndirectObject(t *testing.T) {
	tests := []struct {
		name       string
		number     int
		generation int
		data       []byte
	}{
		{
			name:       "simple object",
			number:     1,
			generation: 0,
			data:       []byte("<< /Type /Catalog >>"),
		},
		{
			name:       "object with generation",
			number:     5,
			generation: 2,
			data:       []byte("[ 1 2 3 4 5 ]"),
		},
		{
			name:       "empty data",
			number:     10,
			generation: 0,
			data:       []byte{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewIndirectObject(tt.number, tt.generation, tt.data)

			if obj == nil {
				t.Fatal("NewIndirectObject returned nil")
			}

			if obj.Number != tt.number {
				t.Errorf("Number = %d, want %d", obj.Number, tt.number)
			}

			if obj.Generation != tt.generation {
				t.Errorf("Generation = %d, want %d", obj.Generation, tt.generation)
			}

			if !bytes.Equal(obj.Data, tt.data) {
				t.Errorf("Data = %v, want %v", obj.Data, tt.data)
			}
		})
	}
}

func TestIndirectObject_WriteTo(t *testing.T) {
	tests := []struct {
		name       string
		number     int
		generation int
		data       []byte
		want       string
	}{
		{
			name:       "catalog object",
			number:     1,
			generation: 0,
			data:       []byte("<< /Type /Catalog /Pages 2 0 R >>"),
			want:       "1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n",
		},
		{
			name:       "pages object",
			number:     2,
			generation: 0,
			data:       []byte("<< /Type /Pages /Kids [3 0 R] /Count 1 >>"),
			want:       "2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n",
		},
		{
			name:       "data with trailing newline",
			number:     3,
			generation: 0,
			data:       []byte("<< /Type /Page >>\n"),
			want:       "3 0 obj\n<< /Type /Page >>\nendobj\n",
		},
		{
			name:       "object with generation",
			number:     5,
			generation: 2,
			data:       []byte("[ 1 2 3 ]"),
			want:       "5 2 obj\n[ 1 2 3 ]\nendobj\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			obj := NewIndirectObject(tt.number, tt.generation, tt.data)
			var buf bytes.Buffer

			n, err := obj.WriteTo(&buf)
			if err != nil {
				t.Fatalf("WriteTo() error = %v", err)
			}

			got := buf.String()
			if got != tt.want {
				t.Errorf("WriteTo() output:\ngot:  %q\nwant: %q", got, tt.want)
			}

			if n != int64(len(tt.want)) {
				t.Errorf("WriteTo() returned %d bytes, want %d", n, len(tt.want))
			}
		})
	}
}

func TestIndirectObject_String(t *testing.T) {
	obj := NewIndirectObject(1, 0, []byte("test data"))
	str := obj.String()

	if !strings.Contains(str, "Object 1 0") {
		t.Errorf("String() should contain 'Object 1 0', got: %s", str)
	}

	if !strings.Contains(str, "9 bytes") {
		t.Errorf("String() should contain '9 bytes', got: %s", str)
	}
}

func TestIndirectObject_WriteTo_EmptyData(t *testing.T) {
	obj := NewIndirectObject(1, 0, []byte{})
	var buf bytes.Buffer

	n, err := obj.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	// Implementation outputs no extra newline for empty data
	want := "1 0 obj\nendobj\n"
	got := buf.String()

	if got != want {
		t.Errorf("WriteTo() with empty data:\ngot:  %q\nwant: %q", got, want)
	}

	if n != int64(len(want)) {
		t.Errorf("WriteTo() returned %d bytes, want %d", n, len(want))
	}
}

func TestIndirectObject_WriteTo_MultilineData(t *testing.T) {
	data := []byte(`<<
  /Type /Page
  /MediaBox [0 0 595 842]
>>`)

	obj := NewIndirectObject(3, 0, data)
	var buf bytes.Buffer

	_, err := obj.WriteTo(&buf)
	if err != nil {
		t.Fatalf("WriteTo() error = %v", err)
	}

	output := buf.String()

	// Check structure
	if !strings.HasPrefix(output, "3 0 obj\n") {
		t.Errorf("Output should start with '3 0 obj\\n'")
	}

	if !strings.HasSuffix(output, "endobj\n") {
		t.Errorf("Output should end with 'endobj\\n'")
	}

	if !strings.Contains(output, "/Type /Page") {
		t.Errorf("Output should contain data")
	}
}
