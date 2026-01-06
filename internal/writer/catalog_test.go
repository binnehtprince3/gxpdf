package writer

import (
	"strings"
	"testing"

	"github.com/coregx/gxpdf/internal/document"
)

func TestCreateCatalog(t *testing.T) {
	tests := []struct {
		name     string
		pagesRef int
		wantType string
		wantRef  string
	}{
		{
			name:     "simple catalog",
			pagesRef: 2,
			wantType: "/Type /Catalog",
			wantRef:  "/Pages 2 0 R",
		},
		{
			name:     "catalog with different pages ref",
			pagesRef: 5,
			wantType: "/Type /Catalog",
			wantRef:  "/Pages 5 0 R",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &PdfWriter{
				nextObjNum: 1,
			}

			doc := document.NewDocument()

			obj := w.createCatalog(tt.pagesRef, doc)

			if obj == nil {
				t.Fatal("createCatalog() returned nil")
			}

			if obj.Number != 1 {
				t.Errorf("Object number = %d, want 1", obj.Number)
			}

			if obj.Generation != 0 {
				t.Errorf("Generation = %d, want 0", obj.Generation)
			}

			data := string(obj.Data)

			if !strings.Contains(data, tt.wantType) {
				t.Errorf("Catalog should contain '%s', got: %s", tt.wantType, data)
			}

			if !strings.Contains(data, tt.wantRef) {
				t.Errorf("Catalog should contain '%s', got: %s", tt.wantRef, data)
			}

			// Check dictionary format
			if !strings.HasPrefix(data, "<<") {
				t.Error("Catalog should start with '<<'")
			}

			if !strings.HasSuffix(data, ">>") {
				t.Error("Catalog should end with '>>'")
			}
		})
	}
}

func TestCreateCatalog_ObjectNumberAllocation(t *testing.T) {
	w := &PdfWriter{
		nextObjNum: 1,
	}

	doc := document.NewDocument()

	// Create first catalog
	obj1 := w.createCatalog(2, doc)
	if obj1.Number != 1 {
		t.Errorf("First catalog object number = %d, want 1", obj1.Number)
	}

	if w.nextObjNum != 2 {
		t.Errorf("After first allocation, nextObjNum = %d, want 2", w.nextObjNum)
	}

	// Create second catalog
	obj2 := w.createCatalog(3, doc)
	if obj2.Number != 2 {
		t.Errorf("Second catalog object number = %d, want 2", obj2.Number)
	}

	if w.nextObjNum != 3 {
		t.Errorf("After second allocation, nextObjNum = %d, want 3", w.nextObjNum)
	}
}

func TestCreateCatalog_ValidDictionary(t *testing.T) {
	w := &PdfWriter{
		nextObjNum: 1,
	}

	doc := document.NewDocument()
	obj := w.createCatalog(2, doc)

	data := string(obj.Data)

	// Check for required PDF dictionary elements
	if !strings.Contains(data, "<<") || !strings.Contains(data, ">>") {
		t.Error("Should contain dictionary delimiters << and >>")
	}

	// Check for proper spacing in dictionary
	if strings.Contains(data, "<<>") {
		t.Error("Dictionary should not be empty")
	}

	// Verify all entries are inside dictionary
	typeIndex := strings.Index(data, "/Type")
	pagesIndex := strings.Index(data, "/Pages")
	openIndex := strings.Index(data, "<<")
	closeIndex := strings.Index(data, ">>")

	if typeIndex < openIndex || typeIndex > closeIndex {
		t.Error("/Type should be inside dictionary")
	}

	if pagesIndex < openIndex || pagesIndex > closeIndex {
		t.Error("/Pages should be inside dictionary")
	}
}
