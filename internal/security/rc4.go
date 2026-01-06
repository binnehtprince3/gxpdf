// Package security provides PDF encryption and security features.
//
// This package implements the PDF Standard Security Handler (Algorithm 2.A)
// as specified in PDF Reference 1.7, Section 3.5.
//
// Supported encryption algorithms:
//   - RC4 with 40-bit keys (PDF 1.1 compatible)
//   - RC4 with 128-bit keys (PDF 1.4 compatible)
//
// The package handles:
//   - User password (opens document with restrictions)
//   - Owner password (opens document with full access)
//   - Permission flags (print, copy, modify, etc.)
package security

import (
	"crypto/md5" //nolint:gosec // MD5 required by PDF Standard Security Handler
	"crypto/rc4" //nolint:gosec // RC4 required by PDF Standard Security Handler
	"fmt"
)

const (
	// PDF 1.7 Standard Security Handler padding string.
	// This is a fixed value from the PDF specification.
	paddingString = "\x28\xBF\x4E\x5E\x4E\x75\x8A\x41\x64\x00\x4E\x56" +
		"\xFF\xFA\x01\x08\x2E\x2E\x00\xB6\xD0\x68\x3E\x80\x2F\x0C" +
		"\xA9\xFE\x64\x53\x69\x7A"

	// filterStandard is the filter name for the Standard Security Handler.
	filterStandard = "Standard"
)

// EncryptionConfig holds the encryption configuration for a PDF document.
type EncryptionConfig struct {
	// UserPassword allows opening the document with restrictions.
	UserPassword string

	// OwnerPassword allows opening the document with full access.
	// If empty, defaults to UserPassword.
	OwnerPassword string

	// Permissions specifies what operations are allowed (print, copy, etc.).
	Permissions Permission

	// KeyLength specifies the encryption key length in bits (40 or 128).
	KeyLength int

	// FileID is the document's unique identifier from the trailer dictionary.
	FileID string
}

// Validate checks if the encryption config is valid.
func (c *EncryptionConfig) Validate() error {
	if c.KeyLength != 40 && c.KeyLength != 128 {
		return fmt.Errorf("key length must be 40 or 128 bits, got %d", c.KeyLength)
	}

	if c.FileID == "" {
		return ErrMissingFileID
	}

	return nil
}

// EncryptionDict represents the PDF encryption dictionary.
type EncryptionDict struct {
	// Filter is always "/Standard" for Standard Security Handler.
	Filter string

	// V is the algorithm version (1 for 40-bit RC4, 2 for 128-bit RC4, 4 for AES-128, 5 for AES-256).
	V int

	// R is the algorithm revision (2 for 40-bit RC4, 3 for 128-bit RC4, 4 for AES-128, 6 for AES-256).
	R int

	// Length is the key length in bits (40, 128, or 256).
	Length int

	// P is the permission flags (32-bit integer).
	P int32

	// O is the owner password hash (32 bytes for RC4, variable for AES).
	O []byte

	// U is the user password hash (32 bytes for RC4, variable for AES).
	U []byte

	// CFM is the crypt filter method (empty for RC4, "AESV2" for AES-128, "AESV3" for AES-256).
	CFM string
}

// RC4Encryptor handles RC4 encryption/decryption for PDF objects.
type RC4Encryptor struct {
	config *EncryptionConfig
	dict   *EncryptionDict
}

// NewRC4Encryptor creates a new RC4 encryptor with the given configuration.
func NewRC4Encryptor(config *EncryptionConfig) (*RC4Encryptor, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	enc := &RC4Encryptor{
		config: config,
		dict:   &EncryptionDict{},
	}

	if err := enc.buildEncryptionDict(); err != nil {
		return nil, fmt.Errorf("build encryption dict: %w", err)
	}

	return enc, nil
}

// buildEncryptionDict creates the encryption dictionary from the config.
func (e *RC4Encryptor) buildEncryptionDict() error {
	e.dict.Filter = filterStandard
	e.dict.Length = e.config.KeyLength

	// Set version and revision based on key length.
	if e.config.KeyLength == 40 {
		e.dict.V = 1
		e.dict.R = 2
	} else {
		e.dict.V = 2
		e.dict.R = 3
	}

	// Set permissions.
	e.dict.P = int32(e.config.Permissions)

	// Use owner password or default to user password.
	ownerPwd := e.config.OwnerPassword
	if ownerPwd == "" {
		ownerPwd = e.config.UserPassword
	}

	// Compute O value (owner password hash).
	o, err := e.computeO(ownerPwd, e.config.UserPassword)
	if err != nil {
		return fmt.Errorf("compute O: %w", err)
	}
	e.dict.O = o

	// Compute U value (user password hash).
	u, err := e.computeU(e.config.UserPassword)
	if err != nil {
		return fmt.Errorf("compute U: %w", err)
	}
	e.dict.U = u

	return nil
}

