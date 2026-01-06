package creator

import (
	"errors"
	"testing"
)

// TestAddBookmark_Success tests adding valid bookmarks.
func TestAddBookmark_Success(t *testing.T) {
	c := New()

	// Add pages.
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}

	// Add top-level bookmark.
	err := c.AddBookmark("Chapter 1", 0, 0)
	if err != nil {
		t.Errorf("AddBookmark failed: %v", err)
	}

	// Add nested bookmark.
	err = c.AddBookmark("Section 1.1", 0, 1)
	if err != nil {
		t.Errorf("AddBookmark failed for nested bookmark: %v", err)
	}

	// Add another top-level bookmark.
	err = c.AddBookmark("Chapter 2", 1, 0)
	if err != nil {
		t.Errorf("AddBookmark failed for second chapter: %v", err)
	}

	// Verify bookmarks were added.
	bookmarks := c.Bookmarks()
	if len(bookmarks) != 3 {
		t.Errorf("Expected 3 bookmarks, got %d", len(bookmarks))
	}
}

// TestAddBookmark_ValidateTitle tests title validation.
func TestAddBookmark_ValidateTitle(t *testing.T) {
	c := New()

	// Add a page.
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}

	// Empty title should fail.
	err := c.AddBookmark("", 0, 0)
	if err == nil {
		t.Error("Expected error for empty title, got nil")
	}
	if !errors.Is(err, ErrEmptyBookmarkTitle) {
		t.Errorf("Expected ErrEmptyBookmarkTitle, got: %v", err)
	}
}

// TestAddBookmark_ValidatePageIndex tests page index validation.
func TestAddBookmark_ValidatePageIndex(t *testing.T) {
	c := New()

	// Add a page.
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}

	tests := []struct {
		name      string
		pageIndex int
		wantErr   bool
	}{
		{
			name:      "Valid page index 0",
			pageIndex: 0,
			wantErr:   false,
		},
		{
			name:      "Valid page index 10",
			pageIndex: 10,
			wantErr:   false,
		},
		{
			name:      "Negative page index",
			pageIndex: -1,
			wantErr:   true,
		},
		{
			name:      "Negative page index -5",
			pageIndex: -5,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.AddBookmark("Test", tt.pageIndex, 0)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBookmark() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidBookmarkPage) {
				t.Errorf("Expected ErrInvalidBookmarkPage, got: %v", err)
			}
		})
	}
}

// TestAddBookmark_ValidateLevel tests level validation.
func TestAddBookmark_ValidateLevel(t *testing.T) {
	c := New()

	// Add a page.
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}

	tests := []struct {
		name    string
		level   int
		wantErr bool
	}{
		{
			name:    "Level 0 - top-level",
			level:   0,
			wantErr: false,
		},
		{
			name:    "Level 1 - child",
			level:   1,
			wantErr: false,
		},
		{
			name:    "Level 5 - deep nesting",
			level:   5,
			wantErr: false,
		},
		{
			name:    "Negative level",
			level:   -1,
			wantErr: true,
		},
		{
			name:    "Negative level -3",
			level:   -3,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := c.AddBookmark("Test", 0, tt.level)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddBookmark() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && !errors.Is(err, ErrInvalidBookmarkLevel) {
				t.Errorf("Expected ErrInvalidBookmarkLevel, got: %v", err)
			}
		})
	}
}

// TestBookmarks_ReturnsCopy tests that Bookmarks() returns a copy.
func TestBookmarks_ReturnsCopy(t *testing.T) {
	c := New()

	// Add a page.
	if _, err := c.NewPage(); err != nil {
		t.Fatalf("Failed to add page: %v", err)
	}

	// Add a bookmark.
	if err := c.AddBookmark("Chapter 1", 0, 0); err != nil {
		t.Fatalf("Failed to add bookmark: %v", err)
	}

	// Get bookmarks.
	bookmarks1 := c.Bookmarks()
	if len(bookmarks1) != 1 {
		t.Fatalf("Expected 1 bookmark, got %d", len(bookmarks1))
	}

	// Modify the returned slice.
	bookmarks1[0].Title = "Modified"

	// Get bookmarks again.
	bookmarks2 := c.Bookmarks()
	if bookmarks2[0].Title != "Chapter 1" {
		t.Error("Bookmarks() did not return a copy (modification affected original)")
	}
}

// TestBookmarks_PreservesOrder tests that bookmarks maintain order.
func TestBookmarks_PreservesOrder(t *testing.T) {
	c := New()

	// Add pages.
	for i := 0; i < 5; i++ {
		if _, err := c.NewPage(); err != nil {
			t.Fatalf("Failed to add page %d: %v", i, err)
		}
	}

	// Add bookmarks in specific order.
	bookmarkTitles := []string{
		"Introduction",
		"Chapter 1",
		"Section 1.1",
		"Section 1.2",
		"Chapter 2",
	}

	for i, title := range bookmarkTitles {
		if err := c.AddBookmark(title, i%5, 0); err != nil {
			t.Fatalf("Failed to add bookmark %q: %v", title, err)
		}
	}

	// Verify order is preserved.
	bookmarks := c.Bookmarks()
	if len(bookmarks) != len(bookmarkTitles) {
		t.Fatalf("Expected %d bookmarks, got %d", len(bookmarkTitles), len(bookmarks))
	}

	for i, expected := range bookmarkTitles {
		if bookmarks[i].Title != expected {
			t.Errorf("Bookmark %d: expected title %q, got %q",
				i, expected, bookmarks[i].Title)
		}
	}
}

