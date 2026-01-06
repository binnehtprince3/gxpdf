package writer_test

import (
	"fmt"

	"github.com/coregx/gxpdf/internal/writer"
)

// ExampleResourceDictionary demonstrates basic usage of ResourceDictionary.
func ExampleResourceDictionary() {
	// Create a new resource dictionary.
	rd := writer.NewResourceDictionary()

	// Add a font resource.
	fontName := rd.AddFont(5) // Font object is at 5 0 R
	fmt.Printf("Font name: %s\n", fontName)

	// Add another font.
	fontName2 := rd.AddFont(6) // Font object is at 6 0 R
	fmt.Printf("Second font name: %s\n", fontName2)

	// Add an image.
	imageName := rd.AddImage(10) // Image object is at 10 0 R
	fmt.Printf("Image name: %s\n", imageName)

	// Get the PDF dictionary.
	fmt.Println(rd.String())

	// Output:
	// Font name: F1
	// Second font name: F2
	// Image name: Im1
	// << /Font << /F1 5 0 R /F2 6 0 R >> /XObject << /Im1 10 0 R >> /ProcSet [/PDF /Text /ImageB /ImageC /ImageI] >>
}

// ExampleResourceDictionary_empty demonstrates an empty resource dictionary.
func ExampleResourceDictionary_empty() {
	rd := writer.NewResourceDictionary()

	// Check if dictionary has resources.
	fmt.Printf("Has resources: %v\n", rd.HasResources())

	// Empty dictionary outputs minimal syntax.
	fmt.Println(rd.String())

	// Output:
	// Has resources: false
	// << >>
}

// ExampleResourceDictionary_fonts demonstrates managing font resources.
func ExampleResourceDictionary_fonts() {
	rd := writer.NewResourceDictionary()

	// Add fonts for a multi-font document.
	helvetica := rd.AddFont(5) // Object 5: Helvetica
	times := rd.AddFont(6)     // Object 6: Times-Roman
	courier := rd.AddFont(7)   // Object 7: Courier

	fmt.Printf("Helvetica: %s, Times: %s, Courier: %s\n", helvetica, times, courier)
	fmt.Println(rd.String())

	// Output:
	// Helvetica: F1, Times: F2, Courier: F3
	// << /Font << /F1 5 0 R /F2 6 0 R /F3 7 0 R >> /ProcSet [/PDF /Text /ImageB /ImageC /ImageI] >>
}

// ExampleResourceDictionary_images demonstrates managing image resources.
func ExampleResourceDictionary_images() {
	rd := writer.NewResourceDictionary()

	// Add images for a document with graphics.
	logo := rd.AddImage(10)    // Object 10: Company logo
	photo := rd.AddImage(11)   // Object 11: Product photo
	diagram := rd.AddImage(12) // Object 12: Technical diagram

	fmt.Printf("Logo: %s, Photo: %s, Diagram: %s\n", logo, photo, diagram)
	fmt.Println(rd.String())

	// Output:
	// Logo: Im1, Photo: Im2, Diagram: Im3
	// << /XObject << /Im1 10 0 R /Im2 11 0 R /Im3 12 0 R >> /ProcSet [/PDF /Text /ImageB /ImageC /ImageI] >>
}

// ExampleResourceDictionary_mixed demonstrates combining different resource types.
func ExampleResourceDictionary_mixed() {
	rd := writer.NewResourceDictionary()

	// Real-world document with fonts, images, and graphics states.
	rd.AddFont(5)       // Helvetica
	rd.AddFont(6)       // Times-Roman
	rd.AddImage(10)     // Logo
	rd.AddExtGState(15) // Transparency state

	fmt.Printf("Has resources: %v\n", rd.HasResources())
	fmt.Println(rd.String())

	// Output:
	// Has resources: true
	// << /Font << /F1 5 0 R /F2 6 0 R >> /XObject << /Im1 10 0 R >> /ExtGState << /GS1 15 0 R >> /ProcSet [/PDF /Text /ImageB /ImageC /ImageI] >>
}
