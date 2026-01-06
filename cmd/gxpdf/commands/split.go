package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/coregx/gxpdf/creator"
	"github.com/spf13/cobra"
)

var splitOutput string

var splitCmd = &cobra.Command{
	Use:   "split FILE PAGES -o OUTPUT",
	Short: "Split PDF by extracting specific pages",
	Long: `Extract specific pages from a PDF file into a new PDF.

Page ranges can be specified as:
  - Single pages: 1, 5, 10
  - Ranges: 1-5, 10-20
  - Combined: 1-3,5,7-9

Examples:
  gxpdf split document.pdf 1-10 -o first_10.pdf
  gxpdf split report.pdf 1,3,5 -o selected.pdf
  gxpdf split book.pdf 1-5,10-15 -o chapters.pdf`,
	Args: cobra.ExactArgs(2),
	RunE: runSplit,
}

func init() {
	splitCmd.Flags().StringVarP(&splitOutput, "output", "o", "", "Output file (required)")
	_ = splitCmd.MarkFlagRequired("output")
}

func runSplit(_ *cobra.Command, args []string) error {
	filePath := args[0]
	pageSpec := args[1]

	pages, err := parsePageSpec(pageSpec)
	if err != nil {
		return fmt.Errorf("invalid page specification: %w", err)
	}

	printVerbosef("Extracting pages %v from %s", pages, filePath)

	// Use merger to extract specific pages.
	merger := creator.NewMerger()
	defer func() { _ = merger.Close() }()

	if err := merger.AddPages(filePath, pages...); err != nil {
		return fmt.Errorf("failed to extract pages: %w", err)
	}

	if err := merger.Write(splitOutput); err != nil {
		return fmt.Errorf("failed to write output: %w", err)
	}

	fmt.Printf("Extracted %d page(s) to %s\n", len(pages), splitOutput)
	return nil
}

// parsePageSpec parses a page specification like "1-5,7,10-12".
func parsePageSpec(spec string) ([]int, error) {
	var pages []int
	parts := strings.Split(spec, ",")

	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		parsed, err := parsePagePart(part)
		if err != nil {
			return nil, err
		}
		pages = append(pages, parsed...)
	}

	if len(pages) == 0 {
		return nil, fmt.Errorf("no pages specified")
	}

	return pages, nil
}

// parsePagePart parses a single part of page spec (either "5" or "1-5").
func parsePagePart(part string) ([]int, error) {
	if !strings.Contains(part, "-") {
		page, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid page number: %s", part)
		}
		return []int{page}, nil
	}
	return parsePageRange(part)
}

// parsePageRange parses a page range like "1-5".
func parsePageRange(part string) ([]int, error) {
	rangeParts := strings.Split(part, "-")
	if len(rangeParts) != 2 {
		return nil, fmt.Errorf("invalid range: %s", part)
	}

	start, err := strconv.Atoi(strings.TrimSpace(rangeParts[0]))
	if err != nil {
		return nil, fmt.Errorf("invalid page number: %s", rangeParts[0])
	}

	end, err := strconv.Atoi(strings.TrimSpace(rangeParts[1]))
	if err != nil {
		return nil, fmt.Errorf("invalid page number: %s", rangeParts[1])
	}

	if start > end {
		return nil, fmt.Errorf("invalid range: start > end in %s", part)
	}

	pages := make([]int, 0, end-start+1)
	for i := start; i <= end; i++ {
		pages = append(pages, i)
	}
	return pages, nil
}
