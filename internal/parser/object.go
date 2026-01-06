// Package primitive implements PDF primitive object types as defined in PDF specification.
//
// PDF specification defines 8 basic object types:
// - Boolean values
// - Integer and Real numbers
// - Strings (literal and hexadecimal)
// - Names
// - Arrays
// - Dictionaries
// - Streams
// - The null object
//
// Reference: PDF 1.7 specification, Section 7.3 "Objects"
package parser

import (
	"fmt"
	"io"
)

// PdfObject is the base interface for all PDF objects.
// All PDF primitive types implement this interface.
type PdfObject interface {
	// String returns a string representation of the object.
	String() string

	// WriteTo writes the PDF representation to w.
	// Returns the number of bytes written and any error.
	WriteTo(w io.Writer) (int64, error)
}

// Type represents the type of a PDF object.
// This is useful for type assertions and debugging.
type Type int

// PDF object type constants.
const (
	// TypeNull represents the null type.
	TypeNull Type = iota
	TypeBoolean
	TypeInteger
	TypeReal
	TypeString
	TypeName
	TypeArray
	TypeDictionary
	TypeStream
	TypeIndirect
	TypeReference
)

// String returns the name of the type.
//
//nolint:cyclop // Switch statement for all types
func (t Type) String() string {
	switch t {
	case TypeNull:
		return "Null"
	case TypeBoolean:
		return "Boolean"
	case TypeInteger:
		return "Integer"
	case TypeReal:
		return "Real"
	case TypeString:
		return "String"
	case TypeName:
		return "Name"
	case TypeArray:
		return "Array"
	case TypeDictionary:
		return "Dictionary"
	case TypeStream:
		return "Stream"
	case TypeIndirect:
		return "Indirect"
	case TypeReference:
		return "Reference"
	default:
		return fmt.Sprintf("Unknown(%d)", t)
	}
}

// TypeOf returns the type of a PDF object.
func TypeOf(obj PdfObject) Type {
	switch obj.(type) {
	case *Null:
		return TypeNull
	case *Boolean:
		return TypeBoolean
	case *Integer:
		return TypeInteger
	case *Real:
		return TypeReal
	case *String:
		return TypeString
	case *Name:
		return TypeName
	case *Array:
		return TypeArray
	case *Dictionary:
		return TypeDictionary
	default:
		return Type(-1)
	}
}

// Clone creates a deep copy of a PDF object.
// This is useful when you need to modify an object without affecting the original.
func Clone(obj PdfObject) PdfObject {
	switch o := obj.(type) {
	case *Null:
		return NewNull()
	case *Boolean:
		return NewBoolean(o.Value())
	case *Integer:
		return NewInteger(o.Value())
	case *Real:
		return NewReal(o.Value())
	case *String:
		return NewString(o.Value())
	case *Name:
		return NewName(o.Value())
	case *Array:
		return o.Clone()
	case *Dictionary:
		return o.Clone()
	default:
		// For unknown types, return nil
		return nil
	}
}

// Resolve resolves indirect references to direct objects.
// For direct objects, it returns the object itself.
//
// Note: Full indirect object support (e.g., "1 0 R" references) will be
// implemented in Phase 2 (PDF Parser) as part of the document reader.
// See SUMMARY.md for the complete roadmap.
func Resolve(obj PdfObject) PdfObject {
	// Phase 1: Direct objects only
	// Phase 2: Will handle indirect references with cross-reference table
	return obj
}
