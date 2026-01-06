package creator

import (
	"testing"
)

func TestNewList_Defaults(t *testing.T) {
	list := NewList()

	if list == nil {
		t.Fatal("NewList() returned nil")
	}

	// Test marker settings.
	assertMarkerType(t, list, BulletMarker)
	assertBulletChar(t, list, "•")

	// Test font settings.
	assertFont(t, list, Helvetica, 12)
	assertColor(t, list, Black)

	// Test spacing settings.
	assertLineSpacing(t, list, 1.2)
	assertIndent(t, list, 20)
	assertMarkerIndent(t, list, 10)

	// Test numbering settings.
	assertStartNumber(t, list, 1)
	assertItemCount(t, list, 0)
}

// Helper functions for assertions.
func assertMarkerType(t *testing.T, list *List, expected MarkerType) {
	t.Helper()
	if list.markerType != expected {
		t.Errorf("Expected marker type %v, got %v", expected, list.markerType)
	}
}

func assertBulletChar(t *testing.T, list *List, expected string) {
	t.Helper()
	if list.bulletChar != expected {
		t.Errorf("Expected bullet char %q, got %q", expected, list.bulletChar)
	}
}

func assertFont(t *testing.T, list *List, expectedFont FontName, expectedSize float64) {
	t.Helper()
	if list.font != expectedFont {
		t.Errorf("Expected font %v, got %v", expectedFont, list.font)
	}
	if list.fontSize != expectedSize {
		t.Errorf("Expected font size %f, got %f", expectedSize, list.fontSize)
	}
}

func assertColor(t *testing.T, list *List, expected Color) {
	t.Helper()
	if list.color != expected {
		t.Errorf("Expected color %v, got %v", expected, list.color)
	}
}

func assertLineSpacing(t *testing.T, list *List, expected float64) {
	t.Helper()
	if list.lineSpacing != expected {
		t.Errorf("Expected line spacing %f, got %f", expected, list.lineSpacing)
	}
}

func assertIndent(t *testing.T, list *List, expected float64) {
	t.Helper()
	if list.indent != expected {
		t.Errorf("Expected indent %f, got %f", expected, list.indent)
	}
}

func assertMarkerIndent(t *testing.T, list *List, expected float64) {
	t.Helper()
	if list.markerIndent != expected {
		t.Errorf("Expected marker indent %f, got %f", expected, list.markerIndent)
	}
}

func assertStartNumber(t *testing.T, list *List, expected int) {
	t.Helper()
	if list.startNumber != expected {
		t.Errorf("Expected start number %d, got %d", expected, list.startNumber)
	}
}

func assertItemCount(t *testing.T, list *List, expected int) {
	t.Helper()
	if len(list.items) != expected {
		t.Errorf("Expected %d items, got %d", expected, len(list.items))
	}
}

func TestNewNumberedList(t *testing.T) {
	list := NewNumberedList()

	if list == nil {
		t.Fatal("NewNumberedList() returned nil")
	}

	if list.markerType != NumberMarker {
		t.Errorf("Expected marker type NumberMarker, got %v", list.markerType)
	}

	if list.numberFormat != NumberFormatArabic {
		t.Errorf("Expected number format Arabic, got %v", list.numberFormat)
	}
}

func TestList_SetMarkerType(t *testing.T) {
	list := NewList()

	result := list.SetMarkerType(NumberMarker)

	if result != list {
		t.Error("SetMarkerType() should return the same list for chaining")
	}

	if list.markerType != NumberMarker {
		t.Errorf("Expected marker type NumberMarker, got %v", list.markerType)
	}
}

func TestList_SetBulletChar(t *testing.T) {
	list := NewList()

	result := list.SetBulletChar("-")

	if result != list {
		t.Error("SetBulletChar() should return the same list for chaining")
	}

	if list.bulletChar != "-" {
		t.Errorf("Expected bullet char '-', got %q", list.bulletChar)
	}
}

