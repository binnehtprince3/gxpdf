package security

import (
	"bytes"
	"testing"
)

func TestEncryptionConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  EncryptionConfig
		wantErr bool
	}{
		{
			name: "valid 40-bit config",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    40,
				FileID:       "test-file-id",
			},
			wantErr: false,
		},
		{
			name: "valid 128-bit config",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    128,
				FileID:       "test-file-id",
			},
			wantErr: false,
		},
		{
			name: "invalid key length",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    256,
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
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNewRC4Encryptor(t *testing.T) {
	tests := []struct {
		name    string
		config  EncryptionConfig
		wantErr bool
	}{
		{
			name: "valid 40-bit encryptor",
			config: EncryptionConfig{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionPrint | PermissionCopy,
				KeyLength:     40,
				FileID:        "test-file-id",
			},
			wantErr: false,
		},
		{
			name: "valid 128-bit encryptor",
			config: EncryptionConfig{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionAll,
				KeyLength:     128,
				FileID:        "test-file-id",
			},
			wantErr: false,
		},
		{
			name: "invalid config",
			config: EncryptionConfig{
				UserPassword: "test",
				KeyLength:    256,
				FileID:       "test-file-id",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enc, err := NewRC4Encryptor(&tt.config)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRC4Encryptor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			verifyEncryptor(t, enc, &tt.config)
		})
	}
}

// verifyEncryptor verifies the encryptor properties.
func verifyEncryptor(t *testing.T, enc *RC4Encryptor, config *EncryptionConfig) {
	t.Helper()

	if enc == nil {
		t.Error("NewRC4Encryptor() returned nil encryptor")
		return
	}

	dict := enc.GetEncryptionDict()
	if dict == nil {
		t.Error("GetEncryptionDict() returned nil")
		return
	}

	verifyEncryptionDict(t, dict, config)
}

// verifyEncryptionDict verifies encryption dictionary properties.
func verifyEncryptionDict(t *testing.T, dict *EncryptionDict, config *EncryptionConfig) {
	t.Helper()

	if dict.Filter != "Standard" {
		t.Errorf("Filter = %v, want %v", dict.Filter, "Standard")
	}

	if dict.Length != config.KeyLength {
		t.Errorf("Length = %v, want %v", dict.Length, config.KeyLength)
	}

	verifyVersionAndRevision(t, dict, config.KeyLength)
	verifyPasswordHashes(t, dict)
	verifyPermissions(t, dict, config.Permissions)
}

// verifyVersionAndRevision verifies V and R values.
func verifyVersionAndRevision(t *testing.T, dict *EncryptionDict, keyLength int) {
	t.Helper()

	expectedV := 1
	expectedR := 2
	if keyLength == 128 {
		expectedV = 2
		expectedR = 3
	}

	if dict.V != expectedV {
		t.Errorf("V = %v, want %v", dict.V, expectedV)
	}

	if dict.R != expectedR {
		t.Errorf("R = %v, want %v", dict.R, expectedR)
	}
}

// verifyPasswordHashes verifies O and U are 32 bytes.
func verifyPasswordHashes(t *testing.T, dict *EncryptionDict) {
	t.Helper()

	if len(dict.O) != 32 {
		t.Errorf("O length = %v, want 32", len(dict.O))
	}

	if len(dict.U) != 32 {
		t.Errorf("U length = %v, want 32", len(dict.U))
	}
}

// verifyPermissions verifies permissions are set correctly.
func verifyPermissions(t *testing.T, dict *EncryptionDict, perms Permission) {
	t.Helper()

	if dict.P != int32(perms) {
		t.Errorf("P = %v, want %v", dict.P, perms)
	}
}

func TestPadPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantLen  int
	}{
		{
			name:     "empty password",
			password: "",
			wantLen:  32,
		},
		{
			name:     "short password",
			password: "test",
			wantLen:  32,
		},
		{
			name:     "32-byte password",
			password: "12345678901234567890123456789012",
			wantLen:  32,
		},
		{
			name:     "long password",
			password: "123456789012345678901234567890123456789012345678901234567890",
			wantLen:  32,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := padPassword(tt.password)
			if len(result) != tt.wantLen {
				t.Errorf("padPassword() length = %v, want %v", len(result), tt.wantLen)
			}

			// Verify password is at the start.
			if len(tt.password) > 0 && len(tt.password) < 32 {
				if !bytes.HasPrefix(result, []byte(tt.password)) {
					t.Error("padPassword() does not start with password")
				}
			}
		})
	}
}

