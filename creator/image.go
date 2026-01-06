package creator

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // Import JPEG decoder
	_ "image/png"  // Import PNG decoder
	"io"
	"os"

	"github.com/coregx/gxpdf/internal/encoding"
)

// Image represents an image that can be embedded in a PDF document.
//
// Currently supports:
//   - JPEG images (RGB and CMYK color spaces)
//   - PNG images (RGB, RGBA, grayscale, paletted)
//
// The image data is stored as:
//   - JPEG: Raw JPEG bytes (DCTDecode)
//   - PNG: Raw pixel data compressed with FlateDecode
//
// For RGBA PNG with transparency, the alpha channel is stored separately
// as an SMask (soft mask) for proper PDF rendering.
//
// Example:
//
//	img, err := creator.LoadImage("photo.jpg")  // or "photo.png"
//	if err != nil {
//	    return err
//	}
//	page.DrawImage(img, 100, 500, 200, 150)
type Image struct {
	// Image format (jpeg or png).
	format string

	// Raw image data (JPEG bytes or compressed PNG pixels).
	data []byte

	// Alpha mask data for RGBA PNG (compressed with FlateDecode).
	alphaMask []byte

	// Image dimensions.
	width  int
	height int

	// Color space information.
	colorSpace ColorSpace

	// Components per pixel (1=Gray, 3=RGB, 4=CMYK).
	components int

	// Bits per component (8 for most images).
	bitsPerComponent int
}

// ColorSpace represents the image color space.
type ColorSpace string

const (
	// ColorSpaceRGB is RGB color space (3 components).
	ColorSpaceRGB ColorSpace = "DeviceRGB"

	// ColorSpaceCMYK is CMYK color space (4 components).
	ColorSpaceCMYK ColorSpace = "DeviceCMYK"

	// ColorSpaceGray is grayscale (1 component).
	ColorSpaceGray ColorSpace = "DeviceGray"
)

// LoadImage loads an image from a file.
//
// Supported formats: JPEG, PNG.
// For JPEG: RGB and CMYK color spaces.
// For PNG: RGB, RGBA (with alpha mask), grayscale, paletted.
//
// Example:
//
//	img, err := creator.LoadImage("photo.jpg")  // or "photo.png"
//	if err != nil {
//	    return err
//	}
func LoadImage(path string) (*Image, error) {
	//nolint:gosec // File path is provided by user, G304 false positive.
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer func() {
		_ = file.Close() // Best effort cleanup.
	}()

	return LoadImageFromReader(file)
}

// LoadImageFromReader loads an image from an io.Reader.
//
// Supported formats: JPEG, PNG.
// This allows loading images from various sources (files, HTTP responses, etc.).
//
// Example:
//
//	resp, _ := http.Get("https://example.com/image.png")
//	defer resp.Body.Close()
//	img, _ := creator.LoadImageFromReader(resp.Body)
func LoadImageFromReader(r io.Reader) (*Image, error) {
	// Read all data first.
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("failed to read image data: %w", err)
	}

	// Detect format.
	format := detectImageFormat(data)
	if format == "" {
		return nil, ErrUnsupportedImageFormat
	}

	// Load based on format.
	switch format {
	case "jpeg":
		return loadJPEG(data)
	case "png":
		return loadPNG(data)
	default:
		return nil, fmt.Errorf("unsupported image format: %s", format)
	}
}

// detectImageFormat detects the image format by checking file header.
func detectImageFormat(data []byte) string {
	if len(data) < 8 {
		return ""
	}

	// Check JPEG signature (0xFF 0xD8 0xFF).
	if data[0] == 0xFF && data[1] == 0xD8 && data[2] == 0xFF {
		return "jpeg"
	}

	// Check PNG signature (0x89 PNG \r\n 0x1A \n).
	if data[0] == 0x89 && data[1] == 'P' && data[2] == 'N' && data[3] == 'G' {
		return "png"
	}

	return ""
}

// loadJPEG loads a JPEG image from raw data.
func loadJPEG(data []byte) (*Image, error) {
	// Decode config to get dimensions.
	cfg, _, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode JPEG: %w", err)
	}

	return &Image{
		format:           "jpeg",
		data:             data,
		width:            cfg.Width,
		height:           cfg.Height,
		colorSpace:       ColorSpaceRGB, // JPEG defaults to RGB.
		components:       3,
		bitsPerComponent: 8,
	}, nil
}

// loadPNG loads a PNG image from raw data.
func loadPNG(data []byte) (*Image, error) {
	// Decode the full PNG image.
	img, err := decodePNGImage(data)
	if err != nil {
		return nil, err
	}

	// Convert PNG to raw pixel data.
	return convertPNGToImage(img)
}

