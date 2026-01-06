// Package detector implements table detection algorithms.
package tabledetect

import (
	"math"
	"sort"

	"github.com/coregx/gxpdf/internal/extractor"
)

// DetectRowsHybrid finds horizontal alignment using HYBRID approach:
// Multiple passes for maximum accuracy (Gap + Overlap + Alignment).
//
// This is the UNIVERSAL solution for any table type.
// Pass 1: Gap detection (finds rows with significant whitespace)
// Pass 2: Overlap detection (finds rows with tight spacing, Tabula-inspired)
// Pass 3: Alignment detection (finds rows by element alignment)
//
// Returns a slice of Y coordinates representing row boundaries,
// sorted bottom to top (in PDF coordinates).
func (wa *DefaultWhitespaceAnalyzer) DetectRowsHybrid(elements []*extractor.TextElement) []float64 {
	if len(elements) == 0 {
		return []float64{}
	}

	var allRows []float64

	// PASS 1: Gap-based detection (existing logic)
	gapRows := wa.detectRowsByGaps(elements)
	allRows = append(allRows, gapRows...)

	// PASS 2: Overlap-based detection (Tabula-inspired)
	overlapRows := wa.detectRowsByOverlap(elements)
	allRows = append(allRows, overlapRows...)

	// PASS 3: Alignment-based detection (existing logic)
	alignmentRows := wa.findHorizontalAlignments(elements)
	allRows = append(allRows, alignmentRows...)

	// Merge and deduplicate
	allRows = wa.uniqueAndSort(allRows, wa.alignmentTolerance)

	// PASS 4: Multi-line cell merger (filter out false boundaries from split cells)
	allRows = wa.filterMultiLineCellBoundaries(elements, allRows)

	return allRows
}

// detectRowsByGaps extracts rows using gap detection (existing logic from DetectRows).
func (wa *DefaultWhitespaceAnalyzer) detectRowsByGaps(elements []*extractor.TextElement) []float64 {
	avgFontSize := wa.calculateAverageFontSize(elements)

	var minGapWidth float64
	if wa.isLatticeMode {
		// Lattice: 2.0x fontSize to ignore within-cell gaps (12-15pt)
		adaptiveMinGapWidth := avgFontSize * 2.0 // 20pt for 10pt font
		minGapWidth = math.Max(wa.minGapWidth, adaptiveMinGapWidth)
	} else {
		// Stream: use SMALLER threshold to detect tight row spacing (gaps < 10pt)
		// Alfa-Bank has very tight spacing between transactions (~5-8pt gaps)
		minGapWidth = avgFontSize * 0.5 // 5pt for 10pt font (UNIVERSAL for tight tables)
	}

	// Use projection profile to find gaps
	horizontalProfile := wa.projectionAnalyzer.AnalyzeHorizontal(elements)
	gaps := wa.projectionAnalyzer.FindSignificantGaps(horizontalProfile, minGapWidth)

	// Extract row boundaries from gaps
	var rows []float64

	// Add bottom edge of content area
	rows = append(rows, horizontalProfile.Min)

	// Add center of each gap as row boundary
	for _, gap := range gaps {
		rows = append(rows, gap.Center())
	}

	// Add top edge of content area
	rows = append(rows, horizontalProfile.Max)

	return rows
}

// detectRowsByOverlap extracts rows using vertical overlap detection (Tabula-inspired).
// Elements with overlap < threshold are considered separate rows.
func (wa *DefaultWhitespaceAnalyzer) detectRowsByOverlap(elements []*extractor.TextElement) []float64 {
	if len(elements) == 0 {
		return []float64{}
	}

	// Sort elements by Y coordinate (bottom to top)
	sorted := make([]*extractor.TextElement, len(elements))
	copy(sorted, elements)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Y < sorted[j].Y
	})

	// Calculate adaptive overlap threshold
	overlapThreshold := wa.calculateOptimalOverlapThreshold(elements)

	// Group elements into rows based on overlap
	var rowGroups [][]*extractor.TextElement
	currentRow := []*extractor.TextElement{sorted[0]}

	for i := 1; i < len(sorted); i++ {
		elem := sorted[i]

		// Check overlap with last element in current row
		last := currentRow[len(currentRow)-1]
		overlap := last.VerticalOverlapRatio(elem)

		if overlap < overlapThreshold {
			// Start new row
			rowGroups = append(rowGroups, currentRow)
			currentRow = []*extractor.TextElement{elem}
		} else {
			// Add to current row
			currentRow = append(currentRow, elem)
		}
	}
	rowGroups = append(rowGroups, currentRow)

	// Extract row boundaries (average Y position of each row)
	var rows []float64
	for _, group := range rowGroups {
		avgY := 0.0
		for _, elem := range group {
			avgY += elem.Y
		}
		avgY /= float64(len(group))
		rows = append(rows, avgY)
	}

	return rows
}