// computeO computes the O value (owner password hash).
//
// Algorithm 3.3 from PDF Reference 1.7:
// 1. Pad owner password to 32 bytes
// 2. Compute MD5 hash
// 3. For 128-bit: iterate MD5 50 times
// 4. Create RC4 key from hash
// 5. Encrypt padded user password
// 6. For 128-bit: iterate with different keys.
func (e *RC4Encryptor) computeO(ownerPwd, userPwd string) ([]byte, error) {
	// Step 1: Pad owner password.
	ownerPadded := padPassword(ownerPwd)

	// Step 2: Compute MD5 hash.
	hash := md5.Sum(ownerPadded) //nolint:gosec // MD5 required by PDF spec

	// Step 3: For 128-bit, iterate MD5 50 times.
	if e.config.KeyLength == 128 {
		for i := 0; i < 50; i++ {
			hash = md5.Sum(hash[:]) //nolint:gosec // MD5 required by PDF spec
		}
	}

	// Step 4: Create RC4 key (first n bytes of hash).
	keyLen := e.config.KeyLength / 8
	encKey := hash[:keyLen]

	// Step 5: Encrypt padded user password.
	userPadded := padPassword(userPwd)
	result := make([]byte, len(userPadded))

	if err := encryptRC4(encKey, userPadded, result); err != nil {
		return nil, err
	}

	// Step 6: For 128-bit, iterate with different keys.
	if e.config.KeyLength == 128 {
		for i := 1; i <= 19; i++ {
			newKey := xorKey(encKey, byte(i))
			if err := encryptRC4(newKey, result, result); err != nil {
				return nil, err
			}
		}
	}

	return result, nil
}

// computeU computes the U value (user password hash).
//
// Algorithm 3.4/3.5 from PDF Reference 1.7:
// 1. Compute encryption key from user password
// 2. For 40-bit: encrypt padding string with key
// 3. For 128-bit: encrypt MD5(padding + file ID) with key, iterate 20 times.
func (e *RC4Encryptor) computeU(userPwd string) ([]byte, error) {
	// Compute encryption key.
	encKey := e.computeEncryptionKey(userPwd)

	if e.config.KeyLength == 40 {
		// Algorithm 3.4: Encrypt padding string.
		result := make([]byte, 32)
		if err := encryptRC4(encKey, []byte(paddingString), result); err != nil {
			return nil, err
		}
		return result, nil
	}

	// Algorithm 3.5: For 128-bit.
	// Step 1: MD5(padding + file ID).
	h := md5.New() //nolint:gosec // MD5 required by PDF spec
	h.Write([]byte(paddingString))
	h.Write([]byte(e.config.FileID))
	hash := h.Sum(nil)

	// Step 2: Encrypt with key.
	result := make([]byte, len(hash))
	if err := encryptRC4(encKey, hash, result); err != nil {
		return nil, err
	}

	// Step 3: Iterate with different keys.
	for i := 1; i <= 19; i++ {
		newKey := xorKey(encKey, byte(i))
		if err := encryptRC4(newKey, result, result); err != nil {
			return nil, err
		}
	}

	// Step 4: Pad to 32 bytes.
	fullResult := make([]byte, 32)
	copy(fullResult, result)

	return fullResult, nil
}

// computeEncryptionKey computes the encryption key from user password.
//
// Algorithm 3.2 from PDF Reference 1.7:
// 1. Pad password to 32 bytes
// 2. MD5(password + O + P + file ID)
// 3. For 128-bit: iterate MD5 50 times
// 4. Return first n bytes as key.
func (e *RC4Encryptor) computeEncryptionKey(userPwd string) []byte {
	// Step 1: Pad password.
	padded := padPassword(userPwd)

	// Step 2: Build MD5 input.
	h := md5.New() //nolint:gosec // MD5 required by PDF spec
	h.Write(padded)
	h.Write(e.dict.O)
	h.Write(int32ToBytes(e.dict.P))
	h.Write([]byte(e.config.FileID))
	hash := h.Sum(nil)

	// Step 3: For 128-bit, iterate MD5 50 times.
	if e.config.KeyLength == 128 {
		for i := 0; i < 50; i++ {
			hashArray := md5.Sum(hash[:e.config.KeyLength/8]) //nolint:gosec // MD5 required by PDF spec
			hash = hashArray[:]
		}
	}

	// Step 4: Return first n bytes.
	return hash[:e.config.KeyLength/8]
}

// GetEncryptionDict returns the encryption dictionary.
func (e *RC4Encryptor) GetEncryptionDict() *EncryptionDict {
	return e.dict
}

// Helper functions.

// padPassword pads a password to 32 bytes using the PDF padding string.
func padPassword(password string) []byte {
	result := make([]byte, 32)
	pwdBytes := []byte(password)

	if len(pwdBytes) >= 32 {
		copy(result, pwdBytes[:32])
	} else {
		copy(result, pwdBytes)
		copy(result[len(pwdBytes):], paddingString)
	}

	return result
}

// encryptRC4 encrypts data with RC4 cipher.
func encryptRC4(key, data, result []byte) error {
	cipher, err := rc4.NewCipher(key) //nolint:gosec // RC4 required by PDF spec
	if err != nil {
		return fmt.Errorf("create RC4 cipher: %w", err)
	}

	cipher.XORKeyStream(result, data)
	return nil
}

// xorKey XORs each byte of key with value.
func xorKey(key []byte, value byte) []byte {
	result := make([]byte, len(key))
	for i, b := range key {
		result[i] = b ^ value
	}
	return result
}

// int32ToBytes converts int32 to little-endian byte array.
func int32ToBytes(n int32) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
	}
}
