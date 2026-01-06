package creator

import (
	"testing"
)

func TestNewChapter(t *testing.T) {
	ch := NewChapter("Introduction")

	if ch.Title() != "Introduction" {
		t.Errorf("Expected title 'Introduction', got '%s'", ch.Title())
	}

	if len(ch.Number()) != 0 {
		t.Errorf("Expected empty number for new chapter, got %v", ch.Number())
	}

	if ch.PageIndex() != -1 {
		t.Errorf("Expected page index -1 for new chapter, got %d", ch.PageIndex())
	}
}

func TestChapterAdd(t *testing.T) {
	ch := NewChapter("Test")
	p := NewParagraph("Test paragraph")

	err := ch.Add(p)
	if err != nil {
		t.Errorf("Failed to add paragraph: %v", err)
	}

	if len(ch.Content()) != 1 {
		t.Errorf("Expected 1 content item, got %d", len(ch.Content()))
	}
}

func TestChapterAddNil(t *testing.T) {
	ch := NewChapter("Test")
	err := ch.Add(nil)

	if err == nil {
		t.Error("Expected error when adding nil drawable")
	}
}

func TestChapterNumbering(t *testing.T) {
	// Create chapter with manual numbering
	ch := NewChapter("Chapter 1")
	ch.assignNumbers([]int{}, 0)

	if ch.NumberString() != "1" {
		t.Errorf("Expected number '1', got '%s'", ch.NumberString())
	}

	// Create sub-chapter
	sub := ch.NewSubChapter("Section 1.1")

	if sub.NumberString() != "1.1" {
		t.Errorf("Expected number '1.1', got '%s'", sub.NumberString())
	}

	// Create sub-sub-chapter
	subsub := sub.NewSubChapter("Subsection 1.1.1")

	if subsub.NumberString() != "1.1.1" {
		t.Errorf("Expected number '1.1.1', got '%s'", subsub.NumberString())
	}
}

func TestChapterFullTitle(t *testing.T) {
	ch := NewChapter("Introduction")
	ch.assignNumbers([]int{}, 0)

	expected := "1 Introduction"
	if ch.FullTitle() != expected {
		t.Errorf("Expected full title '%s', got '%s'", expected, ch.FullTitle())
	}

	// Test with numbering disabled
	style := ch.Style()
	style.ShowNumber = false
	ch.SetStyle(style)

	if ch.FullTitle() != "Introduction" {
		t.Errorf("Expected full title 'Introduction', got '%s'", ch.FullTitle())
	}
}

func TestChapterLevel(t *testing.T) {
	ch := NewChapter("Chapter")
	ch.assignNumbers([]int{}, 0)

	if ch.Level() != 1 {
		t.Errorf("Expected level 1, got %d", ch.Level())
	}

	sub := ch.NewSubChapter("Section")
	if sub.Level() != 2 {
		t.Errorf("Expected level 2, got %d", sub.Level())
	}

	subsub := sub.NewSubChapter("Subsection")
	if subsub.Level() != 3 {
		t.Errorf("Expected level 3, got %d", subsub.Level())
	}
}

func TestChapterSubChapters(t *testing.T) {
	ch := NewChapter("Chapter")
	ch.assignNumbers([]int{}, 0)

	sec1 := ch.NewSubChapter("Section 1")
	sec2 := ch.NewSubChapter("Section 2")

	subChapters := ch.SubChapters()
	if len(subChapters) != 2 {
		t.Errorf("Expected 2 sub-chapters, got %d", len(subChapters))
	}

	if subChapters[0] != sec1 {
		t.Error("First sub-chapter doesn't match")
	}

	if subChapters[1] != sec2 {
		t.Error("Second sub-chapter doesn't match")
	}
}

