// Package valueobjects contains value objects used across the domain.
package types

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

// Image represents an extracted PDF image with metadata.
//
// This is a value object containing image data and properties.
// Images are immutable once created.
//
// Example:
//
//	img := NewImage(data, width, height, "DeviceRGB", 8, "/DCTDecode")
//	img.SaveToFile("output.jpg")
//	goImg, _ := img.ToGoImage()
type Image struct {
	// Raw image data (decoded pixel data or original compressed data)
	data []byte

	// Image metadata
	width            int
	height           int
	colorSpace       string // DeviceRGB, DeviceGray, DeviceCMYK, Indexed, etc.
	bitsPerComponent int
	filter           string // Original filter: /DCTDecode, /FlateDecode, etc.

	// Additional metadata
	name string // XObject name (e.g., "/Im1")
}

// NewImage creates a new Image value object.
//
// Parameters:
//   - data: Raw image data (decoded or compressed)
//   - width: Image width in pixels
//   - height: Image height in pixels
//   - colorSpace: PDF color space name
//   - bitsPerComponent: Bits per color component (typically 8)
//   - filter: Original PDF filter (e.g., "/DCTDecode")
//
// Returns error if parameters are invalid.
func NewImage(data []byte, width, height int, colorSpace string, bitsPerComponent int, filter string) (*Image, error) {
	if len(data) == 0 {
		return nil, ErrEmptyImageData
	}
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("%w: width=%d, height=%d", ErrInvalidImageDimensions, width, height)
	}
	if bitsPerComponent <= 0 || bitsPerComponent > 16 {
		return nil, fmt.Errorf("%w: %d", ErrInvalidBitsPerComponent, bitsPerComponent)
	}

	return &Image{
		data:             data,
		width:            width,
		height:           height,
		colorSpace:       colorSpace,
		bitsPerComponent: bitsPerComponent,
		filter:           filter,
	}, nil
}

// Data returns the raw image data.
func (img *Image) Data() []byte {
	// Return copy to maintain immutability
	result := make([]byte, len(img.data))
	copy(result, img.data)
	return result
}

// Width returns the image width in pixels.
func (img *Image) Width() int {
	return img.width
}

// Height returns the image height in pixels.
func (img *Image) Height() int {
	return img.height
}

// ColorSpace returns the PDF color space name.
func (img *Image) ColorSpace() string {
	return img.colorSpace
}

// BitsPerComponent returns bits per color component.
func (img *Image) BitsPerComponent() int {
	return img.bitsPerComponent
}

// Filter returns the original PDF filter.
func (img *Image) Filter() string {
	return img.filter
}

// Name returns the XObject name.
func (img *Image) Name() string {
	return img.name
}

// SetName sets the XObject name (for internal use).
func (img *Image) SetName(name string) {
	img.name = name
}

// SaveToFile saves the image to a file.
//
// The file format is determined by the extension:
//   - .jpg, .jpeg: JPEG format (best for DCTDecode images)
//   - .png: PNG format (best for lossless images)
//
// For DCTDecode images, the original data is saved directly.
// For other formats, the image is encoded to the target format.
//
// Parameters:
//   - path: Output file path
//
// Returns error if file cannot be written.
func (img *Image) SaveToFile(path string) error {
	ext := strings.ToLower(filepath.Ext(path))

	// For JPEG images from DCTDecode, save original data directly
	if (ext == ".jpg" || ext == ".jpeg") && img.filter == "/DCTDecode" {
		return img.saveRaw(path)
	}

	// For PNG or other formats, convert to Go image first
	goImg, err := img.ToGoImage()
	if err != nil {
		return fmt.Errorf("failed to convert to Go image: %w", err)
	}

	// Save based on extension
	switch ext {
	case ".jpg", ".jpeg":
		return img.saveAsJPEG(path, goImg)
	case ".png":
		return img.saveAsPNG(path, goImg)
	default:
		// Default to PNG for unknown extensions
		return img.saveAsPNG(path, goImg)
	}
}

