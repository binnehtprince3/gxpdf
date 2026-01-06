package extractor

import (
	"bytes"
	"fmt"

	"github.com/coregx/gxpdf/internal/parser"
)

// Operator represents a PDF content stream operator with its operands.
//
// PDF content streams consist of a sequence of operators and their operands.
// The general format is:
//
//	operand1 operand2 ... operandN operator
//
// For example:
//   - "100 200 Td" - Move text position to (100, 200)
//   - "(Hello) Tj" - Show text "Hello"
//   - "/F1 12 Tf" - Set font F1 with size 12
//
// Reference: PDF 1.7 specification, Section 7.8.2 (Content Streams).
type Operator struct {
	Name     string             // Operator name (e.g., "Tj", "Tm", "BT")
	Operands []parser.PdfObject // Operands for the operator
}

// NewOperator creates a new Operator with the given name and operands.
func NewOperator(name string, operands []parser.PdfObject) *Operator {
	return &Operator{
		Name:     name,
		Operands: operands,
	}
}

// String returns a string representation of the operator.
func (op *Operator) String() string {
	return fmt.Sprintf("Operator{%s, operands=%d}", op.Name, len(op.Operands))
}

// ContentParser parses PDF content streams into operators.
//
// Content streams contain a sequence of operators that describe page graphics and text.
// The parser reads the stream and extracts operators with their operands.
//
// Example content stream:
//
//	BT
//	  /F1 12 Tf
//	  100 200 Td
//	  (Hello, World!) Tj
//	ET
//
// This would be parsed into operators:
//   - Operator{Name: "BT"}
//   - Operator{Name: "Tf", Operands: ["/F1", 12]}
//   - Operator{Name: "Td", Operands: [100, 200]}
//   - Operator{Name: "Tj", Operands: ["(Hello, World!)"]}
//   - Operator{Name: "ET"}
//
// Reference: PDF 1.7 specification, Section 7.8 (Content Streams).
type ContentParser struct {
	lexer *parser.Lexer
}

// NewContentParser creates a new ContentParser for the given content stream.
func NewContentParser(content []byte) *ContentParser {
	// Create lexer from bytes by wrapping in bytes.Reader
	lexer := parser.NewLexer(bytes.NewReader(content))
	return &ContentParser{
		lexer: lexer,
	}
}

// ParseOperators parses all operators from the content stream.
//
// Returns a slice of operators in the order they appear in the stream.
// Returns error if parsing fails.
//
// Content streams are sequences of objects followed by operators (keywords).
// Example: "100 200 Td" means: push 100, push 200, execute Td operator.
func (cp *ContentParser) ParseOperators() ([]*Operator, error) {
	var operators []*Operator
	var operandStack []parser.PdfObject

	for {
		// Get next token
		token, err := cp.lexer.NextToken()
		if err != nil {
			// Error occurred during tokenization
			return operators, err
		}
		if token.Type == parser.TokenEOF {
			// End of stream reached
			break
		}

		// Check if token is an operator (keyword)
		if token.Type == parser.TokenKeyword {
			// Create operator with current operand stack
			op := NewOperator(token.Value, operandStack)
			operators = append(operators, op)

			// Clear operand stack
			operandStack = nil
		} else {
			// Token is an operand, convert it to an object
			obj, err := cp.tokenToObject(token)
			if err != nil {
				return nil, fmt.Errorf("failed to parse operand: %w", err)
			}
			operandStack = append(operandStack, obj)
		}
	}

	return operators, nil
}

