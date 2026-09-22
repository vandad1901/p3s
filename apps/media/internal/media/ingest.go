package media

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"io"
	"os"

	"image/gif"
	"image/jpeg"
	"image/png"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"golang.org/x/image/draw"
	_ "golang.org/x/image/webp"
)

const (
	maxDimension  = 10_000
	maxPixelCount = 40_000_000
)

const (
	smallWidth  = 320
	mediumWidth = 640
	largeWidth  = 1280
)

const (
	imageTypeJPEG = "jpeg"
	imageTypePNG  = "png"
	imageTypeGIF  = "gif"
	imageTypeWEBP = "webp"
)

const (
	jpegQuality = 80
)

const s3Bucket = "p3s-upload-bucket"

func (s *Service) HandleMedia(ctx context.Context, key string) error {
	queuedAt, err := s.takeLease(ctx, key)
	if err != nil {
		return fmt.Errorf("creating media: %w", err)
	}

	res, err := s.s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s3Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("getting media from S3: %w", err)
	}

	derivatives, err := IngestMedia(ctx, res.Body, *res.ContentLength)
	if err != nil {
		return err
	}

	err = res.Body.Close()
	if err != nil {
		return fmt.Errorf("closing media body: %w", err)
	}

	for _, derivative := range derivatives {
		_, err = derivative.file.Seek(0, io.SeekStart)
		if err != nil {
			return fmt.Errorf("seeking derivative file: %w", err)
		}

		_, err = s.s3Client.PutObject(ctx, &s3.PutObjectInput{
			Bucket: aws.String(s3Bucket),
			Key:    aws.String(fmt.Sprintf("%s.%s", key, derivative.Extension)),
			Body:   derivative.file,
		})
		if err != nil {
			return fmt.Errorf("putting derivative to S3: %w", err)
		}
	}

	_, err = s.changeStatus(ctx, key, queuedAt, MediaIngestStatusIngested)
	if err != nil {
		return fmt.Errorf("changing media status: %w", err)
	}

	return nil
}

func checkContext(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err() //nolint:wrapcheck
	default:
		return nil
	}
}

func IngestMedia(ctx context.Context, file io.Reader, contentLength int64,
) ([]derivative, error) {
	tempFile, err := createTempFile(ctx, file, contentLength)
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}

	img, formatStr, err := parseAndValidateImage(ctx, tempFile)
	if err != nil {
		return nil, err
	}

	err = tempFile.Close()
	if err != nil {
		return nil, fmt.Errorf("closing temp file: %w", err)
	}

	err = os.Remove(tempFile.Name())
	if err != nil {
		return nil, fmt.Errorf("removing temp file: %w", err)
	}

	derivatives, err := createDerivatives(ctx, img, formatStr)
	if err != nil {
		return nil, fmt.Errorf("deriving images: %w", err)
	}

	return derivatives, nil
}

func createTempFile(ctx context.Context, file io.Reader, contentLength int64) (*os.File, error) {
	err := checkContext(ctx)
	if err != nil {
		return nil, err
	}

	tempFile, err := os.CreateTemp("", "media-*")
	if err != nil {
		return nil, fmt.Errorf("creating temp file: %w", err)
	}

	_, err = io.CopyN(tempFile, file, contentLength)
	if err != nil {
		return nil, fmt.Errorf("filling temp file: %w", err)
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("rewinding temp file: %w", err)
	}

	return tempFile, nil
}

func parseAndValidateImage(ctx context.Context, tempFile *os.File) (image.Image, string, error) {
	err := checkContext(ctx)
	if err != nil {
		return nil, "", err
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, "", fmt.Errorf("seeking temp file: %w", err)
	}

	img, _, err := image.Decode(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("decoding image: %w", err)
	}

	_, err = tempFile.Seek(0, io.SeekStart)
	if err != nil {
		return nil, "", fmt.Errorf("seeking temp file: %w", err)
	}

	cfg, format, err := image.DecodeConfig(tempFile)
	if err != nil {
		return nil, "", fmt.Errorf("decoding image config: %w", err)
	}

	err = validateImage(cfg, format)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil
}

func validateImage(cfg image.Config, format string) error {
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return errInvalidDimension
	}

	if cfg.Width > maxDimension || cfg.Height > maxDimension {
		return errInvalidDimension
	}

	if cfg.Width*cfg.Height > maxPixelCount {
		return errInvalidDimension
	}

	switch format {
	case imageTypeJPEG, imageTypePNG, imageTypeGIF, imageTypeWEBP:
	default:
		return errInvalidFormat
	}

	return nil
}

func createDerivatives(ctx context.Context, img image.Image, formatStr string) ([]derivative, error) {
	err := checkContext(ctx)
	if err != nil {
		return nil, err
	}

	derivatives := []derivative{}
	srcWidth := img.Bounds().Dx()

	sizes := []int{smallWidth, mediumWidth, largeWidth}
	for _, width := range sizes {
		resizedImg := img
		outWidth := srcWidth

		if srcWidth > width {
			resizedImg = resizeImage(img, width)
			outWidth = width
		}

		encodedImage, ext, err := encodeForStorage(resizedImg, formatStr)
		if err != nil {
			return nil, fmt.Errorf("encoding for storage: %w", err)
		}

		derivatives = append(derivatives, derivative{
			file:      bytes.NewReader(encodedImage.Bytes()),
			Extension: fmt.Sprintf("%d.%s", outWidth, ext),
		})

		if srcWidth <= width {
			break
		}
	}

	return derivatives, nil
}

func encodeForStorage(resizedImg image.Image, formatStr string) (bytes.Buffer, string, error) {
	var (
		encodedImage bytes.Buffer
	)

	switch formatStr {
	case imageTypePNG:
		err := png.Encode(&encodedImage, resizedImg)
		if err != nil {
			return bytes.Buffer{}, "", fmt.Errorf("encoding to png: %w", err)
		}

		return encodedImage, "png", nil
	case imageTypeGIF:
		err := gif.Encode(&encodedImage, resizedImg, nil)
		if err != nil {
			return bytes.Buffer{}, "", fmt.Errorf("encoding to gif: %w", err)
		}

		return encodedImage, "gif", nil
	default:
		err := jpeg.Encode(&encodedImage, resizedImg, &jpeg.Options{Quality: jpegQuality})
		if err != nil {
			return bytes.Buffer{}, "", fmt.Errorf("encoding to jpeg: %w", err)
		}

		return encodedImage, "jpg", nil
	}
}

func resizeImage(src image.Image, width int) image.Image {
	srcBounds := src.Bounds()
	height := max(srcBounds.Dy()*width/srcBounds.Dx(), 1)
	dst := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.BiLinear.Scale(dst, dst.Bounds(), src, srcBounds, draw.Over, nil)

	return dst
}
