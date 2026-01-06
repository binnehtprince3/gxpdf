package types

import (
	"bytes"
	"image"
	"image/jpeg"
	"os"
	"testing"
)

func TestNewImage(t *testing.T) {
	tests := []struct {
		name             string
		data             []byte
		width            int
		height           int
		colorSpace       string
		bitsPerComponent int
		filter           string
		wantErr          bool
		expectedErrType  error
	}{
		{
			name:             "valid RGB image",
			data:             []byte{255, 0, 0, 0, 255, 0, 0, 0, 255}, // 3 RGB pixels
			width:            3,
			height:           1,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 8,
			filter:           "/FlateDecode",
			wantErr:          false,
		},
		{
			name:             "valid JPEG image",
			data:             createTestJPEG(t),
			width:            10,
			height:           10,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 8,
			filter:           "/DCTDecode",
			wantErr:          false,
		},
		{
			name:             "empty data",
			data:             []byte{},
			width:            10,
			height:           10,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 8,
			filter:           "/DCTDecode",
			wantErr:          true,
			expectedErrType:  ErrEmptyImageData,
		},
		{
			name:             "invalid width",
			data:             []byte{1, 2, 3},
			width:            0,
			height:           10,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 8,
			filter:           "/DCTDecode",
			wantErr:          true,
			expectedErrType:  ErrInvalidImageDimensions,
		},
		{
			name:             "invalid height",
			data:             []byte{1, 2, 3},
			width:            10,
			height:           -1,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 8,
			filter:           "/DCTDecode",
			wantErr:          true,
			expectedErrType:  ErrInvalidImageDimensions,
		},
		{
			name:             "invalid bits per component",
			data:             []byte{1, 2, 3},
			width:            10,
			height:           10,
			colorSpace:       "DeviceRGB",
			bitsPerComponent: 0,
			filter:           "/DCTDecode",
			wantErr:          true,
			expectedErrType:  ErrInvalidBitsPerComponent,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := NewImage(tt.data, tt.width, tt.height, tt.colorSpace, tt.bitsPerComponent, tt.filter)

			if tt.wantErr {
				if err == nil {
					t.Errorf("NewImage() expected error, got nil")
					return
				}
				if tt.expectedErrType != nil && !bytes.Contains([]byte(err.Error()), []byte(tt.expectedErrType.Error())) {
					t.Errorf("NewImage() error = %v, want error containing %v", err, tt.expectedErrType)
				}
				return
			}

			if err != nil {
				t.Errorf("NewImage() unexpected error: %v", err)
				return
			}

			if img == nil {
				t.Error("NewImage() returned nil image")
				return
			}

			if img.Width() != tt.width {
				t.Errorf("Image.Width() = %d, want %d", img.Width(), tt.width)
			}

			if img.Height() != tt.height {
				t.Errorf("Image.Height() = %d, want %d", img.Height(), tt.height)
			}

			if img.ColorSpace() != tt.colorSpace {
				t.Errorf("Image.ColorSpace() = %s, want %s", img.ColorSpace(), tt.colorSpace)
			}

			if img.BitsPerComponent() != tt.bitsPerComponent {
				t.Errorf("Image.BitsPerComponent() = %d, want %d", img.BitsPerComponent(), tt.bitsPerComponent)
			}

			if img.Filter() != tt.filter {
				t.Errorf("Image.Filter() = %s, want %s", img.Filter(), tt.filter)
			}
		})
	}
}

func TestImage_SaveToFile(t *testing.T) {
	// Create test JPEG data
	jpegData := createTestJPEG(t)

	// Create test RGB data
	rgbData := make([]byte, 10*10*3) // 10x10 RGB image
	for i := range rgbData {
		rgbData[i] = byte(i % 256)
	}

	tests := []struct {
		name    string
		img     *Image
		path    string
		wantErr bool
	}{
		{
			name: "save JPEG as JPEG",
			img: &Image{
				data:             jpegData,
				width:            10,
				height:           10,
				colorSpace:       "DeviceRGB",
				bitsPerComponent: 8,
				filter:           "/DCTDecode",
			},
			path:    "test_output.jpg",
			wantErr: false,
		},
		{
			name: "save RGB as PNG",
			img: &Image{
				data:             rgbData,
				width:            10,
				height:           10,
				colorSpace:       "DeviceRGB",
				bitsPerComponent: 8,
				filter:           "/FlateDecode",
			},
			path:    "test_output.png",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up test file if it exists
			defer os.Remove(tt.path)

			err := tt.img.SaveToFile(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.SaveToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				// Verify file was created
				if _, err := os.Stat(tt.path); os.IsNotExist(err) {
					t.Errorf("SaveToFile() did not create file %s", tt.path)
				}
			}
		})
	}
}

