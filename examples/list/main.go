package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create a new PDF document.
	c := creator.New()
	c.SetTitle("List Examples")
	c.SetAuthor("gxpdf")

	// Create a page.
	page, err := c.NewPage()
	if err != nil {
		log.Fatalf("Failed to create page: %v", err)
	}

	// Example 1: Simple bullet list.
	bulletList := creator.NewList()
	bulletList.SetBulletChar("•")
	bulletList.SetFont(creator.Helvetica, 12)
	bulletList.Add("First item")
	bulletList.Add("Second item")
	bulletList.Add("Third item with longer text that will demonstrate text wrapping capabilities")
	bulletList.Add("Fourth item")

	err = page.Draw(bulletList)
	if err != nil {
		log.Fatalf("Failed to draw bullet list: %v", err)
	}

	// Add some spacing.
	ctx := page.GetLayoutContext()
	ctx.MoveCursor(0, 20)

	// Example 2: Numbered list with Arabic numerals.
	numberedList := creator.NewNumberedList()
	numberedList.SetNumberFormat(creator.NumberFormatArabic)
	numberedList.SetFont(creator.HelveticaBold, 12)
	numberedList.Add("Step one: Prepare the environment")
	numberedList.Add("Step two: Configure settings")
	numberedList.Add("Step three: Run the application")

	err = page.Draw(numberedList)
	if err != nil {
		log.Fatalf("Failed to draw numbered list: %v", err)
	}

	ctx.MoveCursor(0, 20)

	// Example 3: Nested list.
	parentList := creator.NewList()
	parentList.SetBulletChar("■")
	parentList.SetFont(creator.Helvetica, 12)
	parentList.Add("Main topic 1")

	// Create sublist.
	subList1 := creator.NewList()
	subList1.SetBulletChar("◆")
	subList1.SetIndent(15)
	subList1.SetFont(creator.Helvetica, 11)
	subList1.Add("Sub-topic 1.1")
	subList1.Add("Sub-topic 1.2")

	parentList.AddSubList(subList1)
	parentList.Add("Main topic 2")

	// Create another sublist.
	subList2 := creator.NewList()
	subList2.SetBulletChar("◆")
	subList2.SetIndent(15)
	subList2.SetFont(creator.Helvetica, 11)
	subList2.Add("Sub-topic 2.1")
	subList2.Add("Sub-topic 2.2")
	subList2.Add("Sub-topic 2.3")

	parentList.AddSubList(subList2)
	parentList.Add("Main topic 3")

	err = page.Draw(parentList)
	if err != nil {
		log.Fatalf("Failed to draw nested list: %v", err)
	}

	ctx.MoveCursor(0, 20)

	// Example 4: Alphabetic list.
	alphaList := creator.NewNumberedList()
	alphaList.SetNumberFormat(creator.NumberFormatLowerAlpha)
	alphaList.SetFont(creator.TimesRoman, 12)
	alphaList.Add("Option alpha")
	alphaList.Add("Option beta")
	alphaList.Add("Option gamma")

	err = page.Draw(alphaList)
	if err != nil {
		log.Fatalf("Failed to draw alphabetic list: %v", err)
	}

	ctx.MoveCursor(0, 20)

	// Example 5: Roman numerals list.
	romanList := creator.NewNumberedList()
	romanList.SetNumberFormat(creator.NumberFormatUpperRoman)
	romanList.SetFont(creator.TimesBold, 12)
	romanList.Add("Part one")
	romanList.Add("Part two")
	romanList.Add("Part three")

	err = page.Draw(romanList)
	if err != nil {
		log.Fatalf("Failed to draw roman list: %v", err)
	}

	ctx.MoveCursor(0, 20)

	// Example 6: Custom styling.
	customList := creator.NewList()
	customList.SetBulletChar("→")
	customList.SetFont(creator.Courier, 11)
	customList.SetColor(creator.Color{R: 0.2, G: 0.4, B: 0.8}) // Blue
	customList.SetLineSpacing(1.5)
	customList.SetIndent(25)
	customList.SetMarkerIndent(15)
	customList.Add("Custom styled item")
	customList.Add("With blue color and custom spacing")
	customList.Add("And Courier font")

	err = page.Draw(customList)
	if err != nil {
		log.Fatalf("Failed to draw custom list: %v", err)
	}

	// Save the PDF.
	err = c.WriteToFile("list_example.pdf")
	if err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Println("PDF created successfully: list_example.pdf")
}
