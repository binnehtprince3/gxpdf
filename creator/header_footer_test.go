package creator

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBlock_Creation(t *testing.T) {
	tests := []struct {
		name   string
		width  float64
		height float64
	}{
		{"small block", 100, 50},
		{"large block", 500, 200},
		{"zero dimensions", 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			block := NewBlock(tt.width, tt.height)

			assert.NotNil(t, block)
			assert.Equal(t, tt.width, block.Width())
			assert.Equal(t, tt.height, block.Height())
			assert.Empty(t, block.GetDrawables())
		})
	}
}

func TestBlock_SetDimensions(t *testing.T) {
	block := NewBlock(100, 50)

	block.SetWidth(200)
	block.SetHeight(100)

	assert.Equal(t, 200.0, block.Width())
	assert.Equal(t, 100.0, block.Height())
}

func TestBlock_Margins(t *testing.T) {
	block := NewBlock(500, 100)

	// Initially no margins.
	margins := block.Margins()
	assert.Equal(t, 0.0, margins.Top)
	assert.Equal(t, 0.0, margins.Right)
	assert.Equal(t, 0.0, margins.Bottom)
	assert.Equal(t, 0.0, margins.Left)

	// Set margins.
	block.SetMargins(Margins{Top: 10, Right: 20, Bottom: 10, Left: 20})

	margins = block.Margins()
	assert.Equal(t, 10.0, margins.Top)
	assert.Equal(t, 20.0, margins.Right)
	assert.Equal(t, 10.0, margins.Bottom)
	assert.Equal(t, 20.0, margins.Left)
}

func TestBlock_ContentDimensions(t *testing.T) {
	block := NewBlock(500, 100)
	block.SetMargins(Margins{Top: 10, Right: 20, Bottom: 10, Left: 20})

	// Content width = 500 - 20 - 20 = 460.
	assert.Equal(t, 460.0, block.ContentWidth())

	// Content height = 100 - 10 - 10 = 80.
	assert.Equal(t, 80.0, block.ContentHeight())
}

func TestBlock_Draw(t *testing.T) {
	block := NewBlock(500, 100)

	// Drawing nil should fail.
	err := block.Draw(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "nil")

	// Drawing a paragraph should work.
	p := NewParagraph("Test")
	err = block.Draw(p)
	require.NoError(t, err)

	drawables := block.GetDrawables()
	assert.Len(t, drawables, 1)
	assert.Equal(t, p, drawables[0].Drawable)
}

func TestBlock_DrawAt(t *testing.T) {
	block := NewBlock(500, 100)

	// Drawing nil should fail.
	err := block.DrawAt(nil, 10, 20)
	assert.Error(t, err)

	// Drawing at specific position.
	p := NewParagraph("Test")
	err = block.DrawAt(p, 50, 25)
	require.NoError(t, err)

	drawables := block.GetDrawables()
	assert.Len(t, drawables, 1)
	assert.Equal(t, 50.0, drawables[0].X)
	assert.Equal(t, 25.0, drawables[0].Y)
}

func TestBlock_Clear(t *testing.T) {
	block := NewBlock(500, 100)

	// Add some drawables.
	_ = block.Draw(NewParagraph("Test 1"))
	_ = block.Draw(NewParagraph("Test 2"))
	assert.Len(t, block.GetDrawables(), 2)

	// Clear.
	block.Clear()
	assert.Empty(t, block.GetDrawables())
}

func TestBlock_Cursor(t *testing.T) {
	block := NewBlock(500, 100)

	// Initial cursor at origin.
	block.SetCursor(50, 25)
	block.MoveCursor(10, 5)

	// Draw should use cursor position (adjusted for margins).
	p := NewParagraph("Test")
	_ = block.Draw(p)

	drawables := block.GetDrawables()
	assert.Len(t, drawables, 1)
	// Position should be cursor + margins (0 in this case).
	assert.Equal(t, 60.0, drawables[0].X)
	assert.Equal(t, 30.0, drawables[0].Y)
}

