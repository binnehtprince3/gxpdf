package parser

import "fmt"

// TokenType represents the type of a PDF token.
type TokenType int

// Token types recognized by the PDF lexer.
// Based on PDF 1.7 specification, Section 7.2 (Lexical Conventions).
const (
	// TokenError represents an error during tokenization.
	TokenError TokenType = iota

	// TokenEOF represents end of file.
	TokenEOF

	// Basic types.
	TokenInteger   // 123, -456, +789
	TokenReal      // 3.14, -2.5, .5
	TokenString    // (Hello) - literal string
	TokenHexString // <48656C6C6F> - hexadecimal string
	TokenName      // /Type, /Page, /Name#20With#20Spaces
	TokenBoolean   // true, false
	TokenNull      // null

	// Keywords.
	TokenKeyword // obj, endobj, stream, endstream, xref, trailer, startxref, R, n, f

	// Delimiters.
	TokenArrayStart // [
	TokenArrayEnd   // ]
	TokenDictStart  // <<
	TokenDictEnd    // >>
)

// PDF keyword string constants.
const (
	KeywordObj       = "obj"
	KeywordEndobj    = "endobj"
	KeywordStream    = "stream"
	KeywordEndstream = "endstream"
	KeywordXref      = "xref"
	KeywordTrailer   = "trailer"
	KeywordStartxref = "startxref"
)

// String returns the string representation of a TokenType.
//
//nolint:cyclop // Simple switch for token type names.
func (t TokenType) String() string {
	switch t {
	case TokenError:
		return "ERROR"
	case TokenEOF:
		return "EOF"
	case TokenInteger:
		return "INTEGER"
	case TokenReal:
		return "REAL"
	case TokenString:
		return "STRING"
	case TokenHexString:
		return "HEX_STRING"
	case TokenName:
		return "NAME"
	case TokenBoolean:
		return "BOOLEAN"
	case TokenNull:
		return "NULL"
	case TokenKeyword:
		return "KEYWORD"
	case TokenArrayStart:
		return "ARRAY_START"
	case TokenArrayEnd:
		return "ARRAY_END"
	case TokenDictStart:
		return "DICT_START"
	case TokenDictEnd:
		return "DICT_END"
	default:
		return "UNKNOWN"
	}
}

// Token represents a single token from the PDF byte stream.
type Token struct {
	Type   TokenType // Type of the token
	Value  string    // String value of the token
	Line   int       // Line number (1-based)
	Column int       // Column number (1-based)
}

// String returns a string representation of the token for debugging.
func (t Token) String() string {
	return fmt.Sprintf("%s(%q) at %d:%d", t.Type, t.Value, t.Line, t.Column)
}

// NewToken creates a new token with the given type, value, line, and column.
func NewToken(typ TokenType, value string, line, column int) Token {
	return Token{
		Type:   typ,
		Value:  value,
		Line:   line,
		Column: column,
	}
}

// ErrorToken creates an error token with the given message.
func ErrorToken(msg string, line, column int) Token {
	return Token{
		Type:   TokenError,
		Value:  msg,
		Line:   line,
		Column: column,
	}
}

// EOFToken creates an end-of-file token.
func EOFToken(line, column int) Token {
	return Token{
		Type:   TokenEOF,
		Value:  "",
		Line:   line,
		Column: column,
	}
}

// IsKeyword checks if a string is a PDF keyword.
func IsKeyword(s string) bool {
	switch s {
	case KeywordObj, KeywordEndobj, KeywordStream, KeywordEndstream,
		KeywordXref, KeywordTrailer, KeywordStartxref, "R", "n", "f":
		return true
	default:
		return false
	}
}

// IsContentStreamOperator checks if a string is a PDF content stream operator.
// These are operators used in content streams for graphics and text operations.
//
// Reference: PDF 1.7 specification, Appendix A (Operator Summary).
func IsContentStreamOperator(s string) bool {
	switch s {
	// Text object operators (Section 9.4)
	case "BT", "ET":
		return true

	// Text state operators (Section 9.3)
	case "Tc", "Tw", "Tz", "TL", "Tf", "Tr", "Ts":
		return true

	// Text positioning operators (Section 9.4.2)
	case "Td", "TD", "Tm", "T*":
		return true

	// Text showing operators (Section 9.4.3)
	case "Tj", "TJ", "'", "\"":
		return true

	// Graphics state operators (Section 8.4.4)
	case "q", "Q", "cm", "w", "J", "j", "M", "d", "ri", "i", "gs":
		return true

	// Path construction operators (Section 8.5.2)
	case "m", "l", "c", "v", "y", "h", "re":
		return true

	// Path painting operators (Section 8.5.3)
	case "S", "s", "f", "F", "f*", "B", "B*", "b", "b*", "n":
		return true

	// Clipping path operators (Section 8.5.4)
	case "W", "W*":
		return true

	// Color operators (Section 8.6)
	case "CS", "cs", "SC", "SCN", "sc", "scn", "G", "g", "RG", "rg", "K", "k":
		return true

	// Shading operators (Section 8.7.4.3)
	case "sh":
		return true

	// Inline image operators (Section 8.9.7)
	case "BI", "ID", "EI":
		return true

	// XObject operators (Section 8.8)
	case "Do":
		return true

	// Marked content operators (Section 14.6)
	case "MP", "DP", "BMC", "BDC", "EMC":
		return true

	// Compatibility operators (Section 9.9)
	case "BX", "EX":
		return true

	default:
		return false
	}
}
