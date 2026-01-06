package creator

import (
	"fmt"

	"github.com/coregx/gxpdf/creator/forms"
	"github.com/coregx/gxpdf/internal/document"
)

// convertFieldToDomain converts a creator form field to a domain form field.
//
// This function handles the conversion from the high-level creator API
// to the low-level domain model.
//
// Supported field types:
//   - *forms.TextField -> domain.FormField with type "Tx"
//
// Returns ErrUnsupportedFieldType if the field type is not recognized.
func convertFieldToDomain(field interface{}) (*document.FormField, error) {
	switch f := field.(type) {
	case *forms.TextField:
		return convertTextFieldToDomain(f)
	default:
		return nil, fmt.Errorf("%w: %T", ErrUnsupportedFieldType, field)
	}
}

// convertTextFieldToDomain converts a creator TextField to a domain FormField.
func convertTextFieldToDomain(tf *forms.TextField) (*document.FormField, error) {
	// Validate the text field before conversion
	if err := tf.Validate(); err != nil {
		return nil, fmt.Errorf("text field validation failed: %w", err)
	}

	// Create domain form field with type "Tx" (Text)
	field := document.NewFormField("Tx", tf.Name(), tf.Rect())

	// Set value and default value
	field.SetValue(tf.Value().(string))
	field.SetDefaultValue(tf.DefaultValue().(string))

	// Set flags
	field.SetFlags(tf.Flags())

	// Build default appearance string (/DA)
	appearance := buildAppearanceString(tf.FontName(), tf.FontSize(), tf.TextColor())
	field.SetAppearance(appearance)

	// Set border color if present
	if bc := tf.BorderColor(); bc != nil {
		field.SetBorderColor(bc[0], bc[1], bc[2])
	}

	// Set fill color if present
	if fc := tf.FillColor(); fc != nil {
		field.SetFillColor(fc[0], fc[1], fc[2])
	}

	// Set max length if specified
	if tf.MaxLength() > 0 {
		field.SetMaxLength(tf.MaxLength())
	}

	return field, nil
}

// buildAppearanceString builds the PDF default appearance string (/DA).
//
// The default appearance string specifies the font and color for text fields.
//
// Format: "/<FontName> <FontSize> Tf <r> <g> <b> rg"
//
// Example: "/Helv 12 Tf 0 0 0 rg" = Helvetica 12pt, black
func buildAppearanceString(fontName string, fontSize float64, color [3]float64) string {
	// Map font names to PDF font names
	pdfFontName := mapFontNameToPDF(fontName)

	// Build appearance string
	// Format: /<FontName> <FontSize> Tf <r> <g> <b> rg
	return fmt.Sprintf("/%s %.2f Tf %.3f %.3f %.3f rg",
		pdfFontName, fontSize, color[0], color[1], color[2])
}

// mapFontNameToPDF maps creator font names to PDF font names.
//
// PDF requires specific font names in the default appearance string.
func mapFontNameToPDF(fontName string) string {
	// Standard 14 fonts mapping
	switch fontName {
	case "Helvetica":
		return "Helv"
	case "Courier":
		return "Cour"
	case "Times-Roman", "Times":
		return "TiRo"
	default:
		// Use first 4 characters as fallback
		if len(fontName) > 4 {
			return fontName[:4]
		}
		return fontName
	}
}
