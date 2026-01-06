package creator

// TextStyle defines styling for a text chunk.
//
// TextStyle combines font, size, and color into a single style definition
// that can be applied to text chunks in a StyledParagraph.
//
// Example:
//
//	boldRed := TextStyle{
//	    Font:  HelveticaBold,
//	    Size:  14,
//	    Color: Red,
//	}
type TextStyle struct {
	// Font is the font to use (one of the Standard 14 fonts).
	Font FontName

	// Size is the font size in points.
	Size float64

	// Color is the text color (RGB, 0.0 to 1.0 range).
	Color Color
}

// DefaultTextStyle returns the default text style.
//
// Default style:
//   - Font: Helvetica
//   - Size: 12pt
//   - Color: Black
func DefaultTextStyle() TextStyle {
	return TextStyle{
		Font:  Helvetica,
		Size:  12,
		Color: Black,
	}
}
