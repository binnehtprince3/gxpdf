// Package main demonstrates the use of complex vector shapes in PDF creation.
//
// This example shows how to use:
// - Polygon (closed shapes with N vertices)
// - Polyline (open paths with N vertices)
// - Ellipse (not just circles)
// - Bezier curves (complex smooth curves)
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new PDF creator.
	c := creator.New()
	c.SetTitle("Complex Vector Shapes Demo")
	c.SetAuthor("GxPDF Examples")

	// Create a page.
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("failed to create page: %v", err)
	}

	// Draw title.
	if err := page.AddText("Complex Vector Shapes", 200, 750, creator.HelveticaBold, 18); err != nil {
		log.Fatalf("failed to add title: %v", err)
	}

	// Draw all shape examples.
	if err := drawPolygonExamples(page); err != nil {
		log.Fatalf("failed to draw polygons: %v", err)
	}

	if err := drawPolylineExamples(page); err != nil {
		log.Fatalf("failed to draw polylines: %v", err)
	}

	if err := drawEllipseExamples(page); err != nil {
		log.Fatalf("failed to draw ellipses: %v", err)
	}

	if err := drawBezierExamples(page); err != nil {
		log.Fatalf("failed to draw beziers: %v", err)
	}

	// Add footer.
	if err := page.AddText("Generated with GxPDF Complex Shapes API", 150, 50, creator.Helvetica, 10); err != nil {
		log.Fatalf("failed to add footer: %v", err)
	}

	// Write PDF to file.
	outputPath := "complex_shapes_demo.pdf"
	if err := c.WriteToFile(outputPath); err != nil {
		log.Fatalf("failed to write PDF: %v", err)
	}

	fmt.Printf("Successfully created %s\n", outputPath)
	fmt.Println("The PDF demonstrates:")
	fmt.Println("  - Polygons: Triangle, Pentagon, Star")
	fmt.Println("  - Polylines: Zigzag, Dashed Wave")
	fmt.Println("  - Ellipses: Horizontal, Vertical, Circle")
	fmt.Println("  - Bezier: S-curve, Multi-segment, Closed & Filled")
}

// drawPolygonExamples draws triangle, pentagon, and star polygons.
func drawPolygonExamples(page *creator.Page) error {
	// Triangle.
	if err := page.AddText("Polygon: Triangle", 50, 700, creator.Helvetica, 10); err != nil {
		return err
	}
	triangleVertices := []creator.Point{{X: 100, Y: 650}, {X: 150, Y: 600}, {X: 50, Y: 600}}
	if err := page.DrawPolygon(triangleVertices, &creator.PolygonOptions{
		StrokeColor: &creator.Black, StrokeWidth: 2.0, FillColor: &creator.Blue,
	}); err != nil {
		return err
	}

	// Pentagon.
	if err := page.AddText("Polygon: Pentagon", 200, 700, creator.Helvetica, 10); err != nil {
		return err
	}
	pentagonVertices := []creator.Point{
		{X: 250, Y: 650}, {X: 290, Y: 630}, {X: 280, Y: 590}, {X: 220, Y: 590}, {X: 210, Y: 630},
	}
	if err := page.DrawPolygon(pentagonVertices, &creator.PolygonOptions{
		StrokeColor: &creator.Black, StrokeWidth: 1.5, FillColor: &creator.Green,
	}); err != nil {
		return err
	}

	// Star.
	if err := page.AddText("Polygon: Star", 350, 700, creator.Helvetica, 10); err != nil {
		return err
	}
	starVertices := []creator.Point{
		{X: 400, Y: 650}, {X: 415, Y: 620}, {X: 435, Y: 615}, {X: 420, Y: 595}, {X: 425, Y: 570},
		{X: 400, Y: 585}, {X: 375, Y: 570}, {X: 380, Y: 595}, {X: 365, Y: 615}, {X: 385, Y: 620},
	}
	if err := page.DrawPolygon(starVertices, &creator.PolygonOptions{
		StrokeColor: &creator.Black, StrokeWidth: 1.0, FillColor: &creator.Yellow,
	}); err != nil {
		return err
	}

	return nil
}

