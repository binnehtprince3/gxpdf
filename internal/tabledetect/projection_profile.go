// Package detector implements table detection algorithms.
package tabledetect

import (
	"fmt"
	"math"

	"github.com/coregx/gxpdf/internal/extractor"
)

// ProjectionProfile represents text density distribution across a dimension.
//
// A projection profile is a histogram that shows the distribution of text
// along an axis. This is used for stream mode table detection to find
// row and column boundaries based on whitespace.
//
// Algorithm inspired by Nurminen's thesis and tabula-java's detection algorithms.
// Reference: http://dspace.cc.tut.fi/dpub/bitstream/handle/123456789/21520/Nurminen.pdf
type ProjectionProfile struct {
	Bins    []float64 // Density values for each bin
	BinSize float64   // Size of each bin (in points)
	Min     float64   // Minimum coordinate
	Max     float64   // Maximum coordinate
}

// NewProjectionProfile creates a new ProjectionProfile.
func NewProjectionProfile(bins []float64, binSize, min, max float64) *ProjectionProfile {
	return &ProjectionProfile{
		Bins:    bins,
		BinSize: binSize,
		Min:     min,
		Max:     max,
	}
}

// BinCount returns the number of bins in the profile.
func (pp *ProjectionProfile) BinCount() int {
	return len(pp.Bins)
}

// GetDensity returns the density value at the given coordinate.
func (pp *ProjectionProfile) GetDensity(coord float64) float64 {
	if coord < pp.Min || coord > pp.Max {
		return 0
	}

	binIndex := int((coord - pp.Min) / pp.BinSize)
	if binIndex < 0 || binIndex >= len(pp.Bins) {
		return 0
	}

	return pp.Bins[binIndex]
}

// String returns a string representation of the projection profile.
func (pp *ProjectionProfile) String() string {
	return fmt.Sprintf("ProjectionProfile{bins=%d, binSize=%.2f, range=[%.2f, %.2f]}",
		len(pp.Bins), pp.BinSize, pp.Min, pp.Max)
}

// Gap represents a whitespace gap in a projection profile.
//
// Gaps indicate potential row or column boundaries in tables.
type Gap struct {
	Start float64 // Start coordinate of gap
	End   float64 // End coordinate of gap
	Width float64 // Width of gap
}

// NewGap creates a new Gap.
func NewGap(start, end float64) Gap {
	return Gap{
		Start: start,
		End:   end,
		Width: end - start,
	}
}

// Center returns the center coordinate of the gap.
func (g Gap) Center() float64 {
	return (g.Start + g.End) / 2
}

// String returns a string representation of the gap.
func (g Gap) String() string {
	return fmt.Sprintf("Gap{start=%.2f, end=%.2f, width=%.2f}", g.Start, g.End, g.Width)
}

// DefaultProjectionAnalyzer analyzes text distribution to find whitespace gaps.
//
// This is the default implementation of the ProjectionAnalyzer interface.
// It is used for stream mode table detection, where tables don't have
// visible borders and we need to infer structure from text positioning.
//
// Algorithm inspired by tabula-java's BasicExtractionAlgorithm.
// Reference: tabula-java/technology/tabula/extractors/BasicExtractionAlgorithm.java
type DefaultProjectionAnalyzer struct {
	binSize   float64 // Granularity of analysis (in points)
	threshold float64 // Minimum density to consider as text
}

// NewDefaultProjectionAnalyzer creates a new DefaultProjectionAnalyzer with default settings.
func NewDefaultProjectionAnalyzer() *DefaultProjectionAnalyzer {
	return &DefaultProjectionAnalyzer{
		binSize:   2.0,  // 2 points per bin (~0.7mm)
		threshold: 0.01, // Very low threshold for presence
	}
}

// NewProjectionAnalyzer creates a new DefaultProjectionAnalyzer with default settings.
// Deprecated: Use NewDefaultProjectionAnalyzer instead. Kept for backward compatibility.
func NewProjectionAnalyzer() *DefaultProjectionAnalyzer {
	return NewDefaultProjectionAnalyzer()
}

// WithBinSize sets the bin size.
func (pa *DefaultProjectionAnalyzer) WithBinSize(size float64) *DefaultProjectionAnalyzer {
	pa.binSize = size
	return pa
}

// WithThreshold sets the density threshold.
func (pa *DefaultProjectionAnalyzer) WithThreshold(threshold float64) *DefaultProjectionAnalyzer {
	pa.threshold = threshold
	return pa
}

// AnalyzeHorizontal creates a horizontal projection profile from text elements.
//
// The horizontal profile shows text density by Y coordinate (vertical distribution).
// High density = lots of text at that Y position
// Low density = whitespace (potential row boundary)
//
// Returns an array where index represents Y position, value represents text density.
func (pa *DefaultProjectionAnalyzer) AnalyzeHorizontal(elements []*extractor.TextElement) *ProjectionProfile {
	if len(elements) == 0 {
		return NewProjectionProfile([]float64{}, pa.binSize, 0, 0)
	}

	// Find Y coordinate range
	minY := elements[0].Y
	maxY := elements[0].Top()

	for _, elem := range elements {
		if elem.Y < minY {
			minY = elem.Y
		}
		if elem.Top() > maxY {
			maxY = elem.Top()
		}
	}

	// Create bins
	binCount := int(math.Ceil((maxY-minY)/pa.binSize)) + 1
	bins := make([]float64, binCount)

	// Fill bins with text density
	for _, elem := range elements {
		// Calculate which bins this element covers
		startBin := int((elem.Y - minY) / pa.binSize)
		endBin := int((elem.Top() - minY) / pa.binSize)

		// Ensure bins are within range
		if startBin < 0 {
			startBin = 0
		}
		if endBin >= binCount {
			endBin = binCount - 1
		}

		// Add density to bins
		// Density is proportional to text width
		density := elem.Width
		for i := startBin; i <= endBin; i++ {
			bins[i] += density
		}
	}

	return NewProjectionProfile(bins, pa.binSize, minY, maxY)
}

