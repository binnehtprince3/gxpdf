package creator

import (
	"fmt"
	"strings"

	"github.com/coregx/gxpdf/internal/fonts"
)

// MarkerType defines the type of list marker.
type MarkerType int

const (
	// BulletMarker represents a bullet list (•, -, *, etc.).
	BulletMarker MarkerType = iota

	// NumberMarker represents a numbered list (1., 2., 3. or a), b), c) etc.).
	NumberMarker
)

// NumberFormat defines the numbering style.
type NumberFormat int

const (
	// NumberFormatArabic uses Arabic numerals (1. 2. 3.).
	NumberFormatArabic NumberFormat = iota

	// NumberFormatLowerAlpha uses lowercase letters (a. b. c.).
	NumberFormatLowerAlpha

	// NumberFormatUpperAlpha uses uppercase letters (A. B. C.).
	NumberFormatUpperAlpha

	// NumberFormatLowerRoman uses lowercase Roman numerals (i. ii. iii.).
	NumberFormatLowerRoman

	// NumberFormatUpperRoman uses uppercase Roman numerals (I. II. III.).
	NumberFormatUpperRoman
)

// ListItem represents a single item in a list.
type ListItem struct {
	text    string
	subList *List // nested list (optional)
}

// NewListItem creates a new list item with the given text.
func NewListItem(text string) ListItem {
	return ListItem{text: text}
}

// NewListItemWithSubList creates a new list item with a nested list.
func NewListItemWithSubList(text string, subList *List) ListItem {
	return ListItem{text: text, subList: subList}
}

// List represents a bullet or numbered list.
//
// Lists support nested sublists, custom markers, and automatic text wrapping.
//
// Example:
//
//	list := NewList()
//	list.SetBulletChar("•")
//	list.Add("First item")
//	list.Add("Second item with long text that will wrap")
//	page.Draw(list)
type List struct {
	items        []ListItem
	markerType   MarkerType
	numberFormat NumberFormat
	bulletChar   string // default "•"
	font         FontName
	fontSize     float64
	color        Color
	lineSpacing  float64 // multiplier (1.0 = normal)
	indent       float64 // indentation per level
	markerIndent float64 // space between marker and text
	startNumber  int     // starting number for numbered lists
}

// NewList creates a new bullet list with default settings.
//
// Default settings:
//   - Marker: Bullet ("•")
//   - Font: Helvetica, 12pt
//   - Color: Black
//   - Line spacing: 1.2 (120%)
//   - Indent: 20pt per level
//   - Marker indent: 10pt
//   - Start number: 1
func NewList() *List {
	return &List{
		items:        make([]ListItem, 0),
		markerType:   BulletMarker,
		numberFormat: NumberFormatArabic,
		bulletChar:   "•",
		font:         Helvetica,
		fontSize:     12,
		color:        Black,
		lineSpacing:  1.2,
		indent:       20,
		markerIndent: 10,
		startNumber:  1,
	}
}

// NewNumberedList creates a new numbered list with default settings.
//
// Default settings:
//   - Marker: Number (1. 2. 3.)
//   - Font: Helvetica, 12pt
//   - Color: Black
//   - Line spacing: 1.2 (120%)
//   - Indent: 20pt per level
//   - Marker indent: 10pt
//   - Start number: 1
func NewNumberedList() *List {
	list := NewList()
	list.markerType = NumberMarker
	return list
}

// SetMarkerType sets the marker type (bullet or number).
// Returns the list for method chaining.
func (l *List) SetMarkerType(t MarkerType) *List {
	l.markerType = t
	return l
}

// SetNumberFormat sets the numbering format for numbered lists.
// Returns the list for method chaining.
func (l *List) SetNumberFormat(f NumberFormat) *List {
	l.numberFormat = f
	return l
}

// SetBulletChar sets the bullet character for bullet lists.
// Returns the list for method chaining.
func (l *List) SetBulletChar(char string) *List {
	l.bulletChar = char
	return l
}

// SetFont sets the font and size for the list.
// Returns the list for method chaining.
func (l *List) SetFont(font FontName, size float64) *List {
	l.font = font
	l.fontSize = size
	return l
}

