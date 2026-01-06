package creator

import (
	"fmt"
	"strings"

	"github.com/coregx/gxpdf/internal/fonts"
)

// TextChunk represents a piece of text with specific styling.
//
// TextChunks are combined to create a StyledParagraph with multiple
// text styles within the same paragraph.
//
// Example:
//
//	chunk := TextChunk{
//	    Text:  "Hello World",
//	    Style: DefaultTextStyle(),
//	}
type TextChunk struct {
	// Text is the text content.
	Text string

	// Style is the styling to apply to this chunk.
	Style TextStyle
}

// StyledParagraph is a paragraph with multiple text styles.
//
// Unlike a simple Paragraph that has one style for all text,
// StyledParagraph allows mixing different fonts, sizes, and colors
// within the same paragraph while maintaining proper text wrapping
// and alignment.
//
// Example:
//
//	sp := NewStyledParagraph()
//	sp.Append("This is ")
//	sp.AppendStyled("bold", TextStyle{Font: HelveticaBold, Size: 12, Color: Black})
//	sp.Append(" and this is ")
//	sp.AppendStyled("red", TextStyle{Font: Helvetica, Size: 12, Color: Red})
//	sp.Append(" text.")
//	sp.SetAlignment(AlignJustify)
//	page.Draw(sp)
type StyledParagraph struct {
	chunks      []TextChunk
	alignment   Alignment
	lineSpacing float64
}

// styledWord represents a word with its style and measured width.
type styledWord struct {
	text  string
	style TextStyle
	width float64
}

// styledLine represents a line of styled words with their metrics.
type styledLine struct {
	words        []styledWord
	totalWidth   float64
	maxAscender  float64 // Maximum ascender in font units
	maxDescender float64 // Maximum descender in font units (negative)
}

// NewStyledParagraph creates a new styled paragraph.
//
// Default settings:
//   - Alignment: Left
//   - Line spacing: 1.2 (120%)
//   - No text chunks initially
func NewStyledParagraph() *StyledParagraph {
	return &StyledParagraph{
		chunks:      make([]TextChunk, 0, 4),
		alignment:   AlignLeft,
		lineSpacing: 1.2,
	}
}

// Append adds text using the default style.
// Returns the paragraph for method chaining.
func (sp *StyledParagraph) Append(text string) *StyledParagraph {
	return sp.AppendStyled(text, DefaultTextStyle())
}

// AppendStyled adds text with a specific style.
// Returns the paragraph for method chaining.
func (sp *StyledParagraph) AppendStyled(text string, style TextStyle) *StyledParagraph {
	sp.chunks = append(sp.chunks, TextChunk{
		Text:  text,
		Style: style,
	})
	return sp
}

// SetAlignment sets the text alignment.
// Returns the paragraph for method chaining.
func (sp *StyledParagraph) SetAlignment(a Alignment) *StyledParagraph {
	sp.alignment = a
	return sp
}

// SetLineSpacing sets the line spacing multiplier.
// 1.0 = single spacing, 1.5 = 150% spacing, 2.0 = double spacing.
// Returns the paragraph for method chaining.
func (sp *StyledParagraph) SetLineSpacing(spacing float64) *StyledParagraph {
	sp.lineSpacing = spacing
	return sp
}

// Height calculates the total height of the styled paragraph when rendered.
func (sp *StyledParagraph) Height(ctx *LayoutContext) float64 {
	if len(sp.chunks) == 0 {
		return 0
	}

	lines := sp.wrapText(ctx.AvailableWidth())
	if len(lines) == 0 {
		return 0
	}

	var totalHeight float64
	for _, line := range lines {
		lineHeight := sp.calculateLineHeight(line)
		totalHeight += lineHeight
	}

	return totalHeight
}

// Draw renders the styled paragraph on the page at the current cursor position.
func (sp *StyledParagraph) Draw(ctx *LayoutContext, page *Page) error {
	if len(sp.chunks) == 0 {
		return nil
	}

	lines := sp.wrapText(ctx.AvailableWidth())

	for _, line := range lines {
		if err := sp.drawLine(ctx, page, line); err != nil {
			return err
		}

		lineHeight := sp.calculateLineHeight(line)
		ctx.CursorY += lineHeight
	}

	return nil
}

// drawLine renders a single line of styled words.
func (sp *StyledParagraph) drawLine(ctx *LayoutContext, page *Page, line styledLine) error {
	x := sp.calculateLineX(ctx, line)
	y := ctx.CurrentPDFY()

	for _, word := range line.words {
		// Calculate baseline position for this word.
		// Use the word's own ascender to position baseline correctly.
		metrics := fonts.GetMetrics(string(word.style.Font))
		if metrics == nil {
			return fmt.Errorf("font metrics not found for font: %s", word.style.Font)
		}

		// Baseline Y = current Y - (ascender * size / 1000).
		ascenderPoints := float64(metrics.GetAscender()) * word.style.Size / 1000.0
		baselineY := y - ascenderPoints

		err := page.AddTextColor(word.text, x, baselineY, word.style.Font, word.style.Size, word.style.Color)
		if err != nil {
			return fmt.Errorf("failed to add text: %w", err)
		}

		// Advance X by word width.
		x += word.width
	}

	return nil
}

