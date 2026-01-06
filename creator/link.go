package creator

// LinkStyle defines the visual style for a link.
//
// This controls how clickable links appear in the PDF.
// By default, links are displayed in blue with underline (like in web browsers).
//
// Example:
//
//	// Default style (blue, underlined)
//	style := DefaultLinkStyle()
//
//	// Custom style (red, bold, no underline)
//	custom := LinkStyle{
//	    Font:      HelveticaBold,
//	    Size:      14,
//	    Color:     Red,
//	    Underline: false,
//	}
type LinkStyle struct {
	// Font to use for the link text.
	Font FontName

	// Size of the font in points.
	Size float64

	// Color of the link text (RGB, 0.0 to 1.0 range).
	Color Color

	// Underline indicates whether to draw a line under the link text.
	// Default: true (like web browsers).
	Underline bool
}

// DefaultLinkStyle returns the default link style.
//
// Default style:
//   - Font: Helvetica
//   - Size: 12pt
//   - Color: Blue (0, 0, 1)
//   - Underline: true
//
// This matches the typical appearance of links in web browsers.
//
// Example:
//
//	style := DefaultLinkStyle()
//	page.AddLinkStyled("Click here", "https://example.com", 100, 700, style)
func DefaultLinkStyle() LinkStyle {
	return LinkStyle{
		Font:      Helvetica,
		Size:      12,
		Color:     Blue,
		Underline: true,
	}
}
