package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new document.
	c := creator.New()

	// Add a page.
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	// Get layout context.
	ctx := page.GetLayoutContext()

	// Create a styled paragraph with multiple styles.
	sp := creator.NewStyledParagraph()

	// Add text with default style.
	sp.Append("This is a ")

	// Add bold text.
	sp.AppendStyled("bold", creator.TextStyle{
		Font:  creator.HelveticaBold,
		Size:  12,
		Color: creator.Black,
	})

	// Add more default text.
	sp.Append(" word, and this is ")

	// Add red text.
	sp.AppendStyled("red", creator.TextStyle{
		Font:  creator.Helvetica,
		Size:  12,
		Color: creator.Red,
	})

	// Add more default text.
	sp.Append(" text, and this is ")

	// Add blue italic text.
	sp.AppendStyled("blue italic", creator.TextStyle{
		Font:  creator.HelveticaOblique,
		Size:  12,
		Color: creator.Blue,
	})

	// Add final text.
	sp.Append(" text. This is a longer paragraph that demonstrates text wrapping across multiple styles when the available width is limited.")

	// Set alignment and line spacing.
	sp.SetAlignment(creator.AlignJustify)
	sp.SetLineSpacing(1.5)

	// Draw the paragraph.
	if err := page.Draw(sp); err != nil {
		log.Fatalf("Failed to draw paragraph: %v", err)
	}

	// Add some spacing.
	ctx.CursorY += 20

	// Create another styled paragraph with larger text.
	sp2 := creator.NewStyledParagraph()
	sp2.AppendStyled("Large Red Title", creator.TextStyle{
		Font:  creator.HelveticaBold,
		Size:  18,
		Color: creator.Red,
	})
	sp2.SetAlignment(creator.AlignCenter)

	if err := page.Draw(sp2); err != nil {
		log.Fatalf("Failed to draw title: %v", err)
	}

	ctx.CursorY += 15

	// Create a paragraph with mixed sizes.
	sp3 := creator.NewStyledParagraph()
	sp3.Append("This paragraph has ")
	sp3.AppendStyled("LARGE", creator.TextStyle{
		Font:  creator.HelveticaBold,
		Size:  16,
		Color: creator.Black,
	})
	sp3.Append(" and ")
	sp3.AppendStyled("small", creator.TextStyle{
		Font:  creator.Helvetica,
		Size:  8,
		Color: creator.DarkGray,
	})
	sp3.Append(" text mixed together.")

	if err := page.Draw(sp3); err != nil {
		log.Fatalf("Failed to draw mixed paragraph: %v", err)
	}

	// Save to file.
	if err := c.WriteToFile("styled_paragraph_example.pdf"); err != nil {
		log.Fatalf("Failed to save PDF: %v", err)
	}

	fmt.Println("PDF created successfully: styled_paragraph_example.pdf")
}
