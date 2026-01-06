package security

import (
	"bytes"
	"crypto/aes"
	"testing"
)

//nolint:funlen // Test table requires more lines for clarity
func TestNewAESEncryptor(t *testing.T) {
	tests := []struct {
		name    string
		config  EncryptionConfig
		wantErr bool
		wantV   int
		wantR   int
		wantCFM string
	}{
		{
			name: "valid AES-128 encryptor",
			config: EncryptionConfig{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionPrint | PermissionCopy,
				KeyLength:     128,
				FileID:        "test-file-id",
			},
			wantErr: false,
			wantV:   4,
			wantR:   4,
			wantCFM: "AESV2",
		},
		{
			name: "valid AES-256 encryptor",
			config: EncryptionConfig{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionAll,
				KeyLength:     256,
				FileID:        "test-file-id",
			},
			wantErr: false,
			wantV:   5,
			wantR:   6,
			wantCFM: "AESV3",
		},
		{
			name: "invalid key length",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    192, // Invalid for AES in PDF
				FileID:       "test-file-id",
			},
			wantErr: true,
		},
		{
			name: "missing file ID",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    128,
				FileID:       "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewAESEncryptor(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewAESEncryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Verify encryptor was created.
			if enc == nil {
				t.Error("NewAESEncryptor() returned nil encryptor")
				return
			}

			dict := enc.GetEncryptionDict()
			if dict == nil {
				t.Error("GetEncryptionDict() returned nil")
				return
			}

			verifyAESEncryptionDict(t, dict, &tt.config, tt.wantV, tt.wantR, tt.wantCFM)
		})
	}
}

// verifyAESEncryptionDict verifies AES encryption dictionary properties.
func verifyAESEncryptionDict(t *testing.T, dict *EncryptionDict, config *EncryptionConfig, wantV, wantR int, wantCFM string) {
	t.Helper()

	if dict.Filter != "Standard" {
		t.Errorf("Filter = %v, want %v", dict.Filter, "Standard")
	}

	if dict.Length != config.KeyLength {
		t.Errorf("Length = %v, want %v", dict.Length, config.KeyLength)
	}

	if dict.V != wantV {
		t.Errorf("V = %v, want %v", dict.V, wantV)
	}

	if dict.R != wantR {
		t.Errorf("R = %v, want %v", dict.R, wantR)
	}

	if dict.CFM != wantCFM {
		t.Errorf("CFM = %v, want %v", dict.CFM, wantCFM)
	}

	verifyAESPasswordHashes(t, dict, config.KeyLength)
	verifyPermissions(t, dict, config.Permissions)
}

// verifyAESPasswordHashes verifies O and U values for AES.
func verifyAESPasswordHashes(t *testing.T, dict *EncryptionDict, keyLength int) {
	t.Helper()

	if len(dict.O) == 0 {
		t.Error("O is empty")
	}

	if len(dict.U) == 0 {
		t.Error("U is empty")
	}

	// For AES-256, O and U should be 40 bytes (32-byte hash + 8-byte salt).
	if keyLength == 256 {
		if len(dict.O) != 40 {
			t.Errorf("O length = %v, want 40 for AES-256", len(dict.O))
		}
		if len(dict.U) != 40 {
			t.Errorf("U length = %v, want 40 for AES-256", len(dict.U))
		}
	}
}

func TestAddPKCS7Padding(t *testing.T) {
	tests := []struct {
		name      string
		data      []byte
		blockSize int
		wantLen   int
		wantLast  byte
	}{
		{"13 bytes needs 3 padding", make([]byte, 13), 16, 16, 0x03},
		{"16 bytes needs 16 padding (full block)", make([]byte, 16), 16, 32, 0x10},
		{"empty data needs 16 padding", []byte{}, 16, 16, 0x10},
		{"1 byte needs 15 padding", []byte{0x01}, 16, 16, 0x0F},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := addPKCS7Padding(tt.data, tt.blockSize)

			if len(result) != tt.wantLen {
				t.Errorf("addPKCS7Padding() length = %v, want %v", len(result), tt.wantLen)
			}

			if result[len(result)-1] != tt.wantLast {
				t.Errorf("last byte = 0x%02X, want 0x%02X", result[len(result)-1], tt.wantLast)
			}

			// Verify all padding bytes and original data.
			verifyPKCS7Padding(t, result, tt.data, tt.wantLast)
		})
	}
}

// verifyPKCS7Padding verifies padding bytes and original data.
func verifyPKCS7Padding(t *testing.T, result, original []byte, wantLast byte) {
	t.Helper()

	paddingLen := int(wantLast)
	for i := len(result) - paddingLen; i < len(result); i++ {
		if result[i] != wantLast {
			t.Errorf("padding byte at %d = 0x%02X, want 0x%02X", i, result[i], wantLast)
		}
	}

	if !bytes.Equal(result[:len(original)], original) {
		t.Error("original data was modified")
	}
}

