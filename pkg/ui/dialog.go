package ui

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// OpenFilePicker launches a native GUI file selection dialog (zenity / kdialog).
// Returns the selected absolute path or empty string if cancelled/unavailable.
func OpenFilePicker(title string, fileFilter string) (string, error) {
	// Check if DISPLAY or WAYLAND_DISPLAY is set (GUI environment)
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", nil // Headless / terminal-only
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. Try Zenity (GNOME / standard GTK)
	if _, err := exec.LookPath("zenity"); err == nil {
		args := []string{
			"--file-selection",
			"--title=" + title,
		}
		if fileFilter != "" {
			args = append(args, "--file-filter="+fileFilter)
		}
		cmd := exec.CommandContext(ctx, "zenity", args...)
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
	}

	// 2. Try Kdialog (KDE / Qt)
	if _, err := exec.LookPath("kdialog"); err == nil {
		cmd := exec.CommandContext(ctx, "kdialog", "--getopenfilename", ".", "*.pdf *.PDF")
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
	}

	return "", nil
}

// OpenDirectoryPicker launches a native GUI folder selection dialog.
func OpenDirectoryPicker(title string, initialDir string) (string, error) {
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. Try Zenity
	if _, err := exec.LookPath("zenity"); err == nil {
		args := []string{
			"--file-selection",
			"--directory",
			"--title=" + title,
		}
		if initialDir != "" {
			args = append(args, "--filename="+initialDir+"/")
		}
		cmd := exec.CommandContext(ctx, "zenity", args...)
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
	}

	// 2. Try Kdialog
	if _, err := exec.LookPath("kdialog"); err == nil {
		cmd := exec.CommandContext(ctx, "kdialog", "--getexistingdirectory", initialDir)
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
	}

	return "", nil
}
