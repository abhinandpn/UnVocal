package handlers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func HandleKaraoke(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseMultipartForm(500 << 20); err != nil {
		http.Error(w, "Error parsing form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if header.Size > 500<<20 {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	jobID := generateID()

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Error getting working directory", http.StatusInternalServerError)
		return
	}

	tmpInput := filepath.Join(cwd, "tmp_input")
	tmpOutput := filepath.Join(cwd, "tmp_output", jobID)

	if err := os.MkdirAll(tmpInput, os.ModePerm); err != nil {
		http.Error(w, "Error creating input directory", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(tmpOutput, os.ModePerm); err != nil {
		http.Error(w, "Error creating output directory", http.StatusInternalServerError)
		return
	}

	// Create unique filename
	safeFilename := jobID + "_" + header.Filename

	inputFilePath := filepath.Join(
		tmpInput,
		safeFilename,
	)

	out, err := os.Create(inputFilePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}

	defer os.Remove(inputFilePath)
	defer os.RemoveAll(tmpOutput)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		30*time.Minute,
	)
	defer cancel()

	pythonPath := filepath.Join(cwd, "venv", "bin", "python")

	if runtime.GOOS == "windows" {
		pythonPath = filepath.Join(
			cwd,
			"venv",
			"Scripts",
			"python.exe",
		)
	}

	cmd := exec.CommandContext(
		ctx,
		pythonPath,
		"-m",
		"demucs",
		"--two-stems=vocals",
		"-o",
		tmpOutput,
		inputFilePath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		http.Error(w, "Error processing audio", http.StatusInternalServerError)
		return
	}

	originalName := strings.TrimSuffix(
		safeFilename,
		filepath.Ext(safeFilename),
	)

	targetAudio := filepath.Join(
		tmpOutput,
		"htdemucs",
		originalName,
		"no_vocals.wav",
	)

	if _, err := os.Stat(targetAudio); os.IsNotExist(err) {
		http.Error(w, "Processed file not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set(
		"Content-Disposition",
		"attachment; filename=\"karaoke.wav\"",
	)

	http.ServeFile(w, r, targetAudio)
}