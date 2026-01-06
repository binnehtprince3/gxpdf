package creator

import (
	"errors"
	"fmt"
)

// Bookmark represents a PDF bookmark (also known as outline item).
//
// Bookmarks provide a navigational tree structure in PDF documents.
// They allow users to jump to specific pages and create a table of contents.
//
// Example:
//
//	bookmark := Bookmark{
//	    Title:     "Chapter 1",
//	    PageIndex: 0,  // First page (0-based)
//	    Level:     0,  // Top-level bookmark
//	}
type Bookmark struct {
	// Title is the text displayed in the bookmark tree.
	Title string

	// PageIndex is the target page (0-based index).
	// 0 = first page, 1 = second page, etc.
	PageIndex int

	// Level is the nesting level in the bookmark hierarchy.
	// 0 = top-level, 1 = child of top-level, 2 = grandchild, etc.
	Level int
}

// AddBookmark adds a bookmark to the document.
//
// Bookmarks create a navigational tree structure (outline) in the PDF.
// Use Level to create nested bookmarks (chapters → sections → subsections).
//
// Parameters:
//   - title: Text to display in the bookmark tree
//   - pageIndex: Target page (0-based: 0 = first page, 1 = second, etc.)
//   - level: Nesting level (0 = top-level, 1 = child, 2 = grandchild, etc.)
//
// Returns an error if the parameters are invalid.
//
// Example:
//
//	c := creator.New()
//	page1, _ := c.NewPage()
//	page2, _ := c.NewPage()
//	page3, _ := c.NewPage()
//
//	// Add top-level bookmarks
//	c.AddBookmark("Chapter 1", 0, 0)  // Points to page 1
//	c.AddBookmark("Chapter 2", 2, 0)  // Points to page 3
//
//	// Add nested bookmarks (sections under Chapter 1)
//	c.AddBookmark("Section 1.1", 0, 1)  // Child of Chapter 1
//	c.AddBookmark("Section 1.2", 1, 1)  // Child of Chapter 1
//
//	c.WriteToFile("document.pdf")
func (c *Creator) AddBookmark(title string, pageIndex int, level int) error {
	// Validate title.
	if title == "" {
		return ErrEmptyBookmarkTitle
	}

	// Validate page index.
	if pageIndex < 0 {
		return fmt.Errorf("%w: pageIndex must be >= 0, got %d",
			ErrInvalidBookmarkPage, pageIndex)
	}

	// Validate level.
	if level < 0 {
		return fmt.Errorf("%w: level must be >= 0, got %d",
			ErrInvalidBookmarkLevel, level)
	}

	// Create bookmark.
	bookmark := Bookmark{
		Title:     title,
		PageIndex: pageIndex,
		Level:     level,
	}

	// Add to document's bookmark list.
	c.bookmarks = append(c.bookmarks, bookmark)

	return nil
}

// Bookmarks returns a copy of all bookmarks in the document.
//
// The returned slice is a copy, so modifications won't affect the document.
// Bookmarks are returned in the order they were added.
//
// Example:
//
//	c.AddBookmark("Chapter 1", 0, 0)
//	c.AddBookmark("Section 1.1", 0, 1)
//
//	bookmarks := c.Bookmarks()
//	fmt.Printf("Document has %d bookmarks\n", len(bookmarks))
func (c *Creator) Bookmarks() []Bookmark {
	// Return a copy to prevent external modifications.
	result := make([]Bookmark, len(c.bookmarks))
	copy(result, c.bookmarks)
	return result
}

// Bookmark-related errors.
var (
	// ErrEmptyBookmarkTitle is returned when bookmark title is empty.
	ErrEmptyBookmarkTitle = errors.New("bookmark title cannot be empty")

	// ErrInvalidBookmarkPage is returned when page index is invalid.
	ErrInvalidBookmarkPage = errors.New("invalid bookmark page index")

	// ErrInvalidBookmarkLevel is returned when level is invalid.
	ErrInvalidBookmarkLevel = errors.New("invalid bookmark level")
)
