# UnVocal

UnVocal is a Go-based tool that automates the process of extracting vocals from video files using the Demucs music source separation library.

## Features

- Automated setup verification for required dependencies
- Vocal extraction from MP4 video files
- Uses Demucs for high-quality source separation
- Simple command-line interface

## Prerequisites

Before running UnVocal, ensure you have the following installed:

- **Go** (version 1.25.0 or later)
- **Python 3.11** (with pip)
- **Demucs** (music source separation library)
- **FFmpeg** (for audio/video processing)

### Installation Instructions

#### macOS (using Homebrew)
```bash
# Install Python 3.11
brew install python@3.11

# Install FFmpeg
brew install ffmpeg

# Install Demucs
python3 -m pip install demucs
```

#### Other Platforms
Please refer to the official documentation for installing Python, FFmpeg, and Demucs on your platform.

## Installation

1. Clone the repository:
```bash
git clone https://github.com/abhinandpn/UnVocal.git
cd UnVocal
```

2. Install Go dependencies:
```bash
go mod tidy
```

## Usage

1. Place your MP4 video file in the `mp4/` directory (e.g., `mp4/videoplayback.mp4`)

2. Run the application:
```bash
go run main.go
```

The program will:
- Check that all required dependencies are installed
- Run Demucs to extract vocals from the video
- Save the separated audio files to the `output/` directory

## Project Structure

```
UnVocal/
├── main.go           # Main application entry point
├── setup/
│   └── setup.go      # Dependency checking and setup
├── mp4/              # Input video files directory
├── output/           # Output directory for separated audio
├── go.mod            # Go module file
└── README.md         # This file
```

## How It Works

1. **Setup Check**: The application verifies that Go, Python, pip, Demucs, and FFmpeg are installed
2. **Vocal Separation**: Uses Demucs with the `--two-stems=vocals` option to separate vocals from the input video
3. **Output**: Saves the extracted vocals and other stems to the output directory

## Dependencies

- [Demucs](https://github.com/facebookresearch/demucs) - Music source separation library
- [FFmpeg](https://ffmpeg.org/) - Audio/video processing tool
- Go standard library

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is open source. Please check the license file for details.

## Support

If you encounter any issues, please open an issue on the GitHub repository.
