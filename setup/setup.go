package system

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

	fmt.Println("OS:", runtime.GOOS)

	checks := []struct {
		Name string
		Cmd  string
		Args []string
	}{
		{"Go", "go", []string{"version"}},
		{"Python", "python3", []string{"--version"}},
		{"pip", "pip3", []string{"--version"}},
		{"FFmpeg", "ffmpeg", []string{"-version"}},
	}

	for _, check := range checks {
		if err := checkCommand(check.Cmd, check.Args...); err != nil {
			return fmt.Errorf("%s not installed", check.Name)
		}

		fmt.Printf("%s: ✅ Installed\n", check.Name)
	}

	fmt.Println("=== Runtime OK ===")

	return nil
}