// Package main provides the gxpdf command-line interface.
//
// gxpdf is a powerful PDF processing tool that provides table extraction,
// text extraction, PDF manipulation, and more.
//
// Usage:
//
//	gxpdf [command] [flags]
//
// Available Commands:
//
//	tables      Extract tables from PDF (100% accuracy on bank statements)
//	text        Extract text from PDF
//	info        Display PDF metadata and information
//	merge       Merge multiple PDF files
//	split       Split PDF into separate files
//	encrypt     Encrypt PDF with password
//	decrypt     Decrypt password-protected PDF
//	version     Print version information
//
// Use "gxpdf [command] --help" for more information about a command.
package main

import (
	"os"

	"github.com/coregx/gxpdf/cmd/gxpdf/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		os.Exit(1)
	}
}