// SetColor sets the text color.
// Returns the list for method chaining.
func (l *List) SetColor(c Color) *List {
	l.color = c
	return l
}

// SetLineSpacing sets the line spacing multiplier.
// 1.0 = single spacing, 1.5 = 150% spacing, 2.0 = double spacing.
// Returns the list for method chaining.
func (l *List) SetLineSpacing(spacing float64) *List {
	l.lineSpacing = spacing
	return l
}

// SetIndent sets the indentation per nesting level.
// Returns the list for method chaining.
func (l *List) SetIndent(indent float64) *List {
	l.indent = indent
	return l
}

// SetMarkerIndent sets the space between marker and text.
// Returns the list for method chaining.
func (l *List) SetMarkerIndent(indent float64) *List {
	l.markerIndent = indent
	return l
}

// SetStartNumber sets the starting number for numbered lists.
// Returns the list for method chaining.
func (l *List) SetStartNumber(n int) *List {
	l.startNumber = n
	return l
}

// Add adds a text item to the list.
// Returns the list for method chaining.
func (l *List) Add(text string) *List {
	l.items = append(l.items, NewListItem(text))
	return l
}

// AddItem adds a ListItem to the list.
// Returns the list for method chaining.
func (l *List) AddItem(item ListItem) *List {
	l.items = append(l.items, item)
	return l
}

// AddSubList adds a nested sublist as the last item.
// Returns the list for method chaining.
func (l *List) AddSubList(subList *List) *List {
	if len(l.items) > 0 {
		lastIdx := len(l.items) - 1
		l.items[lastIdx].subList = subList
	}
	return l
}

// Height calculates the total height of the list when rendered.
func (l *List) Height(ctx *LayoutContext) float64 {
	return l.calculateHeight(ctx, 0)
}

// Draw renders the list on the page at the current cursor position.
func (l *List) Draw(ctx *LayoutContext, page *Page) error {
	return l.draw(ctx, page, 0)
}

// calculateHeight calculates the total height of the list at a given nesting level.
func (l *List) calculateHeight(ctx *LayoutContext, level int) float64 {
	lineHeight := l.calculateLineHeight()
	totalHeight := 0.0

	currentIndent := float64(level) * l.indent
	availableWidth := ctx.AvailableWidth() - currentIndent - l.markerIndent

	for _, item := range l.items {
		// Calculate height for item text.
		itemHeight := l.calculateItemHeight(item.text, availableWidth, lineHeight)
		totalHeight += itemHeight

		// Add height for sublist.
		if item.subList != nil {
			subHeight := item.subList.calculateHeight(ctx, level+1)
			totalHeight += subHeight
		}
	}

	return totalHeight
}

// draw renders the list at a given nesting level.
func (l *List) draw(ctx *LayoutContext, page *Page, level int) error {
	lineHeight := l.calculateLineHeight()
	currentIndent := float64(level) * l.indent
	availableWidth := ctx.AvailableWidth() - currentIndent - l.markerIndent

	for idx, item := range l.items {
		// Draw marker.
		marker := l.getMarker(idx)
		markerX := ctx.ContentLeft() + currentIndent
		markerY := ctx.CurrentPDFY() - l.fontSize

		err := page.AddTextColor(marker, markerX, markerY, l.font, l.fontSize, l.color)
		if err != nil {
			return fmt.Errorf("failed to draw list marker: %w", err)
		}

		// Draw item text.
		textX := markerX + l.markerIndent
		lines := l.wrapText(item.text, availableWidth)

		for lineIdx, line := range lines {
			textY := ctx.CurrentPDFY() - l.fontSize
			err := page.AddTextColor(line, textX, textY, l.font, l.fontSize, l.color)
			if err != nil {
				return fmt.Errorf("failed to draw list item text: %w", err)
			}

			// Move to next line only if there are more lines.
			if lineIdx < len(lines)-1 {
				ctx.CursorY += lineHeight
			}
		}

		// Move to next item.
		ctx.CursorY += lineHeight

		// Draw sublist.
		if item.subList != nil {
			err := item.subList.draw(ctx, page, level+1)
			if err != nil {
				return fmt.Errorf("failed to draw sublist: %w", err)
			}
		}
	}

	return nil
}

