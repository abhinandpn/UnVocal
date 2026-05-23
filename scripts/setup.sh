#!/bin/bash

echo "=== UnVocal Setup ==="

# Create venv
python3 -m venv .venv

# Activate venv
source .venv/bin/activate

# Upgrade pip
pip install --upgrade pip

# Install demucs
pip install demucs

# Install ffmpeg (macOS)
brew install ffmpeg

# Create storage folders
mkdir -p storage/input
mkdir -p storage/output
mkdir -p storage/temp
mkdir -p logs

echo "✅ Setup completed"