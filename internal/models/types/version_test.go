package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewVersion tests version creation with valid inputs.
func TestNewVersion(t *testing.T) {
	tests := []struct {
		name  string
		major int
		minor int
		want  string
	}{
		{"PDF 1.0", 1, 0, "1.0"},
		{"PDF 1.4", 1, 4, "1.4"},
		{"PDF 1.7", 1, 7, "1.7"},
		{"PDF 2.0", 2, 0, "2.0"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersion(tt.major, tt.minor)
			require.NoError(t, err)
			assert.Equal(t, tt.major, v.Major())
			assert.Equal(t, tt.minor, v.Minor())
			assert.Equal(t, tt.want, v.String())
		})
	}
}

// TestNewVersion_Errors tests version creation with invalid inputs.
func TestNewVersion_Errors(t *testing.T) {
	tests := []struct {
		name      string
		major     int
		minor     int
		wantError error
	}{
		{"negative major", -1, 0, ErrInvalidVersion},
		{"negative minor", 1, -1, ErrInvalidVersion},
		{"negative both", -1, -1, ErrInvalidVersion},
		{"unsupported major version", 3, 0, ErrInvalidVersion},
		{"unsupported major version 4", 4, 0, ErrInvalidVersion},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersion(tt.major, tt.minor)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, tt.wantError), "expected error to be %v, got %v", tt.wantError, err)
			assert.Equal(t, Version{}, v)
		})
	}
}

// TestParseVersion tests parsing PDF version strings.
func TestParseVersion(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantMajor int
		wantMinor int
	}{
		{"simple 1.4", "1.4", 1, 4},
		{"simple 1.7", "1.7", 1, 7},
		{"simple 2.0", "2.0", 2, 0},
		{"with PDF prefix", "PDF-1.7", 1, 7},
		{"with spaces", " 1.4 ", 1, 4},
		{"with PDF prefix and spaces", " PDF-1.7 ", 1, 7},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseVersion(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.wantMajor, v.Major())
			assert.Equal(t, tt.wantMinor, v.Minor())
		})
	}
}

// TestParseVersion_Errors tests parsing invalid version strings.
func TestParseVersion_Errors(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"no dot", "14"},
		{"too many dots", "1.4.3"},
		{"empty string", ""},
		{"non-numeric major", "a.4"},
		{"non-numeric minor", "1.b"},
		{"negative in string", "-1.4"},
		{"missing minor", "1."},
		{"missing major", ".4"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := ParseVersion(tt.input)
			assert.Error(t, err)
			assert.True(t, errors.Is(err, ErrInvalidVersion))
			assert.Equal(t, Version{}, v)
		})
	}
}

// TestVersion_String tests version string representation.
func TestVersion_String(t *testing.T) {
	tests := []struct {
		name  string
		major int
		minor int
		want  string
	}{
		{"1.0", 1, 0, "1.0"},
		{"1.7", 1, 7, "1.7"},
		{"2.0", 2, 0, "2.0"},
		{"1.10", 1, 10, "1.10"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, err := NewVersion(tt.major, tt.minor)
			require.NoError(t, err)
			assert.Equal(t, tt.want, v.String())
		})
	}
}

// TestVersion_AtLeast tests version comparison.
func TestVersion_AtLeast(t *testing.T) {
	tests := []struct {
		name      string
		version   Version
		testMajor int
		testMinor int
		want      bool
	}{
		{"1.7 >= 1.4", PDF17, 1, 4, true},
		{"1.7 >= 1.7", PDF17, 1, 7, true},
		{"1.7 >= 2.0", PDF17, 2, 0, false},
		{"2.0 >= 1.7", PDF20, 1, 7, true},
		{"1.4 >= 1.5", PDF14, 1, 5, false},
		{"1.5 >= 1.4", PDF15, 1, 4, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.version.AtLeast(tt.testMajor, tt.testMinor)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestVersion_Equals tests version equality.
func TestVersion_Equals(t *testing.T) {
	tests := []struct {
		name string
		v1   Version
		v2   Version
		want bool
	}{
		{"equal 1.7", PDF17, PDF17, true},
		{"equal 2.0", PDF20, PDF20, true},
		{"different major", PDF14, PDF20, false},
		{"different minor", PDF14, PDF15, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Equals(tt.v2)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestVersion_Compare tests version comparison.
func TestVersion_Compare(t *testing.T) {
	tests := []struct {
		name string
		v1   Version
		v2   Version
		want int
	}{
		{"equal", PDF17, PDF17, 0},
		{"less major", PDF14, PDF20, -1},
		{"greater major", PDF20, PDF14, 1},
		{"less minor", PDF14, PDF15, -1},
		{"greater minor", PDF15, PDF14, 1},
		{"1.0 < 1.7", PDF10, PDF17, -1},
		{"2.0 > 1.7", PDF20, PDF17, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.v1.Compare(tt.v2)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestPredefinedVersions tests that predefined versions are valid.
func TestPredefinedVersions(t *testing.T) {
	tests := []struct {
		name  string
		v     Version
		major int
		minor int
	}{
		{"PDF10", PDF10, 1, 0},
		{"PDF11", PDF11, 1, 1},
		{"PDF12", PDF12, 1, 2},
		{"PDF13", PDF13, 1, 3},
		{"PDF14", PDF14, 1, 4},
		{"PDF15", PDF15, 1, 5},
		{"PDF16", PDF16, 1, 6},
		{"PDF17", PDF17, 1, 7},
		{"PDF20", PDF20, 2, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.major, tt.v.Major())
			assert.Equal(t, tt.minor, tt.v.Minor())
		})
	}
}

// TestVersion_Immutability tests that Version is immutable.
func TestVersion_Immutability(t *testing.T) {
	v1 := PDF17
	v2 := v1 // Copy

	// Verify they are equal
	assert.True(t, v1.Equals(v2))

	// Create a new version (we cannot modify v2 since fields are private)
	v3, err := NewVersion(2, 0)
	require.NoError(t, err)

	// v1 and v2 should still be equal and unaffected
	assert.True(t, v1.Equals(v2))
	assert.False(t, v1.Equals(v3))
}

// BenchmarkNewVersion benchmarks version creation.
func BenchmarkNewVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewVersion(1, 7)
	}
}

// BenchmarkParseVersion benchmarks version parsing.
func BenchmarkParseVersion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ParseVersion("PDF-1.7")
	}
}

// BenchmarkVersion_Compare benchmarks version comparison.
func BenchmarkVersion_Compare(b *testing.B) {
	v1 := PDF17
	v2 := PDF14
	for i := 0; i < b.N; i++ {
		_ = v1.Compare(v2)
	}
}
