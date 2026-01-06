package creator

import (
	"strings"
	"testing"
)

func TestNewTOC(t *testing.T) {
	toc := NewTOC()

	if toc.Title() != "Table of Contents" {
		t.Errorf("Expected default title 'Table of Contents', got '%s'", toc.Title())
	}

	if !toc.ShowPageNumbers() {
		t.Error("Expected page numbers to be shown by default")
	}

	if toc.Leader() != "." {
		t.Errorf("Expected default leader '.', got '%s'", toc.Leader())
	}
}

func TestTOCSetTitle(t *testing.T) {
	toc := NewTOC()
	toc.SetTitle("Contents")

	if toc.Title() != "Contents" {
		t.Errorf("Expected title 'Contents', got '%s'", toc.Title())
	}
}

func TestTOCSetShowPageNumbers(t *testing.T) {
	toc := NewTOC()

	toc.SetShowPageNumbers(false)
	if toc.ShowPageNumbers() {
		t.Error("Expected page numbers to be hidden")
	}

	toc.SetShowPageNumbers(true)
	if !toc.ShowPageNumbers() {
		t.Error("Expected page numbers to be shown")
	}
}

func TestTOCSetLeader(t *testing.T) {
	toc := NewTOC()

	toc.SetLeader("-")
	if toc.Leader() != "-" {
		t.Errorf("Expected leader '-', got '%s'", toc.Leader())
	}
}

func TestDefaultTOCStyle(t *testing.T) {
	style := DefaultTOCStyle()

	if style.TitleFont != HelveticaBold {
		t.Errorf("Expected title font HelveticaBold, got %v", style.TitleFont)
	}

	if style.TitleSize != 24 {
		t.Errorf("Expected title size 24, got %.1f", style.TitleSize)
	}

	if style.EntryFont != Helvetica {
		t.Errorf("Expected entry font Helvetica, got %v", style.EntryFont)
	}

	if style.EntrySize != 12 {
		t.Errorf("Expected entry size 12, got %.1f", style.EntrySize)
	}

	if style.IndentPerLevel != 20 {
		t.Errorf("Expected indent 20, got %.1f", style.IndentPerLevel)
	}
}

func TestTOCSetStyle(t *testing.T) {
	toc := NewTOC()

	customStyle := TOCStyle{
		TitleFont:       TimesBold,
		TitleSize:       28,
		TitleColor:      Blue,
		EntryFont:       TimesRoman,
		EntrySize:       11,
		EntryColor:      DarkGray,
		IndentPerLevel:  25,
		LineSpacing:     1.3,
		SpaceAfterTitle: 25,
	}

	toc.SetStyle(customStyle)

	style := toc.Style()
	if style.TitleFont != TimesBold {
		t.Error("Custom title font not set correctly")
	}

	if style.TitleSize != 28 {
		t.Error("Custom title size not set correctly")
	}

	if style.IndentPerLevel != 25 {
		t.Error("Custom indent not set correctly")
	}
}

func TestTOCGetEntries(t *testing.T) {
	toc := NewTOC()

	// Create chapter hierarchy
	ch1 := NewChapter("Chapter 1")
	ch1.assignNumbers([]int{}, 0)
	ch1.setPageIndex(1)

	sec1 := ch1.NewSubChapter("Section 1.1")
	sec1.setPageIndex(2)

	ch2 := NewChapter("Chapter 2")
	ch2.assignNumbers([]int{}, 1)
	ch2.setPageIndex(3)

	// Set chapters in TOC
	toc.setChapters([]*Chapter{ch1, ch2})

	// Get entries
	entries := toc.GetEntries()

	// Should have 3 entries: Chapter 1, Section 1.1, Chapter 2
	if len(entries) != 3 {
		t.Fatalf("Expected 3 entries, got %d", len(entries))
	}

	// Verify entry structure
	verifyTOCEntry(t, entries[0], "Chapter 1", "1", 1, 1)
	verifyTOCEntry(t, entries[1], "Section 1.1", "1.1", 2, 2)
	verifyTOCEntry(t, entries[2], "Chapter 2", "2", 1, 3)
}

// verifyTOCEntry is a helper to verify TOC entry properties.
func verifyTOCEntry(t *testing.T, entry TOCEntry, title, number string, level, pageIndex int) {
	t.Helper()
	if entry.Title != title {
		t.Errorf("Expected entry title '%s', got '%s'", title, entry.Title)
	}
	if entry.Number != number {
		t.Errorf("Expected entry number '%s', got '%s'", number, entry.Number)
	}
	if entry.Level != level {
		t.Errorf("Expected entry level %d, got %d", level, entry.Level)
	}
	if entry.PageIndex != pageIndex {
		t.Errorf("Expected page index %d, got %d", pageIndex, entry.PageIndex)
	}
}

