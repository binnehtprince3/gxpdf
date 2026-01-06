package gxpdf

import (
	"image"

	"github.com/coregx/gxpdf/internal/models/types"
)

// Image represents an image extracted from a PDF.
//
// This is a thin wrapper around the internal Image value object,
// providing a clean public API.
//
// Example:
//
//	images := doc.GetImages()
//	for i, img := range images {
//	    fmt.Printf("Image %d: %dx%d\n", i, img.Width(), img.Height())
//	    img.SaveToFile(fmt.Sprintf("image_%d.jpg", i))
//	}
type Image struct {
	internal *types.Image
}

// Width returns the image width in pixels.
func (img *Image) Width() int {
	return img.internal.Width()
}

// Height returns the image height in pixels.
func (img *Image) Height() int {
	return img.internal.Height()
}

// ColorSpace returns the PDF color space name.
//
// Common values: "DeviceRGB", "DeviceGray", "DeviceCMYK", "Indexed"
func (img *Image) ColorSpace() string {
	return img.internal.ColorSpace()
}

// BitsPerComponent returns bits per color component (typically 8).
func (img *Image) BitsPerComponent() int {
	return img.internal.BitsPerComponent()
}

// Filter returns the original PDF filter used for compression.
//
// Common values: "/DCTDecode" (JPEG), "/FlateDecode" (zlib)
func (img *Image) Filter() string {
	return img.internal.Filter()
}

// Name returns the XObject name from the PDF (e.g., "/Im1").
func (img *Image) Name() string {
	return img.internal.Name()
}

// SaveToFile saves the image to a file.
//
// The file format is determined by the extension:
//   - .jpg, .jpeg: JPEG format (best for DCTDecode images)
//   - .png: PNG format (best for lossless images)
//
// For DCTDecode (JPEG) images, the original data is saved directly
// without re-encoding, preserving quality.
//
// Example:
//
//	err := img.SaveToFile("extracted_image.jpg")
//	if err != nil {
//	    log.Fatal(err)
//	}
func (img *Image) SaveToFile(path string) error {
	return img.internal.SaveToFile(path)
}

// ToGoImage converts the image to Go's standard image.Image.
//
// This is useful for further processing with Go's image libraries.
//
// Example:
//
//	goImg, err := img.ToGoImage()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Process with Go image libraries
//	resized := resize.Resize(100, 100, goImg, resize.Lanczos3)
func (img *Image) ToGoImage() (image.Image, error) {
	return img.internal.ToGoImage()
}

// String returns a string representation of the image.
func (img *Image) String() string {
	return img.internal.String()
}
