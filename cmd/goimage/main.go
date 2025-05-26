package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/nirdeo/goimage/pkg/processor"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

var (
	// Flags
	inputFile    string
	inputDir     string
	outputFile   string
	outputDir    string
	format       string
	quality      int
	workers      int
	recursive    bool
	skipExisting bool
)

// Job represents a single image processing job
type Job struct {
	InputPath  string
	OutputPath string
	Format     string
	Quality    int
}

// Result represents the result of a processing job
type Result struct {
	Job      Job
	Success  bool
	Error    error
	Duration time.Duration
}

// ProgressBarWrapper wraps progressbar.ProgressBar to implement ProgressTracker
type ProgressBarWrapper struct {
	bar *progressbar.ProgressBar
}

// Add implements ProgressTracker interface
func (p *ProgressBarWrapper) Add(n int) error {
	return p.bar.Add(n)
}

// Describe implements ProgressTracker interface
func (p *ProgressBarWrapper) Describe(desc string) {
	p.bar.Describe(desc)
}

// Worker processes jobs from the jobs channel
func worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		start := time.Now()

		// Get the appropriate processor for the format
		proc, err := processor.GetProcessor(job.Format)
		if err != nil {
			results <- Result{
				Job:      job,
				Success:  false,
				Error:    err,
				Duration: time.Since(start),
			}
			continue
		}

		// Create a simple progress tracker (no visual bar for parallel processing)
		bar := &processor.SimpleProgressTracker{}

		// Process the image
		err = proc.Process(job.InputPath, job.OutputPath, job.Quality, bar)

		results <- Result{
			Job:      job,
			Success:  err == nil,
			Error:    err,
			Duration: time.Since(start),
		}
	}
}

// getSupportedImageExtensions returns a list of supported image file extensions
func getSupportedImageExtensions() []string {
	return []string{".jpg", ".jpeg", ".png", ".webp"}
}

// isImageFile checks if a file has a supported image extension
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	for _, supportedExt := range getSupportedImageExtensions() {
		if ext == supportedExt {
			return true
		}
	}
	return false
}

// collectImageFiles collects all image files from input directory
func collectImageFiles(inputDir string, recursive bool) ([]string, error) {
	var imageFiles []string

	if recursive {
		err := filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && isImageFile(info.Name()) {
				imageFiles = append(imageFiles, path)
			}
			return nil
		})
		return imageFiles, err
	} else {
		entries, err := os.ReadDir(inputDir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() && isImageFile(entry.Name()) {
				imageFiles = append(imageFiles, filepath.Join(inputDir, entry.Name()))
			}
		}
		return imageFiles, nil
	}
}

// generateOutputPath generates the output path for a given input file
func generateOutputPath(inputPath, outputDir, format string) string {
	inputBase := filepath.Base(inputPath)
	inputExt := filepath.Ext(inputBase)
	inputName := inputBase[:len(inputBase)-len(inputExt)]

	var ext string
	switch format {
	case "jpeg", "jpg", "mozjpeg":
		ext = ".jpg"
	case "webp":
		ext = ".webp"
	case "png", "optipng":
		ext = ".png"
	default:
		ext = ".png"
	}

	if outputDir != "" {
		return filepath.Join(outputDir, inputName+"_compressed"+ext)
	}

	// If no output directory, place in same directory as input
	inputDir := filepath.Dir(inputPath)
	return filepath.Join(inputDir, inputName+"_compressed"+ext)
}