func TestTOCBuildEntryWithLeader(t *testing.T) {
	toc := NewTOC()

	// Create a minimal layout context for testing
	ctx := &LayoutContext{
		PageWidth:  595, // A4 width
		PageHeight: 842,
		Margins: Margins{
			Top:    72,
			Right:  72,
			Bottom: 72,
			Left:   72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	// Build entry with leader
	entry := toc.buildEntryWithLeader("Introduction", 1, 12, 0, ctx)

	// Entry should contain title, leaders, and page number
	if !strings.Contains(entry, "Introduction") {
		t.Error("Entry should contain title")
	}

	if !strings.Contains(entry, "1") {
		t.Error("Entry should contain page number")
	}

	// Should have leader dots
	if !strings.Contains(entry, "...") {
		t.Error("Entry should contain leader dots")
	}
}

func TestTOCEmptyChapters(t *testing.T) {
	toc := NewTOC()
	toc.setChapters([]*Chapter{})

	entries := toc.GetEntries()
	if len(entries) != 0 {
		t.Errorf("Expected 0 entries for empty chapters, got %d", len(entries))
	}
}

func TestTOCDeepNesting(t *testing.T) {
	toc := NewTOC()

	// Create deep hierarchy
	ch := NewChapter("Chapter 1")
	ch.assignNumbers([]int{}, 0)

	sec := ch.NewSubChapter("Section 1.1")
	subsec := sec.NewSubChapter("Subsection 1.1.1")
	_ = subsec.NewSubChapter("Subsubsection 1.1.1.1")

	toc.setChapters([]*Chapter{ch})

	entries := toc.GetEntries()

	// Should have 4 entries with increasing levels
	if len(entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(entries))
	}

	// Check levels
	expectedLevels := []int{1, 2, 3, 4}
	for i, entry := range entries {
		if entry.Level != expectedLevels[i] {
			t.Errorf("Entry %d: expected level %d, got %d",
				i, expectedLevels[i], entry.Level)
		}
	}

	// Check deepest numbering
	if entries[3].Number != "1.1.1.1" {
		t.Errorf("Expected deepest number '1.1.1.1', got '%s'",
			entries[3].Number)
	}
}

func TestTOCMultipleTopLevelChapters(t *testing.T) {
	toc := NewTOC()

	ch1 := NewChapter("Chapter 1")
	ch1.assignNumbers([]int{}, 0)
	ch1.setPageIndex(1)

	ch2 := NewChapter("Chapter 2")
	ch2.assignNumbers([]int{}, 1)
	ch2.setPageIndex(5)

	ch3 := NewChapter("Chapter 3")
	ch3.assignNumbers([]int{}, 2)
	ch3.setPageIndex(10)

	toc.setChapters([]*Chapter{ch1, ch2, ch3})

	entries := toc.GetEntries()

	if len(entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(entries))
	}

	// Check page indices
	expectedPages := []int{1, 5, 10}
	for i, entry := range entries {
		if entry.PageIndex != expectedPages[i] {
			t.Errorf("Entry %d: expected page %d, got %d",
				i, expectedPages[i], entry.PageIndex)
		}
	}
}

func TestTOCLeaderCustomization(t *testing.T) {
	toc := NewTOC()

	// Test different leader characters
	leaders := []string{".", "-", "_", "*"}

	for _, leader := range leaders {
		toc.SetLeader(leader)
		if toc.Leader() != leader {
			t.Errorf("Expected leader '%s', got '%s'", leader, toc.Leader())
		}
	}
}

func TestTOCHeightCalculation(t *testing.T) {
	toc := NewTOC()

	// Create chapters
	ch1 := NewChapter("Chapter 1")
	ch1.assignNumbers([]int{}, 0)

	ch2 := NewChapter("Chapter 2")
	ch2.assignNumbers([]int{}, 1)

	toc.setChapters([]*Chapter{ch1, ch2})

	ctx := &LayoutContext{
		PageWidth:  595,
		PageHeight: 842,
		Margins: Margins{
			Top:    72,
			Right:  72,
			Bottom: 72,
			Left:   72,
		},
		CursorX: 72,
		CursorY: 0,
	}

	height := toc.Height(ctx)

	// Height should be positive
	if height <= 0 {
		t.Errorf("Expected positive height, got %.1f", height)
	}

	// Height should include title, space, and entries
	style := toc.Style()
	minHeight := style.TitleSize*1.2 + style.SpaceAfterTitle + 2*style.EntrySize*style.LineSpacing
	if height < minHeight {
		t.Errorf("Height %.1f is less than expected minimum %.1f", height, minHeight)
	}
}
