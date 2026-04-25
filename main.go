package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/abhinandpn/UnVocal/setup"
)

func main() {

	if err := setup.Setup(); err != nil {
		fmt.Println("❌ Setup failed:")
		fmt.Println(err)
		return
	}

	fmt.Println("✅ Setup completed successfully")

	start := time.Now() // ⏱ start time

	input := "/Users/abhinanpn/Desktop/UnVocal/mp4/videoplayback.mp4"
	output := "/Users/abhinanpn/Desktop/UnVocal/output"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(
		ctx,
		"/Users/abhinanpn/Desktop/UnVocal/venv/bin/python",
		"-m", "demucs",
		"--two-stems=vocals",
		"-o", output,
		input,
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("Running Demucs...")

	err := cmd.Run()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	elapsed := time.Since(start) // ⏱ end time

	fmt.Println("Done!")
	fmt.Println("Execution time :", elapsed)
}