func TestRemovePKCS7Padding(t *testing.T) {
	tests := []struct {
		name    string
		data    []byte
		want    []byte
		wantErr bool
	}{
		{
			name:    "valid padding (3 bytes)",
			data:    []byte{0x01, 0x02, 0x03, 0x03, 0x03},
			want:    []byte{0x01, 0x02},
			wantErr: false,
		},
		{
			name:    "valid padding (16 bytes - full block)",
			data:    append(make([]byte, 16), 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10, 0x10),
			want:    make([]byte, 16),
			wantErr: false,
		},
		{
			name:    "invalid padding (incorrect bytes)",
			data:    []byte{0x01, 0x02, 0x03, 0x03, 0x04},
			wantErr: true,
		},
		{
			name:    "invalid padding (too large)",
			data:    []byte{0x01, 0x02, 0x20},
			wantErr: true,
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := removePKCS7Padding(tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("removePKCS7Padding() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			if !bytes.Equal(result, tt.want) {
				t.Errorf("removePKCS7Padding() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestEncryptDecryptAES(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		data []byte
	}{
		{
			name: "16-byte key (AES-128)",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
			data: []byte("Hello, World!"),
		},
		{
			name: "32-byte key (AES-256)",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20},
			data: []byte("This is a test message for AES encryption."),
		},
		{
			name: "empty data",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
			data: []byte{},
		},
		{
			name: "exact block size (16 bytes)",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
			data: []byte("1234567890123456"), // Exactly 16 bytes
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt.
			encrypted, err := encryptAES(tt.key, tt.data)
			if err != nil {
				t.Fatalf("encryptAES() error = %v", err)
			}

			// Verify encrypted is longer (IV + data + padding).
			if len(encrypted) < aes.BlockSize {
				t.Error("encrypted data is too short")
			}

			// Verify encrypted is different from original (unless empty).
			if len(tt.data) > 0 && bytes.Equal(encrypted[aes.BlockSize:], tt.data) {
				t.Error("encryptAES() did not change the data")
			}

			// Decrypt.
			decrypted, err := decryptAES(tt.key, encrypted)
			if err != nil {
				t.Fatalf("decryptAES() error = %v", err)
			}

			// Verify decrypted matches original.
			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("decryptAES() mismatch:\ngot:  %v\nwant: %v", decrypted, tt.data)
			}
		})
	}
}

func TestEncryptDecryptAESWithEncryptor(t *testing.T) {
	tests := []struct {
		name      string
		keyLength int
		data      []byte
	}{
		{
			name:      "AES-128 with short data",
			keyLength: 128,
			data:      []byte("Hello, World!"),
		},
		{
			name:      "AES-256 with long data",
			keyLength: 256,
			data:      []byte("This is a much longer test message for AES-256 encryption that should span multiple blocks."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &EncryptionConfig{
				UserPassword:  "testuser",
				OwnerPassword: "testowner",
				Permissions:   PermissionAll,
				KeyLength:     tt.keyLength,
				FileID:        "test-file-id-12345",
			}

			enc, err := NewAESEncryptor(config)
			if err != nil {
				t.Fatalf("NewAESEncryptor() error = %v", err)
			}

			// Encrypt.
			encrypted, err := enc.EncryptData(tt.data)
			if err != nil {
				t.Fatalf("EncryptData() error = %v", err)
			}

			// Decrypt.
			decrypted, err := enc.DecryptData(encrypted)
			if err != nil {
				t.Fatalf("DecryptData() error = %v", err)
			}

			// Verify.
			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("Encrypt/Decrypt mismatch:\ngot:  %v\nwant: %v", decrypted, tt.data)
			}
		})
	}
}

func TestDecryptAESErrors(t *testing.T) {
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{
			name:    "data too short (less than block size)",
			data:    []byte{0x01, 0x02, 0x03},
			wantErr: true,
		},
		{
			name:    "data not multiple of block size",
			data:    append(make([]byte, aes.BlockSize), 0x01, 0x02, 0x03), // 16 + 3 = 19 bytes
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := decryptAES(key, tt.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("decryptAES() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Benchmark tests.

func BenchmarkNewAESEncryptor128(b *testing.B) {
	config := &EncryptionConfig{
		UserPassword:  "user123",
		OwnerPassword: "owner123",
		Permissions:   PermissionAll,
		KeyLength:     128,
		FileID:        "test-file-id",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewAESEncryptor(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewAESEncryptor256(b *testing.B) {
	config := &EncryptionConfig{
		UserPassword:  "user123",
		OwnerPassword: "owner123",
		Permissions:   PermissionAll,
		KeyLength:     256,
		FileID:        "test-file-id",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewAESEncryptor(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncryptAES128(b *testing.B) {
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10}
	data := []byte("This is a test message for AES encryption benchmark.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encryptAES(key, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncryptAES256(b *testing.B) {
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1A, 0x1B, 0x1C, 0x1D, 0x1E, 0x1F, 0x20}
	data := []byte("This is a test message for AES encryption benchmark.")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := encryptAES(key, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}
