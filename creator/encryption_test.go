package creator

import (
	"errors"
	"testing"

	"github.com/coregx/gxpdf/internal/security"
)

//nolint:funlen // Test table requires more lines for clarity
func TestCreator_SetEncryption(t *testing.T) {
	tests := []struct {
		name    string
		opts    EncryptionOptions
		wantErr bool
	}{
		{
			name: "valid 40-bit encryption",
			opts: EncryptionOptions{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionPrint | PermissionCopy,
				KeyLength:     40,
			},
			wantErr: false,
		},
		{
			name: "valid 128-bit encryption",
			opts: EncryptionOptions{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionAll,
				KeyLength:     128,
			},
			wantErr: false,
		},
		{
			name: "default key length",
			opts: EncryptionOptions{
				UserPassword:  "user123",
				OwnerPassword: "owner123",
				Permissions:   PermissionPrint,
				// KeyLength not set, should default to 128
			},
			wantErr: false,
		},
		{
			name: "invalid key length (now maps to AES-256, so valid)",
			opts: EncryptionOptions{
				UserPassword: "user123",
				KeyLength:    256,
			},
			wantErr: false,
		},
		{
			name: "no passwords",
			opts: EncryptionOptions{
				Permissions: PermissionNone,
				KeyLength:   128,
			},
			wantErr: false, // Empty password is allowed (document has restrictions but no password)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.SetEncryption(tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetEncryption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			verifyCreatorEncryption(t, c, tt.opts)
		})
	}
}

// verifyCreatorEncryption verifies encryption options were stored correctly.
func verifyCreatorEncryption(t *testing.T, c *Creator, opts EncryptionOptions) {
	t.Helper()

	if c.encryptionOpts == nil {
		t.Error("SetEncryption() did not store options")
		return
	}

	// Verify algorithm was set (either from Algorithm or KeyLength).
	if c.encryptionOpts.Algorithm == 0 {
		t.Error("Algorithm should be set to a non-zero value")
	}

	// Verify permissions were stored.
	if c.encryptionOpts.Permissions != opts.Permissions {
		t.Errorf("Permissions = %v, want %v", c.encryptionOpts.Permissions, opts.Permissions)
	}
}

func TestPermissionConstants(t *testing.T) {
	// Verify that creator package exports permission constants correctly.
	tests := []struct {
		name         string
		creatorPerm  Permission
		securityPerm security.Permission
	}{
		{"Print", PermissionPrint, security.PermissionPrint},
		{"Modify", PermissionModify, security.PermissionModify},
		{"Copy", PermissionCopy, security.PermissionCopy},
		{"Annotate", PermissionAnnotate, security.PermissionAnnotate},
		{"FillForms", PermissionFillForms, security.PermissionFillForms},
		{"Extract", PermissionExtract, security.PermissionExtract},
		{"Assemble", PermissionAssemble, security.PermissionAssemble},
		{"PrintHighQuality", PermissionPrintHighQuality, security.PermissionPrintHighQuality},
		{"All", PermissionAll, security.PermissionAll},
		{"None", PermissionNone, security.PermissionNone},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.creatorPerm != tt.securityPerm {
				t.Errorf("Permission %s mismatch: creator=%v, security=%v",
					tt.name, tt.creatorPerm, tt.securityPerm)
			}
		})
	}
}

func TestEncryptionOptions_PermissionCombinations(t *testing.T) {
	tests := []struct {
		name  string
		perms Permission
		want  []Permission
	}{
		{
			name:  "single permission",
			perms: PermissionPrint,
			want:  []Permission{PermissionPrint},
		},
		{
			name:  "multiple permissions",
			perms: PermissionPrint | PermissionCopy,
			want:  []Permission{PermissionPrint, PermissionCopy},
		},
		{
			name:  "all permissions",
			perms: PermissionAll,
			want: []Permission{
				PermissionPrint,
				PermissionModify,
				PermissionCopy,
				PermissionAnnotate,
				PermissionFillForms,
				PermissionExtract,
				PermissionAssemble,
				PermissionPrintHighQuality,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.SetEncryption(EncryptionOptions{
				UserPassword: "test",
				Permissions:  tt.perms,
				KeyLength:    128,
			})

			if err != nil {
				t.Fatalf("SetEncryption() error = %v", err)
			}

			// Verify all expected permissions are present.
			for _, perm := range tt.want {
				if !c.encryptionOpts.Permissions.Has(perm) {
					t.Errorf("Missing permission: %v", perm)
				}
			}
		})
	}
}