func TestImage_ToGoImage(t *testing.T) {
	tests := []struct {
		name       string
		img        *Image
		wantWidth  int
		wantHeight int
		wantErr    bool
	}{
		{
			name: "JPEG to Go image",
			img: &Image{
				data:             createTestJPEG(t),
				width:            10,
				height:           10,
				colorSpace:       "DeviceRGB",
				bitsPerComponent: 8,
				filter:           "/DCTDecode",
			},
			wantWidth:  10,
			wantHeight: 10,
			wantErr:    false,
		},
		{
			name: "RGB to Go image",
			img: &Image{
				data:             make([]byte, 5*5*3), // 5x5 RGB
				width:            5,
				height:           5,
				colorSpace:       "DeviceRGB",
				bitsPerComponent: 8,
				filter:           "/FlateDecode",
			},
			wantWidth:  5,
			wantHeight: 5,
			wantErr:    false,
		},
		{
			name: "Grayscale to Go image",
			img: &Image{
				data:             make([]byte, 8*8), // 8x8 grayscale
				width:            8,
				height:           8,
				colorSpace:       "DeviceGray",
				bitsPerComponent: 8,
				filter:           "/FlateDecode",
			},
			wantWidth:  8,
			wantHeight: 8,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goImg, err := tt.img.ToGoImage()
			if (err != nil) != tt.wantErr {
				t.Errorf("Image.ToGoImage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if goImg == nil {
					t.Error("ToGoImage() returned nil image")
					return
				}

				bounds := goImg.Bounds()
				if bounds.Dx() != tt.wantWidth {
					t.Errorf("ToGoImage() width = %d, want %d", bounds.Dx(), tt.wantWidth)
				}

				if bounds.Dy() != tt.wantHeight {
					t.Errorf("ToGoImage() height = %d, want %d", bounds.Dy(), tt.wantHeight)
				}
			}
		})
	}
}

func TestImage_Equals(t *testing.T) {
	data1 := []byte{1, 2, 3, 4, 5, 6}
	data2 := []byte{1, 2, 3, 4, 5, 6}
	data3 := []byte{7, 8, 9, 10, 11, 12}

	img1, _ := NewImage(data1, 2, 1, "DeviceRGB", 8, "/FlateDecode")
	img2, _ := NewImage(data2, 2, 1, "DeviceRGB", 8, "/FlateDecode")
	img3, _ := NewImage(data3, 2, 1, "DeviceRGB", 8, "/FlateDecode")
	img4, _ := NewImage(data1, 3, 1, "DeviceRGB", 8, "/FlateDecode") // Different dimensions

	tests := []struct {
		name  string
		img1  *Image
		img2  *Image
		equal bool
	}{
		{
			name:  "equal images",
			img1:  img1,
			img2:  img2,
			equal: true,
		},
		{
			name:  "different data",
			img1:  img1,
			img2:  img3,
			equal: false,
		},
		{
			name:  "different dimensions",
			img1:  img1,
			img2:  img4,
			equal: false,
		},
		{
			name:  "nil comparison",
			img1:  img1,
			img2:  nil,
			equal: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.img1.Equals(tt.img2); got != tt.equal {
				t.Errorf("Image.Equals() = %v, want %v", got, tt.equal)
			}
		})
	}
}

func TestImage_String(t *testing.T) {
	img, _ := NewImage([]byte{1, 2, 3}, 1, 1, "DeviceRGB", 8, "/DCTDecode")
	str := img.String()

	if str == "" {
		t.Error("Image.String() returned empty string")
	}

	// Check that the string contains key information
	expectedParts := []string{"1x1", "DeviceRGB", "8 bits", "/DCTDecode", "3 bytes"}
	for _, part := range expectedParts {
		if !bytes.Contains([]byte(str), []byte(part)) {
			t.Errorf("Image.String() = %q, should contain %q", str, part)
		}
	}
}

func TestImage_SetName(t *testing.T) {
	img, _ := NewImage([]byte{1, 2, 3}, 1, 1, "DeviceRGB", 8, "/DCTDecode")

	if img.Name() != "" {
		t.Errorf("Image.Name() = %q, want empty string", img.Name())
	}

	img.SetName("/Im1")

	if img.Name() != "/Im1" {
		t.Errorf("Image.Name() = %q, want /Im1", img.Name())
	}
}

// createTestJPEG creates a small test JPEG image.
func createTestJPEG(t *testing.T) []byte {
	t.Helper()

	// Create a simple 10x10 red image
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			img.Set(x, y, image.NewRGBA64(image.Rect(0, 0, 1, 1)).At(0, 0))
		}
	}

	// Encode to JPEG
	var buf bytes.Buffer
	if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 90}); err != nil {
		t.Fatalf("Failed to create test JPEG: %v", err)
	}

	return buf.Bytes()
}
