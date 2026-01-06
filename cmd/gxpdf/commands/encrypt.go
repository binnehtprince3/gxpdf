package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	encryptPassword  string
	encryptOwner     string
	encryptAlgorithm string
	encryptOutput    string
)

var encryptCmd = &cobra.Command{
	Use:   "encrypt FILE -p PASSWORD -o OUTPUT",
	Short: "Encrypt PDF with password protection",
	Long: `Encrypt a PDF file with password protection.

Supports:
  - AES-256 encryption (default, most secure)
  - AES-128 encryption
  - RC4 encryption (legacy compatibility)

You can set both user password (to open) and owner password (to edit).

Examples:
  gxpdf encrypt secret.pdf -p mypassword -o encrypted.pdf
  gxpdf encrypt doc.pdf -p user123 --owner admin456 -o protected.pdf
  gxpdf encrypt legacy.pdf -p pass --algorithm rc4 -o encrypted.pdf`,
	Args: cobra.ExactArgs(1),
	RunE: runEncrypt,
}

func init() {
	encryptCmd.Flags().StringVarP(&encryptPassword, "password", "p", "", "User password (required)")
	encryptCmd.Flags().StringVar(&encryptOwner, "owner", "", "Owner password (optional)")
	encryptCmd.Flags().StringVar(&encryptAlgorithm, "algorithm", "aes128", "Encryption: aes256, aes128, rc4")
	encryptCmd.Flags().StringVarP(&encryptOutput, "output", "o", "", "Output file (required)")
	_ = encryptCmd.MarkFlagRequired("password")
	_ = encryptCmd.MarkFlagRequired("output")
}

func runEncrypt(_ *cobra.Command, _ []string) error {
	// TODO: Implement encryption for existing PDFs
	// Currently encryption is available through the Creator API for new documents
	return fmt.Errorf("encrypt command not yet implemented for CLI\n\nUse the Creator API for document encryption:\n  c := creator.New()\n  c.SetEncryption(creator.EncryptionOptions{...})")
}