// calculateLineHeight calculates the height of a line.
// Uses the maximum ascender and descender across all words in the line.
func (sp *StyledParagraph) calculateLineHeight(line styledLine) float64 {
	// Line height = (maxAscender - maxDescender) * lineSpacing.
	// maxDescender is negative, so this becomes maxAscender + abs(maxDescender).
	lineHeightPoints := (line.maxAscender - line.maxDescender) * sp.lineSpacing
	return lineHeightPoints
}

// calculateLineX calculates the X position for a line based on alignment.
func (sp *StyledParagraph) calculateLineX(ctx *LayoutContext, line styledLine) float64 {
	availableWidth := ctx.AvailableWidth()

	switch sp.alignment {
	case AlignCenter:
		return ctx.ContentLeft() + (availableWidth-line.totalWidth)/2
	case AlignRight:
		return ctx.ContentRight() - line.totalWidth
	case AlignJustify, AlignLeft:
		return ctx.ContentLeft()
	default:
		return ctx.ContentLeft()
	}
}

// wrapText breaks the text into lines that fit within the given width.
func (sp *StyledParagraph) wrapText(availableWidth float64) []styledLine {
	if len(sp.chunks) == 0 {
		return []styledLine{}
	}

	// Split all chunks into styled words.
	words := sp.splitChunksIntoWords()
	if len(words) == 0 {
		return []styledLine{}
	}

	// Build lines by accumulating words.
	return sp.buildLines(words, availableWidth)
}

// splitChunksIntoWords splits all text chunks into individual styled words.
func (sp *StyledParagraph) splitChunksIntoWords() []styledWord {
	var words []styledWord

	for _, chunk := range sp.chunks {
		if chunk.Text == "" {
			continue
		}

		// Split chunk into words.
		chunkWords := strings.Fields(chunk.Text)

		// Create styled words.
		for _, wordText := range chunkWords {
			width := fonts.MeasureString(string(chunk.Style.Font), wordText, chunk.Style.Size)

			// Add space width if not the first word.
			if len(words) > 0 {
				spaceWidth := fonts.MeasureString(string(chunk.Style.Font), " ", chunk.Style.Size)
				width += spaceWidth
				wordText = " " + wordText
			}

			words = append(words, styledWord{
				text:  wordText,
				style: chunk.Style,
				width: width,
			})
		}
	}

	return words
}

// buildLines groups styled words into lines that fit the available width.
func (sp *StyledParagraph) buildLines(words []styledWord, availableWidth float64) []styledLine {
	var lines []styledLine
	var currentLine styledLine

	for i, word := range words {
		// Check if adding this word exceeds available width.
		newWidth := currentLine.totalWidth + word.width

		if newWidth > availableWidth && len(currentLine.words) > 0 {
			// Finalize current line and start a new one.
			lines = append(lines, currentLine)

			// Start new line with this word.
			// If it's not the first word overall, remove leading space.
			wordText := word.text
			if i > 0 && strings.HasPrefix(wordText, " ") {
				wordText = strings.TrimPrefix(wordText, " ")
				// Recalculate width without leading space.
				word.width = fonts.MeasureString(string(word.style.Font), wordText, word.style.Size)
				word.text = wordText
			}

			currentLine = sp.createLine(word)
		} else {
			// Add to current line.
			sp.addWordToLine(&currentLine, word)
		}
	}

	// Add the last line.
	if len(currentLine.words) > 0 {
		lines = append(lines, currentLine)
	}

	return lines
}

// createLine creates a new line with the given word.
func (sp *StyledParagraph) createLine(word styledWord) styledLine {
	metrics := fonts.GetMetrics(string(word.style.Font))
	if metrics == nil {
		// Fallback to default metrics if not found.
		return styledLine{
			words:        []styledWord{word},
			totalWidth:   word.width,
			maxAscender:  word.style.Size * 0.75,  // Approximate.
			maxDescender: -word.style.Size * 0.25, // Approximate.
		}
	}

	ascender := float64(metrics.GetAscender()) * word.style.Size / 1000.0
	descender := float64(metrics.GetDescender()) * word.style.Size / 1000.0

	return styledLine{
		words:        []styledWord{word},
		totalWidth:   word.width,
		maxAscender:  ascender,
		maxDescender: descender,
	}
}

// addWordToLine adds a word to an existing line.
func (sp *StyledParagraph) addWordToLine(line *styledLine, word styledWord) {
	line.words = append(line.words, word)
	line.totalWidth += word.width

	// Update line metrics.
	metrics := fonts.GetMetrics(string(word.style.Font))
	if metrics != nil {
		ascender := float64(metrics.GetAscender()) * word.style.Size / 1000.0
		descender := float64(metrics.GetDescender()) * word.style.Size / 1000.0

		if ascender > line.maxAscender {
			line.maxAscender = ascender
		}
		if descender < line.maxDescender {
			line.maxDescender = descender
		}
	}
}
