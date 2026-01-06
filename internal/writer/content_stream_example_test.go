package writer_test

import (
	"fmt"

	"github.com/coregx/gxpdf/internal/writer"
)

// ExampleContentStreamWriter_simpleText demonstrates creating a simple text content stream.
func ExampleContentStreamWriter_simpleText() {
	csw := writer.NewContentStreamWriter()

	// Create text content.
	csw.BeginText()
	csw.SetFont("Helvetica", 12.0)
	csw.MoveTextPosition(100.0, 700.0)
	csw.ShowText("Hello World")
	csw.EndText()

	// Get the content stream.
	fmt.Println(csw.String())
	// Output:
	// BT
	// /Helvetica 12.00 Tf
	// 100.00 700.00 Td
	// (Hello World) Tj
	// ET
}

// ExampleContentStreamWriter_rectangle demonstrates drawing a rectangle.
func ExampleContentStreamWriter_rectangle() {
	csw := writer.NewContentStreamWriter()

	// Save state, set stroke color, draw rectangle.
	csw.SaveState()
	csw.SetStrokeColorRGB(1.0, 0.0, 0.0) // Red
	csw.SetLineWidth(2.0)
	csw.Rectangle(50.0, 50.0, 200.0, 100.0)
	csw.Stroke()
	csw.RestoreState()

	fmt.Println(csw.String())
	// Output:
	// q
	// 1.00 0.00 0.00 RG
	// 2.00 w
	// 50.00 50.00 200.00 100.00 re
	// S
	// Q
}

// ExampleContentStreamWriter_combined demonstrates combining text and graphics.
func ExampleContentStreamWriter_combined() {
	csw := writer.NewContentStreamWriter()

	// Draw a filled rectangle.
	csw.SaveState()
	csw.SetFillColorRGB(0.9, 0.9, 1.0) // Light blue
	csw.Rectangle(40.0, 40.0, 220.0, 80.0)
	csw.Fill()
	csw.RestoreState()

	// Add text on top.
	csw.BeginText()
	csw.SetFont("Times-Roman", 14.0)
	csw.MoveTextPosition(50.0, 60.0)
	csw.ShowText("Inside Box")
	csw.EndText()

	// Verify content was created.
	if csw.Len() > 0 {
		fmt.Println("Content stream created successfully")
	}
	// Output:
	// Content stream created successfully
}

// ExampleContentStreamWriter_compression demonstrates content stream compression.
func ExampleContentStreamWriter_compression() {
	csw := writer.NewContentStreamWriter()

	// Create some content.
	csw.BeginText()
	csw.SetFont("Courier", 10.0)
	csw.MoveTextPosition(50.0, 750.0)
	csw.ShowText("This text will be compressed")
	csw.EndText()

	// Compress the content.
	compressed, err := csw.Compress()
	if err != nil {
		fmt.Printf("Compression error: %v\n", err)
		return
	}

	// Show that compression worked.
	if len(compressed) > 0 && csw.Len() > 0 {
		fmt.Println("Content stream compressed successfully")
	}
	// Output:
	// Content stream compressed successfully
}
