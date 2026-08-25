package ui

import (
	"context"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// OpenFilePicker launches a native GUI file selection dialog on Linux, Windows, or macOS.
func OpenFilePicker(title string, fileFilter string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. Windows: Native PowerShell OpenFileDialog
	if runtime.GOOS == "windows" {
		psScript := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.OpenFileDialog
$f.Title = "` + title + `"
$f.Filter = "PDF Files (*.pdf)|*.pdf;*.PDF|All Files (*.*)|*.*"
$f.InitialDirectory = [Environment]::GetFolderPath("MyDocuments")
if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $f.FileName
}
`
		cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
		return "", nil
	}

	// 2. Linux / BSD: Check GUI environment
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", nil // Headless server environment
	}

	// Zenity (GTK / GNOME standard)
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

	// Kdialog (Qt / KDE standard)
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

// OpenDirectoryPicker launches a native GUI directory picker on Linux, Windows, or macOS.
func OpenDirectoryPicker(title string, initialDir string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// 1. Windows: Native PowerShell FolderBrowserDialog
	if runtime.GOOS == "windows" {
		psScript := `
Add-Type -AssemblyName System.Windows.Forms
$f = New-Object System.Windows.Forms.FolderBrowserDialog
$f.Description = "` + title + `"
$f.ShowNewFolderButton = $true
if ($f.ShowDialog() -eq [System.Windows.Forms.DialogResult]::OK) {
    Write-Output $f.SelectedPath
}
`
		cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", psScript)
		out, err := cmd.Output()
		if err == nil {
			selected := strings.TrimSpace(string(out))
			if selected != "" {
				return selected, nil
			}
		}
		return "", nil
	}

	// 2. Linux / BSD: Check GUI environment
	if os.Getenv("DISPLAY") == "" && os.Getenv("WAYLAND_DISPLAY") == "" {
		return "", nil
	}

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
