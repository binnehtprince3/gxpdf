package creator

import (
	"fmt"
	"strings"

	"github.com/coregx/gxpdf/internal/fonts"
)

// TOC represents a Table of Contents for the document.
//
// The TOC is automatically generated from document chapters and includes
// clickable links to each chapter/section.
//
// Example:
//
//	c := creator.New()
//	c.EnableTOC()
//	// Add chapters...
//	c.WriteToFile("document.pdf")  // TOC is automatically generated
type TOC struct {
	// Title of the TOC (default: "Table of Contents")
	title string

	// Chapters to include in the TOC
	chapters []*Chapter

	// Style options
	style TOCStyle

	// Whether to show page numbers
	showPageNumbers bool

	// Leader character between title and page number (default: ".")
	leader string
}

// TOCStyle defines the visual style for the Table of Contents.
type TOCStyle struct {
	// TitleFont for the "Table of Contents" heading
	TitleFont FontName

	// TitleSize for the heading
	TitleSize float64

	// TitleColor for the heading
	TitleColor Color

	// EntryFont for TOC entries
	EntryFont FontName

	// EntrySize for TOC entries (base size, reduced for sub-levels)
	EntrySize float64

	// EntryColor for TOC entries
	EntryColor Color

	// IndentPerLevel is the indentation in points for each sub-level
	IndentPerLevel float64

	// LineSpacing between TOC entries
	LineSpacing float64

	// SpaceAfterTitle is the space after the TOC heading
	SpaceAfterTitle float64

	// LeaderFont for the leader dots (...)
	LeaderFont FontName

	// LeaderSize for the leader dots
	LeaderSize float64

	// LeaderColor for the leader dots
	LeaderColor Color
}

// NewTOC creates a new Table of Contents.
//
// The TOC will be populated with chapters when the document is rendered.
//
// Example:
//
//	toc := NewTOC()
//	toc.SetTitle("Contents")
func NewTOC() *TOC {
	return &TOC{
		title:           "Table of Contents",
		chapters:        make([]*Chapter, 0),
		style:           DefaultTOCStyle(),
		showPageNumbers: true,
		leader:          ".",
	}
}

// DefaultTOCStyle returns the default TOC style.
//
// Default style:
//   - TitleFont: HelveticaBold, 24pt, Black
//   - EntryFont: Helvetica, 12pt, Black
//   - IndentPerLevel: 20pt
//   - LineSpacing: 1.5
//   - SpaceAfterTitle: 20pt
//   - LeaderFont: Helvetica, 12pt, Gray
func DefaultTOCStyle() TOCStyle {
	return TOCStyle{
		TitleFont:       HelveticaBold,
		TitleSize:       24,
		TitleColor:      Black,
		EntryFont:       Helvetica,
		EntrySize:       12,
		EntryColor:      Black,
		IndentPerLevel:  20,
		LineSpacing:     1.5,
		SpaceAfterTitle: 20,
		LeaderFont:      Helvetica,
		LeaderSize:      12,
		LeaderColor:     Gray,
	}
}

// SetTitle sets the TOC heading title.
func (t *TOC) SetTitle(title string) {
	t.title = title
}

// Title returns the TOC heading title.
func (t *TOC) Title() string {
	return t.title
}

// SetStyle sets the TOC visual style.
func (t *TOC) SetStyle(style TOCStyle) {
	t.style = style
}

// Style returns the current TOC style.
func (t *TOC) Style() TOCStyle {
	return t.style
}

// SetShowPageNumbers enables/disables page number display.
func (t *TOC) SetShowPageNumbers(show bool) {
	t.showPageNumbers = show
}

// ShowPageNumbers returns whether page numbers are displayed.
func (t *TOC) ShowPageNumbers() bool {
	return t.showPageNumbers
}

// SetLeader sets the leader character (default: ".").
//
// Example:
//
//	toc.SetLeader(".")  // Introduction .............. 1
//	toc.SetLeader("-")  // Introduction ------------- 1
func (t *TOC) SetLeader(leader string) {
	t.leader = leader
}

// Leader returns the current leader character.
func (t *TOC) Leader() string {
	return t.leader
}

// setChapters sets the chapters to include in the TOC.
//
// This is called internally by the Creator when rendering.
func (t *TOC) setChapters(chapters []*Chapter) {
	t.chapters = chapters
}

// Height calculates the total height needed for the TOC.
func (t *TOC) Height(ctx *LayoutContext) float64 {
	height := t.style.TitleSize * 1.2 // Title height
	height += t.style.SpaceAfterTitle

	// Calculate height for each chapter entry
	lineHeight := t.style.EntrySize * t.style.LineSpacing
	for _, ch := range t.chapters {
		allChapters := ch.GetAllChapters()
		height += float64(len(allChapters)) * lineHeight
	}

	return height
}

// Draw renders the Table of Contents.
func (t *TOC) Draw(ctx *LayoutContext, page *Page) error {
	// Draw TOC title
	if err := t.drawTitle(ctx, page); err != nil {
		return fmt.Errorf("failed to draw TOC title: %w", err)
	}

	// Draw TOC entries
	for _, ch := range t.chapters {
		if err := t.drawChapterEntries(ctx, page, ch); err != nil {
			return fmt.Errorf("failed to draw TOC entries: %w", err)
		}
	}

	return nil
}