// decodePNGImage decodes PNG data to an image.Image.
func decodePNGImage(data []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode PNG: %w", err)
	}
	return img, nil
}

// convertPNGToImage converts a PNG image to our Image structure.
func convertPNGToImage(img image.Image) (*Image, error) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Detect color model and convert accordingly.
	switch img.ColorModel() {
	case color.RGBAModel:
		return convertRGBAPNG(img, width, height)
	case color.NRGBAModel:
		return convertRGBAPNG(img, width, height)
	case color.GrayModel:
		return convertGrayPNG(img, width, height)
	default:
		// For paletted and other formats, convert to RGB.
		return convertGenericPNG(img, width, height)
	}
}

// convertRGBAPNG converts an RGBA PNG image (with alpha channel).
func convertRGBAPNG(img image.Image, width, height int) (*Image, error) {
	// Extract RGB and alpha channels separately.
	rgbData, alphaData := extractRGBAndAlpha(img, width, height)

	// Compress both with FlateDecode.
	compressedRGB, err := compressData(rgbData)
	if err != nil {
		return nil, fmt.Errorf("failed to compress RGB data: %w", err)
	}

	var compressedAlpha []byte
	if alphaData != nil {
		compressedAlpha, err = compressData(alphaData)
		if err != nil {
			return nil, fmt.Errorf("failed to compress alpha data: %w", err)
		}
	}

	return &Image{
		format:           "png",
		data:             compressedRGB,
		alphaMask:        compressedAlpha,
		width:            width,
		height:           height,
		colorSpace:       ColorSpaceRGB,
		components:       3,
		bitsPerComponent: 8,
	}, nil
}

// convertGrayPNG converts a grayscale PNG image.
func convertGrayPNG(img image.Image, width, height int) (*Image, error) {
	// Extract grayscale data.
	grayData := extractGrayscale(img, width, height)

	// Compress with FlateDecode.
	compressed, err := compressData(grayData)
	if err != nil {
		return nil, fmt.Errorf("failed to compress grayscale data: %w", err)
	}

	return &Image{
		format:           "png",
		data:             compressed,
		width:            width,
		height:           height,
		colorSpace:       ColorSpaceGray,
		components:       1,
		bitsPerComponent: 8,
	}, nil
}

// convertGenericPNG converts paletted and other PNG formats to RGB.
func convertGenericPNG(img image.Image, width, height int) (*Image, error) {
	// Convert to RGB.
	rgbData := extractRGB(img, width, height)

	// Compress with FlateDecode.
	compressed, err := compressData(rgbData)
	if err != nil {
		return nil, fmt.Errorf("failed to compress RGB data: %w", err)
	}

	return &Image{
		format:           "png",
		data:             compressed,
		width:            width,
		height:           height,
		colorSpace:       ColorSpaceRGB,
		components:       3,
		bitsPerComponent: 8,
	}, nil
}

// extractRGBAndAlpha extracts RGB and alpha from RGBA image.
func extractRGBAndAlpha(img image.Image, width, height int) ([]byte, []byte) {
	rgbData := make([]byte, width*height*3)
	alphaData := make([]byte, width*height)
	hasAlpha := false

	idx := 0
	alphaIdx := 0
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			rgbData[idx] = byte(r >> 8)
			rgbData[idx+1] = byte(g >> 8)
			rgbData[idx+2] = byte(b >> 8)
			idx += 3

			alphaValue := byte(a >> 8)
			alphaData[alphaIdx] = alphaValue
			alphaIdx++

			if alphaValue != 255 {
				hasAlpha = true
			}
		}
	}

	// Return nil for alpha if fully opaque.
	if !hasAlpha {
		return rgbData, nil
	}
	return rgbData, alphaData
}

// extractRGB extracts RGB data from any image format.
func extractRGB(img image.Image, width, height int) []byte {
	rgbData := make([]byte, width*height*3)
	idx := 0
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			rgbData[idx] = byte(r >> 8)
			rgbData[idx+1] = byte(g >> 8)
			rgbData[idx+2] = byte(b >> 8)
			idx += 3
		}
	}

	return rgbData
}

// extractGrayscale extracts grayscale data from a grayscale image.
func extractGrayscale(img image.Image, width, height int) []byte {
	grayData := make([]byte, width*height)
	idx := 0
	bounds := img.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			gray := color.GrayModel.Convert(img.At(x, y)).(color.Gray)
			grayData[idx] = gray.Y
			idx++
		}
	}

	return grayData
}

