package parser

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// ============================================================================
// Null
// ============================================================================

// Null represents the PDF null object.
type Null struct{}

// NewNull creates a new Null object.
func NewNull() *Null {
	return &Null{}
}

// String returns "null".
func (n *Null) String() string {
	return "null"
}

// WriteTo writes "null" to w.
func (n *Null) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte("null"))
	return int64(written), err
}

// ============================================================================
// Boolean
// ============================================================================

// Boolean represents a PDF boolean value (true or false).
type Boolean struct {
	value bool
}

// NewBoolean creates a new Boolean object.
func NewBoolean(value bool) *Boolean {
	return &Boolean{value: value}
}

// Value returns the boolean value.
func (b *Boolean) Value() bool {
	return b.value
}

// String returns "true" or "false".
func (b *Boolean) String() string {
	if b.value {
		return "true"
	}
	return "false"
}

// WriteTo writes "true" or "false" to w.
func (b *Boolean) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte(b.String()))
	return int64(written), err
}

// ============================================================================
// Integer
// ============================================================================

// Integer represents a PDF integer object.
// PDF integers are signed 32-bit values.
type Integer struct {
	value int64
}

// NewInteger creates a new Integer object.
func NewInteger(value int64) *Integer {
	return &Integer{value: value}
}

// Value returns the integer value.
func (i *Integer) Value() int64 {
	return i.value
}

// Int returns the value as int (may overflow on 32-bit systems).
func (i *Integer) Int() int {
	return int(i.value)
}

// String returns the string representation of the integer.
func (i *Integer) String() string {
	return strconv.FormatInt(i.value, 10)
}

// WriteTo writes the integer to w.
func (i *Integer) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte(i.String()))
	return int64(written), err
}

// ============================================================================
// Real
// ============================================================================

// Real represents a PDF real (floating-point) number.
type Real struct {
	value float64
}

// NewReal creates a new Real object.
func NewReal(value float64) *Real {
	return &Real{value: value}
}

// Value returns the float64 value.
func (r *Real) Value() float64 {
	return r.value
}

// String returns the string representation of the real number.
// Uses minimal precision to keep PDF files compact.
func (r *Real) String() string {
	// Remove trailing zeros and decimal point if integer
	s := strconv.FormatFloat(r.value, 'f', -1, 64)

	// Handle special cases
	if strings.Contains(s, ".") {
		// Remove trailing zeros after decimal point
		s = strings.TrimRight(s, "0")
		s = strings.TrimRight(s, ".")
	}

	return s
}

// WriteTo writes the real number to w.
func (r *Real) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte(r.String()))
	return int64(written), err
}

// ============================================================================
// String
// ============================================================================

// String represents a PDF string object.
// PDF strings can be literal strings (text) or hexadecimal strings (<hex>).
type String struct {
	value []byte
	isHex bool // true if hexadecimal string
}

// NewString creates a new literal String object.
func NewString(value string) *String {
	return &String{
		value: []byte(value),
		isHex: false,
	}
}

// NewStringBytes creates a new literal String from bytes.
func NewStringBytes(value []byte) *String {
	return &String{
		value: value,
		isHex: false,
	}
}

// NewHexString creates a new hexadecimal String object.
func NewHexString(value string) *String {
	return &String{
		value: []byte(value),
		isHex: true,
	}
}

// Value returns the string value as a Go string.
func (s *String) Value() string {
	return string(s.value)
}

// Bytes returns the raw bytes.
func (s *String) Bytes() []byte {
	return s.value
}

// IsHex returns true if this is a hexadecimal string.
func (s *String) IsHex() bool {
	return s.isHex
}

// String returns the string representation.
// For debugging purposes, not PDF format.
func (s *String) String() string {
	if s.isHex {
		return fmt.Sprintf("<%x>", s.value)
	}
	return fmt.Sprintf("(%s)", s.escapeLiteral())
}

// WriteTo writes the PDF representation to w.
func (s *String) WriteTo(w io.Writer) (int64, error) {
	var buf bytes.Buffer

	if s.isHex {
		// Hexadecimal string: <48656C6C6F>
		buf.WriteByte('<')
		buf.WriteString(fmt.Sprintf("%X", s.value))
		buf.WriteByte('>')
	} else {
		// Literal string: (Hello)
		buf.WriteByte('(')
		buf.WriteString(s.escapeLiteral())
		buf.WriteByte(')')
	}

	written, err := w.Write(buf.Bytes())
	return int64(written), err
}

// escapeLiteral escapes special characters in literal strings.
func (s *String) escapeLiteral() string {
	var buf bytes.Buffer
	for _, b := range s.value {
		switch b {
		case '\\', '(', ')':
			buf.WriteByte('\\')
			buf.WriteByte(b)
		case '\n':
			buf.WriteString("\\n")
		case '\r':
			buf.WriteString("\\r")
		case '\t':
			buf.WriteString("\\t")
		default:
			buf.WriteByte(b)
		}
	}
	return buf.String()
}

// ============================================================================
// Name
// ============================================================================

// Name represents a PDF name object.
// Names are unique identifiers and always start with '/'.
type Name struct {
	value string
}

// NewName creates a new Name object.
// The leading '/' is added automatically if not present.
func NewName(value string) *Name {
	// Remove leading '/' if present (we'll add it back in String())
	value = strings.TrimPrefix(value, "/")
	return &Name{value: value}
}

// Value returns the name without the leading '/'.
func (n *Name) Value() string {
	return n.value
}

// String returns the name with leading '/'.
func (n *Name) String() string {
	return "/" + n.escape()
}

// WriteTo writes the name to w.
func (n *Name) WriteTo(w io.Writer) (int64, error) {
	written, err := w.Write([]byte(n.String()))
	return int64(written), err
}

// escape escapes special characters in names.
// Characters outside 33-126 (! to ~) except # must be written as #XX.
//
//nolint:cyclop // Multiple characters need escaping
func (n *Name) escape() string {
	var buf bytes.Buffer
	for _, r := range n.value {
		// Characters that need escaping
		if r < 33 || r > 126 || r == '#' || r == '/' || r == '(' || r == ')' ||
			r == '<' || r == '>' || r == '[' || r == ']' || r == '{' || r == '}' ||
			r == '%' {
			fmt.Fprintf(&buf, "#%02X", r)
		} else {
			buf.WriteRune(r)
		}
	}
	return buf.String()
}

// Equals checks if two names are equal.
func (n *Name) Equals(other *Name) bool {
	if other == nil {
		return false
	}
	return n.value == other.value
}
