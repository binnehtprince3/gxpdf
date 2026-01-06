# DCTDecode Filter Example

This example demonstrates how the DCTDecode (JPEG) filter is automatically handled when reading PDF streams.

## Overview

The DCTDecode filter is used in PDF files to compress image data using JPEG compression. When the PDF parser encounters a stream with DCTDecode filter, it automatically:

1. Detects the filter from the stream dictionary
2. Extracts decode parameters (if present)
3. Decodes the JPEG data to raw RGB or grayscale pixels
4. Returns the decoded data for further processing

## Integration

The DCTDecode filter is integrated into the stream decoding pipeline at:
- `internal/infrastructure/parser/reader.go` - Main stream decoder
- `internal/infrastructure/encoding/dct.go` - DCT decoder implementation

### Automatic Decoding

When reading a PDF with JPEG images:

```go
reader := parser.NewReader("document.pdf")
err := reader.Open()
if err != nil {
    log.Fatal(err)
}
defer reader.Close()

// Get image stream object
imageStream := ... // Get from PDF structure

// Decode automatically handles DCTDecode filter
decodedData, err := reader.DecodeStream(imageStream)
if err != nil {
    log.Fatal(err)
}

// decodedData now contains raw RGB/grayscale pixels
```

## Supported Features

### Color Transforms

The decoder supports the ColorTransform parameter:
- `0` = No color transform
- `1` = YCbCr to RGB (default for color images)

Example stream dictionary with ColorTransform:
```
<<
  /Type /XObject
  /Subtype /Image
  /Width 800
  /Height 600
  /ColorSpace /DeviceRGB
  /BitsPerComponent 8
  /Filter /DCTDecode
  /DecodeParms << /ColorTransform 1 >>
>>
```

### Supported Image Types

- RGB images (3 components, 8 bits per component)
- Grayscale images (1 component, 8 bits per component)
- YCbCr images (automatically converted to RGB)

## Output Format

Decoded data format:
- **RGB**: 3 bytes per pixel (R, G, B), width × height × 3 bytes total
- **Grayscale**: 1 byte per pixel, width × height bytes total

## Testing

Comprehensive tests are available in:
- `internal/infrastructure/encoding/dct_test.go` - DCT decoder unit tests
- `internal/infrastructure/parser/stream_test.go` - Stream integration tests

Run tests:
```bash
go test ./internal/infrastructure/encoding/...
go test ./internal/infrastructure/parser/... -run=TestStreamDecoder
```

## Implementation Details

### Filter Detection

The parser automatically detects DCTDecode from the stream's `/Filter` entry:

```go
// Single filter
/Filter /DCTDecode

// Multiple filters (array)
/Filter [ /DCTDecode ]
```

### Parameter Extraction

Decode parameters are extracted from `/DecodeParms`:

```go
if parmsDict, ok := decodeParmsObj.(*primitive.Dictionary); ok {
    if ctObj := parmsDict.Get("ColorTransform"); ctObj != nil {
        colorTransform = ctObj.Value()
    }
}
```

### Decoding Pipeline

1. **Filter Detection**: Extract filter name from stream dictionary
2. **Parameter Extraction**: Get ColorTransform and other parameters
3. **Decoder Creation**: Create DCTDecoder with extracted parameters
4. **JPEG Decoding**: Use Go's standard `image/jpeg` package
5. **Pixel Extraction**: Convert image to raw RGB/grayscale bytes

## Performance

- **Zero-copy**: Minimal memory allocations
- **Fast**: Uses Go's optimized JPEG decoder
- **Memory Efficient**: Streams large images without loading entire PDF

## Future Enhancements

- Support for CMYK color space
- Support for ICC color profiles
- Progressive JPEG support
- Better error recovery for corrupted JPEG data

## References

- PDF 1.7 Specification, Section 7.4.8 (DCTDecode Filter)
- JPEG Standard (ISO/IEC 10918-1)
- Go image/jpeg package documentation
