# GoImage

GoImage is a command-line tool for image compression, inspired by Squoosh and Rimage, but implemented in Go. It supports multiple image formats including PNG, JPEG, WebP, OptiPNG, and MozJPEG.

## Features

- Compress images to various formats (PNG, JPEG, WebP)
- Support for advanced compression algorithms (OptiPNG, MozJPEG)
- **Parallel batch processing** with worker pools (4-8 workers based on CPU cores)
- **Single file and directory processing** modes
- **Recursive directory processing** for nested folder structures
- Simple command-line interface
- Configurable compression quality
- Progress bar for visual feedback during processing
- **Skip existing files** option for incremental processing
- **Automatic worker scaling** based on system capabilities

## Installation

### Prerequisites

- Go 1.18 or higher
- For WebP conversion, the WebP binaries will be automatically downloaded on first use
- For OptiPNG optimization: Install OptiPNG (`optipng` command must be available in PATH)
- For MozJPEG optimization: Install MozJPEG (`cjpeg` command must be available in PATH)

#### Installing OptiPNG and MozJPEG

**Windows:**
- OptiPNG: Download from [http://optipng.sourceforge.net/](http://optipng.sourceforge.net/) or use `choco install optipng`
- MozJPEG: Download from [https://github.com/mozilla/mozjpeg/releases](https://github.com/mozilla/mozjpeg/releases)

**macOS:**
```bash
# Using Homebrew
brew install optipng mozjpeg
```

**Ubuntu/Debian:**
```bash
# OptiPNG
sudo apt install optipng

# MozJPEG (may need to build from source or use alternative repos)
```

**Note:** If these tools are not installed, GoImage will automatically fall back to basic compression methods.

### Building from source

#### Windows (PowerShell)

```powershell
git clone https://github.com/Nirdeo/goimage.git
cd goimage
go build -o goimage.exe .\cmd\goimage
```

#### macOS/Linux (Bash/Zsh)

```bash
git clone https://github.com/Nirdeo/goimage.git
cd goimage
go build -o goimage ./cmd/goimage
```

## Usage

### Windows (PowerShell)

```powershell
# Single file processing
.\goimage.exe --input image.png --output compressed.webp
.\goimage.exe -i image.png -o compressed.webp

# Specify format and quality
.\goimage.exe --input image.jpg --format webp --quality 80
.\goimage.exe -i image.jpg -f webp -q 80

# Use OptiPNG for better PNG compression
.\goimage.exe --input image.png --format optipng
.\goimage.exe -i image.png -f optipng

# Use MozJPEG for better JPEG compression
.\goimage.exe --input image.png --format mozjpeg --quality 85
.\goimage.exe -i image.png -f mozjpeg -q 85

# Batch processing with parallel workers
.\goimage.exe --input-dir .\images --output-dir .\compressed --format webp --workers 6
.\goimage.exe --input-dir .\images --format optipng --workers 4 --quality 90

# Recursive processing with auto worker scaling
.\goimage.exe --input-dir .\photos --format webp --recursive --workers 0 --quality 85

# Skip existing files (incremental processing)
.\goimage.exe --input-dir .\images --format mozjpeg --skip-existing --quality 90

# Show help
.\goimage.exe --help
```

### macOS/Linux (Bash/Zsh)

```bash
# Basic usage
./goimage --input image.png --output compressed.webp
# or using short flags
./goimage -i image.png -o compressed.webp

# Specify format and quality
./goimage --input image.jpg --format webp --quality 80
# or using short flags
./goimage -i image.jpg -f webp -q 80

# Use OptiPNG for better PNG compression
./goimage --input image.png --format optipng
# or using short flags
./goimage -i image.png -f optipng

# Use MozJPEG for better JPEG compression
./goimage --input image.png --format mozjpeg --quality 85
# or using short flags
./goimage -i image.png -f mozjpeg -q 85

# Show help
./goimage --help
```

### Command-line options

#### Single file processing
- `--input`, `-i`: Input image file path
- `--output`, `-o`: Output image file path (if not specified, will use input filename with appropriate extension)

#### Batch processing
- `--input-dir`: Input directory containing images for batch processing
- `--output-dir`: Output directory for batch processing (if not specified, outputs will be placed next to input files)
- `--recursive`: Process images in subdirectories recursively
- `--skip-existing`: Skip processing if output file already exists
- `--workers`, `-w`: Number of worker goroutines (0 = auto: 4-8 based on CPU cores)

#### Common options
- `--format`, `-f`: Output format: png, jpeg, webp, optipng, mozjpeg (default: auto - based on output file extension)
- `--quality`, `-q`: Compression quality (0-100, higher is better quality, default: 75)
- `--help`, `-h`: Show help

## Current Status

This project is fully functional with the following features:
- Basic PNG compression
- Basic JPEG compression
- WebP compression
- OptiPNG integration (requires optipng binary)
- MozJPEG integration (requires mozjpeg cjpeg binary)
- Progress bar for visual feedback during processing
- Automatic fallback to basic compression when external tools are not available

## License

This project is licensed under the MIT License - see the LICENSE file for details.