func TestEncryptionOptions_PasswordHandling(t *testing.T) {
	tests := []struct {
		name          string
		userPassword  string
		ownerPassword string
	}{
		{
			name:          "both passwords set",
			userPassword:  "user123",
			ownerPassword: "owner123",
		},
		{
			name:          "only user password",
			userPassword:  "user123",
			ownerPassword: "",
		},
		{
			name:          "empty passwords",
			userPassword:  "",
			ownerPassword: "",
		},
		{
			name:          "long passwords",
			userPassword:  "this-is-a-very-long-user-password-for-testing",
			ownerPassword: "this-is-a-very-long-owner-password-for-testing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.SetEncryption(EncryptionOptions{
				UserPassword:  tt.userPassword,
				OwnerPassword: tt.ownerPassword,
				Permissions:   PermissionAll,
				KeyLength:     128,
			})

			if err != nil {
				t.Fatalf("SetEncryption() error = %v", err)
			}

			// Verify passwords were stored.
			if c.encryptionOpts.UserPassword != tt.userPassword {
				t.Errorf("UserPassword = %v, want %v",
					c.encryptionOpts.UserPassword, tt.userPassword)
			}

			if c.encryptionOpts.OwnerPassword != tt.ownerPassword {
				t.Errorf("OwnerPassword = %v, want %v",
					c.encryptionOpts.OwnerPassword, tt.ownerPassword)
			}
		})
	}
}

func TestEncryptionOptions_KeyLengthDefault(t *testing.T) {
	c := New()

	// Set encryption without specifying key length or algorithm.
	err := c.SetEncryption(EncryptionOptions{
		UserPassword: "test",
		Permissions:  PermissionPrint,
	})

	if err != nil {
		t.Fatalf("SetEncryption() error = %v", err)
	}

	// Verify default algorithm is AES-128.
	if c.encryptionOpts.Algorithm != EncryptionAES128 {
		t.Errorf("Default Algorithm = %v, want %v", c.encryptionOpts.Algorithm, EncryptionAES128)
	}
}

// Example usage tests.

func ExampleCreator_SetEncryption() {
	c := New()

	// Enable encryption with user and owner passwords.
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword:  "userpass",
		OwnerPassword: "ownerpass",
		Permissions:   PermissionPrint | PermissionCopy,
		KeyLength:     128,
	})

	// Document will be encrypted when written to file.
	// (WriteToFile not shown here as it requires more setup)
}

func ExampleEncryptionOptions_permissions() {
	c := New()

	// Allow only printing, no modifications.
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword: "secret",
		Permissions:  PermissionPrint,
		KeyLength:    128,
	})

	// Allow print and copy, but no modify.
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword: "secret",
		Permissions:  PermissionPrint | PermissionCopy,
		KeyLength:    128,
	})

	// Allow all operations.
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword: "secret",
		Permissions:  PermissionAll,
		KeyLength:    128,
	})
}

// AES Encryption Tests.

//nolint:funlen // Test table requires more lines for clarity
func TestSetEncryptionAES(t *testing.T) {
	tests := []struct {
		name    string
		opts    EncryptionOptions
		wantErr bool
	}{
		{
			name: "AES-128 encryption",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionPrint | PermissionCopy,
				Algorithm:     EncryptionAES128,
			},
			wantErr: false,
		},
		{
			name: "AES-256 encryption",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionAll,
				Algorithm:     EncryptionAES256,
			},
			wantErr: false,
		},
		{
			name: "RC4-128 encryption (legacy)",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionPrint,
				Algorithm:     EncryptionRC4_128,
			},
			wantErr: false,
		},
		{
			name: "RC4-40 encryption (legacy)",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionNone,
				Algorithm:     EncryptionRC4_40,
			},
			wantErr: false,
		},
		{
			name: "default algorithm (should use AES-128)",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionPrint | PermissionCopy,
			},
			wantErr: false,
		},
		{
			name: "backward compatibility - KeyLength 128",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionPrint,
				KeyLength:     128,
			},
			wantErr: false,
		},
		{
			name: "backward compatibility - KeyLength 256",
			opts: EncryptionOptions{
				UserPassword:  "userpass",
				OwnerPassword: "ownerpass",
				Permissions:   PermissionPrint,
				KeyLength:     256,
			},
			wantErr: false,
		},
		{
			name: "invalid algorithm",
			opts: EncryptionOptions{
				UserPassword: "userpass",
				Algorithm:    EncryptionAlgorithm(999),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			err := c.SetEncryption(tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("SetEncryption() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if c.encryptionOpts == nil {
					t.Error("encryptionOpts should be set")
					return
				}

				// Verify algorithm is set correctly.
				if c.encryptionOpts.Algorithm == 0 {
					t.Error("Algorithm should be set to a non-zero value")
				}
			}
		})
	}
}