// saveRaw saves the raw image data to file.
func (img *Image) saveRaw(path string) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	if _, err := file.Write(img.data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

// saveAsJPEG saves the image as JPEG.
func (img *Image) saveAsJPEG(path string, goImg image.Image) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	opts := &jpeg.Options{Quality: 90}
	if err := jpeg.Encode(file, goImg, opts); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

// saveAsPNG saves the image as PNG.
func (img *Image) saveAsPNG(path string, goImg image.Image) (err error) {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer func() {
		if closeErr := file.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close file: %w", closeErr)
		}
	}()

	if err := png.Encode(file, goImg); err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// ToGoImage converts the image to Go's image.Image.
//
// This is useful for further processing or saving in different formats.
//
// Returns error if conversion fails.
func (img *Image) ToGoImage() (image.Image, error) {
	// For DCTDecode (JPEG), decode directly
	if img.filter == "/DCTDecode" {
		return img.decodeJPEG()
	}

	// For other formats, build image from raw pixel data
	return img.buildGoImage()
}

// decodeJPEG decodes JPEG data to image.Image.
func (img *Image) decodeJPEG() (image.Image, error) {
	goImg, err := jpeg.Decode(bytes.NewReader(img.data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode JPEG: %w", err)
	}
	return goImg, nil
}

// buildGoImage builds a Go image from raw pixel data.
func (img *Image) buildGoImage() (image.Image, error) {
	// Determine color model based on color space
	switch img.colorSpace {
	case "DeviceGray", "/DeviceGray":
		return img.buildGrayImage()
	case "DeviceRGB", "/DeviceRGB":
		return img.buildRGBImage()
	case "DeviceCMYK", "/DeviceCMYK":
		return nil, fmt.Errorf("CMYK to RGB conversion not yet implemented")
	default:
		return nil, fmt.Errorf("unsupported color space: %s", img.colorSpace)
	}
}

// buildGrayImage builds a grayscale image.
func (img *Image) buildGrayImage() (image.Image, error) {
	expectedLen := img.width * img.height
	if len(img.data) < expectedLen {
		return nil, fmt.Errorf("insufficient data for grayscale image: expected %d bytes, got %d",
			expectedLen, len(img.data))
	}

	goImg := image.NewGray(image.Rect(0, 0, img.width, img.height))
	copy(goImg.Pix, img.data[:expectedLen])

	return goImg, nil
}

// buildRGBImage builds an RGB image.
func (img *Image) buildRGBImage() (image.Image, error) {
	expectedLen := img.width * img.height * 3
	if len(img.data) < expectedLen {
		return nil, fmt.Errorf("insufficient data for RGB image: expected %d bytes, got %d",
			expectedLen, len(img.data))
	}

	goImg := image.NewRGBA(image.Rect(0, 0, img.width, img.height))

	// Convert RGB to RGBA
	srcIdx := 0
	dstIdx := 0
	for y := 0; y < img.height; y++ {
		for x := 0; x < img.width; x++ {
			goImg.Pix[dstIdx] = img.data[srcIdx]     // R
			goImg.Pix[dstIdx+1] = img.data[srcIdx+1] // G
			goImg.Pix[dstIdx+2] = img.data[srcIdx+2] // B
			goImg.Pix[dstIdx+3] = 255                // A (opaque)
			srcIdx += 3
			dstIdx += 4
		}
	}

	return goImg, nil
}

// Equals checks if two images are equal.
//
// Two images are equal if they have the same dimensions, color space,
// bits per component, and data.
func (img *Image) Equals(other *Image) bool {
	if other == nil {
		return false
	}

	if img.width != other.width || img.height != other.height {
		return false
	}

	if img.colorSpace != other.colorSpace {
		return false
	}

	if img.bitsPerComponent != other.bitsPerComponent {
		return false
	}

	if !bytes.Equal(img.data, other.data) {
		return false
	}

	return true
}

// String returns a string representation of the image.
func (img *Image) String() string {
	return fmt.Sprintf("Image{%dx%d, %s, %d bits, filter=%s, size=%d bytes}",
		img.width, img.height, img.colorSpace, img.bitsPerComponent, img.filter, len(img.data))
}

// Domain errors
var (
	// ErrEmptyImageData is returned when image data is empty.
	ErrEmptyImageData = errors.New("image data cannot be empty")

	// ErrInvalidImageDimensions is returned when image dimensions are invalid.
	ErrInvalidImageDimensions = errors.New("image dimensions must be positive")

	// ErrInvalidBitsPerComponent is returned when bits per component is invalid.
	ErrInvalidBitsPerComponent = errors.New("bits per component must be between 1 and 16")
)
