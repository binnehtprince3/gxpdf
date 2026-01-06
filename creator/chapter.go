package creator

import (
	"errors"
	"fmt"
	"strings"
)

// Chapter represents a document chapter with title and content.
//
// Chapters provide hierarchical document structure with automatic numbering.
// They can contain sub-chapters (sections, subsections, etc.) up to any depth.
//
// Example:
//
//	ch := NewChapter("Introduction")
//	ch.Add(NewParagraph("This is the introduction..."))
//
//	sec := ch.NewSubChapter("Background")
//	sec.Add(NewParagraph("Background information..."))
type Chapter struct {
	// Title of the chapter
	title string

	// Number components (e.g., [1, 2, 3] for "1.2.3")
	number []int

	// Content elements (paragraphs, tables, etc.)
	content []Drawable

	// Sub-chapters (sections, subsections, etc.)
	subChapters []*Chapter

	// Parent chapter (nil for top-level)
	parent *Chapter

	// Page index where chapter starts (set during rendering)
	pageIndex int

	// Style options
	style ChapterStyle
}

// ChapterStyle defines the visual style for chapter headings.
type ChapterStyle struct {
	// Font for the chapter title
	Font FontName

	// FontSize for the chapter title
	FontSize float64

	// Color for the chapter title
	Color Color

	// SpaceBefore is the vertical space before the chapter heading
	SpaceBefore float64

	// SpaceAfter is the vertical space after the chapter heading
	SpaceAfter float64

	// ShowNumber indicates whether to show chapter numbers
	ShowNumber bool

	// NumberSeparator is the separator between number components (default: ".")
	NumberSeparator string
}

// NewChapter creates a new top-level chapter with the given title.
//
// The chapter uses default styling and can have content and sub-chapters added.
//
// Example:
//
//	ch := NewChapter("Chapter 1: Introduction")
//	ch.Add(NewParagraph("Welcome to this document..."))
func NewChapter(title string) *Chapter {
	return &Chapter{
		title:       title,
		number:      []int{},
		content:     make([]Drawable, 0),
		subChapters: make([]*Chapter, 0),
		parent:      nil,
		pageIndex:   -1,
		style:       DefaultChapterStyle(),
	}
}

// DefaultChapterStyle returns the default chapter heading style.
//
// Default style:
//   - Font: HelveticaBold
//   - FontSize: 18pt (scaled down for sub-chapters)
//   - Color: Black
//   - SpaceBefore: 20pt
//   - SpaceAfter: 10pt
//   - ShowNumber: true
//   - NumberSeparator: "."
func DefaultChapterStyle() ChapterStyle {
	return ChapterStyle{
		Font:            HelveticaBold,
		FontSize:        18,
		Color:           Black,
		SpaceBefore:     20,
		SpaceAfter:      10,
		ShowNumber:      true,
		NumberSeparator: ".",
	}
}

// Title returns the chapter title.
func (c *Chapter) Title() string {
	return c.title
}

// SetTitle sets the chapter title.
func (c *Chapter) SetTitle(title string) {
	c.title = title
}

// Number returns the chapter number components.
//
// For example, section 1.2.3 returns []int{1, 2, 3}.
func (c *Chapter) Number() []int {
	return c.number
}

// NumberString returns the formatted chapter number.
//
// For example: "1", "1.2", "1.2.3".
func (c *Chapter) NumberString() string {
	if len(c.number) == 0 {
		return ""
	}
	parts := make([]string, len(c.number))
	for i, n := range c.number {
		parts[i] = fmt.Sprintf("%d", n)
	}
	return strings.Join(parts, c.style.NumberSeparator)
}

// FullTitle returns the title with number prefix if numbering is enabled.
//
// For example: "1.2 Background" or just "Background" if ShowNumber is false.
func (c *Chapter) FullTitle() string {
	if !c.style.ShowNumber || len(c.number) == 0 {
		return c.title
	}
	return c.NumberString() + " " + c.title
}

// Level returns the nesting level of the chapter.
//
// Top-level chapters return 0, sub-chapters return 1, etc.
func (c *Chapter) Level() int {
	return len(c.number)
}

// SetStyle sets the chapter heading style.
func (c *Chapter) SetStyle(style ChapterStyle) {
	c.style = style
}

// Style returns the current chapter heading style.
func (c *Chapter) Style() ChapterStyle {
	return c.style
}

// Add adds a content element (paragraph, table, etc.) to the chapter.
//
// Example:
//
//	ch.Add(NewParagraph("This is the introduction..."))
//	ch.Add(NewTable())
func (c *Chapter) Add(d Drawable) error {
	if d == nil {
		return errors.New("cannot add nil drawable to chapter")
	}
	c.content = append(c.content, d)
	return nil
}

