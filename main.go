package main

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/abhinandpn/UnVocal/setup"
)

func generateID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func main() {
	if err := setup.Setup(); err != nil {
		log.Fatalf("❌ Setup failed: %v", err)
	}
	log.Println("✅ Setup completed successfully")

	http.HandleFunc("/api/separate-audio", handleSeparateAudio)

	port := "8080"
	log.Printf("Server listening on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}

func handleSeparateAudio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form (max 500MB)
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

	// Validate size
	if header.Size > 500<<20 {
		http.Error(w, "File too large", http.StatusBadRequest)
		return
	}

	// Track type
	track := r.FormValue("track")
	if track != "no_vocals" {
		track = "vocals"
	}

	jobID := generateID()

	cwd, err := os.Getwd()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	tmpInput := filepath.Join(cwd, "tmp_input")
	tmpOutput := filepath.Join(cwd, "tmp_output", jobID)

	os.MkdirAll(tmpInput, os.ModePerm)
	os.MkdirAll(tmpOutput, os.ModePerm)

	// Use original filename (IMPORTANT)
	inputFilePath := filepath.Join(tmpInput, header.Filename)

	out, err := os.Create(inputFilePath)
	if err != nil {
		http.Error(w, "Error saving file", http.StatusInternalServerError)
		return
	}

	if _, err := io.Copy(out, file); err != nil {
		out.Close()
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
	out.Close()

	defer os.Remove(inputFilePath)
	defer os.RemoveAll(tmpOutput)

	log.Printf("Processing job %s...", jobID)
	start := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// Cross-platform python path
	pythonPath := filepath.Join(cwd, "venv", "bin", "python")
	if runtime.GOOS == "windows" {
		pythonPath = filepath.Join(cwd, "venv", "Scripts", "python.exe")
	}

	cmd := exec.CommandContext(
		ctx,
		pythonPath,
		"-m", "demucs",
		"--two-stems=vocals",
		"-o", tmpOutput,
		inputFilePath,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Demucs error: %v", err)
		http.Error(w, "Error processing audio", http.StatusInternalServerError)
		return
	}

	elapsed := time.Since(start)
	log.Printf("Job %s completed in %v", jobID, elapsed)

	// Correct output path
	originalName := strings.TrimSuffix(header.Filename, filepath.Ext(header.Filename))
	targetAudio := filepath.Join(tmpOutput, "htdemucs", originalName, track+".wav")

	if _, err := os.Stat(targetAudio); os.IsNotExist(err) {
		log.Printf("Output not found: %s", targetAudio)
		http.Error(w, "Processed file not found", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "audio/wav")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s.wav\"", track))

	http.ServeFile(w, r, targetAudio)
}