// getMarker returns the marker string for the item at the given index.
func (l *List) getMarker(index int) string {
	switch l.markerType {
	case BulletMarker:
		return l.bulletChar
	case NumberMarker:
		return l.formatNumber(l.startNumber + index)
	default:
		return l.bulletChar
	}
}

// formatNumber formats a number according to the number format.
func (l *List) formatNumber(n int) string {
	switch l.numberFormat {
	case NumberFormatArabic:
		return l.formatArabic(n)
	case NumberFormatLowerAlpha:
		return l.formatAlpha(n, false)
	case NumberFormatUpperAlpha:
		return l.formatAlpha(n, true)
	case NumberFormatLowerRoman:
		return toRoman(n, false) + "."
	case NumberFormatUpperRoman:
		return toRoman(n, true) + "."
	default:
		return l.formatArabic(n)
	}
}

// formatArabic formats a number as Arabic numerals.
func (l *List) formatArabic(n int) string {
	return fmt.Sprintf("%d.", n)
}

// formatAlpha formats a number as alphabetic characters (a-z or A-Z).
func (l *List) formatAlpha(n int, upper bool) string {
	if n < 1 || n > 26 {
		return l.formatArabic(n)
	}

	base := 'a'
	if upper {
		base = 'A'
	}
	return string(base+rune(n)-1) + "."
}

// calculateLineHeight returns the height of one line.
func (l *List) calculateLineHeight() float64 {
	return l.fontSize * l.lineSpacing
}

// calculateItemHeight calculates the height of a single item with text wrapping.
func (l *List) calculateItemHeight(
	text string,
	availableWidth float64,
	lineHeight float64,
) float64 {
	lines := l.wrapText(text, availableWidth)
	return float64(len(lines)) * lineHeight
}

// wrapText breaks the text into lines that fit within the given width.
func (l *List) wrapText(text string, availableWidth float64) []string {
	if text == "" {
		return []string{}
	}

	words := strings.Fields(text)
	if len(words) == 0 {
		return []string{}
	}

	spaceWidth := fonts.MeasureString(string(l.font), " ", l.fontSize)

	var lines []string
	var currentLine []string
	var currentWidth float64

	for _, word := range words {
		wordWidth := fonts.MeasureString(string(l.font), word, l.fontSize)

		// Check if adding this word exceeds available width.
		newWidth := currentWidth + wordWidth
		if len(currentLine) > 0 {
			newWidth += spaceWidth
		}

		if newWidth > availableWidth && len(currentLine) > 0 {
			// Start a new line.
			lines = append(lines, strings.Join(currentLine, " "))
			currentLine = []string{word}
			currentWidth = wordWidth
		} else {
			// Add to current line.
			currentLine = append(currentLine, word)
			if len(currentLine) > 1 {
				currentWidth += spaceWidth
			}
			currentWidth += wordWidth
		}
	}

	// Add the last line.
	if len(currentLine) > 0 {
		lines = append(lines, strings.Join(currentLine, " "))
	}

	return lines
}

// toRoman converts an integer to Roman numerals.
func toRoman(n int, upper bool) string {
	if n <= 0 || n > 3999 {
		return fmt.Sprintf("%d", n)
	}

	var (
		values = []int{1000, 900, 500, 400, 100, 90, 50, 40, 10, 9, 5, 4, 1}
		lower  = []string{
			"m", "cm", "d", "cd", "c", "xc", "l", "xl", "x", "ix", "v", "iv", "i",
		}
		upperSymbols = []string{
			"M", "CM", "D", "CD", "C", "XC", "L", "XL", "X", "IX", "V", "IV", "I",
		}
	)

	symbols := lower
	if upper {
		symbols = upperSymbols
	}

	var result strings.Builder
	for i, value := range values {
		for n >= value {
			result.WriteString(symbols[i])
			n -= value
		}
	}

	return result.String()
}
