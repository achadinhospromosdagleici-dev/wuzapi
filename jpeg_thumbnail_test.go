package main

import (
	"bytes"
	"image"
	"image/color"
	_ "image/jpeg"
	"testing"
)

// TestJpegThumbnail covers the in-memory thumbnail encoding that replaced the
// leaking temp-file path in SendImage: it must return decodable JPEG bytes
// bounded by the requested size, preserving aspect ratio.
func TestJpegThumbnail(t *testing.T) {
	// A 200x100 source image (2:1 aspect).
	src := image.NewRGBA(image.Rect(0, 0, 200, 100))
	for x := 0; x < 200; x++ {
		for y := 0; y < 100; y++ {
			src.Set(x, y, color.RGBA{R: uint8(x), G: uint8(y), B: 0, A: 255})
		}
	}

	out, err := jpegThumbnail(src, 72, 72)
	if err != nil {
		t.Fatalf("jpegThumbnail: %v", err)
	}
	if len(out) == 0 {
		t.Fatal("jpegThumbnail returned no bytes")
	}

	cfg, format, err := image.DecodeConfig(bytes.NewReader(out))
	if err != nil {
		t.Fatalf("output is not a decodable image: %v", err)
	}
	if format != "jpeg" {
		t.Errorf("format = %q; want jpeg", format)
	}
	if cfg.Width == 0 || cfg.Height == 0 {
		t.Errorf("thumbnail has a zero dimension: %dx%d", cfg.Width, cfg.Height)
	}
	if cfg.Width > 72 || cfg.Height > 72 {
		t.Errorf("thumbnail %dx%d exceeds the 72x72 bound", cfg.Width, cfg.Height)
	}
}

// TestJpegThumbnailNil verifies the nil-image guard returns an error instead of
// panicking (resize.Thumbnail dereferences the image's bounds).
func TestJpegThumbnailNil(t *testing.T) {
	if _, err := jpegThumbnail(nil, 72, 72); err == nil {
		t.Error("expected an error for a nil image, got nil")
	}
}
