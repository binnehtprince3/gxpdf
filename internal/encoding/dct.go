// Package encoding implements PDF stream encoding and decoding filters.
package encoding

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
)

// DCTDecoder implements DCTDecode (JPEG) stream decompression.
//
// DCTDecode is used for JPEG-compressed image data in PDF files.
// The decoder converts JPEG data to raw RGB/Gray pixel data.
//
// Reference: PDF 1.7 specification, Section 7.4.8 (DCTDecode Filter).
type DCTDecoder struct {
	// ColorTransform specifies color transformation:
	// 0 = no transform
	// 1 = YCbCr to RGB (default for RGB images)
	ColorTransform int
}

// DCTResult contains decoded image data and metadata.
type DCTResult struct {
	// Data contains raw pixel data in RGB (3 bytes per pixel) or Gray (1 byte per pixel).
	Data []byte

	// Width is the image width in pixels.
	Width int

	// Height is the image height in pixels.
	Height int

	// Components is the number of color components (1 for grayscale, 3 for RGB).
	Components int

	// BitsPerComponent is always 8 for JPEG.
	BitsPerComponent int
}

// NewDCTDecoder creates a new DCT (JPEG) decoder.
func NewDCTDecoder() *DCTDecoder {
	return &DCTDecoder{
		ColorTransform: 1, // Default: YCbCr to RGB
	}
}

// NewDCTDecoderWithParams creates a DCT decoder with specific parameters.
func NewDCTDecoderWithParams(colorTransform int) *DCTDecoder {
	return &DCTDecoder{
		ColorTransform: colorTransform,
	}
}

// Decode decompresses JPEG-encoded data to raw pixels.
//
// Returns raw RGB or grayscale pixel data.
// For RGB images: 3 bytes per pixel (R, G, B).
// For grayscale images: 1 byte per pixel.
//
// Parameters:
//   - data: JPEG-compressed image data
//
// Returns: Raw pixel data bytes, or error if decoding fails.
func (d *DCTDecoder) Decode(data []byte) ([]byte, error) {
	result, err := d.DecodeWithMetadata(data)
	if err != nil {
		return nil, err
	}
	return result.Data, nil
}

// DecodeWithMetadata decodes JPEG data and returns both pixel data and metadata.
//
// This is useful when you need to know the image dimensions and color space.
func (d *DCTDecoder) DecodeWithMetadata(data []byte) (*DCTResult, error) {
	// Decode JPEG using Go's standard library.
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode JPEG: %w", err)
	}

	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	// Determine color model and extract pixels.
	switch img := img.(type) {
	case *image.Gray:
		return d.extractGray(img, width, height)
	case *image.YCbCr:
		return d.extractRGB(img, width, height)
	case *image.RGBA:
		return d.extractFromRGBA(img, width, height)
	case *image.NRGBA:
		return d.extractFromNRGBA(img, width, height)
	default:
		// Fallback: convert any image to RGB.
		return d.extractGeneric(img, width, height)
	}
}

// extractGray extracts pixel data from a grayscale image.
func (d *DCTDecoder) extractGray(img *image.Gray, width, height int) (*DCTResult, error) {
	// For grayscale, we can use the pixel data directly.
	data := make([]byte, width*height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			data[y*width+x] = img.GrayAt(x, y).Y
		}
	}

	return &DCTResult{
		Data:             data,
		Width:            width,
		Height:           height,
		Components:       1,
		BitsPerComponent: 8,
	}, nil
}

// extractRGB extracts RGB pixel data from a YCbCr image.
func (d *DCTDecoder) extractRGB(img *image.YCbCr, width, height int) (*DCTResult, error) {
	data := make([]byte, width*height*3)
	idx := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[idx] = byte(r >> 8)
			data[idx+1] = byte(g >> 8)
			data[idx+2] = byte(b >> 8)
			idx += 3
		}
	}

	return &DCTResult{
		Data:             data,
		Width:            width,
		Height:           height,
		Components:       3,
		BitsPerComponent: 8,
	}, nil
}

