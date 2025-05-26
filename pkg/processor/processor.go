package processor

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/nickalie/go-webpbin"

	// Register WebP decoder
	_ "golang.org/x/image/webp"
)

// Initialize webpbin library
func init() {
	// Check for unsupported platforms
	webpbin.DetectUnsupportedPlatforms()
}

// ProgressTracker defines the interface for progress tracking
type ProgressTracker interface {
	Add(int) error
	Describe(string)
}

// Processor defines the interface for image processing
type Processor interface {
	Process(inputPath, outputPath string, quality int, bar ProgressTracker) error
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

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// PNGProcessor handles PNG image processing
type PNGProcessor struct{}

// Process implements the Processor interface for PNG
func (p *PNGProcessor) Process(inputPath, outputPath string, quality int, bar ProgressTracker) error {
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
func (p *OptiPNGProcessor) Process(inputPath, outputPath string, quality int, bar ProgressTracker) error {
	bar.Describe("Loading image...")
	// Load the input image first to validate it
	_, err := loadImage(inputPath)
	if err != nil {
		return err
	}
	bar.Add(20) // 20% progress after loading

	// Create a temporary PNG file if input is not PNG
	tempPNG := ""
	inputExt := strings.ToLower(filepath.Ext(inputPath))
	if inputExt != ".png" {
		bar.Describe("Converting to PNG format...")
		// Create temporary PNG file
		tempPNG = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + "_temp.png"
		pngProcessor := &PNGProcessor{}
		if err := pngProcessor.Process(inputPath, tempPNG, quality, bar); err != nil {
			return err
		}
		inputPath = tempPNG
		defer func() {
			if tempPNG != "" {
				os.Remove(tempPNG)
			}
		}()
	} else {
		bar.Add(30) // Skip conversion step
	}

	bar.Describe("Optimizing with OptiPNG...")

	// Check if optipng is available
	_, err = exec.LookPath("optipng")
	if err != nil {
		// Fallback to basic PNG processing
		bar.Describe("OptiPNG not found, using basic PNG compression...")
		if tempPNG == "" {
			pngProcessor := &PNGProcessor{}
			if err := pngProcessor.Process(inputPath, outputPath, quality, bar); err != nil {
				return err
			}
		} else {
			// Copy temp PNG to output
			if err := copyFile(tempPNG, outputPath); err != nil {
				return err
			}
			bar.Add(50)
		}
		return nil
	}

	// Determine optimization level based on quality
	optimizationLevel := "2" // Default level
	if quality >= 90 {
		optimizationLevel = "7" // Maximum optimization
	} else if quality >= 70 {
		optimizationLevel = "5"
	} else if quality >= 50 {
		optimizationLevel = "3"
	}

	// Run optipng command
	args := []string{
		"-o" + optimizationLevel,
		"-quiet",
		"-out", outputPath,
		inputPath,
	}

	cmd := exec.Command("optipng", args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("optipng optimization failed: %v", err)
	}

	bar.Add(50) // 100% progress after optimization
	bar.Describe("OptiPNG optimization completed")
	return nil
}

// JPEGProcessor handles JPEG image processing
type JPEGProcessor struct{}

// Process implements the Processor interface for JPEG
func (p *JPEGProcessor) Process(inputPath, outputPath string, quality int, bar ProgressTracker) error {
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
func (p *MozJPEGProcessor) Process(inputPath, outputPath string, quality int, bar ProgressTracker) error {
	bar.Describe("Loading image...")
	// Load the input image first to validate it
	img, err := loadImage(inputPath)
	if err != nil {
		return err
	}
	bar.Add(20) // 20% progress after loading

	// Create a temporary JPEG file if input is not JPEG
	tempJPEG := ""
	inputExt := strings.ToLower(filepath.Ext(inputPath))
	if inputExt != ".jpg" && inputExt != ".jpeg" {
		bar.Describe("Converting to JPEG format...")
		// Create temporary JPEG file
		tempJPEG = strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + "_temp.jpg"
		jpegProcessor := &JPEGProcessor{}
		if err := jpegProcessor.Process(inputPath, tempJPEG, quality, bar); err != nil {
			return err
		}
		inputPath = tempJPEG
		defer func() {
			if tempJPEG != "" {
				os.Remove(tempJPEG)
			}
		}()
	} else {
		bar.Add(30) // Skip conversion step
	}

	bar.Describe("Optimizing with MozJPEG...")

	// Check if cjpeg (mozjpeg) is available
	cjpegPath, err := exec.LookPath("cjpeg")
	if err != nil {
		// Try alternative name
		cjpegPath, err = exec.LookPath("mozjpeg-cjpeg")
		if err != nil {
			// Fallback to basic JPEG processing
			bar.Describe("MozJPEG not found, using basic JPEG compression...")
			if tempJPEG == "" {
				jpegProcessor := &JPEGProcessor{}
				if err := jpegProcessor.Process(inputPath, outputPath, quality, bar); err != nil {
					return err
				}
			} else {
				// Copy temp JPEG to output
				if err := copyFile(tempJPEG, outputPath); err != nil {
					return err
				}
				bar.Add(50)
			}
			return nil
		}
	}

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Prepare cjpeg arguments
	args := []string{
		"-quality", fmt.Sprintf("%d", quality),
		"-optimize",
		"-progressive",
		"-outfile", outputPath,
	}

	// If we have the original image in memory and it's not a JPEG, encode directly
	if tempJPEG == "" && (inputExt != ".jpg" && inputExt != ".jpeg") {
		// Create a temporary input file for cjpeg (it needs a file)
		tempInput := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + "_input_temp.ppm"
		defer os.Remove(tempInput)

		// Save as PPM format for cjpeg input
		ppmFile, err := os.Create(tempInput)
		if err != nil {
			return err
		}

		// Write PPM header and data
		bounds := img.Bounds()
		fmt.Fprintf(ppmFile, "P6\n%d %d\n255\n", bounds.Dx(), bounds.Dy())

		// Write pixel data
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := img.At(x, y).RGBA()
				ppmFile.Write([]byte{byte(r >> 8), byte(g >> 8), byte(b >> 8)})
			}
		}
		ppmFile.Close()

		args = append(args, tempInput)
	} else {
		args = append(args, inputPath)
	}

	// Run cjpeg command
	cmd := exec.Command(cjpegPath, args...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("mozjpeg optimization failed: %v", err)
	}

	bar.Add(50) // 100% progress after optimization
	bar.Describe("MozJPEG optimization completed")
	return nil
}

// WebPProcessor handles WebP image processing
type WebPProcessor struct{}

// Process implements the Processor interface for WebP
func (p *WebPProcessor) Process(inputPath, outputPath string, quality int, bar ProgressTracker) error {

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

// SimpleProgressTracker is a minimal progress tracker for parallel processing
type SimpleProgressTracker struct {
	current     int
	description string
}

// Add simulates progress addition (interface compatibility)
func (s *SimpleProgressTracker) Add(n int) error {
	s.current += n
	return nil
}

// Describe sets the current description (interface compatibility)
func (s *SimpleProgressTracker) Describe(desc string) {
	s.description = desc
}