func TestList_SetNumberFormat(t *testing.T) {
	tests := []struct {
		name   string
		format NumberFormat
	}{
		{"Arabic", NumberFormatArabic},
		{"LowerAlpha", NumberFormatLowerAlpha},
		{"UpperAlpha", NumberFormatUpperAlpha},
		{"LowerRoman", NumberFormatLowerRoman},
		{"UpperRoman", NumberFormatUpperRoman},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewNumberedList()
			result := list.SetNumberFormat(tt.format)

			if result != list {
				t.Error("SetNumberFormat() should return the same list for chaining")
			}

			if list.numberFormat != tt.format {
				t.Errorf("Expected number format %v, got %v", tt.format, list.numberFormat)
			}
		})
	}
}

func TestList_SetFont(t *testing.T) {
	list := NewList()

	result := list.SetFont(TimesBold, 14)

	if result != list {
		t.Error("SetFont() should return the same list for chaining")
	}

	if list.font != TimesBold {
		t.Errorf("Expected font TimesBold, got %v", list.font)
	}

	if list.fontSize != 14 {
		t.Errorf("Expected font size 14, got %f", list.fontSize)
	}
}

func TestList_SetColor(t *testing.T) {
	list := NewList()

	result := list.SetColor(Red)

	if result != list {
		t.Error("SetColor() should return the same list for chaining")
	}

	if list.color != Red {
		t.Errorf("Expected color Red, got %v", list.color)
	}
}

func TestList_SetLineSpacing(t *testing.T) {
	list := NewList()

	result := list.SetLineSpacing(1.5)

	if result != list {
		t.Error("SetLineSpacing() should return the same list for chaining")
	}

	if list.lineSpacing != 1.5 {
		t.Errorf("Expected line spacing 1.5, got %f", list.lineSpacing)
	}
}

func TestList_SetIndent(t *testing.T) {
	list := NewList()

	result := list.SetIndent(30)

	if result != list {
		t.Error("SetIndent() should return the same list for chaining")
	}

	if list.indent != 30 {
		t.Errorf("Expected indent 30, got %f", list.indent)
	}
}

func TestList_SetMarkerIndent(t *testing.T) {
	list := NewList()

	result := list.SetMarkerIndent(15)

	if result != list {
		t.Error("SetMarkerIndent() should return the same list for chaining")
	}

	if list.markerIndent != 15 {
		t.Errorf("Expected marker indent 15, got %f", list.markerIndent)
	}
}

func TestList_SetStartNumber(t *testing.T) {
	list := NewNumberedList()

	result := list.SetStartNumber(5)

	if result != list {
		t.Error("SetStartNumber() should return the same list for chaining")
	}

	if list.startNumber != 5 {
		t.Errorf("Expected start number 5, got %d", list.startNumber)
	}
}

func TestList_Add(t *testing.T) {
	list := NewList()

	result := list.Add("First item")

	if result != list {
		t.Error("Add() should return the same list for chaining")
	}

	if len(list.items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list.items))
	}

	if list.items[0].text != "First item" {
		t.Errorf("Expected item text 'First item', got %q", list.items[0].text)
	}

	if list.items[0].subList != nil {
		t.Error("Expected no sublist")
	}

	// Add more items.
	list.Add("Second item").Add("Third item")

	if len(list.items) != 3 {
		t.Fatalf("Expected 3 items, got %d", len(list.items))
	}
}

func TestList_AddItem(t *testing.T) {
	list := NewList()
	item := NewListItem("Custom item")

	result := list.AddItem(item)

	if result != list {
		t.Error("AddItem() should return the same list for chaining")
	}

	if len(list.items) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(list.items))
	}

	if list.items[0].text != "Custom item" {
		t.Errorf("Expected item text 'Custom item', got %q", list.items[0].text)
	}
}

func TestList_AddSubList(t *testing.T) {
	list := NewList()
	list.Add("Parent item")

	subList := NewList()
	subList.Add("Sub item 1").Add("Sub item 2")

	result := list.AddSubList(subList)

	if result != list {
		t.Error("AddSubList() should return the same list for chaining")
	}

	if len(list.items) != 1 {
		t.Fatalf("Expected 1 parent item, got %d", len(list.items))
	}

	if list.items[0].subList == nil {
		t.Fatal("Expected sublist to be attached")
	}

	if len(list.items[0].subList.items) != 2 {
		t.Errorf("Expected 2 sub-items, got %d", len(list.items[0].subList.items))
	}
}

