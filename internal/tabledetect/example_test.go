package tabledetect_test

import (
	"fmt"

	"github.com/coregx/gxpdf/internal/extractor"
	"github.com/coregx/gxpdf/internal/tabledetect"
)

// Example_basicUsage demonstrates basic table detection with default implementations.
func Example_basicUsage() {
	// Create detector with default implementations
	tableDetector := tabledetect.NewDefaultTableDetector()

	// Create sample text elements (from PDF extraction)
	textElements := []*extractor.TextElement{
		extractor.NewTextElement("Name", 10, 100, 50, 10, "/F1", 10),
		extractor.NewTextElement("Age", 110, 100, 30, 10, "/F1", 10),
		extractor.NewTextElement("John", 10, 80, 50, 10, "/F1", 10),
		extractor.NewTextElement("25", 110, 80, 30, 10, "/F1", 10),
	}

	// No graphics (will use stream mode)
	graphics := []*extractor.GraphicsElement{}

	// Detect tables
	tables, err := tableDetector.DetectTables(textElements, graphics)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Found %d table(s)\n", len(tables))
	for _, tbl := range tables {
		fmt.Printf("Table: %dx%d, Method: %s\n",
			tbl.RowCount(), tbl.ColumnCount(), tbl.Method.String())
	}
}

// Example_customImplementation demonstrates using custom implementations.
func Example_customImplementation() {
	// Create custom ruling detector with different settings
	customRuling := tabledetect.NewDefaultRulingLineDetector().
		WithMinLineLength(15.0).
		WithTolerance(3.0)

	// Create custom whitespace analyzer
	customWhitespace := tabledetect.NewDefaultWhitespaceAnalyzer().
		WithMinGapWidth(12.0).
		WithAlignmentTolerance(3.0)

	// Create detector with custom implementations
	tableDetector := tabledetect.NewTableDetectorWithDeps(
		customRuling,
		customWhitespace,
		tabledetect.NewDefaultGridBuilder(),
	)

	fmt.Printf("Created detector with custom settings: %T\n", tableDetector)
	// Output: Created detector with custom settings: *tabledetect.DefaultTableDetector
}

// Example_dependencyInjection demonstrates dependency injection for testing.
func Example_dependencyInjection() {
	// In production code, use default implementations
	productionDetector := tabledetect.NewDefaultTableDetector()

	// In test code, inject mocks for isolated testing
	// (This would typically be in a _test.go file)
	type mockRuling struct {
		tabledetect.RulingLineDetector
	}
	type mockWhitespace struct {
		tabledetect.WhitespaceAnalyzer
	}
	type mockGrid struct {
		tabledetect.GridBuilder
	}

	testDetector := tabledetect.NewTableDetectorWithDeps(
		&mockRuling{},
		&mockWhitespace{},
		&mockGrid{},
	)

	fmt.Printf("Production detector: %T\n", productionDetector)
	fmt.Printf("Test detector: %T\n", testDetector)
	// Output:
	// Production detector: *tabledetect.DefaultTableDetector
	// Test detector: *tabledetect.DefaultTableDetector
}

// Example_builderPattern demonstrates the builder pattern for configuration.
func Example_builderPattern() {
	// Use builder pattern for step-by-step configuration
	detector := tabledetect.NewDefaultTableDetector().
		WithRulingDetector(
			tabledetect.NewDefaultRulingLineDetector().WithMinLineLength(20.0),
		).
		WithWhitespaceAnalyzer(
			tabledetect.NewDefaultWhitespaceAnalyzer().WithMinGapWidth(15.0),
		).
		WithGridBuilder(
			tabledetect.NewDefaultGridBuilder().WithTolerance(2.5),
		)

	fmt.Printf("Configured detector: %T\n", detector)
	// Output: Configured detector: *tabledetect.DefaultTableDetector
}

// Example_backwardCompatibility demonstrates backward compatibility.
func Example_backwardCompatibility() {
	// Old code still works with deprecated constructors
	detector1 := tabledetect.NewTableDetector()

	// New code uses explicit constructors
	detector2 := tabledetect.NewDefaultTableDetector()

	// Both produce the same type
	fmt.Printf("Old constructor: %T\n", detector1)
	fmt.Printf("New constructor: %T\n", detector2)
	// Output:
	// Old constructor: *tabledetect.DefaultTableDetector
	// New constructor: *tabledetect.DefaultTableDetector
}

// Example_modeDetection demonstrates automatic mode detection.
func Example_modeDetection() {
	detector := tabledetect.NewDefaultTableDetector()

	// Scenario 1: PDF with ruling lines
	graphicsWithLines := []*extractor.GraphicsElement{
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(100, 0),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 50),
				extractor.NewPoint(100, 50),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(0, 0),
				extractor.NewPoint(0, 50),
			},
		},
		{
			Type: extractor.GraphicsTypeLine,
			Points: []extractor.Point{
				extractor.NewPoint(100, 0),
				extractor.NewPoint(100, 50),
			},
		},
	}

	mode1 := detector.DetectMode([]*extractor.TextElement{}, graphicsWithLines)
	fmt.Printf("With ruling lines: %s mode\n", mode1.String())

	// Scenario 2: PDF without ruling lines
	mode2 := detector.DetectMode([]*extractor.TextElement{}, []*extractor.GraphicsElement{})
	fmt.Printf("Without ruling lines: %s mode\n", mode2.String())

	// Output:
	// With ruling lines: Lattice mode
	// Without ruling lines: Stream mode
}

// Example_extensibility demonstrates how to extend with custom algorithms.
func Example_extensibility() {
	// Hypothetical custom implementation (not included in package)
	type MLBasedRulingDetector struct {
		tabledetect.RulingLineDetector
		// Custom fields for ML model...
	}

	// Easy to plug in custom implementations
	// detector := NewTableDetectorWithDeps(
	//     &MLBasedRulingDetector{},  // Custom ML-based detector
	//     NewDefaultWhitespaceAnalyzer(),
	//     NewDefaultGridBuilder(),
	// )

	fmt.Println("Custom implementations can be easily integrated via interfaces")
	// Output: Custom implementations can be easily integrated via interfaces
}
