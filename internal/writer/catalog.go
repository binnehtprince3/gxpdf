package writer

import (
	"bytes"
	"fmt"

	"github.com/coregx/gxpdf/internal/document"
)

// createCatalog creates the PDF Catalog object (document root).
//
// The Catalog is the root of the document object hierarchy and
// contains references to other objects like the Pages tree.
//
// Format:
//
//	<< /Type /Catalog /Pages N 0 R >>
//
// Parameters:
//   - pagesRef: Object number of the Pages root object
//   - doc: Document for additional catalog entries (metadata, etc.)
//
// Returns:
//
//	The Catalog indirect object
func (w *PdfWriter) createCatalog(pagesRef int, doc *document.Document) *IndirectObject {
	catalogNum := w.allocateObjNum()

	var catalog bytes.Buffer
	catalog.WriteString("<<")
	catalog.WriteString(" /Type /Catalog")
	catalog.WriteString(fmt.Sprintf(" /Pages %d 0 R", pagesRef))

	// Add optional entries
	// TODO: Add more catalog entries as needed:
	// - /PageLayout (SinglePage, OneColumn, etc.)
	// - /PageMode (UseNone, UseOutlines, UseThumbs, FullScreen)
	// - /Outlines (bookmarks)
	// - /Names (named destinations)
	// - /OpenAction (action to perform when document is opened)

	catalog.WriteString(" >>")

	return NewIndirectObject(catalogNum, 0, catalog.Bytes())
}
