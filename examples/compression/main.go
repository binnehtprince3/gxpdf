// Package main demonstrates PDF stream compression using FlateDecode.
package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/coregx/gxpdf/internal/writer"
)

func main() {
	fmt.Println("=== PDF Stream Compression Demo ===")

	// Example 1: Basic compression
	basicCompressionExample()

	// Example 2: Different compression levels
	compressionLevelsExample()

	// Example 3: Compression ratio estimation
	compressionRatioExample()
}

func basicCompressionExample() {
	fmt.Println("1. Basic Compression")
	fmt.Println("--------------------")

	// Create sample PDF content stream
	content := []byte("BT /F1 12 Tf 100 700 Td (Hello World) Tj ET\n")
	fmt.Printf("Original size: %d bytes\n", len(content))

	// Compress using default compression
	compressed, err := writer.CompressStream(content, writer.DefaultCompression)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Compressed size: %d bytes\n", len(compressed))

	// Decompress to verify
	decompressed, err := writer.DecompressStream(compressed)
	if err != nil {
		log.Fatal(err)
	}

	// Verify round-trip
	if bytes.Equal(decompressed, content) {
		fmt.Println("✓ Round-trip successful!")
	} else {
		fmt.Println("✗ Round-trip failed!")
	}

	fmt.Println()
}

func compressionLevelsExample() {
	fmt.Println("2. Compression Levels Comparison")
	fmt.Println("---------------------------------")

	// Create larger content for better comparison
	var content []byte
	for i := 0; i < 100; i++ {
		content = append(content, []byte("BT /F1 12 Tf 100 700 Td (Lorem ipsum dolor sit amet) Tj ET\n")...)
	}

	fmt.Printf("Original size: %d bytes\n\n", len(content))

	levels := []struct {
		name  string
		level writer.CompressionLevel
	}{
		{"No Compression", writer.NoCompression},
		{"Best Speed", writer.BestSpeed},
		{"Default", writer.DefaultCompression},
		{"Best Compression", writer.BestCompression},
	}

	for _, l := range levels {
		compressed, err := writer.CompressStream(content, l.level)
		if err != nil {
			log.Fatal(err)
		}

		ratio := float64(len(compressed)) / float64(len(content)) * 100
		fmt.Printf("%s:\n", l.name)
		fmt.Printf("  Size: %d bytes (%.1f%% of original)\n", len(compressed), ratio)
		fmt.Printf("  Savings: %d bytes\n", len(content)-len(compressed))
		fmt.Println()
	}
}

func compressionRatioExample() {
	fmt.Println("3. Compression Ratio Estimation")
	fmt.Println("--------------------------------")

	samples := []struct {
		name string
		data []byte
	}{
		{
			name: "Highly repetitive (all 'A's)",
			data: []byte(string(make([]byte, 1000)) + "AAAAAAAAAA"),
		},
		{
			name: "PDF operators (typical content)",
			data: []byte("BT /F1 12 Tf 100 700 Td (Hello World) Tj ET\n" +
				"BT /F1 12 Tf 100 680 Td (Second line) Tj ET\n" +
				"BT /F1 12 Tf 100 660 Td (Third line) Tj ET\n"),
		},
		{
			name: "Mixed text",
			data: []byte("Lorem ipsum dolor sit amet, consectetur adipiscing elit."),
		},
	}

	for _, s := range samples {
		ratio := writer.EstimateCompressionRatio(s.data)
		shouldCompress := writer.ShouldCompress(s.data)

		fmt.Printf("%s:\n", s.name)
		fmt.Printf("  Size: %d bytes\n", len(s.data))
		fmt.Printf("  Estimated ratio: %.1f%%\n", ratio*100)
		fmt.Printf("  Recommendation: ")
		if shouldCompress {
			fmt.Printf("✓ Compress (saves ~%.0f bytes)\n", float64(len(s.data))*(1-ratio))
		} else {
			fmt.Println("✗ Skip compression (too small)")
		}
		fmt.Println()
	}
}
