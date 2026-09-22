package media_test

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/binary"
	"hash/crc32"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"testing"

	"github.com/vandad1901/p3s/apps/media/internal/media"
)

//go:embed testdata/valid.webp
var validWebP []byte

func blankImage(w, h int) *image.RGBA {
	return image.NewRGBA(image.Rect(0, 0, w, h))
}

func encodeJPEG(t *testing.T, w, h int) []byte {
	t.Helper()

	var buf bytes.Buffer

	err := jpeg.Encode(&buf, blankImage(w, h), &jpeg.Options{Quality: 80})
	if err != nil {
		t.Fatalf("encode jpeg: %v", err)
	}

	return buf.Bytes()
}

func encodePNG(t *testing.T, w, h int) []byte {
	t.Helper()

	var buf bytes.Buffer

	err := png.Encode(&buf, blankImage(w, h))
	if err != nil {
		t.Fatalf("encode png: %v", err)
	}

	return buf.Bytes()
}

func bombPNG(t *testing.T, claimedW, claimedH uint32) []byte {
	t.Helper()

	data := encodePNG(t, 1000, 1000)

	binary.BigEndian.PutUint32(data[16:20], claimedW)
	binary.BigEndian.PutUint32(data[20:24], claimedH)

	crc := crc32.ChecksumIEEE(data[12:29])
	binary.BigEndian.PutUint32(data[29:33], crc)

	return data
}

func encodeGIF(t *testing.T, w, h int) []byte {
	t.Helper()

	var buf bytes.Buffer

	err := gif.Encode(&buf, blankImage(w, h), nil)
	if err != nil {
		t.Fatalf("encode gif: %v", err)
	}

	return buf.Bytes()
}

func truncate(b []byte, n int) []byte {
	if n >= len(b) {
		return nil
	}

	return b[:len(b)-n]
}

func TestIngestMedia(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		data    []byte
		wantErr bool
	}{
		{"valid jpeg", encodeJPEG(t, 100, 100), false},
		{"valid png", encodePNG(t, 100, 100), false},
		{"valid gif", encodeGIF(t, 100, 100), false},
		{"valid webp", validWebP, false},

		{"at width limit", encodePNG(t, 10000, 1), false},
		{"over width limit", encodePNG(t, 10001, 1), true},

		{"at height limit", encodePNG(t, 1, 10000), false},
		{"over height limit", encodePNG(t, 1, 10001), true},

		{"at pixel count limit", encodePNG(t, 8000, 5000), false},
		{"over pixel count limit", encodePNG(t, 6351, 6301), true},

		{"png bomb", bombPNG(t, 40000, 40000), true},

		{"truncated jpeg 1", truncate(encodeJPEG(t, 200, 200), 50), true},
		{"truncated jpeg 2", truncate(encodeJPEG(t, 200, 200), 500), true},

		{"truncated PNG 1", truncate(encodePNG(t, 200, 200), 50), true},
		{"truncated PNG 2", truncate(encodePNG(t, 200, 200), 500), true},

		{"truncated GIF 1", truncate(encodeGIF(t, 200, 200), 50), true},
		{"truncated GIF 2", truncate(encodeGIF(t, 200, 200), 500), true},

		{"truncated WEBP 1", truncate(validWebP, 50), true},
		{"truncated WEBP 2", truncate(validWebP, 500), true},

		{"non-image bytes", []byte("this is not an image, just text"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := media.IngestMedia(context.Background(), bytes.NewReader(tt.data), int64(len(tt.data)))
			if (err != nil) != tt.wantErr {
				t.Errorf("IngestMedia() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
