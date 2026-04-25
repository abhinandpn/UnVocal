package setup

import (
	"fmt"
	"os/exec"
	"runtime"
)

// run command and check
func checkCommand(name string, args ...string) bool {
	cmd := exec.Command(name, args...)
	err := cmd.Run()
	return err == nil
}

func Setup() error {
	fmt.Println("=== UnVocal Setup Check ===")

	// 1. OS Info
	fmt.Println("OS:", runtime.GOOS)

	// 2. Check Go
	if checkCommand("go", "version") {
		fmt.Println("Go: ✅ Installed")
	} else {
		fmt.Println("Go: ❌ Not installed")
	}

	// 3. Check Python
	if checkCommand("python3", "--version") {
		fmt.Println("Python: ✅ Installed")
	} else {
		fmt.Println("Python: ❌ Not installed")
		fmt.Println("Install: brew install python@3.11")
	}

	// 4. Check pip
	if checkCommand("pip3", "--version") {
		fmt.Println("pip: ✅ Installed")
	} else {
		fmt.Println("pip: ❌ Not installed")
	}

	// 5. Check demucs
	if checkCommand("python3", "-m", "demucs", "--help") {
		fmt.Println("Demucs: ✅ Installed")
	} else {
		fmt.Println("Demucs: ❌ Not installed")
		fmt.Println("Install: python3 -m pip install demucs")
	}

	// 6. Check ffmpeg
	if checkCommand("ffmpeg", "-version") {
		fmt.Println("FFmpeg: ✅ Installed")
	} else {
		fmt.Println("FFmpeg: ❌ Not installed")
		fmt.Println("Install: brew install ffmpeg")
	}

	fmt.Println("=== Done ===")
	return nil
}