// extractFromRGBA extracts RGB pixel data from an RGBA image.
func (d *DCTDecoder) extractFromRGBA(img *image.RGBA, width, height int) (*DCTResult, error) {
	data := make([]byte, width*height*3)
	idx := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.RGBAAt(x, y)
			data[idx] = c.R
			data[idx+1] = c.G
			data[idx+2] = c.B
			idx += 3
		}
	}

	return &DCTResult{
		Data:             data,
		Width:            width,
		Height:           height,
		Components:       3,
		BitsPerComponent: 8,
	}, nil
}

// extractFromNRGBA extracts RGB pixel data from an NRGBA image.
func (d *DCTDecoder) extractFromNRGBA(img *image.NRGBA, width, height int) (*DCTResult, error) {
	data := make([]byte, width*height*3)
	idx := 0

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			c := img.NRGBAAt(x, y)
			data[idx] = c.R
			data[idx+1] = c.G
			data[idx+2] = c.B
			idx += 3
		}
	}

	return &DCTResult{
		Data:             data,
		Width:            width,
		Height:           height,
		Components:       3,
		BitsPerComponent: 8,
	}, nil
}

// extractGeneric extracts RGB pixel data from any image type.
func (d *DCTDecoder) extractGeneric(img image.Image, width, height int) (*DCTResult, error) {
	data := make([]byte, width*height*3)
	idx := 0

	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			data[idx] = byte(r >> 8)
			data[idx+1] = byte(g >> 8)
			data[idx+2] = byte(b >> 8)
			idx += 3
		}
	}

	return &DCTResult{
		Data:             data,
		Width:            width,
		Height:           height,
		Components:       3,
		BitsPerComponent: 8,
	}, nil
}

// DecodeToImage decodes JPEG data to a Go image.Image.
//
// This is useful when you need the image in Go's standard format
// for further processing or saving in a different format.
func (d *DCTDecoder) DecodeToImage(data []byte) (image.Image, error) {
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		return nil, fmt.Errorf("failed to decode JPEG: %w", err)
	}
	return img, nil
}

// Encode compresses raw pixel data to JPEG format.
//
// This enables creating JPEG streams for PDF writing.
//
// Parameters:
//   - data: Raw RGB pixel data (3 bytes per pixel)
//   - width: Image width in pixels
//   - height: Image height in pixels
//   - quality: JPEG quality (1-100, 0 for default 75)
//
// Returns: JPEG-compressed data, or error if encoding fails.
func (d *DCTDecoder) Encode(data []byte, width, height, quality int) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 75 // Default quality
	}

	// Validate data length.
	expectedLen := width * height * 3
	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length: expected %d bytes for %dx%d RGB image, got %d",
			expectedLen, width, height, len(data))
	}

	// Create image from raw RGB data.
	img := image.NewRGBA(image.Rect(0, 0, width, height))
	idx := 0
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: data[idx],
				G: data[idx+1],
				B: data[idx+2],
				A: 255,
			})
			idx += 3
		}
	}

	// Encode to JPEG.
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return buf.Bytes(), nil
}

// EncodeGray compresses grayscale pixel data to JPEG format.
//
// Parameters:
//   - data: Raw grayscale pixel data (1 byte per pixel)
//   - width: Image width in pixels
//   - height: Image height in pixels
//   - quality: JPEG quality (1-100, 0 for default 75)
//
// Returns: JPEG-compressed data, or error if encoding fails.
func (d *DCTDecoder) EncodeGray(data []byte, width, height, quality int) ([]byte, error) {
	if quality <= 0 || quality > 100 {
		quality = 75 // Default quality
	}

	// Validate data length.
	expectedLen := width * height
	if len(data) != expectedLen {
		return nil, fmt.Errorf("invalid data length: expected %d bytes for %dx%d grayscale image, got %d",
			expectedLen, width, height, len(data))
	}

	// Create grayscale image from raw data.
	img := image.NewGray(image.Rect(0, 0, width, height))
	copy(img.Pix, data)

	// Encode to JPEG.
	var buf bytes.Buffer
	opts := &jpeg.Options{Quality: quality}
	if err := jpeg.Encode(&buf, img, opts); err != nil {
		return nil, fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return buf.Bytes(), nil
}
