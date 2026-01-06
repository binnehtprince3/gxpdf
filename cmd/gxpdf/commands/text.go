package commands

import (
	"fmt"
	"os"

	"github.com/coregx/gxpdf"
	"github.com/spf13/cobra"
)

var (
	textPage   int
	textOutput string
)

var textCmd = &cobra.Command{
	Use:   "text FILE",
	Short: "Extract text from PDF",
	Long: `Extract text content from PDF files.

Extracts all text while preserving reading order. Supports extracting
from specific pages or the entire document.

Examples:
  gxpdf text document.pdf
  gxpdf text report.pdf --page 1
  gxpdf text book.pdf -o extracted.txt`,
	Args: cobra.ExactArgs(1),
	RunE: runText,
}

func init() {
	textCmd.Flags().IntVarP(&textPage, "page", "p", 0, "Extract from specific page (0 = all)")
	textCmd.Flags().StringVarP(&textOutput, "output", "o", "", "Output file (default: stdout)")
}

func runText(_ *cobra.Command, args []string) error {
	filePath := args[0]

	doc, err := gxpdf.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open PDF: %w", err)
	}
	defer func() { _ = doc.Close() }()

	out, cleanup, err := getTextOutput()
	if err != nil {
		return err
	}
	if cleanup != nil {
		defer cleanup()
	}

	if textPage > 0 {
		return extractSinglePage(doc, out)
	}
	return extractAllPages(doc, out)
}

func getTextOutput() (*os.File, func(), error) {
	if textOutput != "" {
		f, err := os.Create(textOutput) //nolint:gosec // G304: User-specified output file
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create output file: %w", err)
		}
		return f, func() { _ = f.Close() }, nil
	}
	return os.Stdout, nil, nil
}

func extractSinglePage(doc *gxpdf.Document, out *os.File) error {
	if textPage > doc.PageCount() {
		return fmt.Errorf("page %d does not exist (document has %d pages)", textPage, doc.PageCount())
	}
	text, err := doc.ExtractTextFromPage(textPage)
	if err != nil {
		return fmt.Errorf("failed to extract text from page %d: %w", textPage, err)
	}
	_, _ = fmt.Fprintln(out, text)
	return nil
}

//nolint:unparam // Returns nil for consistency with extractSinglePage.
func extractAllPages(doc *gxpdf.Document, out *os.File) error {
	for pageNum := 1; pageNum <= doc.PageCount(); pageNum++ {
		text, err := doc.ExtractTextFromPage(pageNum)
		if err != nil {
			printVerbosef("Warning: failed to extract text from page %d: %v", pageNum, err)
			continue
		}
		if pageNum > 1 {
			_, _ = fmt.Fprintln(out)
			_, _ = fmt.Fprintf(out, "--- Page %d ---\n", pageNum)
		}
		_, _ = fmt.Fprintln(out, text)
	}
	return nil
}