func TestList_MethodChaining(t *testing.T) {
	list := NewList().
		SetBulletChar("-").
		SetFont(Courier, 10).
		SetColor(Blue).
		SetLineSpacing(1.5).
		SetIndent(25).
		SetMarkerIndent(12).
		Add("Item 1").
		Add("Item 2")

	if list.bulletChar != "-" {
		t.Errorf("Expected bullet char '-', got %q", list.bulletChar)
	}

	if list.font != Courier {
		t.Errorf("Expected font Courier, got %v", list.font)
	}

	if list.fontSize != 10 {
		t.Errorf("Expected font size 10, got %f", list.fontSize)
	}

	if list.color != Blue {
		t.Errorf("Expected color Blue, got %v", list.color)
	}

	if list.lineSpacing != 1.5 {
		t.Errorf("Expected line spacing 1.5, got %f", list.lineSpacing)
	}

	if list.indent != 25 {
		t.Errorf("Expected indent 25, got %f", list.indent)
	}

	if list.markerIndent != 12 {
		t.Errorf("Expected marker indent 12, got %f", list.markerIndent)
	}

	if len(list.items) != 2 {
		t.Errorf("Expected 2 items, got %d", len(list.items))
	}
}

func TestList_GetMarker_Bullet(t *testing.T) {
	list := NewList()

	marker := list.getMarker(0)
	if marker != "•" {
		t.Errorf("Expected marker '•', got %q", marker)
	}

	list.SetBulletChar("*")
	marker = list.getMarker(0)
	if marker != "*" {
		t.Errorf("Expected marker '*', got %q", marker)
	}
}

func TestList_GetMarker_Number(t *testing.T) {
	list := NewNumberedList()

	tests := []struct {
		index    int
		expected string
	}{
		{0, "1."},
		{1, "2."},
		{2, "3."},
	}

	for _, tt := range tests {
		marker := list.getMarker(tt.index)
		if marker != tt.expected {
			t.Errorf("For index %d, expected marker %q, got %q",
				tt.index, tt.expected, marker)
		}
	}
}

func TestList_FormatNumber_Arabic(t *testing.T) {
	list := NewNumberedList()
	list.SetNumberFormat(NumberFormatArabic)

	tests := []struct {
		number   int
		expected string
	}{
		{1, "1."},
		{2, "2."},
		{10, "10."},
		{99, "99."},
	}

	for _, tt := range tests {
		result := list.formatNumber(tt.number)
		if result != tt.expected {
			t.Errorf("For number %d, expected %q, got %q",
				tt.number, tt.expected, result)
		}
	}
}

func TestList_FormatNumber_LowerAlpha(t *testing.T) {
	list := NewNumberedList()
	list.SetNumberFormat(NumberFormatLowerAlpha)

	tests := []struct {
		number   int
		expected string
	}{
		{1, "a."},
		{2, "b."},
		{26, "z."},
		{27, "27."}, // Out of range, fallback to numeric
	}

	for _, tt := range tests {
		result := list.formatNumber(tt.number)
		if result != tt.expected {
			t.Errorf("For number %d, expected %q, got %q",
				tt.number, tt.expected, result)
		}
	}
}

func TestList_FormatNumber_UpperAlpha(t *testing.T) {
	list := NewNumberedList()
	list.SetNumberFormat(NumberFormatUpperAlpha)

	tests := []struct {
		number   int
		expected string
	}{
		{1, "A."},
		{2, "B."},
		{26, "Z."},
		{27, "27."}, // Out of range, fallback to numeric
	}

	for _, tt := range tests {
		result := list.formatNumber(tt.number)
		if result != tt.expected {
			t.Errorf("For number %d, expected %q, got %q",
				tt.number, tt.expected, result)
		}
	}
}

