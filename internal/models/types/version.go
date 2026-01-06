package types

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var (
	// ErrInvalidVersion is returned when a PDF version string is malformed.
	ErrInvalidVersion = errors.New("invalid PDF version")
)

// Version represents a PDF version (e.g., 1.4, 1.7, 2.0).
// This is a Value Object - immutable and compared by value.
type Version struct {
	major int
	minor int
}

// NewVersion creates a new Version from major and minor numbers.
func NewVersion(major, minor int) (Version, error) {
	if major < 0 || minor < 0 {
		return Version{}, fmt.Errorf("%w: version cannot be negative", ErrInvalidVersion)
	}
	if major > 2 {
		return Version{}, fmt.Errorf("%w: unsupported major version %d", ErrInvalidVersion, major)
	}
	return Version{major, minor}, nil
}

// ParseVersion parses a PDF version string (e.g., "1.4", "PDF-1.7", "2.0").
func ParseVersion(s string) (Version, error) {
	// Trim spaces first, then remove "PDF-" prefix if present
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "PDF-")

	parts := strings.Split(s, ".")
	if len(parts) != 2 {
		return Version{}, fmt.Errorf("%w: expected format 'X.Y', got '%s'", ErrInvalidVersion, s)
	}

	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return Version{}, fmt.Errorf("%w: invalid major version '%s'", ErrInvalidVersion, parts[0])
	}

	minor, err := strconv.Atoi(parts[1])
	if err != nil {
		return Version{}, fmt.Errorf("%w: invalid minor version '%s'", ErrInvalidVersion, parts[1])
	}

	return NewVersion(major, minor)
}

// String returns the version as a string (e.g., "1.7").
func (v Version) String() string {
	return fmt.Sprintf("%d.%d", v.major, v.minor)
}

// Major returns the major version number.
func (v Version) Major() int {
	return v.major
}

// Minor returns the minor version number.
func (v Version) Minor() int {
	return v.minor
}

// AtLeast checks if this version is at least the specified version.
func (v Version) AtLeast(major, minor int) bool {
	if v.major > major {
		return true
	}
	if v.major == major {
		return v.minor >= minor
	}
	return false
}

// Equals checks if two versions are equal.
func (v Version) Equals(other Version) bool {
	return v.major == other.major && v.minor == other.minor
}

// Compare compares two versions.
// Returns -1 if v < other, 0 if v == other, 1 if v > other.
func (v Version) Compare(other Version) int {
	if v.major < other.major {
		return -1
	}
	if v.major > other.major {
		return 1
	}
	if v.minor < other.minor {
		return -1
	}
	if v.minor > other.minor {
		return 1
	}
	return 0
}

// Common PDF versions.
var (
	// PDF10 represents PDF 1.0 (1993).
	PDF10 = Version{1, 0}

	// PDF11 represents PDF 1.1 (1996).
	PDF11 = Version{1, 1}

	// PDF12 represents PDF 1.2 (1996).
	PDF12 = Version{1, 2}

	// PDF13 represents PDF 1.3 (2000).
	PDF13 = Version{1, 3}

	// PDF14 represents PDF 1.4 (2001).
	PDF14 = Version{1, 4}

	// PDF15 represents PDF 1.5 (2003).
	PDF15 = Version{1, 5}

	// PDF16 represents PDF 1.6 (2004).
	PDF16 = Version{1, 6}

	// PDF17 represents PDF 1.7 (2008) - ISO 32000-1.
	PDF17 = Version{1, 7}

	// PDF20 represents PDF 2.0 (2017) - ISO 32000-2.
	PDF20 = Version{2, 0}
)