func TestChapterGetAllChapters(t *testing.T) {
	// Create chapter hierarchy:
	// 1. Chapter 1
	//    1.1 Section 1
	//        1.1.1 Subsection 1
	//    1.2 Section 2
	ch := NewChapter("Chapter 1")
	ch.assignNumbers([]int{}, 0)

	sec1 := ch.NewSubChapter("Section 1")
	subsec := sec1.NewSubChapter("Subsection 1")
	sec2 := ch.NewSubChapter("Section 2")

	all := ch.GetAllChapters()

	// Should include: Chapter 1, Section 1, Subsection 1, Section 2
	expected := []*Chapter{ch, sec1, subsec, sec2}
	if len(all) != len(expected) {
		t.Errorf("Expected %d chapters, got %d", len(expected), len(all))
	}

	for i, exp := range expected {
		if all[i] != exp {
			t.Errorf("Chapter at index %d doesn't match", i)
		}
	}
}

func TestChapterSetPageIndex(t *testing.T) {
	ch := NewChapter("Test")

	if ch.PageIndex() != -1 {
		t.Errorf("Expected initial page index -1, got %d", ch.PageIndex())
	}

	ch.setPageIndex(5)

	if ch.PageIndex() != 5 {
		t.Errorf("Expected page index 5, got %d", ch.PageIndex())
	}
}

func TestChapterStyleInheritance(t *testing.T) {
	ch := NewChapter("Chapter")
	ch.assignNumbers([]int{}, 0)

	// Check parent style
	parentFontSize := ch.Style().FontSize

	// Create sub-chapter
	sub := ch.NewSubChapter("Section")

	// Sub-chapter should have smaller font size
	subFontSize := sub.Style().FontSize
	if subFontSize >= parentFontSize {
		t.Errorf("Expected sub-chapter font size (%.1f) to be smaller than parent (%.1f)",
			subFontSize, parentFontSize)
	}

	// Font size should be 85% of parent
	expectedSize := parentFontSize * 0.85
	if subFontSize != expectedSize {
		t.Errorf("Expected sub-chapter font size %.1f, got %.1f",
			expectedSize, subFontSize)
	}
}

func TestChapterNumberSeparator(t *testing.T) {
	ch := NewChapter("Test")
	ch.assignNumbers([]int{}, 0)
	sub := ch.NewSubChapter("Sub")

	// Default separator is "."
	if sub.NumberString() != "1.1" {
		t.Errorf("Expected '1.1', got '%s'", sub.NumberString())
	}

	// Change separator
	style := sub.Style()
	style.NumberSeparator = "-"
	sub.SetStyle(style)

	if sub.NumberString() != "1-1" {
		t.Errorf("Expected '1-1', got '%s'", sub.NumberString())
	}
}

func TestDefaultChapterStyle(t *testing.T) {
	style := DefaultChapterStyle()

	if style.Font != HelveticaBold {
		t.Errorf("Expected font HelveticaBold, got %v", style.Font)
	}

	if style.FontSize != 18 {
		t.Errorf("Expected font size 18, got %.1f", style.FontSize)
	}

	if style.Color != Black {
		t.Error("Expected color Black")
	}

	if !style.ShowNumber {
		t.Error("Expected ShowNumber to be true")
	}

	if style.NumberSeparator != "." {
		t.Errorf("Expected separator '.', got '%s'", style.NumberSeparator)
	}
}

func TestChapterParent(t *testing.T) {
	ch := NewChapter("Parent")
	ch.assignNumbers([]int{}, 0)

	if ch.Parent() != nil {
		t.Error("Expected top-level chapter to have nil parent")
	}

	sub := ch.NewSubChapter("Child")

	if sub.Parent() != ch {
		t.Error("Expected sub-chapter parent to be the main chapter")
	}
}

func TestMultipleSubChapterNumbering(t *testing.T) {
	ch := NewChapter("Chapter")
	ch.assignNumbers([]int{}, 0)

	sec1 := ch.NewSubChapter("Section 1")
	sec2 := ch.NewSubChapter("Section 2")
	sec3 := ch.NewSubChapter("Section 3")

	if sec1.NumberString() != "1.1" {
		t.Errorf("Expected '1.1', got '%s'", sec1.NumberString())
	}

	if sec2.NumberString() != "1.2" {
		t.Errorf("Expected '1.2', got '%s'", sec2.NumberString())
	}

	if sec3.NumberString() != "1.3" {
		t.Errorf("Expected '1.3', got '%s'", sec3.NumberString())
	}
}
