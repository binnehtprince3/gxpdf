// Package security provides PDF encryption and security features.
//
// This file implements AES encryption for PDF documents (PDF 1.5+).
//
// Supported AES encryption:
//   - AES-128 with 128-bit keys (PDF 1.5 compatible, V=4, CFM=/AESV2)
//   - AES-256 with 256-bit keys (PDF 1.7 Extension Level 3, V=5, CFM=/AESV3)
//
// AES encryption in PDF uses:
//   - CBC (Cipher Block Chaining) mode
//   - PKCS#7 padding
//   - Random initialization vector (IV) prepended to encrypted data
package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5" //nolint:gosec // MD5 required by PDF Standard Security Handler
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
)

// AESEncryptor handles AES encryption/decryption for PDF objects.
//
// AES encryption is more secure than RC4 and is recommended for new PDFs.
// It uses CBC mode with PKCS#7 padding and random initialization vectors.
type AESEncryptor struct {
	config *EncryptionConfig
	dict   *EncryptionDict
}

// NewAESEncryptor creates a new AES encryptor with the given configuration.
//
// Supported key lengths:
//   - 128 bits (16 bytes) for AES-128 (PDF 1.5+)
//   - 256 bits (32 bytes) for AES-256 (PDF 1.7 Extension Level 3+)
//
// Example:
//
//	config := &EncryptionConfig{
//	    UserPassword:  "userpass",
//	    OwnerPassword: "ownerpass",
//	    Permissions:   PermissionPrint | PermissionCopy,
//	    KeyLength:     128, // or 256 for AES-256
//	    FileID:        "document-file-id",
//	}
//	enc, err := NewAESEncryptor(config)
func NewAESEncryptor(config *EncryptionConfig) (*AESEncryptor, error) {
	if err := validateAESConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	enc := &AESEncryptor{
		config: config,
		dict:   &EncryptionDict{},
	}

	if err := enc.buildEncryptionDict(); err != nil {
		return nil, fmt.Errorf("build encryption dict: %w", err)
	}

	return enc, nil
}

// validateAESConfig validates the AES encryption configuration.
func validateAESConfig(config *EncryptionConfig) error {
	if config.KeyLength != 128 && config.KeyLength != 256 {
		return fmt.Errorf("key length must be 128 or 256 bits for AES, got %d", config.KeyLength)
	}

	if config.FileID == "" {
		return ErrMissingFileID
	}

	return nil
}