func TestMapKeyLengthToAlgorithm(t *testing.T) {
	tests := []struct {
		name      string
		keyLength int
		want      EncryptionAlgorithm
	}{
		{
			name:      "40-bit maps to RC4_40",
			keyLength: 40,
			want:      EncryptionRC4_40,
		},
		{
			name:      "128-bit maps to RC4_128",
			keyLength: 128,
			want:      EncryptionRC4_128,
		},
		{
			name:      "256-bit maps to AES256",
			keyLength: 256,
			want:      EncryptionAES256,
		},
		{
			name:      "invalid key length defaults to AES128",
			keyLength: 192,
			want:      EncryptionAES128,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapKeyLengthToAlgorithm(tt.keyLength)
			if result != tt.want {
				t.Errorf("mapKeyLengthToAlgorithm() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestValidateAlgorithm(t *testing.T) {
	tests := []struct {
		name      string
		algorithm EncryptionAlgorithm
		wantErr   bool
	}{
		{
			name:      "RC4_40 is valid",
			algorithm: EncryptionRC4_40,
			wantErr:   false,
		},
		{
			name:      "RC4_128 is valid",
			algorithm: EncryptionRC4_128,
			wantErr:   false,
		},
		{
			name:      "AES128 is valid",
			algorithm: EncryptionAES128,
			wantErr:   false,
		},
		{
			name:      "AES256 is valid",
			algorithm: EncryptionAES256,
			wantErr:   false,
		},
		{
			name:      "invalid algorithm",
			algorithm: EncryptionAlgorithm(999),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateAlgorithm(tt.algorithm)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateAlgorithm() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify error type if expected.
			if tt.wantErr && !errors.Is(err, security.ErrUnsupportedVersion) {
				t.Errorf("validateAlgorithm() error = %v, want %v", err, security.ErrUnsupportedVersion)
			}
		})
	}
}

func TestEncryptionAlgorithmConstants(t *testing.T) {
	// Verify that algorithm constants have expected values.
	if EncryptionRC4_40 != 0 {
		t.Errorf("EncryptionRC4_40 = %d, want 0", EncryptionRC4_40)
	}

	if EncryptionRC4_128 != 1 {
		t.Errorf("EncryptionRC4_128 = %d, want 1", EncryptionRC4_128)
	}

	if EncryptionAES128 != 2 {
		t.Errorf("EncryptionAES128 = %d, want 2", EncryptionAES128)
	}

	if EncryptionAES256 != 3 {
		t.Errorf("EncryptionAES256 = %d, want 3", EncryptionAES256)
	}
}

func TestEncryptionOptionsDefaults(t *testing.T) {
	c := New()

	// Test with minimal options (should use defaults).
	err := c.SetEncryption(EncryptionOptions{
		UserPassword: "test",
	})

	if err != nil {
		t.Fatalf("SetEncryption() error = %v", err)
	}

	// Verify default algorithm is AES-128.
	if c.encryptionOpts.Algorithm != EncryptionAES128 {
		t.Errorf("Default algorithm = %v, want %v", c.encryptionOpts.Algorithm, EncryptionAES128)
	}
}

func TestEncryptionBackwardCompatibility(t *testing.T) {
	// Test that old API using KeyLength still works.
	c := New()

	err := c.SetEncryption(EncryptionOptions{
		UserPassword: "test",
		KeyLength:    128,
	})

	if err != nil {
		t.Fatalf("SetEncryption() error = %v", err)
	}

	// KeyLength 128 should map to RC4_128 for backward compatibility.
	if c.encryptionOpts.Algorithm != EncryptionRC4_128 {
		t.Errorf("KeyLength 128 mapped to %v, want %v", c.encryptionOpts.Algorithm, EncryptionRC4_128)
	}
}

func TestEncryptionOptionsWithAllPermissions(t *testing.T) {
	c := New()

	err := c.SetEncryption(EncryptionOptions{
		UserPassword:  "user",
		OwnerPassword: "owner",
		Permissions: PermissionPrint |
			PermissionModify |
			PermissionCopy |
			PermissionAnnotate |
			PermissionFillForms |
			PermissionExtract |
			PermissionAssemble |
			PermissionPrintHighQuality,
		Algorithm: EncryptionAES256,
	})

	if err != nil {
		t.Fatalf("SetEncryption() error = %v", err)
	}

	// Verify all permissions are set.
	expectedPerms := PermissionPrint |
		PermissionModify |
		PermissionCopy |
		PermissionAnnotate |
		PermissionFillForms |
		PermissionExtract |
		PermissionAssemble |
		PermissionPrintHighQuality

	if c.encryptionOpts.Permissions != expectedPerms {
		t.Errorf("Permissions = %v, want %v", c.encryptionOpts.Permissions, expectedPerms)
	}

	// Should be equal to PermissionAll.
	if c.encryptionOpts.Permissions != PermissionAll {
		t.Error("All permissions != PermissionAll")
	}
}

// Example usage with AES.

func ExampleCreator_SetEncryption_aes128() {
	c := New()

	// Enable AES-128 encryption (recommended).
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword:  "userpass",
		OwnerPassword: "ownerpass",
		Permissions:   PermissionPrint | PermissionCopy,
		Algorithm:     EncryptionAES128,
	})

	// Document will be encrypted when written to file.
}

func ExampleCreator_SetEncryption_aes256() {
	c := New()

	// Enable AES-256 encryption (most secure).
	_ = c.SetEncryption(EncryptionOptions{
		UserPassword:  "userpass",
		OwnerPassword: "ownerpass",
		Permissions:   PermissionAll,
		Algorithm:     EncryptionAES256,
	})

	// Document will be encrypted when written to file.
}