func TestBlock_GetLayoutContext(t *testing.T) {
	block := NewBlock(500, 100)
	block.SetMargins(Margins{Top: 10, Right: 20, Bottom: 10, Left: 20})

	ctx := block.GetLayoutContext()

	assert.Equal(t, 500.0, ctx.PageWidth)
	assert.Equal(t, 100.0, ctx.PageHeight)
	assert.Equal(t, 20.0, ctx.Margins.Left)
	assert.Equal(t, 20.0, ctx.CursorX)
	assert.Equal(t, 0.0, ctx.CursorY)
}

func TestHeaderFunctionArgs(t *testing.T) {
	block := NewBlock(500, 50)
	args := HeaderFunctionArgs{
		PageNum:    1,
		TotalPages: 5,
		PageWidth:  595,
		PageHeight: 842,
		Block:      block,
	}

	assert.Equal(t, 1, args.PageNum)
	assert.Equal(t, 5, args.TotalPages)
	assert.Equal(t, 595.0, args.PageWidth)
	assert.Equal(t, 842.0, args.PageHeight)
	assert.Same(t, block, args.Block)
}

func TestFooterFunctionArgs(t *testing.T) {
	block := NewBlock(500, 30)
	args := FooterFunctionArgs{
		PageNum:    3,
		TotalPages: 10,
		PageWidth:  612,
		PageHeight: 792,
		Block:      block,
	}

	assert.Equal(t, 3, args.PageNum)
	assert.Equal(t, 10, args.TotalPages)
	assert.Equal(t, 612.0, args.PageWidth)
	assert.Equal(t, 792.0, args.PageHeight)
	assert.Same(t, block, args.Block)
}

func TestCreator_SetHeaderFunc(t *testing.T) {
	c := New()

	// Initially nil.
	assert.Nil(t, c.headerFunc)

	// Set header function.
	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
		// Function body for test.
		_ = args.PageNum
	})

	assert.NotNil(t, c.headerFunc)
}

func TestCreator_SetFooterFunc(t *testing.T) {
	c := New()

	// Initially nil.
	assert.Nil(t, c.footerFunc)

	// Set footer function.
	c.SetFooterFunc(func(args FooterFunctionArgs) {
		// Function body for test.
		_ = args.PageNum
	})

	assert.NotNil(t, c.footerFunc)
}

func TestCreator_HeaderHeight(t *testing.T) {
	c := New()

	// Default height.
	assert.Equal(t, DefaultHeaderHeight, c.HeaderHeight())

	// Set custom height.
	c.SetHeaderHeight(40)
	assert.Equal(t, 40.0, c.HeaderHeight())
}

func TestCreator_FooterHeight(t *testing.T) {
	c := New()

	// Default height.
	assert.Equal(t, DefaultFooterHeight, c.FooterHeight())

	// Set custom height.
	c.SetFooterHeight(25)
	assert.Equal(t, 25.0, c.FooterHeight())
}

func TestCreator_SkipHeaderOnFirstPage(t *testing.T) {
	c := New()

	// Default is false.
	assert.False(t, c.SkipHeaderOnFirstPage())

	// Set to true.
	c.SetSkipHeaderOnFirstPage(true)
	assert.True(t, c.SkipHeaderOnFirstPage())

	// Set back to false.
	c.SetSkipHeaderOnFirstPage(false)
	assert.False(t, c.SkipHeaderOnFirstPage())
}

func TestCreator_SkipFooterOnFirstPage(t *testing.T) {
	c := New()

	// Default is false.
	assert.False(t, c.SkipFooterOnFirstPage())

	// Set to true.
	c.SetSkipFooterOnFirstPage(true)
	assert.True(t, c.SkipFooterOnFirstPage())
}