// drawTitle renders the TOC heading.
func (t *TOC) drawTitle(ctx *LayoutContext, page *Page) error {
	title := NewParagraph(t.title)
	title.SetFont(t.style.TitleFont, t.style.TitleSize)
	title.SetColor(t.style.TitleColor)
	title.SetAlignment(AlignCenter)

	if err := title.Draw(ctx, page); err != nil {
		return err
	}

	ctx.MoveCursor(0, t.style.SpaceAfterTitle)
	return nil
}

// drawChapterEntries renders TOC entries for a chapter and sub-chapters.
func (t *TOC) drawChapterEntries(ctx *LayoutContext, page *Page, ch *Chapter) error {
	allChapters := ch.GetAllChapters()
	for _, chapter := range allChapters {
		if err := t.drawEntry(ctx, page, chapter); err != nil {
			return err
		}
	}
	return nil
}

// drawEntry renders a single TOC entry with link.
func (t *TOC) drawEntry(ctx *LayoutContext, page *Page, ch *Chapter) error {
	level := ch.Level()
	indent := float64(level) * t.style.IndentPerLevel

	// Calculate entry font size (reduce for deeper levels)
	fontSize := t.style.EntrySize
	if level > 0 {
		fontSize *= (1.0 - float64(level)*0.05) // 5% smaller per level
		if fontSize < 8 {
			fontSize = 8 // Minimum font size
		}
	}

	// Calculate positions
	x := ctx.ContentLeft() + indent
	y := ctx.CurrentPDFY()

	// Build entry text
	entryText := ch.FullTitle()

	// Add page number and leader if enabled
	if t.showPageNumbers && ch.PageIndex() >= 0 {
		pageNum := ch.PageIndex() + 1 // Convert to 1-based
		entryText = t.buildEntryWithLeader(entryText, pageNum, fontSize, indent, ctx)
	}

	// Create link style matching entry style
	linkStyle := LinkStyle{
		Font:      t.style.EntryFont,
		Size:      fontSize,
		Color:     t.style.EntryColor,
		Underline: false, // No underline for TOC entries
	}

	// Add as internal link to the chapter's page
	if ch.PageIndex() >= 0 {
		if err := page.addLinkWithStyle(entryText, "", ch.PageIndex(), true, x, y, linkStyle); err != nil {
			return err
		}
	} else {
		// If page not set yet, just render as text
		if err := page.AddTextColor(entryText, x, y, t.style.EntryFont, fontSize, t.style.EntryColor); err != nil {
			return err
		}
	}

	// Advance cursor
	lineHeight := fontSize * t.style.LineSpacing
	ctx.MoveCursor(0, lineHeight)

	return nil
}

// buildEntryWithLeader builds a TOC entry with leader dots and page number.
//
// Example: "Introduction .............. 1".
func (t *TOC) buildEntryWithLeader(title string, pageNum int, fontSize, indent float64, ctx *LayoutContext) string {
	// Measure title width
	titleWidth := fonts.MeasureString(string(t.style.EntryFont), title, fontSize)

	// Measure page number width
	pageStr := fmt.Sprintf("%d", pageNum)
	pageWidth := fonts.MeasureString(string(t.style.EntryFont), pageStr, fontSize)

	// Measure leader character width
	leaderWidth := fonts.MeasureString(string(t.style.LeaderFont), t.leader, t.style.LeaderSize)

	// Calculate available space for leader
	availableWidth := ctx.AvailableWidth() - indent - titleWidth - pageWidth - 10 // 10pt padding

	// Calculate number of leader characters
	numLeaders := int(availableWidth / leaderWidth)
	if numLeaders < 3 {
		numLeaders = 3 // Minimum 3 leaders
	}

	// Build entry with leaders and page number
	leaders := strings.Repeat(t.leader, numLeaders)
	return fmt.Sprintf("%s %s %s", title, leaders, pageStr)
}

// TOCEntry represents a single entry in the Table of Contents.
//
// This is used internally for rendering and testing.
type TOCEntry struct {
	// Title of the chapter/section
	Title string

	// Number of the chapter (e.g., "1.2.3")
	Number string

	// Level of nesting (0 = top-level)
	Level int

	// PageIndex where the chapter starts (0-based)
	PageIndex int
}

// GetEntries returns all TOC entries in document order.
//
// This is useful for custom TOC rendering or testing.
func (t *TOC) GetEntries() []TOCEntry {
	entries := make([]TOCEntry, 0)
	for _, ch := range t.chapters {
		entries = append(entries, t.getChapterEntries(ch)...)
	}
	return entries
}

// getChapterEntries recursively builds TOC entries for a chapter.
func (t *TOC) getChapterEntries(ch *Chapter) []TOCEntry {
	allChapters := ch.GetAllChapters()
	entries := make([]TOCEntry, 0, len(allChapters))

	for _, chapter := range allChapters {
		entry := TOCEntry{
			Title:     chapter.Title(),
			Number:    chapter.NumberString(),
			Level:     chapter.Level(),
			PageIndex: chapter.PageIndex(),
		}
		entries = append(entries, entry)
	}

	return entries
}
