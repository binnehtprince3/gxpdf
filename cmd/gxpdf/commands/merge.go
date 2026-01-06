package commands

import (
	"fmt"

	"github.com/coregx/gxpdf/creator"
	"github.com/spf13/cobra"
)

var mergeOutput string

var mergeCmd = &cobra.Command{
	Use:   "merge FILE1 FILE2 [FILE...] -o OUTPUT",
	Short: "Merge multiple PDF files into one",
	Long: `Merge multiple PDF files into a single PDF document.

The pages from each input file are appended in order.

Examples:
  gxpdf merge doc1.pdf doc2.pdf -o combined.pdf
  gxpdf merge chapter1.pdf chapter2.pdf chapter3.pdf -o book.pdf`,
	Args: cobra.MinimumNArgs(2),
	RunE: runMerge,
}

func init() {
	mergeCmd.Flags().StringVarP(&mergeOutput, "output", "o", "", "Output file (required)")
	_ = mergeCmd.MarkFlagRequired("output")
}

func runMerge(_ *cobra.Command, args []string) error {
	printVerbosef("Merging %d files into %s", len(args), mergeOutput)

	merger := creator.NewMerger()
	defer func() { _ = merger.Close() }()

	for _, file := range args {
		printVerbosef("Adding: %s", file)
		if err := merger.AddAllPages(file); err != nil {
			return fmt.Errorf("failed to add %s: %w", file, err)
		}
	}

	if err := merger.Write(mergeOutput); err != nil {
		return fmt.Errorf("failed to write merged PDF: %w", err)
	}

	fmt.Printf("Merged %d files into %s\n", len(args), mergeOutput)
	return nil
}
