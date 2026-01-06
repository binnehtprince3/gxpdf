# Image Extraction Example

This example demonstrates how to extract images from PDF documents using GxPDF.

## Features Demonstrated

1. **Extract all images from a document**
   - Get all images from all pages
   - Save to disk automatically

2. **Extract images from specific page**
   - Target specific pages
   - Control output location

3. **Process image metadata**
   - Get image dimensions
   - Check color space and encoding
   - Access raw image data

## Usage

```bash
# Extract all images from a PDF
go run main.go document.pdf

# Extract to specific directory
go run main.go document.pdf ./extracted_images

# Help
go run main.go
```

## API Reference

### Document-level Extraction

```go
doc, _ := gxpdf.Open("document.pdf")
defer doc.Close()

// Get all images from all pages
images := doc.GetImages()

// With error handling
images, err := doc.GetImagesWithError()
```

### Page-level Extraction

```go
page := doc.Page(0) // First page

// Get all images from this page
images := page.GetImages()

// With error handling
images, err := page.GetImagesWithError()
```

### Image Operations

```go
for _, img := range images {
    // Get metadata
    fmt.Printf("%dx%d, %s\n", img.Width(), img.Height(), img.ColorSpace())

    // Save to file (format determined by extension)
    img.SaveToFile("output.jpg")  // JPEG
    img.SaveToFile("output.png")  // PNG

    // Convert to Go image for processing
    goImg, _ := img.ToGoImage()
}
```

## Supported Image Formats

GxPDF can extract images with the following encodings:

- **DCTDecode** (JPEG) - Direct extraction, no re-encoding
- **FlateDecode** (zlib) - Decompressed and converted
- **Uncompressed** - Direct extraction

## Color Spaces

- **DeviceRGB** - RGB color images
- **DeviceGray** - Grayscale images
- **DeviceCMYK** - CMYK images (conversion to RGB coming soon)
- **Indexed** - Palette-based images (expansion coming soon)

## Output Formats

Images can be saved as:

- **JPEG** (.jpg, .jpeg) - Best for photos and DCTDecode images
- **PNG** (.png) - Best for lossless images and transparency

For DCTDecode images saved as JPEG, the original compressed data is used
directly, preserving quality without re-encoding.

## Example Output

```
=== Example 1: Extract All Images from Document ===
Found 5 images in document
  Saved: image_0.jpg (800x600, DeviceRGB, /DCTDecode)
  Saved: image_1.jpg (1024x768, DeviceRGB, /DCTDecode)
  Saved: image_2.png (200x200, DeviceGray, /FlateDecode)
  Saved: image_3.jpg (640x480, DeviceRGB, /DCTDecode)
  Saved: image_4.png (100x100, DeviceRGB, /FlateDecode)

=== Example 2: Extract Images from Specific Page ===
Found 2 images on page 1
  Saved: page1_image_0.jpg (800x600)
  Saved: page1_image_1.jpg (1024x768)

=== Example 3: Process Images (Metadata Only) ===
Page 1:
  Image 0:
    Name: /Im1
    Dimensions: 800x600 pixels
    Color Space: DeviceRGB
    Bits per Component: 8
    Filter: /DCTDecode
    Go Image Bounds: (0,0)-(800,600)
```

## Notes

- Images are extracted in the order they appear in the PDF's resource dictionary
- Inline images are not yet supported (coming soon)
- CMYK to RGB conversion is not yet implemented
- Indexed color space expansion is not yet implemented

## See Also

- [Image Embedding Example](../image_embedding/) - How to add images to PDFs
- [GxPDF Documentation](../../README.md) - Full library documentation