func TestCreator_HeaderFooter_MultiPage(t *testing.T) {
	c := New()

	// Track which pages have headers/footers.
	headerPages := make([]int, 0)
	footerPages := make([]int, 0)

	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
		headerPages = append(headerPages, args.PageNum)
		p := NewParagraph(fmt.Sprintf("Header - Page %d of %d", args.PageNum, args.TotalPages))
		_ = args.Block.Draw(p)
	})

	c.SetFooterFunc(func(args FooterFunctionArgs) {
		footerPages = append(footerPages, args.PageNum)
		p := NewParagraph(fmt.Sprintf("Footer - Page %d", args.PageNum))
		p.SetAlignment(AlignCenter)
		_ = args.Block.Draw(p)
	})

	// Create 3 pages.
	for i := 0; i < 3; i++ {
		page, err := c.NewPage()
		require.NoError(t, err)
		_ = page.AddText(fmt.Sprintf("Content on page %d", i+1), 100, 400, Helvetica, 12)
	}

	// Collect content to trigger header/footer rendering.
	textContents, _ := c.collectAllPageContents()

	// Verify headers were called for all pages.
	assert.Equal(t, []int{1, 2, 3}, headerPages)
	assert.Equal(t, []int{1, 2, 3}, footerPages)

	// Each page should have text content.
	assert.Len(t, textContents, 3)
}

func TestCreator_HeaderFooter_SkipFirst(t *testing.T) {
	c := New()

	headerPages := make([]int, 0)
	footerPages := make([]int, 0)

	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
		headerPages = append(headerPages, args.PageNum)
		_ = args.Block.Draw(NewParagraph("Header"))
	})

	c.SetFooterFunc(func(args FooterFunctionArgs) {
		footerPages = append(footerPages, args.PageNum)
		_ = args.Block.Draw(NewParagraph("Footer"))
	})

	// Skip header and footer on first page.
	c.SetSkipHeaderOnFirstPage(true)
	c.SetSkipFooterOnFirstPage(true)

	// Create 3 pages.
	for i := 0; i < 3; i++ {
		_, err := c.NewPage()
		require.NoError(t, err)
	}

	// Collect content.
	_, _ = c.collectAllPageContents()

	// Headers/footers should only be on pages 2 and 3.
	assert.Equal(t, []int{2, 3}, headerPages)
	assert.Equal(t, []int{2, 3}, footerPages)
}

func TestCreator_HeaderFooter_TotalPages(t *testing.T) {
	c := New()

	var capturedTotalPages int

	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
		capturedTotalPages = args.TotalPages
	})

	// Create 5 pages.
	for i := 0; i < 5; i++ {
		_, err := c.NewPage()
		require.NoError(t, err)
	}

	// Collect content.
	_, _ = c.collectAllPageContents()

	// Total pages should be 5.
	assert.Equal(t, 5, capturedTotalPages)
}

func TestCreator_HeaderWithAlignment(t *testing.T) {
	c := New()

	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
		// Left aligned title.
		title := NewParagraph("Document Title")
		title.SetFont(HelveticaBold, 10)
		_ = args.Block.Draw(title)

		// Right aligned page number.
		pageNum := NewParagraph(fmt.Sprintf("Page %d", args.PageNum))
		pageNum.SetAlignment(AlignRight)
		_ = args.Block.DrawAt(pageNum, 400, 0)
	})

	_, err := c.NewPage()
	require.NoError(t, err)

	textContents, _ := c.collectAllPageContents()

	// Should have text content.
	assert.NotEmpty(t, textContents[0])
}

func TestDefaultConstants(t *testing.T) {
	assert.Equal(t, 50.0, DefaultHeaderHeight)
	assert.Equal(t, 30.0, DefaultFooterHeight)
}

func TestCreator_NoHeaderFooter(t *testing.T) {
	c := New()

	// Create a page without setting header/footer functions.
	page, err := c.NewPage()
	require.NoError(t, err)
	_ = page.AddText("Hello World", 100, 700, Helvetica, 12)

	// Collect content - should work without header/footer.
	textContents, _ := c.collectAllPageContents()

	// Should have one text operation (the "Hello World").
	assert.Len(t, textContents[0], 1)
	assert.Equal(t, "Hello World", textContents[0][0].Text)
}
