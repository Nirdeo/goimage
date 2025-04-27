package main

import (
	"fmt"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
	"github.com/victordedomenico/goimage/pkg/processor"
	"os"
	"path/filepath"
)

var (
	// Flags
	inputFile  string
	outputFile string
	format     string
	quality    int
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "goimage",
	Short: "GoImage - Image compression tool",
	Long: `GoImage is a command-line tool for image compression, inspired by Squoosh and Rimage, 
but implemented in Go. It supports multiple image formats including PNG, JPEG, WebP, 
OptiPNG, and MozJPEG.`,
	Example: `  goimage --input image.png --output compressed.webp
  goimage --input image.jpg --format webp --quality 80
  goimage --input image.png --format optipng`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate input file
		if inputFile == "" {
			fmt.Fprintf(os.Stderr, "Error: Input file is required\n")
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

		// Process the image
		err = proc.Process(inputFile, outputFile, quality, bar)
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
	// Define flags
	rootCmd.Flags().StringVarP(&inputFile, "input", "i", "", "Input image file path (required)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output image file path (if not specified, will use input filename with appropriate extension)")
	rootCmd.Flags().StringVarP(&format, "format", "f", "auto", "Output format: png, jpeg, webp, optipng, mozjpeg (default: auto - based on output file extension)")
	rootCmd.Flags().IntVarP(&quality, "quality", "q", 75, "Compression quality (0-100, higher is better quality)")
}

func main() {
	Execute()
}