func TestList_FormatNumber_LowerRoman(t *testing.T) {
	list := NewNumberedList()
	list.SetNumberFormat(NumberFormatLowerRoman)

	tests := []struct {
		number   int
		expected string
	}{
		{1, "i."},
		{2, "ii."},
		{3, "iii."},
		{4, "iv."},
		{5, "v."},
		{9, "ix."},
		{10, "x."},
		{49, "xlix."},
		{50, "l."},
	}

	for _, tt := range tests {
		result := list.formatNumber(tt.number)
		if result != tt.expected {
			t.Errorf("For number %d, expected %q, got %q",
				tt.number, tt.expected, result)
		}
	}
}

func TestList_FormatNumber_UpperRoman(t *testing.T) {
	list := NewNumberedList()
	list.SetNumberFormat(NumberFormatUpperRoman)

	tests := []struct {
		number   int
		expected string
	}{
		{1, "I."},
		{2, "II."},
		{3, "III."},
		{4, "IV."},
		{5, "V."},
		{9, "IX."},
		{10, "X."},
		{49, "XLIX."},
		{50, "L."},
	}

	for _, tt := range tests {
		result := list.formatNumber(tt.number)
		if result != tt.expected {
			t.Errorf("For number %d, expected %q, got %q",
				tt.number, tt.expected, result)
		}
	}
}

func TestToRoman(t *testing.T) {
	tests := []struct {
		number int
		upper  bool
		want   string
	}{
		{1, false, "i"},
		{1, true, "I"},
		{4, false, "iv"},
		{4, true, "IV"},
		{5, false, "v"},
		{5, true, "V"},
		{9, false, "ix"},
		{9, true, "IX"},
		{10, false, "x"},
		{10, true, "X"},
		{40, false, "xl"},
		{40, true, "XL"},
		{50, false, "l"},
		{50, true, "L"},
		{90, false, "xc"},
		{90, true, "XC"},
		{100, false, "c"},
		{100, true, "C"},
		{400, false, "cd"},
		{400, true, "CD"},
		{500, false, "d"},
		{500, true, "D"},
		{900, false, "cm"},
		{900, true, "CM"},
		{1000, false, "m"},
		{1000, true, "M"},
		{1994, false, "mcmxciv"},
		{1994, true, "MCMXCIV"},
		{3999, false, "mmmcmxcix"},
		{3999, true, "MMMCMXCIX"},
		{0, false, "0"},       // Out of range
		{4000, false, "4000"}, // Out of range
	}

	for _, tt := range tests {
		got := toRoman(tt.number, tt.upper)
		if got != tt.want {
			t.Errorf("toRoman(%d, %v) = %q, want %q",
				tt.number, tt.upper, got, tt.want)
		}
	}
}

func TestList_TextWrapping(t *testing.T) {
	list := NewList()
	list.Add("This is a very long text that should be wrapped across multiple lines when rendered")

	// Test that wrapping works.
	lines := list.wrapText(
		"This is a very long text that should be wrapped across multiple lines",
		200, // Small width to force wrapping
	)

	if len(lines) <= 1 {
		t.Error("Expected text to wrap into multiple lines")
	}
}

func TestList_WrapText_Empty(t *testing.T) {
	list := NewList()

	lines := list.wrapText("", 500)
	if len(lines) != 0 {
		t.Errorf("Expected 0 lines for empty text, got %d", len(lines))
	}
}

func TestList_WrapText_SingleWord(t *testing.T) {
	list := NewList()

	lines := list.wrapText("Word", 500)
	if len(lines) != 1 {
		t.Errorf("Expected 1 line for single word, got %d", len(lines))
	}

	if lines[0] != "Word" {
		t.Errorf("Expected line 'Word', got %q", lines[0])
	}
}

func TestList_Height(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()

	list := NewList()
	list.Add("Item 1")
	list.Add("Item 2")

	height := list.Height(ctx)
	if height <= 0 {
		t.Error("Expected positive height")
	}

	// Height should be approximately 2 items * line height.
	expectedHeight := 2 * list.calculateLineHeight()
	if height < expectedHeight*0.9 || height > expectedHeight*1.1 {
		t.Errorf("Expected height around %f, got %f", expectedHeight, height)
	}
}

func TestList_ImplementsDrawable(_ *testing.T) {
	var _ Drawable = (*List)(nil)
}

