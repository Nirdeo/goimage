# GoImage

GoImage is a command-line tool for image compression, inspired by Squoosh and Rimage, but implemented in Go. It supports multiple image formats including PNG, JPEG, WebP, OptiPNG, and MozJPEG.

## Features

- Compress images to various formats (PNG, JPEG, WebP)
- Support for advanced compression algorithms (OptiPNG, MozJPEG)
- Simple command-line interface
- Configurable compression quality
- Progress bar for visual feedback during processing

## Installation

### Prerequisites

- Go 1.18 or higher

### Building from source

```bash
git clone https://github.com/victordedomenico/goimage.git
cd goimage
go build -o goimage ./cmd/goimage
```

## Usage

```bash
# Basic usage
goimage --input image.png --output compressed.webp
# or using short flags
goimage -i image.png -o compressed.webp

# Specify format and quality
goimage --input image.jpg --format webp --quality 80
# or using short flags
goimage -i image.jpg -f webp -q 80

# Use OptiPNG for better PNG compression
goimage --input image.png --format optipng
# or using short flags
goimage -i image.png -f optipng

# Use MozJPEG for better JPEG compression
goimage --input image.png --format mozjpeg --quality 85
# or using short flags
goimage -i image.png -f mozjpeg -q 85

# Show help
goimage --help
```

### Command-line options

- `--input`, `-i`: Input image file path (required)
- `--output`, `-o`: Output image file path (if not specified, will use input filename with appropriate extension)
- `--format`, `-f`: Output format: png, jpeg, webp, optipng, mozjpeg (default: auto - based on output file extension)
- `--quality`, `-q`: Compression quality (0-100, higher is better quality, default: 75)
- `--help`, `-h`: Show help

## Current Status

This project is under development. Currently implemented features:
- Basic PNG compression
- Basic JPEG compression
- Progress bar for visual feedback during processing

Coming soon:
- WebP support
- OptiPNG integration
- MozJPEG integration

## License

This project is licensed under the MIT License - see the LICENSE file for details.
