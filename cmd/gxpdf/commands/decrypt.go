package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	decryptPassword string
	decryptOutput   string
)

var decryptCmd = &cobra.Command{
	Use:   "decrypt FILE -p PASSWORD -o OUTPUT",
	Short: "Decrypt password-protected PDF",
	Long: `Decrypt a password-protected PDF file.

Removes encryption from a PDF, creating an unprotected copy.

Examples:
  gxpdf decrypt encrypted.pdf -p mypassword -o decrypted.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: runDecrypt,
}

func init() {
	decryptCmd.Flags().StringVarP(&decryptPassword, "password", "p", "", "Password to decrypt (required)")
	decryptCmd.Flags().StringVarP(&decryptOutput, "output", "o", "", "Output file (required)")
	_ = decryptCmd.MarkFlagRequired("password")
	_ = decryptCmd.MarkFlagRequired("output")
}

func runDecrypt(_ *cobra.Command, _ []string) error {
	// TODO: Implement decryption for password-protected PDFs
	return fmt.Errorf("decrypt command not yet implemented for CLI\n\nPassword-protected PDF reading will be available in a future release")
}
