# Image Embedding Example

This example demonstrates how to embed JPEG images in PDF documents using GxPDF.

## Features Demonstrated

1. **Basic Image Embedding**
   - Load JPEG images from file or io.Reader
   - Draw images at specific positions and sizes

2. **Image Scaling**
   - Fixed size embedding with `DrawImage()`
   - Aspect-ratio-preserving scaling with `DrawImageFit()`

3. **Multiple Images**
   - Embed multiple images on the same page
   - Different sizes and positions

## Usage

```bash
go run main.go
```

This will create `image_embedding_example.pdf` with embedded JPEG images.

## API Reference

### Load Image

```go
// From file
img, err := creator.LoadImage("photo.jpg")

// From io.Reader
file, _ := os.Open("photo.jpg")
img, err := creator.LoadImageFromReader(file)
```

### Draw Image

```go
// Fixed size (may stretch image)
page.DrawImage(img, x, y, width, height)

// Fit to box (maintains aspect ratio)
page.DrawImageFit(img, x, y, maxWidth, maxHeight)
```

## Supported Formats

Currently supported:
- **JPEG** (RGB and CMYK color spaces)

Not yet supported:
- PNG (planned)
- GIF (planned)
- TIFF (planned)

## Technical Details

JPEG images are embedded directly in PDF using DCTDecode filter:
- No re-encoding (preserves original JPEG quality)
- Efficient storage (JPEG compression maintained)
- Fast rendering (PDF readers handle DCTDecode natively)

## Example Output

The generated PDF contains:
- Title text
- Fixed-size image (200×150 pts)
- Aspect-ratio-fit image (150×150 pts box)
- Multiple images at different scales (80×60, 120×90, 160×120 pts)