func TestEncryptRC4(t *testing.T) {
	tests := []struct {
		name string
		key  []byte
		data []byte
	}{
		{
			name: "5-byte key",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05},
			data: []byte("Hello, World!"),
		},
		{
			name: "16-byte key",
			key:  []byte{0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0x0D, 0x0E, 0x0F, 0x10},
			data: []byte("This is a test message for RC4 encryption."),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Encrypt.
			encrypted := make([]byte, len(tt.data))
			if err := encryptRC4(tt.key, tt.data, encrypted); err != nil {
				t.Fatalf("encryptRC4() error = %v", err)
			}

			// Verify encrypted is different from original.
			if bytes.Equal(encrypted, tt.data) {
				t.Error("encryptRC4() did not change the data")
			}

			// Decrypt (RC4 is symmetric).
			decrypted := make([]byte, len(encrypted))
			if err := encryptRC4(tt.key, encrypted, decrypted); err != nil {
				t.Fatalf("encryptRC4() decrypt error = %v", err)
			}

			// Verify decrypted matches original.
			if !bytes.Equal(decrypted, tt.data) {
				t.Errorf("encryptRC4() decrypt mismatch: got %v, want %v", decrypted, tt.data)
			}
		})
	}
}

func TestXorKey(t *testing.T) {
	tests := []struct {
		name  string
		key   []byte
		value byte
		want  []byte
	}{
		{
			name:  "simple XOR",
			key:   []byte{0x01, 0x02, 0x03},
			value: 0xFF,
			want:  []byte{0xFE, 0xFD, 0xFC},
		},
		{
			name:  "XOR with zero",
			key:   []byte{0x01, 0x02, 0x03},
			value: 0x00,
			want:  []byte{0x01, 0x02, 0x03},
		},
		{
			name:  "XOR with 1",
			key:   []byte{0x00, 0x01, 0x02},
			value: 0x01,
			want:  []byte{0x01, 0x00, 0x03},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := xorKey(tt.key, tt.value)
			if !bytes.Equal(result, tt.want) {
				t.Errorf("xorKey() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestInt32ToBytes(t *testing.T) {
	tests := []struct {
		name  string
		value int32
		want  []byte
	}{
		{
			name:  "zero",
			value: 0,
			want:  []byte{0x00, 0x00, 0x00, 0x00},
		},
		{
			name:  "positive",
			value: 0x01020304,
			want:  []byte{0x04, 0x03, 0x02, 0x01}, // Little-endian
		},
		{
			name:  "negative",
			value: -1,
			want:  []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
		{
			name:  "permissions example",
			value: int32(PermissionPrint | PermissionCopy),
			want:  []byte{byte(PermissionPrint | PermissionCopy), 0x00, 0x00, 0x00},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := int32ToBytes(tt.value)
			if !bytes.Equal(result, tt.want) {
				t.Errorf("int32ToBytes() = %v, want %v", result, tt.want)
			}
		})
	}
}

// Benchmark tests.

func BenchmarkNewRC4Encryptor40(b *testing.B) {
	config := &EncryptionConfig{
		UserPassword:  "user123",
		OwnerPassword: "owner123",
		Permissions:   PermissionAll,
		KeyLength:     40,
		FileID:        "test-file-id",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewRC4Encryptor(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkNewRC4Encryptor128(b *testing.B) {
	config := &EncryptionConfig{
		UserPassword:  "user123",
		OwnerPassword: "owner123",
		Permissions:   PermissionAll,
		KeyLength:     128,
		FileID:        "test-file-id",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := NewRC4Encryptor(config)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkEncryptRC4(b *testing.B) {
	key := []byte{0x01, 0x02, 0x03, 0x04, 0x05}
	data := []byte("This is a test message for RC4 encryption benchmark.")
	result := make([]byte, len(data))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if err := encryptRC4(key, data, result); err != nil {
			b.Fatal(err)
		}
	}
}
