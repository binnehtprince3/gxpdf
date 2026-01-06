// Example: Creating a document with chapters and table of contents
//
// This example demonstrates:
// - Creating chapters with automatic numbering
// - Nested sub-chapters (sections, subsections)
// - Automatic TOC generation with links
// - Customizing TOC appearance
//
// Run: go run examples/creator_chapter_toc.go
package main

import (
	"fmt"
	"log"

	"github.com/coregx/gxpdf/creator"
)

func main() {
	// Create new PDF document
	c := creator.New()
	c.SetTitle("My Technical Document")
	c.SetAuthor("John Doe")

	// Enable Table of Contents
	c.EnableTOC()

	// Customize TOC appearance
	toc := c.TOC()
	toc.SetTitle("Contents")
	style := toc.Style()
	style.TitleSize = 28
	style.EntrySize = 11
	toc.SetStyle(style)

	// Create Chapter 1: Introduction
	ch1 := creator.NewChapter("Introduction")
	ch1.Add(creator.NewParagraph("This document demonstrates the chapter and TOC feature of GxPDF."))
	ch1.Add(creator.NewParagraph("Chapters provide document structure with automatic numbering."))

	// Add sub-chapters (sections)
	sec1_1 := ch1.NewSubChapter("Background")
	sec1_1.Add(creator.NewParagraph("This section provides background information."))
	sec1_1.Add(creator.NewParagraph("Sections are automatically numbered as 1.1, 1.2, etc."))

	sec1_2 := ch1.NewSubChapter("Motivation")
	sec1_2.Add(creator.NewParagraph("This section explains the motivation for this work."))

	// Add sub-sub-chapter (subsection)
	subsec1_2_1 := sec1_2.NewSubChapter("Problem Statement")
	subsec1_2_1.Add(creator.NewParagraph("Subsections are numbered as 1.2.1, 1.2.2, etc."))

	// Add chapter to document
	if err := c.AddChapter(ch1); err != nil {
		log.Fatalf("Failed to add chapter 1: %v", err)
	}

	// Create Chapter 2: Methods
	ch2 := creator.NewChapter("Methods")
	ch2.Add(creator.NewParagraph("This chapter describes the methods used in this research."))

	sec2_1 := ch2.NewSubChapter("Experimental Setup")
	sec2_1.Add(creator.NewParagraph("Description of the experimental setup."))

	sec2_2 := ch2.NewSubChapter("Data Collection")
	sec2_2.Add(creator.NewParagraph("Description of data collection procedures."))

	if err := c.AddChapter(ch2); err != nil {
		log.Fatalf("Failed to add chapter 2: %v", err)
	}

	// Create Chapter 3: Results
	ch3 := creator.NewChapter("Results")
	ch3.Add(creator.NewParagraph("This chapter presents the results of the research."))

	sec3_1 := ch3.NewSubChapter("Quantitative Results")
	sec3_1.Add(creator.NewParagraph("Quantitative analysis results."))

	sec3_2 := ch3.NewSubChapter("Qualitative Results")
	sec3_2.Add(creator.NewParagraph("Qualitative analysis results."))

	if err := c.AddChapter(ch3); err != nil {
		log.Fatalf("Failed to add chapter 3: %v", err)
	}

	// Create Chapter 4: Conclusion
	ch4 := creator.NewChapter("Conclusion")
	ch4.Add(creator.NewParagraph("This chapter summarizes the findings and suggests future work."))

	sec4_1 := ch4.NewSubChapter("Summary")
	sec4_1.Add(creator.NewParagraph("Summary of key findings."))

	sec4_2 := ch4.NewSubChapter("Future Work")
	sec4_2.Add(creator.NewParagraph("Suggestions for future research directions."))

	if err := c.AddChapter(ch4); err != nil {
		log.Fatalf("Failed to add chapter 4: %v", err)
	}

	// Write to file
	// Note: TOC and chapters are automatically rendered
	outputPath := "examples/output/chapter_toc_example.pdf"
	if err := c.WriteToFile(outputPath); err != nil {
		log.Fatalf("Failed to write PDF: %v", err)
	}

	fmt.Printf("PDF created successfully: %s\n", outputPath)
	fmt.Println("\nDocument structure:")
	fmt.Println("- Table of Contents")
	fmt.Println("- 1. Introduction")
	fmt.Println("  - 1.1 Background")
	fmt.Println("  - 1.2 Motivation")
	fmt.Println("    - 1.2.1 Problem Statement")
	fmt.Println("- 2. Methods")
	fmt.Println("  - 2.1 Experimental Setup")
	fmt.Println("  - 2.2 Data Collection")
	fmt.Println("- 3. Results")
	fmt.Println("  - 3.1 Quantitative Results")
	fmt.Println("  - 3.2 Qualitative Results")
	fmt.Println("- 4. Conclusion")
	fmt.Println("  - 4.1 Summary")
	fmt.Println("  - 4.2 Future Work")
}