// compressData compresses data using FlateDecode.
func compressData(data []byte) ([]byte, error) {
	encoder := encoding.NewFlateDecoder()
	compressed, err := encoder.Encode(data)
	if err != nil {
		return nil, fmt.Errorf("flate compression failed: %w", err)
	}
	return compressed, nil
}

// Width returns the image width in pixels.
func (img *Image) Width() int {
	return img.width
}

// Height returns the image height in pixels.
func (img *Image) Height() int {
	return img.height
}

// Data returns the raw JPEG data.
//
// This is used internally by the PDF writer to embed the image.
func (img *Image) Data() []byte {
	return img.data
}

// ColorSpace returns the image color space.
func (img *Image) ColorSpace() ColorSpace {
	return img.colorSpace
}

// Format returns the image format (jpeg or png).
func (img *Image) Format() string {
	return img.format
}

// AlphaMask returns the alpha mask data (nil if no transparency).
//
// For RGBA PNG images with transparency, this contains the compressed
// alpha channel data used for the SMask (soft mask) in PDF.
func (img *Image) AlphaMask() []byte {
	return img.alphaMask
}

// HasAlpha returns true if the image has transparency data.
func (img *Image) HasAlpha() bool {
	return img.alphaMask != nil
}

// Components returns the number of color components.
//
// Returns:
//   - 1 for grayscale
//   - 3 for RGB
//   - 4 for CMYK
func (img *Image) Components() int {
	return img.components
}

// BitsPerComponent returns the bits per component (typically 8).
func (img *Image) BitsPerComponent() int {
	return img.bitsPerComponent
}

// DrawImage draws an image at the specified position and size.
//
// The image is scaled to fit the specified width and height.
// No aspect ratio preservation - the image is stretched to fit.
//
// Parameters:
//   - img: The image to draw
//   - x: Horizontal position in points (from left edge)
//   - y: Vertical position in points (from bottom edge)
//   - width: Display width in points
//   - height: Display height in points
//
// Example:
//
//	img, _ := creator.LoadImage("photo.jpg")
//	page.DrawImage(img, 100, 500, 200, 150)
func (p *Page) DrawImage(img *Image, x, y, width, height float64) error {
	// Validate dimensions.
	if width <= 0 || height <= 0 {
		return errors.New("image dimensions must be positive")
	}

	// Store image operation.
	p.graphicsOps = append(p.graphicsOps, GraphicsOperation{
		Type:   GraphicsOpImage,
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Image:  img,
	})

	return nil
}

// DrawImageFit draws an image scaled to fit within the specified dimensions.
//
// The image is scaled to fit within the width/height while maintaining
// its aspect ratio. The image is centered in the available space.
//
// Parameters:
//   - img: The image to draw
//   - x: Horizontal position in points (from left edge)
//   - y: Vertical position in points (from bottom edge)
//   - maxWidth: Maximum width in points
//   - maxHeight: Maximum height in points
//
// Example:
//
//	img, _ := creator.LoadImage("photo.jpg")
//	page.DrawImageFit(img, 100, 500, 200, 200)  // Fit in 200x200 box
func (p *Page) DrawImageFit(img *Image, x, y, maxWidth, maxHeight float64) error {
	// Validate dimensions.
	if maxWidth <= 0 || maxHeight <= 0 {
		return errors.New("image max dimensions must be positive")
	}

	// Calculate scaled dimensions.
	scaledW, scaledH := calculateFitDimensions(
		float64(img.width),
		float64(img.height),
		maxWidth,
		maxHeight,
	)

	// Center the image in available space.
	centerX := x + (maxWidth-scaledW)/2
	centerY := y + (maxHeight-scaledH)/2

	return p.DrawImage(img, centerX, centerY, scaledW, scaledH)
}

// calculateFitDimensions calculates dimensions to fit within max bounds.
//
// Maintains aspect ratio by scaling down the larger dimension.
func calculateFitDimensions(imgW, imgH, maxW, maxH float64) (float64, float64) {
	// Calculate scale factors for width and height.
	scaleW := maxW / imgW
	scaleH := maxH / imgH

	// Use the smaller scale factor (to fit within both constraints).
	scale := scaleW
	if scaleH < scaleW {
		scale = scaleH
	}

	return imgW * scale, imgH * scale
}

// Errors.
var (
	// ErrUnsupportedImageFormat is returned for unsupported image formats.
	ErrUnsupportedImageFormat = errors.New("unsupported image format (supported: JPEG, PNG)")

	// ErrInvalidImageDimensions is returned for zero/negative dimensions.
	ErrInvalidImageDimensions = errors.New("image dimensions must be positive")
)