// processBatch handles batch processing with worker pool
func processBatch(imageFiles []string, outputDir, format string, quality, numWorkers int) {
	// Determine optimal number of workers
	if numWorkers <= 0 {
		numWorkers = runtime.NumCPU()
		if numWorkers > 8 {
			numWorkers = 8 // Cap at 8 workers
		} else if numWorkers < 4 {
			numWorkers = 4 // Minimum 4 workers
		}
	}

	fmt.Printf("Processing %d images with %d workers...\n", len(imageFiles), numWorkers)

	// Create channels
	jobs := make(chan Job, len(imageFiles))
	results := make(chan Result, len(imageFiles))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i+1, jobs, results, &wg)
	}

	// Create output directory if specified
	if outputDir != "" {
		if err := os.MkdirAll(outputDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating output directory: %v\n", err)
			return
		}
	}

	// Send jobs
	go func() {
		defer close(jobs)
		for _, inputPath := range imageFiles {
			outputPath := generateOutputPath(inputPath, outputDir, format)

			// Skip if output already exists and skipExisting is true
			if skipExisting {
				if _, err := os.Stat(outputPath); err == nil {
					fmt.Printf("Skipping %s (output already exists)\n", inputPath)
					continue
				}
			}

			jobs <- Job{
				InputPath:  inputPath,
				OutputPath: outputPath,
				Format:     format,
				Quality:    quality,
			}
		}
	}()

	// Wait for workers to finish
	go func() {
		wg.Wait()
		close(results)
	}()

	// Create overall progress bar
	bar := progressbar.NewOptions(len(imageFiles),
		progressbar.OptionSetDescription("Processing images..."),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Collect results
	var successful, failed int
	var totalDuration time.Duration

	for result := range results {
		bar.Add(1)
		totalDuration += result.Duration

		if result.Success {
			successful++
			fmt.Printf("✓ %s -> %s (%.2fs)\n",
				filepath.Base(result.Job.InputPath),
				filepath.Base(result.Job.OutputPath),
				result.Duration.Seconds())
		} else {
			failed++
			fmt.Printf("✗ %s: %v\n",
				filepath.Base(result.Job.InputPath),
				result.Error)
		}
	}

	fmt.Printf("\nBatch processing completed!\n")
	fmt.Printf("Successful: %d\n", successful)
	fmt.Printf("Failed: %d\n", failed)
	fmt.Printf("Total time: %.2fs\n", totalDuration.Seconds())
	fmt.Printf("Average time per image: %.2fs\n", totalDuration.Seconds()/float64(successful+failed))
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goimage",
	Short: "GoImage - Image compression tool",
	Long: `GoImage is a command-line tool for image compression, inspired by Squoosh and Rimage, 
but implemented in Go. It supports multiple image formats including PNG, JPEG, WebP, 
OptiPNG, and MozJPEG. Now supports batch processing with worker pools for parallel processing.`,
	Example: `  # Single file processing
  goimage --input image.png --output compressed.webp
  goimage --input image.jpg --format webp --quality 80
  
  # Batch processing
  goimage --input-dir ./images --output-dir ./compressed --format webp --workers 6
  goimage --input-dir ./images --format optipng --recursive --skip-existing`,
	Run: func(cmd *cobra.Command, args []string) {
		// Check if batch processing is requested
		if inputDir != "" {
			// Batch processing mode
			if _, err := os.Stat(inputDir); os.IsNotExist(err) {
				fmt.Fprintf(os.Stderr, "Error: Input directory '%s' does not exist\n", inputDir)
				os.Exit(1)
			}

			// Collect image files
			imageFiles, err := collectImageFiles(inputDir, recursive)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error collecting image files: %v\n", err)
				os.Exit(1)
			}

			if len(imageFiles) == 0 {
				fmt.Fprintf(os.Stderr, "No image files found in directory '%s'\n", inputDir)
				os.Exit(1)
			}

			// Determine format if auto
			if format == "auto" {
				format = "webp" // Default for batch processing
			}

			// Validate quality
			if quality < 0 || quality > 100 {
				fmt.Fprintf(os.Stderr, "Error: Quality must be between 0 and 100\n")
				os.Exit(1)
			}

			// Process batch
			processBatch(imageFiles, outputDir, format, quality, workers)
			return
		}

		// Single file processing mode
		if inputFile == "" {
			fmt.Fprintf(os.Stderr, "Error: Either --input or --input-dir is required\n")
			cmd.Help()
			os.Exit(1)
		}

		if _, err := os.Stat(inputFile); os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "Error: Input file '%s' does not exist\n", inputFile)
			os.Exit(1)
		}

		// Set default output file if not specified
		if outputFile == "" {
			ext := ".png"
			switch format {
			case "jpeg", "jpg", "mozjpeg":
				ext = ".jpg"
			case "webp":
				ext = ".webp"
			case "png", "optipng":
				ext = ".png"
			}

			inputBase := filepath.Base(inputFile)
			inputExt := filepath.Ext(inputBase)
			inputName := inputBase[:len(inputBase)-len(inputExt)]
			outputFile = inputName + "_compressed" + ext
		}

		// Determine format from output file extension if set to auto
		if format == "auto" {
			ext := filepath.Ext(outputFile)
			switch ext {
			case ".jpg", ".jpeg":
				format = "jpeg"
			case ".webp":
				format = "webp"
			case ".png":
				format = "png"
			default:
				fmt.Fprintf(os.Stderr, "Error: Could not determine format from output file extension '%s'\n", ext)
				os.Exit(1)
			}
		}

		// Validate quality
		if quality < 0 || quality > 100 {
			fmt.Fprintf(os.Stderr, "Error: Quality must be between 0 and 100\n")
			os.Exit(1)
		}

		fmt.Printf("Processing image:\n")
		fmt.Printf("  Input: %s\n", inputFile)
		fmt.Printf("  Output: %s\n", outputFile)
		fmt.Printf("  Format: %s\n", format)
		fmt.Printf("  Quality: %d\n", quality)

		// Get the appropriate processor for the format
		proc, err := processor.GetProcessor(format)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		// Create a progress bar
		bar := progressbar.NewOptions(100,
			progressbar.OptionSetDescription("Processing..."),
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(50),
			progressbar.OptionShowCount(),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		// Create a wrapper for the progress bar to implement ProgressTracker
		progressWrapper := &ProgressBarWrapper{bar: bar}

		// Process the image
		err = proc.Process(inputFile, outputFile, quality, progressWrapper)
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nError processing image: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\nImage processed successfully!")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define flags for single file processing
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file path")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file path (if not specified, will use input filename with appropriate extension)")

	// Define flags for batch processing
	rootCmd.Flags().StringVar(&inputDir, "input-dir", "", "Input directory containing images for batch processing")
	rootCmd.Flags().StringVar(&outputDir, "output-dir", "", "Output directory for batch processing (if not specified, outputs will be placed next to input files)")
	rootCmd.Flags().BoolVar(&recursive, "recursive", false, "Process images in subdirectories recursively")
	rootCmd.Flags().BoolVar(&skipExisting, "skip-existing", false, "Skip processing if output file already exists")

	// Common flags
	rootCmd.Flags().StringVarP(&format, "format", "f", "auto", "Output format: png, jpeg, webp, optipng, mozjpeg (default: auto - based on output file extension)")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 75, "Compression quality (0-100, higher is better quality)")
	rootCmd.Flags().IntVarP(&workers, "workers", "w", 0, "Number of worker goroutines (0 = auto: 4-8 based on CPU cores)")
}

func main() {
	Execute()
}