// AnalyzeVertical creates a vertical projection profile from text elements.
//
// The vertical profile shows text density by X coordinate (horizontal distribution).
// High density = lots of text at that X position
// Low density = whitespace (potential column boundary)
//
// Returns an array where index represents X position, value represents text density.
func (pa *DefaultProjectionAnalyzer) AnalyzeVertical(elements []*extractor.TextElement) *ProjectionProfile {
	if len(elements) == 0 {
		return NewProjectionProfile([]float64{}, pa.binSize, 0, 0)
	}

	// Find X coordinate range
	minX := elements[0].X
	maxX := elements[0].Right()

	for _, elem := range elements {
		if elem.X < minX {
			minX = elem.X
		}
		if elem.Right() > maxX {
			maxX = elem.Right()
		}
	}

	// Create bins
	binCount := int(math.Ceil((maxX-minX)/pa.binSize)) + 1
	bins := make([]float64, binCount)

	// Fill bins with text density
	for _, elem := range elements {
		// Calculate which bins this element covers
		startBin := int((elem.X - minX) / pa.binSize)
		endBin := int((elem.Right() - minX) / pa.binSize)

		// Ensure bins are within range
		if startBin < 0 {
			startBin = 0
		}
		if endBin >= binCount {
			endBin = binCount - 1
		}

		// Add density to bins
		// Density is proportional to text height
		density := elem.Height
		for i := startBin; i <= endBin; i++ {
			bins[i] += density
		}
	}

	return NewProjectionProfile(bins, pa.binSize, minX, maxX)
}

// FindGaps finds whitespace gaps in a projection profile.
//
// Gaps are regions where text density is below the threshold,
// indicating potential row or column boundaries.
func (pa *DefaultProjectionAnalyzer) FindGaps(profile *ProjectionProfile) []Gap {
	if profile.BinCount() == 0 {
		return []Gap{}
	}

	var gaps []Gap
	var gapStart float64
	inGap := false

	for i, density := range profile.Bins {
		coord := profile.Min + float64(i)*profile.BinSize

		if density <= pa.threshold {
			// Low density - we're in a gap
			if !inGap {
				// Start of new gap
				gapStart = coord
				inGap = true
			}
		} else {
			// High density - not in gap
			if inGap {
				// End of gap
				gapEnd := coord
				gap := NewGap(gapStart, gapEnd)
				gaps = append(gaps, gap)
				inGap = false
			}
		}
	}

	// Close final gap if still open
	if inGap {
		gapEnd := profile.Max
		gap := NewGap(gapStart, gapEnd)
		gaps = append(gaps, gap)
	}

	return gaps
}

// FindSignificantGaps finds gaps that are wide enough to be meaningful.
//
// Returns gaps with width >= minWidth.
func (pa *DefaultProjectionAnalyzer) FindSignificantGaps(profile *ProjectionProfile, minWidth float64) []Gap {
	allGaps := pa.FindGaps(profile)

	var significant []Gap
	for _, gap := range allGaps {
		if gap.Width >= minWidth {
			significant = append(significant, gap)
		}
	}

	return significant
}

// FindPeaks finds peaks (high density regions) in a projection profile.
//
// Peaks indicate where text is concentrated, which can help identify
// rows or columns of text.
type Peak struct {
	Start    float64 // Start coordinate of peak
	End      float64 // End coordinate of peak
	MaxValue float64 // Maximum density in peak
	Center   float64 // Center coordinate of peak
}

// NewPeak creates a new Peak.
func NewPeak(start, end, maxValue float64) Peak {
	return Peak{
		Start:    start,
		End:      end,
		MaxValue: maxValue,
		Center:   (start + end) / 2,
	}
}

// String returns a string representation of the peak.
func (p Peak) String() string {
	return fmt.Sprintf("Peak{start=%.2f, end=%.2f, max=%.2f, center=%.2f}",
		p.Start, p.End, p.MaxValue, p.Center)
}

// FindPeaks finds peaks in a projection profile.
//
// Peaks are regions where density is above a threshold.
func (pa *DefaultProjectionAnalyzer) FindPeaks(profile *ProjectionProfile, threshold float64) []Peak {
	if profile.BinCount() == 0 {
		return []Peak{}
	}

	var peaks []Peak
	var peakStart float64
	var maxInPeak float64
	inPeak := false

	for i, density := range profile.Bins {
		coord := profile.Min + float64(i)*profile.BinSize

		if density > threshold {
			// High density - we're in a peak
			if !inPeak {
				// Start of new peak
				peakStart = coord
				maxInPeak = density
				inPeak = true
			} else {
				// Continue peak
				if density > maxInPeak {
					maxInPeak = density
				}
			}
		} else {
			// Low density - not in peak
			if inPeak {
				// End of peak
				peakEnd := coord
				peak := NewPeak(peakStart, peakEnd, maxInPeak)
				peaks = append(peaks, peak)
				inPeak = false
			}
		}
	}

	// Close final peak if still open
	if inPeak {
		peakEnd := profile.Max
		peak := NewPeak(peakStart, peakEnd, maxInPeak)
		peaks = append(peaks, peak)
	}

	return peaks
}