// tokenToObject converts a token to a PDF object.
//
//nolint:cyclop // Token type checking requires many cases
func (cp *ContentParser) tokenToObject(token parser.Token) (parser.PdfObject, error) {
	switch token.Type {
	case parser.TokenNull:
		return parser.NewNull(), nil

	case parser.TokenBoolean:
		if token.Value == "true" {
			return parser.NewBoolean(true), nil
		}
		return parser.NewBoolean(false), nil

	case parser.TokenInteger:
		// Parse integer from string value
		var val int64
		_, err := fmt.Sscanf(token.Value, "%d", &val)
		if err != nil {
			return nil, fmt.Errorf("invalid integer: %s", token.Value)
		}
		return parser.NewInteger(val), nil

	case parser.TokenReal:
		// Parse float from string value
		var val float64
		_, err := fmt.Sscanf(token.Value, "%f", &val)
		if err != nil {
			return nil, fmt.Errorf("invalid real: %s", token.Value)
		}
		return parser.NewReal(val), nil

	case parser.TokenString, parser.TokenHexString:
		return parser.NewString(token.Value), nil

	case parser.TokenName:
		// Remove leading slash if present
		name := token.Value
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}
		return parser.NewName(name), nil

	case parser.TokenArrayStart:
		// Parse array elements until ArrayEnd token
		return cp.parseArray()

	case parser.TokenArrayEnd:
		// ARRAY_END should only appear inside parseArray(), not at top level
		// This means unbalanced brackets in the content stream
		return nil, fmt.Errorf("unexpected array end token (unbalanced brackets)")

	case parser.TokenDictStart:
		// For dictionaries, parse dict elements until DictEnd token
		return cp.parseDictionary()

	case parser.TokenDictEnd:
		// Similar to ARRAY_END, should only appear inside parseDictionary()
		return nil, fmt.Errorf("unexpected dictionary end token (unbalanced brackets)")

	default:
		return nil, fmt.Errorf("unexpected token type for operand: %v", token.Type)
	}
}

// parseArray parses an array from the content stream.
//
// Assumes TokenArrayStart has already been consumed.
// Reads tokens until TokenArrayEnd is found.
func (cp *ContentParser) parseArray() (parser.PdfObject, error) {
	arr := parser.NewArray()

	for {
		token, err := cp.lexer.NextToken()
		if err != nil {
			return nil, fmt.Errorf("error reading array element: %w", err)
		}

		if token.Type == parser.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF while parsing array")
		}

		if token.Type == parser.TokenArrayEnd {
			// Array complete
			return arr, nil
		}

		// Convert token to object and add to array
		obj, err := cp.tokenToObject(token)
		if err != nil {
			return nil, fmt.Errorf("failed to parse array element: %w", err)
		}

		arr.Append(obj)
	}
}

// parseDictionary parses a dictionary from the content stream.
//
// Assumes TokenDictStart has already been consumed.
// Reads key-value pairs until TokenDictEnd is found.
func (cp *ContentParser) parseDictionary() (parser.PdfObject, error) {
	dict := parser.NewDictionary()

	for {
		// Read key
		keyToken, err := cp.lexer.NextToken()
		if err != nil {
			return nil, fmt.Errorf("error reading dictionary key: %w", err)
		}

		if keyToken.Type == parser.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF while parsing dictionary")
		}

		if keyToken.Type == parser.TokenDictEnd {
			// Dictionary complete
			return dict, nil
		}

		// Key must be a name
		if keyToken.Type != parser.TokenName {
			return nil, fmt.Errorf("dictionary key must be a name, got %v", keyToken.Type)
		}

		keyName := keyToken.Value
		if len(keyName) > 0 && keyName[0] == '/' {
			keyName = keyName[1:] // Remove leading slash
		}

		// Read value
		valueToken, err := cp.lexer.NextToken()
		if err != nil {
			return nil, fmt.Errorf("error reading dictionary value: %w", err)
		}

		if valueToken.Type == parser.TokenEOF {
			return nil, fmt.Errorf("unexpected EOF while reading dictionary value")
		}

		// Convert value token to object
		valueObj, err := cp.tokenToObject(valueToken)
		if err != nil {
			return nil, fmt.Errorf("failed to parse dictionary value: %w", err)
		}

		dict.Set(keyName, valueObj)
	}
}