// drawPolylineExamples draws zigzag and dashed wave polylines.
func drawPolylineExamples(page *creator.Page) error {
	// Zigzag.
	if err := page.AddText("Polyline: Zigzag", 50, 530, creator.Helvetica, 10); err != nil {
		return err
	}
	zigzagVertices := []creator.Point{
		{X: 50, Y: 500}, {X: 70, Y: 480}, {X: 90, Y: 500}, {X: 110, Y: 480}, {X: 130, Y: 500}, {X: 150, Y: 480},
	}
	if err := page.DrawPolyline(zigzagVertices, &creator.PolylineOptions{Color: creator.Red, Width: 2.0}); err != nil {
		return err
	}

	// Dashed wave.
	if err := page.AddText("Polyline: Dashed Wave", 200, 530, creator.Helvetica, 10); err != nil {
		return err
	}
	waveVertices := []creator.Point{
		{X: 200, Y: 490}, {X: 220, Y: 475}, {X: 240, Y: 490}, {X: 260, Y: 505},
		{X: 280, Y: 490}, {X: 300, Y: 475}, {X: 320, Y: 490},
	}
	if err := page.DrawPolyline(waveVertices, &creator.PolylineOptions{
		Color: creator.Blue, Width: 2.5, Dashed: true, DashArray: []float64{5, 3},
	}); err != nil {
		return err
	}

	return nil
}

// drawEllipseExamples draws horizontal ellipse, vertical ellipse, and circle.
func drawEllipseExamples(page *creator.Page) error {
	// Horizontal ellipse.
	if err := page.AddText("Ellipse: Horizontal", 50, 420, creator.Helvetica, 10); err != nil {
		return err
	}
	if err := page.DrawEllipse(100, 360, 60, 30, &creator.EllipseOptions{
		StrokeColor: &creator.Black, StrokeWidth: 1.5, FillColor: &creator.LightGray,
	}); err != nil {
		return err
	}

	// Vertical ellipse.
	if err := page.AddText("Ellipse: Vertical", 200, 420, creator.Helvetica, 10); err != nil {
		return err
	}
	if err := page.DrawEllipse(250, 360, 30, 60, &creator.EllipseOptions{
		StrokeColor: &creator.Red, StrokeWidth: 2.0, FillColor: &creator.Yellow,
	}); err != nil {
		return err
	}

	// Circle (rx = ry).
	if err := page.AddText("Ellipse: Circle", 340, 420, creator.Helvetica, 10); err != nil {
		return err
	}
	if err := page.DrawEllipse(390, 360, 50, 50, &creator.EllipseOptions{
		StrokeColor: &creator.Blue, StrokeWidth: 2.0,
	}); err != nil {
		return err
	}

	return nil
}

// drawBezierExamples draws S-curve, multi-segment curve, and closed filled curve.
func drawBezierExamples(page *creator.Page) error {
	// Simple S-curve.
	if err := page.AddText("Bezier: S-curve", 50, 280, creator.Helvetica, 10); err != nil {
		return err
	}
	sCurveSegments := []creator.BezierSegment{{
		Start: creator.Point{X: 50, Y: 250}, C1: creator.Point{X: 100, Y: 200},
		C2: creator.Point{X: 100, Y: 200}, End: creator.Point{X: 150, Y: 250},
	}}
	if err := page.DrawBezierCurve(sCurveSegments, &creator.BezierOptions{Color: creator.Blue, Width: 2.0}); err != nil {
		return err
	}

	// Multi-segment wave.
	if err := page.AddText("Bezier: Multi-segment", 200, 280, creator.Helvetica, 10); err != nil {
		return err
	}
	multiSegments := []creator.BezierSegment{
		{Start: creator.Point{X: 200, Y: 240}, C1: creator.Point{X: 220, Y: 220}, C2: creator.Point{X: 240, Y: 220}, End: creator.Point{X: 260, Y: 240}},
		{Start: creator.Point{X: 260, Y: 240}, C1: creator.Point{X: 280, Y: 260}, C2: creator.Point{X: 300, Y: 260}, End: creator.Point{X: 320, Y: 240}},
	}
	if err := page.DrawBezierCurve(multiSegments, &creator.BezierOptions{Color: creator.Green, Width: 2.0}); err != nil {
		return err
	}

	// Closed filled shape.
	if err := page.AddText("Bezier: Closed & Filled", 380, 280, creator.Helvetica, 10); err != nil {
		return err
	}
	closedSegments := []creator.BezierSegment{
		{Start: creator.Point{X: 430, Y: 220}, C1: creator.Point{X: 470, Y: 240}, C2: creator.Point{X: 470, Y: 260}, End: creator.Point{X: 430, Y: 280}},
		{Start: creator.Point{X: 430, Y: 280}, C1: creator.Point{X: 390, Y: 260}, C2: creator.Point{X: 390, Y: 240}, End: creator.Point{X: 430, Y: 220}},
	}
	fillColor := creator.Color{R: 1.0, G: 0.8, B: 0.0} // Orange
	if err := page.DrawBezierCurve(closedSegments, &creator.BezierOptions{
		Color: creator.Black, Width: 1.5, Closed: true, FillColor: &fillColor,
	}); err != nil {
		return err
	}

	return nil
}
