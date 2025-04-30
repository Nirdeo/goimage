package processor

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"

	"github.com/nickalie/go-webpbin"
	"github.com/schollz/progressbar/v3"

	// Register WebP decoder
	_ "golang.org/x/image/webp"
)

// Initialize webpbin library
func init() {
	// Check for unsupported platforms
	webpbin.DetectUnsupportedPlatforms()
}

// Processor defines the interface for image processing
type Processor interface {
	Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error
}

// GetProcessor returns the appropriate processor based on the format
func GetProcessor(format string) (Processor, error) {
	switch strings.ToLower(format) {
	case "png":
		return &PNGProcessor{}, nil
	case "optipng":
		return &OptiPNGProcessor{}, nil
	case "jpeg", "jpg":
		return &JPEGProcessor{}, nil
	case "mozjpeg":
		return &MozJPEGProcessor{}, nil
	case "webp":
		return &WebPProcessor{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// loadImage loads an image from the given path
func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// PNGProcessor handles PNG image processing
type PNGProcessor struct{}

// Process implements the Processor interface for PNG
func (p *PNGProcessor) Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error {
	// Load the input image
	bar.Describe("Loading image...")
	img, err := loadImage(inputPath)
	if err != nil {
		return err
	}
	bar.Add(50) // 50% progress after loading

	// Create the output file
	bar.Describe("Creating output file...")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	bar.Add(10) // 60% progress after creating file

	// Encode as PNG
	bar.Describe("Encoding as PNG...")
	encoder := png.Encoder{
		CompressionLevel: png.BestCompression,
	}
	err = encoder.Encode(outFile, img)
	if err != nil {
		return err
	}
	bar.Add(40) // 100% progress after encoding
	return nil
}

// OptiPNGProcessor handles OptiPNG image processing
type OptiPNGProcessor struct{}

// Process implements the Processor interface for OptiPNG
func (p *OptiPNGProcessor) Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error {
	// First, use the standard PNG processor to create a PNG
	bar.Describe("Applying basic PNG compression...")
	pngProcessor := &PNGProcessor{}
	if err := pngProcessor.Process(inputPath, outputPath, quality, bar); err != nil {
		return err
	}

	// TODO: Implement OptiPNG optimization
	// This would typically involve calling the optipng binary
	// For now, we'll just return a message
	bar.Describe("OptiPNG optimization not implemented")
	return errors.New("OptiPNG optimization not yet implemented - basic PNG compression applied")
}

// JPEGProcessor handles JPEG image processing
type JPEGProcessor struct{}

// Process implements the Processor interface for JPEG
func (p *JPEGProcessor) Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error {
	// Load the input image
	bar.Describe("Loading image...")
	img, err := loadImage(inputPath)
	if err != nil {
		return err
	}
	bar.Add(50) // 50% progress after loading

	// Create the output file
	bar.Describe("Creating output file...")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	bar.Add(10) // 60% progress after creating file

	// Encode as JPEG
	bar.Describe("Encoding as JPEG...")
	options := jpeg.Options{
		Quality: quality,
	}
	err = jpeg.Encode(outFile, img, &options)
	if err != nil {
		return err
	}
	bar.Add(40) // 100% progress after encoding
	return nil
}

// MozJPEGProcessor handles MozJPEG image processing
type MozJPEGProcessor struct{}

// Process implements the Processor interface for MozJPEG
func (p *MozJPEGProcessor) Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error {
	// First, use the standard JPEG processor to create a JPEG
	bar.Describe("Applying basic JPEG compression...")
	jpegProcessor := &JPEGProcessor{}
	if err := jpegProcessor.Process(inputPath, outputPath, quality, bar); err != nil {
		return err
	}

	// TODO: Implement MozJPEG optimization
	// This would typically involve calling the mozjpeg binary
	// For now, we'll just return a message
	bar.Describe("MozJPEG optimization not implemented")
	return errors.New("MozJPEG optimization not yet implemented - basic JPEG compression applied")
}

// WebPProcessor handles WebP image processing
type WebPProcessor struct{}

// Process implements the Processor interface for WebP
func (p *WebPProcessor) Process(inputPath, outputPath string, quality int, bar *progressbar.ProgressBar) error {

	bar.Describe("Setting up WebP encoding...")
	// Setup CWebP encoder
	cwebp := webpbin.NewCWebP()
	cwebp.InputFile(inputPath)
	cwebp.OutputFile(outputPath)
	cwebp.Quality(uint(quality))

	bar.Add(20) // 20% progress after setup

	// Show progress
	bar.Describe("Encoding to WebP...")

	// Run the encoder
	err := cwebp.Run()
	if err != nil {
		return fmt.Errorf("failed to encode image as WebP: %v", err)
	}

	bar.Add(80) // 100% progress after encoding
	return nil
}