// Content returns all content elements in the chapter.
func (c *Chapter) Content() []Drawable {
	return c.content
}

// NewSubChapter creates a new sub-chapter (section) under this chapter.
//
// The sub-chapter inherits styling from the parent but with smaller font size.
// Numbering is automatically assigned based on the position.
//
// Example:
//
//	ch := NewChapter("Introduction")
//	sec1 := ch.NewSubChapter("Background")     // 1.1 Background
//	sec2 := ch.NewSubChapter("Motivation")     // 1.2 Motivation
//	subsec := sec1.NewSubChapter("History")    // 1.1.1 History
func (c *Chapter) NewSubChapter(title string) *Chapter {
	sub := NewChapter(title)
	sub.parent = c

	// Inherit style with reduced font size
	sub.style = c.style
	sub.style.FontSize = c.style.FontSize * 0.85 // 15% smaller
	sub.style.SpaceBefore = c.style.SpaceBefore * 0.75
	sub.style.SpaceAfter = c.style.SpaceAfter * 0.75

	// Assign number
	subNumber := len(c.subChapters) + 1
	sub.number = append(append([]int{}, c.number...), subNumber)

	c.subChapters = append(c.subChapters, sub)
	return sub
}

// SubChapters returns all sub-chapters.
func (c *Chapter) SubChapters() []*Chapter {
	return c.subChapters
}

// Parent returns the parent chapter (nil for top-level).
func (c *Chapter) Parent() *Chapter {
	return c.parent
}

// PageIndex returns the page index where the chapter starts.
//
// Returns -1 if not yet rendered.
func (c *Chapter) PageIndex() int {
	return c.pageIndex
}

// setPageIndex sets the page index (used internally during rendering).
func (c *Chapter) setPageIndex(index int) {
	c.pageIndex = index
}

// Height calculates the total height needed to render this chapter.
//
// This includes the heading, all content, and all sub-chapters.
func (c *Chapter) Height(ctx *LayoutContext) float64 {
	height := c.style.SpaceBefore
	height += c.style.FontSize * 1.2 // Heading height
	height += c.style.SpaceAfter

	// Add content height
	for _, d := range c.content {
		height += d.Height(ctx)
	}

	// Add sub-chapters height
	for _, sub := range c.subChapters {
		height += sub.Height(ctx)
	}

	return height
}

// Draw renders the chapter on the page.
//
// This renders:
// 1. Chapter heading (with number if enabled)
// 2. All content elements
// 3. All sub-chapters (recursively).
func (c *Chapter) Draw(ctx *LayoutContext, page *Page) error {
	// Draw chapter heading
	if err := c.drawHeading(ctx, page); err != nil {
		return fmt.Errorf("failed to draw chapter heading: %w", err)
	}

	// Draw content
	for _, d := range c.content {
		if err := d.Draw(ctx, page); err != nil {
			return fmt.Errorf("failed to draw chapter content: %w", err)
		}
	}

	// Draw sub-chapters
	for _, sub := range c.subChapters {
		if err := sub.Draw(ctx, page); err != nil {
			return fmt.Errorf("failed to draw sub-chapter: %w", err)
		}
	}

	return nil
}

// drawHeading renders the chapter heading.
func (c *Chapter) drawHeading(ctx *LayoutContext, page *Page) error {
	// Add space before heading
	ctx.MoveCursor(0, c.style.SpaceBefore)

	// Create heading paragraph
	heading := NewParagraph(c.FullTitle())
	heading.SetFont(c.style.Font, c.style.FontSize)
	heading.SetColor(c.style.Color)

	// Draw heading
	if err := heading.Draw(ctx, page); err != nil {
		return err
	}

	// Add space after heading
	ctx.MoveCursor(0, c.style.SpaceAfter)

	return nil
}

// GetAllChapters returns a flat list of this chapter and all sub-chapters.
//
// The list is in document order (depth-first traversal).
//
// This is useful for building table of contents.
func (c *Chapter) GetAllChapters() []*Chapter {
	result := []*Chapter{c}
	for _, sub := range c.subChapters {
		result = append(result, sub.GetAllChapters()...)
	}
	return result
}

// assignNumbers assigns numbers to this chapter and all sub-chapters.
//
// This is called internally when chapters are added to the document.
func (c *Chapter) assignNumbers(parentNumber []int, index int) {
	// Assign this chapter's number
	c.number = append(append([]int{}, parentNumber...), index+1)

	// Assign sub-chapter numbers recursively
	for i, sub := range c.subChapters {
		sub.assignNumbers(c.number, i)
	}
}