func TestList_Draw(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()

	list := NewList()
	list.Add("Item 1")
	list.Add("Item 2")

	err = list.Draw(ctx, page)
	if err != nil {
		t.Errorf("Draw() failed: %v", err)
	}

	// Check that cursor moved.
	if ctx.CursorY <= 0 {
		t.Error("Expected cursor to move after drawing")
	}
}

func TestList_Draw_NumberedList(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()

	list := NewNumberedList()
	list.Add("First")
	list.Add("Second")
	list.Add("Third")

	err = list.Draw(ctx, page)
	if err != nil {
		t.Errorf("Draw() failed: %v", err)
	}
}

func TestList_Draw_NestedLists(t *testing.T) {
	c := New()
	page, err := c.NewPage()
	if err != nil {
		t.Fatalf("Failed to create page: %v", err)
	}

	ctx := page.GetLayoutContext()

	list := NewList()
	list.Add("Parent item 1")

	subList := NewList()
	subList.SetBulletChar("-")
	subList.Add("Sub item A")
	subList.Add("Sub item B")

	list.AddSubList(subList)
	list.Add("Parent item 2")

	err = list.Draw(ctx, page)
	if err != nil {
		t.Errorf("Draw() failed: %v", err)
	}
}

func TestList_StartNumber(t *testing.T) {
	list := NewNumberedList()
	list.SetStartNumber(5)

	marker := list.getMarker(0)
	if marker != "5." {
		t.Errorf("Expected marker '5.', got %q", marker)
	}

	marker = list.getMarker(1)
	if marker != "6." {
		t.Errorf("Expected marker '6.', got %q", marker)
	}
}

func TestList_NumberFormats(t *testing.T) {
	tests := []struct {
		name         string
		format       NumberFormat
		startNumber  int
		index        int
		expectedMark string
	}{
		{"Arabic_1", NumberFormatArabic, 1, 0, "1."},
		{"Arabic_2", NumberFormatArabic, 1, 1, "2."},
		{"LowerAlpha_a", NumberFormatLowerAlpha, 1, 0, "a."},
		{"LowerAlpha_b", NumberFormatLowerAlpha, 1, 1, "b."},
		{"UpperAlpha_A", NumberFormatUpperAlpha, 1, 0, "A."},
		{"UpperAlpha_B", NumberFormatUpperAlpha, 1, 1, "B."},
		{"LowerRoman_i", NumberFormatLowerRoman, 1, 0, "i."},
		{"LowerRoman_ii", NumberFormatLowerRoman, 1, 1, "ii."},
		{"UpperRoman_I", NumberFormatUpperRoman, 1, 0, "I."},
		{"UpperRoman_II", NumberFormatUpperRoman, 1, 1, "II."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list := NewNumberedList()
			list.SetNumberFormat(tt.format)
			list.SetStartNumber(tt.startNumber)

			marker := list.getMarker(tt.index)
			if marker != tt.expectedMark {
				t.Errorf("Expected marker %q, got %q", tt.expectedMark, marker)
			}
		})
	}
}

func TestNewListItem(t *testing.T) {
	item := NewListItem("Test item")

	if item.text != "Test item" {
		t.Errorf("Expected text 'Test item', got %q", item.text)
	}

	if item.subList != nil {
		t.Error("Expected no sublist")
	}
}

func TestNewListItemWithSubList(t *testing.T) {
	subList := NewList()
	subList.Add("Sub item")

	item := NewListItemWithSubList("Parent", subList)

	if item.text != "Parent" {
		t.Errorf("Expected text 'Parent', got %q", item.text)
	}

	if item.subList == nil {
		t.Fatal("Expected sublist to be set")
	}

	if len(item.subList.items) != 1 {
		t.Errorf("Expected 1 sub-item, got %d", len(item.subList.items))
	}
}

func TestList_CalculateLineHeight(t *testing.T) {
	list := NewList()
	list.SetFont(Helvetica, 12)
	list.SetLineSpacing(1.5)

	expected := 12.0 * 1.5
	got := list.calculateLineHeight()

	if got != expected {
		t.Errorf("Expected line height %f, got %f", expected, got)
	}
}
