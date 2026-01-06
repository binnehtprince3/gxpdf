package creator

import (
	"strings"

	"github.com/coregx/gxpdf/internal/fonts"
)

// Paragraph represents a block of text with automatic word wrapping.
//
// Paragraphs automatically wrap text based on the available width,
// handle alignment, and track their height for layout purposes.
//
// Example:
//
//	p := NewParagraph("This is a long text that will be wrapped automatically.")
//	p.SetFont(Helvetica, 12)
//	p.SetAlignment(AlignJustify)
//	page.Draw(p)
type Paragraph struct {
	text        string
	font        FontName
	fontSize    float64
	color       Color
	alignment   Alignment
	lineSpacing float64 // multiplier (1.0 = normal)
}

// NewParagraph creates a new paragraph with the given text.
//
// Default settings:
//   - Font: Helvetica, 12pt
//   - Color: Black
//   - Alignment: Left
//   - Line spacing: 1.2 (120%)
func NewParagraph(text string) *Paragraph {
	return &Paragraph{
		text:        text,
		font:        Helvetica,
		fontSize:    12,
		color:       Black,
		alignment:   AlignLeft,
		lineSpacing: 1.2,
	}
}

// SetFont sets the font and size for the paragraph.
// Returns the paragraph for method chaining.
func (p *Paragraph) SetFont(font FontName, size float64) *Paragraph {
	p.font = font
	p.fontSize = size
	return p
}

// SetColor sets the text color.
// Returns the paragraph for method chaining.
func (p *Paragraph) SetColor(c Color) *Paragraph {
	p.color = c
	return p
}

// SetAlignment sets the text alignment.
// Returns the paragraph for method chaining.
func (p *Paragraph) SetAlignment(a Alignment) *Paragraph {
	p.alignment = a
	return p
}

// SetLineSpacing sets the line spacing multiplier.
// 1.0 = single spacing, 1.5 = 150% spacing, 2.0 = double spacing.
// Returns the paragraph for method chaining.
func (p *Paragraph) SetLineSpacing(spacing float64) *Paragraph {
	p.lineSpacing = spacing
	return p
}

// Font returns the current font name.
func (p *Paragraph) Font() FontName {
	return p.font
}

// FontSize returns the current font size.
func (p *Paragraph) FontSize() float64 {
	return p.fontSize
}

// Color returns the current text color.
func (p *Paragraph) Color() Color {
	return p.color
}

// Alignment returns the current text alignment.
func (p *Paragraph) Alignment() Alignment {
	return p.alignment
}

// LineSpacing returns the current line spacing multiplier.
func (p *Paragraph) LineSpacing() float64 {
	return p.lineSpacing
}

// Text returns the paragraph text.
func (p *Paragraph) Text() string {
	return p.text
}

// SetText sets the paragraph text.
// Returns the paragraph for method chaining.
func (p *Paragraph) SetText(text string) *Paragraph {
	p.text = text
	return p
}

// Height calculates the total height of the paragraph when rendered.
func (p *Paragraph) Height(ctx *LayoutContext) float64 {
	lines := p.wrapText(ctx.AvailableWidth())
	lineHeight := p.calculateLineHeight()
	return float64(len(lines)) * lineHeight
}

// Draw renders the paragraph on the page at the current cursor position.
func (p *Paragraph) Draw(ctx *LayoutContext, page *Page) error {
	lines := p.wrapText(ctx.AvailableWidth())
	lineHeight := p.calculateLineHeight()

	for _, line := range lines {
		x := p.calculateLineX(ctx, line)
		y := ctx.CurrentPDFY() - p.fontSize // baseline position

		err := page.AddTextColor(line, x, y, p.font, p.fontSize, p.color)
		if err != nil {
			return err
		}

		ctx.CursorY += lineHeight
	}

	return nil
}

// calculateLineHeight returns the height of one line.
func (p *Paragraph) calculateLineHeight() float64 {
	return p.fontSize * p.lineSpacing
}

// calculateLineX calculates the X position for a line based on alignment.
func (p *Paragraph) calculateLineX(ctx *LayoutContext, line string) float64 {
	lineWidth := fonts.MeasureString(string(p.font), line, p.fontSize)
	availableWidth := ctx.AvailableWidth()

	switch p.alignment {
	case AlignCenter:
		return ctx.ContentLeft() + (availableWidth-lineWidth)/2
	case AlignRight:
		return ctx.ContentRight() - lineWidth
	case AlignJustify, AlignLeft:
		return ctx.ContentLeft()
	default:
		return ctx.ContentLeft()
	}
}

// wrapText breaks the text into lines that fit within the given width.
func (p *Paragraph) wrapText(availableWidth float64) []string {
	if p.text == "" {
		return []string{}
	}

	words := strings.Fields(p.text)
	if len(words) == 0 {
		return []string{}
	}

	spaceWidth := fonts.MeasureString(string(p.font), " ", p.fontSize)

	var lines []string
	var currentLine []string
	var currentWidth float64

	for _, word := range words {
		wordWidth := fonts.MeasureString(string(p.font), word, p.fontSize)

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

// WrapTextLines returns the lines after wrapping (for testing/debugging).
func (p *Paragraph) WrapTextLines(availableWidth float64) []string {
	return p.wrapText(availableWidth)
}