// calculateOptimalOverlapThreshold determines the best overlap threshold
// based on table structure analysis (adaptive, no magic numbers).
func (wa *DefaultWhitespaceAnalyzer) calculateOptimalOverlapThreshold(elements []*extractor.TextElement) float64 {
	avgFontSize := wa.calculateAverageFontSize(elements)

	// Calculate average vertical gap between elements
	avgGap := wa.calculateAverageVerticalGap(elements)

	// Heuristic: if gaps are small relative to fontSize → tight spacing → lower threshold
	ratio := avgGap / avgFontSize

	if ratio < 0.3 {
		// Very tight spacing (e.g., Alfa-Bank sparse rows with gaps < 3pt)
		return 0.05 // 5% overlap threshold (stricter than Tabula's 10%)
	} else if ratio < 0.8 {
		// Moderate spacing
		return 0.1 // 10% (Tabula default)
	} else {
		// Wide spacing (e.g., VTB with multi-line cells)
		return 0.15 // 15% (more lenient)
	}
}

// calculateAverageVerticalGap calculates the average vertical gap between consecutive elements.
func (wa *DefaultWhitespaceAnalyzer) calculateAverageVerticalGap(elements []*extractor.TextElement) float64 {
	if len(elements) < 2 {
		return 0.0
	}

	// Sort by Y coordinate
	sorted := make([]*extractor.TextElement, len(elements))
	copy(sorted, elements)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Y < sorted[j].Y
	})

	// Calculate gaps
	totalGap := 0.0
	count := 0
	for i := 1; i < len(sorted); i++ {
		gap := sorted[i].Y - sorted[i-1].Top()
		if gap > 0 { // Only positive gaps
			totalGap += gap
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return totalGap / float64(count)
}

// filterMultiLineCellBoundaries removes false row boundaries created by multi-line cells.
//
// Problem: detectRowsByOverlap() creates separate row boundaries for each line in multi-line cells.
// Example: "31.03.24 Перевод C163103240088312 через Систему быстрых платежей"
//
//	"на +79139919790. Без НДС."
//
// This becomes 2 rows instead of 1 → DUPLICATE!
//
// Solution: Detect and remove boundaries between lines that belong to the same cell:
// - Small Y-gap between rows (< 1.5x fontSize)
// - Second row has NO date column (or continues previous date)
// - Second row starts with lowercase/continuation pattern
func (wa *DefaultWhitespaceAnalyzer) filterMultiLineCellBoundaries(elements []*extractor.TextElement, rowBoundaries []float64) []float64 {
	if len(rowBoundaries) < 3 { // Need at least 3 boundaries to define 2 rows
		return rowBoundaries
	}

	avgFontSize := wa.calculateAverageFontSize(elements)
	maxMultiLineGap := avgFontSize * 1.5 // 15pt for 10pt font

	// Group elements by rows
	type rowInfo struct {
		boundary  float64
		elements  []*extractor.TextElement
		hasDate   bool
		hasAmount bool // Has monetary amount (e.g., "1 650,00" or "-350.50")
		minX      float64
		avgY      float64
		elemCount int // Number of text elements in this row
	}

	var rows []rowInfo
	for i := 0; i < len(rowBoundaries)-1; i++ {
		bottom := rowBoundaries[i]
		top := rowBoundaries[i+1]

		// Find elements in this row
		var rowElems []*extractor.TextElement
		for _, elem := range elements {
			elemCenter := elem.CenterY()
			if elemCenter >= bottom && elemCenter < top {
				rowElems = append(rowElems, elem)
			}
		}

		if len(rowElems) == 0 {
			continue
		}

		// Calculate row statistics
		info := rowInfo{
			boundary:  bottom,
			elements:  rowElems,
			hasDate:   wa.hasDateColumn(rowElems),
			hasAmount: wa.hasAmountColumn(rowElems),
			minX:      wa.findMinX(rowElems),
			elemCount: len(rowElems),
		}

		// Calculate average Y
		totalY := 0.0
		for _, elem := range rowElems {
			totalY += elem.Y
		}
		info.avgY = totalY / float64(len(rowElems))

		rows = append(rows, info)
	}
	// Add top boundary
	rows = append(rows, rowInfo{boundary: rowBoundaries[len(rowBoundaries)-1]})

	// Filter boundaries: keep if NOT a multi-line continuation
	var filtered []float64
	filtered = append(filtered, rows[0].boundary) // Always keep first boundary

	for i := 1; i < len(rows)-1; i++ {
		current := rows[i]
		previous := rows[i-1]

		// Calculate gap between rows
		gap := current.avgY - previous.avgY

		// Detect multi-line continuation - UNIVERSAL RULE:
		//
		// A row is a continuation of the previous row if:
		// 1. Small gap (< 1.5x fontSize) - rows are close together
		// 2. Current row has NO AMOUNT - not a transaction row
		// 3. Current row starts at similar or indented X - not a new column
		//
		// EXAMPLES:
		//
		// Alfa-Bank: Continuation row has NO date and NO amount
		//   Row 1: "31.03.24 Перевод C163103240088312 через СБП -350,00 1 234,56"
		//   Row 2: "на +79139919790. Без НДС." (no date, no amount)
		//
		// Sberbank: Continuation row HAS date but NO amount
		//   Row 1: "16.09.2025 13:58 089743 Прочие расходы -1 650,00 84 095,86"
		//   Row 2: "16.09.2025 WILDBERRIES SBERPAY MOSCOW RUS. Операция" (date, description, NO amount)
		//
		// KEY INSIGHT: Amount column is the BEST discriminator!
		// - Transaction row ALWAYS has amount (debit/credit)
		// - Continuation row NEVER has amount (just description)

		isMultiLineContinuation := gap < maxMultiLineGap &&
			!current.hasAmount && // NO AMOUNT = continuation row
			current.minX >= previous.minX-10 // Allow 10pt tolerance

		if !isMultiLineContinuation {
			// Keep this boundary (start of new transaction)
			filtered = append(filtered, current.boundary)
		}
		// else: SKIP this boundary (it's a multi-line continuation)
	}

	// Always keep last boundary
	filtered = append(filtered, rows[len(rows)-1].boundary)

	return filtered
}

// hasDateColumn checks if the row has a date-like element in the first column.
// Alfa-Bank date format: DD.MM.YY (e.g., "31.03.24")
func (wa *DefaultWhitespaceAnalyzer) hasDateColumn(elements []*extractor.TextElement) bool {
	if len(elements) == 0 {
		return false
	}

	// Find leftmost element (date column is first)
	leftmost := elements[0]
	for _, elem := range elements {
		if elem.X < leftmost.X {
			leftmost = elem
		}
	}

	// Check if it looks like a date: XX.XX.XX pattern
	text := leftmost.Text
	if len(text) < 8 {
		return false
	}

	// Simple pattern: digit.digit.digit
	hasDots := 0
	for _, ch := range text {
		if ch == '.' {
			hasDots++
		}
	}

	return hasDots >= 2 // DD.MM.YY has 2 dots
}

// hasAmountColumn checks if the row has a monetary amount element.
// Amount patterns:
//   - "1 650,00" (Russian format with space separator)
//   - "-350.50" (negative with decimal point)
//   - "1234.56" (decimal point)
//   - "1 234,56" (space + comma)
//   - "-1 650,00" (negative with space)
func (wa *DefaultWhitespaceAnalyzer) hasAmountColumn(elements []*extractor.TextElement) bool {
	for _, elem := range elements {
		text := elem.Text
		if wa.looksLikeAmount(text) {
			return true
		}
	}
	return false
}

// looksLikeAmount checks if text looks like a monetary amount.
func (wa *DefaultWhitespaceAnalyzer) looksLikeAmount(text string) bool {
	if len(text) < 1 {
		return false
	}

	// Remove leading sign if present (+ or -)
	// + indicates credit (positive), - indicates debit (negative)
	if text[0] == '-' || text[0] == '+' {
		text = text[1:]
	}

	// Must have digits
	hasDigit := false
	commaCount := 0
	dotCount := 0

	for _, ch := range text {
		if ch >= '0' && ch <= '9' {
			hasDigit = true
		} else if ch == ',' {
			commaCount++
		} else if ch == '.' {
			dotCount++
		} else if ch == ' ' || ch == '\u00a0' {
			// Space is allowed (thousand separator)
			// \u00a0 = NO-BREAK SPACE (used by some PDF generators)
			continue
		} else {
			// Other characters = not an amount
			return false
		}
	}

	// Amount patterns:
	// - "1 650,00" → 1 comma (decimal)
	// - "1234.56" → 1 dot (decimal)
	// - "1 234,56" → 1 comma (decimal)
	//
	// NOT amount:
	// - "16.09.2025" → 2 dots (date!)
	// - "089743" → 0 separators (ID)
	//
	// Rule: Must have EXACTLY 1 or 2 separators (max)
	// - 1 separator = decimal point/comma
	// - 2 separators = thousand + decimal (e.g., "1,234.56")
	totalSeparators := commaCount + dotCount
	return hasDigit && totalSeparators >= 1 && totalSeparators <= 2 && !(dotCount >= 2)
}

// findMinX returns the minimum X coordinate among elements.
func (wa *DefaultWhitespaceAnalyzer) findMinX(elements []*extractor.TextElement) float64 {
	if len(elements) == 0 {
		return 0
	}

	minX := elements[0].X
	for _, elem := range elements {
		if elem.X < minX {
			minX = elem.X
		}
	}
	return minX
}