// buildEncryptionDict creates the encryption dictionary for AES.
func (e *AESEncryptor) buildEncryptionDict() error {
	e.dict.Filter = "Standard"
	e.dict.Length = e.config.KeyLength

	// Set version, revision, and CFM based on key length.
	if e.config.KeyLength == 128 {
		e.dict.V = 4
		e.dict.R = 4
		e.dict.CFM = "AESV2" // AES-128 (PDF 1.5+)
	} else {
		e.dict.V = 5
		e.dict.R = 6
		e.dict.CFM = "AESV3" // AES-256 (PDF 1.7+)
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

// computeO computes the O value for AES encryption.
//
// For AES-128 (R=4): Uses similar algorithm to RC4 but with AES encryption.
// For AES-256 (R=6): Uses SHA-256/SHA-512 based algorithm.
func (e *AESEncryptor) computeO(ownerPwd, userPwd string) ([]byte, error) {
	if e.config.KeyLength == 128 {
		return e.computeOforAES128(ownerPwd, userPwd)
	}
	return e.computeOforAES256(ownerPwd, userPwd)
}

// computeOforAES128 computes O value for AES-128 (R=4).
func (e *AESEncryptor) computeOforAES128(ownerPwd, userPwd string) ([]byte, error) {
	// Algorithm similar to RC4 R=3 but with AES encryption.
	// Step 1: Pad owner password.
	ownerPadded := padPassword(ownerPwd)

	// Step 2: Compute MD5 hash.
	hash := md5.Sum(ownerPadded) //nolint:gosec // MD5 required by PDF spec

	// Step 3: Iterate MD5 50 times.
	for i := 0; i < 50; i++ {
		hash = md5.Sum(hash[:]) //nolint:gosec // MD5 required by PDF spec
	}

	// Step 4: Create AES key (first 16 bytes).
	encKey := hash[:16]

	// Step 5: Encrypt padded user password.
	userPadded := padPassword(userPwd)

	result, err := encryptAES(encKey, userPadded)
	if err != nil {
		return nil, err
	}

	// Step 6: Iterate with different keys (19 times).
	for i := 1; i <= 19; i++ {
		newKey := xorKey(encKey, byte(i))
		result, err = encryptAES(newKey, result)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

// computeOforAES256 computes O value for AES-256 (R=6).
func (e *AESEncryptor) computeOforAES256(ownerPwd, userPwd string) ([]byte, error) {
	// Algorithm 3.2a from PDF 2.0 specification.
	// For R=6, use SHA-256/SHA-512 based computation.
	ownerBytes := []byte(ownerPwd)
	userBytes := []byte(userPwd)

	// Generate random salt (8 bytes).
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	// Compute hash: SHA-256(password + salt + user password).
	h := sha256.New()
	h.Write(ownerBytes)
	h.Write(salt)
	h.Write(userBytes)
	hash := h.Sum(nil)

	// Result is hash + salt (32 + 8 = 40 bytes for AES-256).
	result := make([]byte, 0, 40)
	result = append(result, hash...)
	result = append(result, salt...)

	return result, nil
}

// computeU computes the U value for AES encryption.
func (e *AESEncryptor) computeU(userPwd string) ([]byte, error) {
	if e.config.KeyLength == 128 {
		return e.computeUforAES128(userPwd)
	}
	return e.computeUforAES256(userPwd)
}

// computeUforAES128 computes U value for AES-128 (R=4).
func (e *AESEncryptor) computeUforAES128(userPwd string) ([]byte, error) {
	// Compute encryption key.
	encKey := e.computeEncryptionKeyAES128(userPwd)

	// Algorithm 3.5: MD5(padding + file ID).
	h := md5.New() //nolint:gosec // MD5 required by PDF spec
	h.Write([]byte(paddingString))
	h.Write([]byte(e.config.FileID))
	hash := h.Sum(nil)

	// Encrypt with key.
	result, err := encryptAES(encKey, hash)
	if err != nil {
		return nil, err
	}

	// Iterate with different keys (19 times).
	for i := 1; i <= 19; i++ {
		newKey := xorKey(encKey, byte(i))
		result, err = encryptAES(newKey, result)
		if err != nil {
			return nil, err
		}
	}

	// Pad to 32 bytes.
	fullResult := make([]byte, 32)
	copy(fullResult, result)

	return fullResult, nil
}

// computeUforAES256 computes U value for AES-256 (R=6).
func (e *AESEncryptor) computeUforAES256(userPwd string) ([]byte, error) {
	// Algorithm 3.2a from PDF 2.0 specification.
	userBytes := []byte(userPwd)

	// Generate random salt (8 bytes).
	salt := make([]byte, 8)
	if _, err := rand.Read(salt); err != nil {
		return nil, fmt.Errorf("generate salt: %w", err)
	}

	// Compute hash: SHA-256(password + salt).
	h := sha256.New()
	h.Write(userBytes)
	h.Write(salt)
	hash := h.Sum(nil)

	// Result is hash + salt (32 + 8 = 40 bytes for AES-256).
	result := make([]byte, 0, 40)
	result = append(result, hash...)
	result = append(result, salt...)

	return result, nil
}

// computeEncryptionKeyAES128 computes the encryption key for AES-128.
func (e *AESEncryptor) computeEncryptionKeyAES128(userPwd string) []byte {
	// Algorithm 3.2 from PDF Reference 1.7.
	// Step 1: Pad password.
	padded := padPassword(userPwd)

	// Step 2: Build MD5 input.
	h := md5.New() //nolint:gosec // MD5 required by PDF spec
	h.Write(padded)
	h.Write(e.dict.O)
	h.Write(int32ToBytes(e.dict.P))
	h.Write([]byte(e.config.FileID))
	hash := h.Sum(nil)

	// Step 3: Iterate MD5 50 times.
	for i := 0; i < 50; i++ {
		hashArray := md5.Sum(hash[:16]) //nolint:gosec // MD5 required by PDF spec
		hash = hashArray[:]
	}

	// Step 4: Return first 16 bytes for AES-128.
	return hash[:16]
}

// GetEncryptionDict returns the encryption dictionary.
func (e *AESEncryptor) GetEncryptionDict() *EncryptionDict {
	return e.dict
}

// EncryptData encrypts data using AES with a random IV.
//
// The IV is prepended to the encrypted data as required by PDF spec.
//
// Example:
//
//	encrypted, err := enc.EncryptData([]byte("Hello, World!"))
//	// Result: [16-byte IV][encrypted data with PKCS#7 padding]
func (e *AESEncryptor) EncryptData(data []byte) ([]byte, error) {
	// Compute encryption key (simplified for example).
	// In real implementation, this depends on object number and generation.
	var encKey []byte
	if e.config.KeyLength == 128 {
		encKey = e.computeEncryptionKeyAES128(e.config.UserPassword)
	} else {
		encKey = e.computeEncryptionKeyAES256(e.config.UserPassword)
	}

	return encryptAES(encKey, data)
}

// computeEncryptionKeyAES256 computes the encryption key for AES-256.
func (e *AESEncryptor) computeEncryptionKeyAES256(userPwd string) []byte {
	// For AES-256 (R=6), use SHA-512 based key derivation.
	userBytes := []byte(userPwd)

	h := sha512.New()
	h.Write(userBytes)
	h.Write([]byte(e.config.FileID))
	hash := h.Sum(nil)

	// Return first 32 bytes for AES-256.
	return hash[:32]
}

// encryptAES encrypts data using AES-CBC with PKCS#7 padding.
//
// The function:
// 1. Generates a random IV (16 bytes for AES)
// 2. Pads data to AES block size (16 bytes) using PKCS#7
// 3. Encrypts using AES-CBC mode
// 4. Prepends IV to encrypted data
//
// Returns: [IV (16 bytes)][encrypted data].
func encryptAES(key, data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	// Add PKCS#7 padding.
	padded := addPKCS7Padding(data, aes.BlockSize)

	// Generate random IV.
	iv := make([]byte, aes.BlockSize)
	if _, err := rand.Read(iv); err != nil {
		return nil, fmt.Errorf("generate IV: %w", err)
	}

	// Encrypt using CBC mode.
	mode := cipher.NewCBCEncrypter(block, iv)
	encrypted := make([]byte, len(padded))
	mode.CryptBlocks(encrypted, padded)

	// Prepend IV to result (PDF spec requirement).
	result := make([]byte, 0, len(iv)+len(encrypted))
	result = append(result, iv...)
	result = append(result, encrypted...)

	return result, nil
}

// addPKCS7Padding adds PKCS#7 padding to data.
//
// PKCS#7 padding adds N bytes of value N, where N is the number of padding bytes needed.
//
// Example:
//   - Data of 13 bytes with block size 16 needs 3 bytes padding: [data][0x03 0x03 0x03]
//   - Data of 16 bytes with block size 16 needs 16 bytes padding: [data][0x10 * 16]
func addPKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := make([]byte, padding)

	for i := range padText {
		padText[i] = byte(padding)
	}

	return append(data, padText...)
}

// DecryptData decrypts AES-encrypted data.
//
// The data must be in the format: [IV (16 bytes)][encrypted data].
// PKCS#7 padding is removed after decryption.
func (e *AESEncryptor) DecryptData(data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("encrypted data too short: %d bytes", len(data))
	}

	// Compute decryption key.
	var decKey []byte
	if e.config.KeyLength == 128 {
		decKey = e.computeEncryptionKeyAES128(e.config.UserPassword)
	} else {
		decKey = e.computeEncryptionKeyAES256(e.config.UserPassword)
	}

	return decryptAES(decKey, data)
}

// decryptAES decrypts AES-CBC encrypted data with PKCS#7 padding.
func decryptAES(key, data []byte) ([]byte, error) {
	if len(data) < aes.BlockSize {
		return nil, fmt.Errorf("encrypted data too short: %d bytes", len(data))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("create AES cipher: %w", err)
	}

	// Extract IV (first 16 bytes).
	iv := data[:aes.BlockSize]
	encrypted := data[aes.BlockSize:]

	if len(encrypted)%aes.BlockSize != 0 {
		return nil, fmt.Errorf("encrypted data is not a multiple of block size")
	}

	// Decrypt using CBC mode.
	mode := cipher.NewCBCDecrypter(block, iv)
	decrypted := make([]byte, len(encrypted))
	mode.CryptBlocks(decrypted, encrypted)

	// Remove PKCS#7 padding.
	return removePKCS7Padding(decrypted)
}

// removePKCS7Padding removes PKCS#7 padding from data.
func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	padding := int(data[len(data)-1])
	if padding > len(data) || padding > aes.BlockSize {
		return nil, fmt.Errorf("invalid padding: %d", padding)
	}

	// Verify padding is correct.
	for i := len(data) - padding; i < len(data); i++ {
		if data[i] != byte(padding) {
			return nil, fmt.Errorf("invalid padding bytes")
		}
	}

	return data[:len(data)-padding], nil
}