// TestBookmark_Structure tests the Bookmark struct.
func TestBookmark_Structure(t *testing.T) {
	c := New()

	// Add pages.
	for i := 0; i < 3; i++ {
		if _, err := c.NewPage(); err != nil {
			t.Fatalf("Failed to add page %d: %v", i, err)
		}
	}

	// Add bookmark with specific values.
	title := "Test Chapter"
	pageIndex := 1
	level := 2

	err := c.AddBookmark(title, pageIndex, level)
	if err != nil {
		t.Fatalf("Failed to add bookmark: %v", err)
	}

	// Verify bookmark structure.
	bookmarks := c.Bookmarks()
	if len(bookmarks) != 1 {
		t.Fatalf("Expected 1 bookmark, got %d", len(bookmarks))
	}

	bookmark := bookmarks[0]
	if bookmark.Title != title {
		t.Errorf("Expected title %q, got %q", title, bookmark.Title)
	}
	if bookmark.PageIndex != pageIndex {
		t.Errorf("Expected pageIndex %d, got %d", pageIndex, bookmark.PageIndex)
	}
	if bookmark.Level != level {
		t.Errorf("Expected level %d, got %d", level, bookmark.Level)
	}
}

// TestAddBookmark_HierarchicalStructure tests nested bookmark hierarchy.
func TestAddBookmark_HierarchicalStructure(t *testing.T) {
	c := New()

	// Add 5 pages.
	for i := 0; i < 5; i++ {
		if _, err := c.NewPage(); err != nil {
			t.Fatalf("Failed to add page %d: %v", i, err)
		}
	}

	// Build hierarchical structure:
	// Chapter 1 (level 0)
	//   Section 1.1 (level 1)
	//     Subsection 1.1.1 (level 2)
	//   Section 1.2 (level 1)
	// Chapter 2 (level 0)

	type bookmark struct {
		title     string
		pageIndex int
		level     int
	}

	hierarchy := []bookmark{
		{"Chapter 1", 0, 0},
		{"Section 1.1", 0, 1},
		{"Subsection 1.1.1", 1, 2},
		{"Section 1.2", 2, 1},
		{"Chapter 2", 3, 0},
	}

	for _, b := range hierarchy {
		err := c.AddBookmark(b.title, b.pageIndex, b.level)
		if err != nil {
			t.Fatalf("Failed to add bookmark %q: %v", b.title, err)
		}
	}

	// Verify hierarchy.
	bookmarks := c.Bookmarks()
	if len(bookmarks) != len(hierarchy) {
		t.Fatalf("Expected %d bookmarks, got %d", len(hierarchy), len(bookmarks))
	}

	for i, expected := range hierarchy {
		if bookmarks[i].Title != expected.title {
			t.Errorf("Bookmark %d: expected title %q, got %q",
				i, expected.title, bookmarks[i].Title)
		}
		if bookmarks[i].PageIndex != expected.pageIndex {
			t.Errorf("Bookmark %d: expected pageIndex %d, got %d",
				i, expected.pageIndex, bookmarks[i].PageIndex)
		}
		if bookmarks[i].Level != expected.level {
			t.Errorf("Bookmark %d: expected level %d, got %d",
				i, expected.level, bookmarks[i].Level)
		}
	}
}

// TestBookmarks_EmptyByDefault tests that new Creator has no bookmarks.
func TestBookmarks_EmptyByDefault(t *testing.T) {
	c := New()

	bookmarks := c.Bookmarks()
	if len(bookmarks) != 0 {
		t.Errorf("Expected 0 bookmarks in new Creator, got %d", len(bookmarks))
	}
}

// TestAddBookmark_MultipleCallsSameLevel tests multiple bookmarks at same level.
func TestAddBookmark_MultipleCallsSameLevel(t *testing.T) {
	c := New()

	// Add pages.
	for i := 0; i < 10; i++ {
		if _, err := c.NewPage(); err != nil {
			t.Fatalf("Failed to add page %d: %v", i, err)
		}
	}

	// Add multiple top-level bookmarks.
	for i := 0; i < 10; i++ {
		title := "Chapter " + string(rune('A'+i))
		if err := c.AddBookmark(title, i, 0); err != nil {
			t.Fatalf("Failed to add bookmark %d: %v", i, err)
		}
	}

	// Verify all were added.
	bookmarks := c.Bookmarks()
	if len(bookmarks) != 10 {
		t.Errorf("Expected 10 bookmarks, got %d", len(bookmarks))
	}

	// All should be level 0.
	for i, b := range bookmarks {
		if b.Level != 0 {
			t.Errorf("Bookmark %d: expected level 0, got %d", i, b.Level)
		}
	}
}
