package setup

import (
	"fmt"
	"os/exec"
	"runtime"
)

func checkCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	return cmd.Run()
}

func Setup() error {
	fmt.Println("=== UnVocal Runtime Check ===")

	// OS
	fmt.Println("OS:", runtime.GOOS)

	// Detect venv python path
	pythonPath := ".venv/bin/python"

	if runtime.GOOS == "windows" {
		pythonPath = ".venv/Scripts/python.exe"
	}

	// Runtime dependency checks
	checks := []struct {
		Name string
		Cmd  string
		Args []string
	}{
		{
			Name: "Go",
			Cmd:  "go",
			Args: []string{"version"},
		},
		{
			Name: "Python",
			Cmd:  pythonPath,
			Args: []string{"--version"},
		},
		{
			Name: "FFmpeg",
			Cmd:  "ffmpeg",
			Args: []string{"-version"},
		},
		{
			Name: "Demucs",
			Cmd:  pythonPath,
			Args: []string{"-c", "import demucs"},
		},
	}

	// Execute checks
	for _, check := range checks {
		if err := checkCommand(check.Cmd, check.Args...); err != nil {
			return fmt.Errorf("%s not installed or not working", check.Name)
		}

		fmt.Printf("%s: ✅ Installed\n", check.Name)
	}

	fmt.Println("=== Runtime OK ===")

	return nil
}
