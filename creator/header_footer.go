package creator

// HeaderFunctionArgs contains information passed to the header function.
//
// This struct provides context about the current page and a Block to draw
// the header content into. The header function is called once for each page.
//
// Example:
//
//	c.SetHeaderFunc(func(args HeaderFunctionArgs) {
//	    p := NewParagraph(fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages))
//	    p.SetAlignment(AlignRight)
//	    args.Block.Draw(p)
//	})
type HeaderFunctionArgs struct {
	// PageNum is the current page number (1-based).
	PageNum int

	// TotalPages is the total number of pages in the document.
	// This is known because headers are rendered after all content is generated.
	TotalPages int

	// PageWidth is the page width in points.
	PageWidth float64

	// PageHeight is the page height in points.
	PageHeight float64

	// Block is the block to draw header content into.
	// The block is positioned at the top of the page within the margins.
	Block *Block
}

// FooterFunctionArgs contains information passed to the footer function.
//
// This struct provides context about the current page and a Block to draw
// the footer content into. The footer function is called once for each page.
//
// Example:
//
//	c.SetFooterFunc(func(args FooterFunctionArgs) {
//	    p := NewParagraph(fmt.Sprintf("Page %d", args.PageNum))
//	    p.SetAlignment(AlignCenter)
//	    args.Block.Draw(p)
//	})
type FooterFunctionArgs struct {
	// PageNum is the current page number (1-based).
	PageNum int

	// TotalPages is the total number of pages in the document.
	// This is known because footers are rendered after all content is generated.
	TotalPages int

	// PageWidth is the page width in points.
	PageWidth float64

	// PageHeight is the page height in points.
	PageHeight float64

	// Block is the block to draw footer content into.
	// The block is positioned at the bottom of the page within the margins.
	Block *Block
}

// HeaderFunc is the function signature for header rendering.
//
// The function receives a HeaderFunctionArgs struct containing page information
// and a Block to draw the header content into. The function should add any
// desired content to the Block using Draw() or DrawAt().
//
// Example:
//
//	var headerFunc HeaderFunc = func(args HeaderFunctionArgs) {
//	    p := NewParagraph("Company Name")
//	    p.SetFont(HelveticaBold, 10)
//	    args.Block.Draw(p)
//	}
type HeaderFunc func(args HeaderFunctionArgs)

// FooterFunc is the function signature for footer rendering.
//
// The function receives a FooterFunctionArgs struct containing page information
// and a Block to draw the footer content into. The function should add any
// desired content to the Block using Draw() or DrawAt().
//
// Example:
//
//	var footerFunc FooterFunc = func(args FooterFunctionArgs) {
//	    text := fmt.Sprintf("Page %d of %d", args.PageNum, args.TotalPages)
//	    p := NewParagraph(text)
//	    p.SetAlignment(AlignCenter)
//	    args.Block.Draw(p)
//	}
type FooterFunc func(args FooterFunctionArgs)

// Default header and footer heights in points.
const (
	// DefaultHeaderHeight is the default height for headers (50 points).
	DefaultHeaderHeight = 50.0

	// DefaultFooterHeight is the default height for footers (30 points).
	DefaultFooterHeight = 30.0
)
